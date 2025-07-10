package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
	"reflect"
	"strconv"

	"github.com/go-faster/errors"

	"github.com/rs/zerolog/log"
)

type MailService struct {
	params MailServiceParams
}

type MailServiceParams struct {
	SMTPFrom string
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	ToMail   string
}

type MailRequest struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Message            string `json:"message"`
	Subject            string `json:"subject"`
	FrcCaptchaSolution string `json:"frc-captcha-solution"`
}

func NewMailService(params MailServiceParams) (*MailService, error) {
	err := ValidateMailParams(params)
	if err != nil {
		return nil, fmt.Errorf("invalid mail service params: %w", err)
	}

	return &MailService{params: params}, nil
}

func ValidateMailParams(s interface{}) error {
	structType := reflect.TypeOf(s)
	if structType.Kind() != reflect.Struct {
		return errors.New("input param should be a struct")
	}

	structVal := reflect.ValueOf(s)
	fieldNum := structVal.NumField()

	for i := range fieldNum {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name
		isSet := field.IsValid() && !field.IsZero()
		if !isSet {
			return fmt.Errorf("%s should be set", fieldName)
		}
	}

	return nil
}

func (m *MailService) Send(mail MailRequest) error {
	to := []string{
		m.params.ToMail,
	}
	auth := smtp.PlainAuth("", m.params.SMTPUser, m.params.SMTPPass, m.params.SMTPHost)

	t, err := template.ParseFiles("template/mail.html")
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	var body *bytes.Buffer
	_, err = fmt.Fprintf(body,
		"Subject: %s \n%s\n\n",
		mail.Subject,
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n")
	if err != nil {
		return fmt.Errorf("writing mail header failed: %w", err)
	}

	err = t.Execute(body, struct {
		Name    string
		Message string
		Email   string
		Subject string
	}{
		Subject: mail.Subject,
		Name:    mail.Name,
		Message: mail.Message,
		Email:   mail.Email,
	})
	if err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	err = smtp.SendMail(
		net.JoinHostPort(
			m.params.SMTPHost,
			strconv.Itoa(m.params.SMTPPort)),
		auth,
		m.params.SMTPFrom,
		to,
		body.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}

	log.Info().Msg("email sent")

	return nil
}

const (
	maxShort = 256
	maxLong  = 80000
)

func (r MailRequest) Validate() error {
	if len(r.Name) == 0 || len(r.Name) > maxShort {
		return errors.New("Invalid name")
	}
	if len(r.Email) == 0 || len(r.Email) > maxShort {
		return errors.New("invalid email length")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email format")
	}
	if len(r.Subject) == 0 || len(r.Subject) > maxShort {
		return errors.New("invalid subject")
	}
	if len(r.Message) == 0 || len(r.Message) > maxLong {
		return errors.New("invalid message")
	}

	if len(r.FrcCaptchaSolution) == 0 || len(r.FrcCaptchaSolution) > maxShort {
		return errors.New("invalid captcha solution")
	}

	return nil
}
