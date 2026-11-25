package service

import "testing"

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
