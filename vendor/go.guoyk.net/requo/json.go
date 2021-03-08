package requo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func JSON(ctx context.Context, method string, url string, in interface{}, out interface{}, fs ...OptionsFunc) (err error) {
	var body io.Reader
	if in != nil {
		var buf []byte
		if buf, err = json.Marshal(in); err != nil {
			return
		}
		body = bytes.NewReader(buf)
	}
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, method, url, body); err != nil {
		return
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if out != nil {
		req.Header.Set("Accept", "application/json")
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
		err = json.Unmarshal(buf, out)
	}
	return
}

func JSONGet(ctx context.Context, url string, out interface{}, fs ...OptionsFunc) error {
	return JSON(ctx, http.MethodGet, url, nil, out, fs...)
}

func JSONPost(ctx context.Context, url string, in, out interface{}, fs ...OptionsFunc) error {
	return JSON(ctx, http.MethodPost, url, in, out, fs...)
}

func JSONPut(ctx context.Context, url string, in, out interface{}, fs ...OptionsFunc) error {
	return JSON(ctx, http.MethodPut, url, in, out, fs...)
}

func JSONPatch(ctx context.Context, url string, in, out interface{}, fs ...OptionsFunc) error {
	return JSON(ctx, http.MethodPatch, url, in, out, fs...)
}

func JSONDelete(ctx context.Context, url string, out interface{}, fs ...OptionsFunc) error {
	return JSON(ctx, http.MethodDelete, url, nil, out, fs...)
}
