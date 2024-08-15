package handler

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go-sdk"
	"github.com/rs/zerolog/log"

	"github.com/davidramiro/go-form-mailer/api"
	"github.com/davidramiro/go-form-mailer/internal/service"
	"github.com/go-faster/errors"
)

type FormHandler struct {
	mutex       sync.Mutex
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

func (f *FormHandler) FormPost(ctx context.Context, req *api.FormData) (api.FormPostRes, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	log.Info().Interface("request", req).Msg("incoming form request")

	if req.Name == "" ||
		req.Subject == "" ||
		req.Message == "" ||
		req.Email == "" ||
		req.CaptchaSolution == "" {
		return &api.FormPostBadRequest{}, errors.New("missing field(s)")
	}

	solution := req.CaptchaSolution
	shouldAccept, err := f.frcClient.CheckCaptchaSolution(ctx, solution)
	if err != nil {
		if errors.Is(err, friendlycaptcha.ErrVerificationFailedDueToClientError) {
			log.Error().Err(err).Msg("frc client misconfigured")
			return &api.FormPostInternalServerError{}, errors.New("captcha client error")
		} else if errors.Is(err, friendlycaptcha.ErrVerificationRequest) {
			log.Error().Err(err).Msg("frc client api error")
			return &api.FormPostInternalServerError{}, errors.New("captcha api unreachable")
		}
	}

	if !shouldAccept {
		return &api.FormPostBadRequest{}, errors.New("captcha solution is invalid")
	}

	err = f.mailService.Send(*req)
	if err != nil {
		return &api.FormPostInternalServerError{}, errors.New("SMTP error")
	}

	return &api.FormPostOK{Message: "Message sent. I will get back to you asap!"}, nil
}

func (f *FormHandler) NewError(_ context.Context, err error) *api.ResponseStatusCode {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return &api.ResponseStatusCode{
		StatusCode: http.StatusBadRequest,
		Response: api.Response{
			Message: fmt.Sprintf("Failed sending message: %s", err.Error()),
		},
	}
}
