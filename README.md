# go-form-mailer

Migrating my website recently from plain HTML/JS/PHP to Hugo, a static site generator, I wanted to ditch my old
form handler that was built on PHPMailer. This service listens to a POST request with URL encoded form data to
send a templated email via SMTP. It uses FriendlyCaptcha for verification.

## Setup

- Rename `config.sample.toml` to `config.toml`
- Fill in SMTP auth info under `[smtp]`
- Create a FriendlyCaptcha API and site key
- Fill in FriendlyCaptcha credentials under `[frc]`
- Run server with `go run .`