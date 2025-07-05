package handler

import (
	"encoding/json"
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
		f.respond(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	req := service.MailRequest{
		Name:               r.Form.Get("name"),
		Email:              r.Form.Get("email"),
		Message:            r.Form.Get("message"),
		Subject:            r.Form.Get("subject"),
		FrcCaptchaSolution: r.Form.Get("frc-captcha-solution"),
	}

	if !req.IsComplete() {
		f.respond(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	log.Info().Interface("MailRequest", req).Msg("incoming form MailRequest")

	solution := req.FrcCaptchaSolution
	shouldAccept, err := f.frcClient.CheckCaptchaSolution(r.Context(), solution)
	if err != nil {
		log.Error().Err(err).Msg("captcha check error")
		f.respond(w, "Captcha error", http.StatusInternalServerError)
		return
	}

	if !shouldAccept {
		f.respond(w, "Invalid captcha", http.StatusBadRequest)
		return
	}

	err = f.mailService.Send(req)
	if err != nil {
		log.Error().Err(err).Msg("smtp error")
		f.respond(w, "Captcha error", http.StatusInternalServerError)
		return
	}

	f.respond(w, "Message sent. I will get back to you asap!", http.StatusOK)
}

type response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (f *FormHandler) respond(w http.ResponseWriter, msg string, statusCode int) {
	res := response{
		Message: msg,
		Success: http.StatusOK == statusCode,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonRes)
	if err != nil {
		log.Warn().Err(err).Msg("error writing response")
	}
}
