package requo

import (
	"fmt"
	"net/http"
)

var (
	DefaultClient *http.Client
)

func init() {
	DefaultClient = &http.Client{}
}

type Options struct {
	Client             *http.Client
	Request            *http.Request
	AllowedStatusCodes []int
}

type OptionsFunc func(opts *Options)

type HTTPError struct {
	Status     string
	StatusCode int
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", h.StatusCode, h.Status)
}
