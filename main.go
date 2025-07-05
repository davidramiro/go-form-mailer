package main

import (
	"fmt"
	"net/http"

	friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go-sdk"

	"github.com/davidramiro/go-form-mailer/internal/handler"
	"github.com/davidramiro/go-form-mailer/internal/service"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	log.Info().Msg("startup, reading config...")

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error reading config file")
	}

	log.Info().Msg("initializing mail service")

	msParams := service.MailServiceParams{
		SMTPFrom: viper.GetString("smtp.from"),
		SMTPHost: viper.GetString("smtp.host"),
		SMTPPort: viper.GetInt("smtp.port"),
		SMTPUser: viper.GetString("smtp.user"),
		SMTPPass: viper.GetString("smtp.password"),
		ToMail:   viper.GetString("smtp.to"),
	}
	ms, err := service.NewMailService(msParams)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mail service")
	}

	log.Info().Msg("initializing form handler")

	frcClient := friendlycaptcha.NewClient(viper.GetString("frc.apiKey"), viper.GetString("frc.siteKey"))

	fs, err := handler.NewFormHandler(ms, frcClient)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating form handler")
	}

	log.Info().Msg("spinning up server")

	m := http.NewServeMux()
	m.Handle("/form", fs)

	srv := &http.Server{
		ReadTimeout:  viper.GetDuration("server.timeout"),
		WriteTimeout: viper.GetDuration("server.timeout"),
		Handler:      m,
		Addr:         fmt.Sprintf(":%d", viper.GetInt("server.port")),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
