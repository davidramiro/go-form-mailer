// Code generated by ogen, DO NOT EDIT.

package api

import (
	"mime"
	"net/http"

	"github.com/go-faster/errors"
	"go.uber.org/multierr"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

func (s *Server) decodeFormPostRequest(r *http.Request) (
	req *FormData,
	close func() error,
	rerr error,
) {
	var closers []func() error
	close = func() error {
		var merr error
		// Close in reverse order, to match defer behavior.
		for i := len(closers) - 1; i >= 0; i-- {
			c := closers[i]
			merr = multierr.Append(merr, c())
		}
		return merr
	}
	defer func() {
		if rerr != nil {
			rerr = multierr.Append(rerr, close())
		}
	}()
	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return req, close, errors.Wrap(err, "parse media type")
	}
	switch {
	case ct == "application/x-www-form-urlencoded":
		if r.ContentLength == 0 {
			return req, close, validate.ErrBodyRequired
		}
		form, err := ht.ParseForm(r)
		if err != nil {
			return req, close, errors.Wrap(err, "parse form")
		}

		var request FormData
		q := uri.NewQueryDecoder(form)
		{
			cfg := uri.QueryParameterDecodingConfig{
				Name:    "name",
				Style:   uri.QueryStyleForm,
				Explode: true,
			}
			if err := q.HasParam(cfg); err == nil {
				if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					request.Name = c
					return nil
				}); err != nil {
					return req, close, errors.Wrap(err, "decode \"name\"")
				}
			} else {
				return req, close, errors.Wrap(err, "query")
			}
		}
		{
			cfg := uri.QueryParameterDecodingConfig{
				Name:    "email",
				Style:   uri.QueryStyleForm,
				Explode: true,
			}
			if err := q.HasParam(cfg); err == nil {
				if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					request.Email = c
					return nil
				}); err != nil {
					return req, close, errors.Wrap(err, "decode \"email\"")
				}
				if err := func() error {
					if err := (validate.String{
						MinLength:    0,
						MinLengthSet: false,
						MaxLength:    0,
						MaxLengthSet: false,
						Email:        true,
						Hostname:     false,
						Regex:        nil,
					}).Validate(string(request.Email)); err != nil {
						return errors.Wrap(err, "string")
					}
					return nil
				}(); err != nil {
					return req, close, errors.Wrap(err, "validate")
				}
			} else {
				return req, close, errors.Wrap(err, "query")
			}
		}
		{
			cfg := uri.QueryParameterDecodingConfig{
				Name:    "subject",
				Style:   uri.QueryStyleForm,
				Explode: true,
			}
			if err := q.HasParam(cfg); err == nil {
				if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					request.Subject = c
					return nil
				}); err != nil {
					return req, close, errors.Wrap(err, "decode \"subject\"")
				}
			} else {
				return req, close, errors.Wrap(err, "query")
			}
		}
		{
			cfg := uri.QueryParameterDecodingConfig{
				Name:    "message",
				Style:   uri.QueryStyleForm,
				Explode: true,
			}
			if err := q.HasParam(cfg); err == nil {
				if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					request.Message = c
					return nil
				}); err != nil {
					return req, close, errors.Wrap(err, "decode \"message\"")
				}
			} else {
				return req, close, errors.Wrap(err, "query")
			}
		}
		{
			cfg := uri.QueryParameterDecodingConfig{
				Name:    "frc-captcha-solution",
				Style:   uri.QueryStyleForm,
				Explode: true,
			}
			if err := q.HasParam(cfg); err == nil {
				if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					request.FrcMinusCaptchaMinusSolution = c
					return nil
				}); err != nil {
					return req, close, errors.Wrap(err, "decode \"frc-captcha-solution\"")
				}
			} else {
				return req, close, errors.Wrap(err, "query")
			}
		}
		return &request, close, nil
	default:
		return req, close, validate.InvalidContentType(ct)
	}
}
