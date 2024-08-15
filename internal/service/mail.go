package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"reflect"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/davidramiro/go-form-mailer/api"
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
		return fmt.Errorf("input param should be a struct")
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

func (m *MailService) Send(mail api.FormData) error {
	to := []string{
		m.params.ToMail,
	}
	auth := smtp.PlainAuth("", m.params.SMTPUser, m.params.SMTPPass, m.params.SMTPHost)

	t, err := template.ParseFiles("template/mail.html")
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	mailHeader := fmt.Sprintf("Subject: %s \n%s\n\n", mail.Subject, mimeHeaders)
	body.Write([]byte(mailHeader))

	err = t.Execute(&body, struct {
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
