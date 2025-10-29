package handler

import (
	"fmt"
	"net/http"

	friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go-sdk"
	"github.com/rs/zerolog/log"

	"github.com/davidramiro/go-form-mailer/internal/service"
	"github.com/go-faster/errors"
)

type FormHandler struct {
	mailService service.MailService
	frcClient   friendlycaptcha.Client
}

func NewFormHandler(mailService *service.MailService, frcClient friendlycaptcha.Client) (*FormHandler, error) {
	if mailService == nil {
		return nil, errors.New("no MailService passed to FormHandler")
	}

	return &FormHandler{
		mailService: *mailService,
		frcClient:   frcClient,
	}, nil
}

func (f *FormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		f.respond(w, errors.New("Error parsing form"))
		return
	}

	req := service.MailRequest{
		Name:               r.Form.Get("name"),
		Email:              r.Form.Get("email"),
		Message:            r.Form.Get("message"),
		Subject:            r.Form.Get("subject"),
		FrcCaptchaSolution: r.Form.Get("frc-captcha-solution"),
	}

	if err := req.Validate(); err != nil {
		f.respond(w, errors.New("Validation failed: "+err.Error()))
		return
	}

	log.Info().Interface("MailRequest", req).Msg("incoming form MailRequest")

	solution := req.FrcCaptchaSolution
	shouldAccept, err := f.frcClient.CheckCaptchaSolution(r.Context(), solution)
	if err != nil {
		log.Error().Err(err).Msg("captcha check error")
		f.respond(w, errors.New("Captcha error"))
		return
	}

	if !shouldAccept {
		f.respond(w, errors.New("Invalid captcha"))
		return
	}

	err = f.mailService.Send(req)
	if err != nil {
		log.Error().Err(err).Msg("smtp error")
		f.respond(w, errors.New("Error sending mail"))
		return
	}

	f.respond(w, nil)
}

const (
	responseHTMLTemplate = `<div class="col-span-full" id="alert">
<div class="flex rounded-md bg-primary-100 px-4 py-3 dark:bg-primary-900">
<span class="pe-3 text-primary-400">
<span class="icon inline-block align-text-middle">
<svg aria-hidden="true" class="hi-svg-inline" fill="currentcolor" height="1em" id="mdi-information-outline" viewBox="0 0 24 24" width="1em">
<path d="M11 9h2V7H11m1 13c-4.41.0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8M12 2A10 10 0 002 12 10 10 0 0012 22 10 10 0 0022 12 10 10 0 0012 2M11 17h2V11H11v6z">
</path>
</svg>
</span> 
</span>
<span class="dark:text-neutral-300" id="alert-message">%s</span></div></div>`
)

func (f *FormHandler) respond(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusOK)
	var response string
	if err != nil {
		response = fmt.Sprintf(responseHTMLTemplate, "An error occured: "+err.Error())
	} else {
		response = fmt.Sprintf(responseHTMLTemplate, "Message has been sent. I will get back to you asap!")
	}

	_, err = w.Write([]byte(response))
	if err != nil {
		log.Err(err).Msg("error writing response")
	}
}
