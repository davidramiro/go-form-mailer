package service

import (
	"strings"
	"testing"
)

func TestNl2br(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "single line", in: "hello", want: "hello"},
		{name: "LF", in: "a\nb", want: "a<br>b"},
		{name: "CRLF", in: "a\r\nb", want: "a<br>b"},
		{name: "blank line", in: "a\n\nb", want: "a<br><br>b"},
		{name: "escapes HTML", in: "<b>hi</b>", want: "&lt;b&gt;hi&lt;/b&gt;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(nl2br(tt.in)); got != tt.want {
				t.Errorf("nl2br(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateSubject(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		wantErr bool
	}{
		{name: "normal", subject: "Hello there"},
		{name: "empty", subject: "", wantErr: true},
		{name: "too long", subject: strings.Repeat("a", maxShort+1), wantErr: true},
		{name: "LF", subject: "Hello\nBcc: attacker@example.com", wantErr: true},
		{name: "CRLF", subject: "Hello\r\nBcc: attacker@example.com", wantErr: true},
		{name: "tab", subject: "Hello\tWorld", wantErr: true},
		{name: "NUL", subject: "Hello\x00World", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSubject(tt.subject)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSubject(%q) error = %v, wantErr %v", tt.subject, err, tt.wantErr)
			}
		})
	}
}

func TestMailRequestValidateRejectsHeaderInjection(t *testing.T) {
	req := MailRequest{
		Name:               "John Doe",
		Email:              "john@example.com",
		Subject:            "Hello\r\nBcc: attacker@example.com",
		Message:            "Hello",
		FrcCaptchaResponse: "solution",
	}

	if err := req.Validate(); err == nil {
		t.Fatal("Validate() accepted a subject containing CRLF")
	}
}
