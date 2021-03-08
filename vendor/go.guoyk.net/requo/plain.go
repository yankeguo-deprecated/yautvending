package requo

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Plain(ctx context.Context, method string, url string, in *string, out *string, fs ...OptionsFunc) (err error) {
	var body io.Reader
	if in != nil {
		body = strings.NewReader(*in)
	}
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, method, url, body); err != nil {
		return
	}
	opts := &Options{
		Client:             DefaultClient,
		Request:            req,
		AllowedStatusCodes: []int{http.StatusOK},
	}
	for _, fn := range fs {
		fn(opts)
	}
	var resp *http.Response
	if resp, err = opts.Client.Do(opts.Request); err != nil {
		return
	}
	defer resp.Body.Close()

	for _, code := range opts.AllowedStatusCodes {
		if code == resp.StatusCode {
			goto allowed
		}
	}

	err = &HTTPError{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
	return

allowed:
	if out != nil {
		var buf []byte
		if buf, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		*out = string(buf)
	}
	return
}

func PlainGet(ctx context.Context, url string, out *string, fs ...OptionsFunc) error {
	return Plain(ctx, http.MethodGet, url, nil, out, fs...)
}

func PlainPost(ctx context.Context, url string, in, out *string, fs ...OptionsFunc) error {
	return Plain(ctx, http.MethodPost, url, in, out, fs...)
}

func PlainPut(ctx context.Context, url string, in, out *string, fs ...OptionsFunc) error {
	return Plain(ctx, http.MethodPut, url, in, out, fs...)
}

func PlainPatch(ctx context.Context, url string, in, out *string, fs ...OptionsFunc) error {
	return Plain(ctx, http.MethodPatch, url, in, out, fs...)
}

func PlainDelete(ctx context.Context, url string, out *string, fs ...OptionsFunc) error {
	return Plain(ctx, http.MethodDelete, url, nil, out, fs...)
}
