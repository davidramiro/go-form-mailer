package main

import (
	"fmt"
	"net/http"

	friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go-sdk"

	_ "github.com/ogen-go/ogen/gen"

	"github.com/davidramiro/go-form-mailer/api"
	"github.com/davidramiro/go-form-mailer/internal/handler"
	"github.com/davidramiro/go-form-mailer/internal/service"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen -package api --clean spec/openapi.yaml

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

	srv, err := api.NewServer(fs)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating server")
	}

	listenAddress := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	log.Fatal().
		Err(http.ListenAndServe(listenAddress, srv)).
		Msg("form mailer server closed")
}
