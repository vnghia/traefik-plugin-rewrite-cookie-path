package traefik_plugin_rewrite_cookie_path // nolint

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

const setCookieHeader string = "Set-Cookie"

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Name        string `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`
	Regex       string `json:"regex,omitempty" toml:"regex,omitempty" yaml:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty" toml:"replacement,omitempty" yaml:"replacement,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	Rewrites []Rewrite `json:"rewrites,omitempty" toml:"rewrites,omitempty" yaml:"rewrites,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type rewrite struct {
	name        string
	regex       *regexp.Regexp
	replacement string
}

type rewriteBody struct {
	name     string
	next     http.Handler
	rewrites []rewrite
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	rewrites := make([]rewrite, len(config.Rewrites))

	for i, rewriteConfig := range config.Rewrites {
		regex, err := regexp.Compile(rewriteConfig.Regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", rewriteConfig.Regex, err)
		}

		rewrites[i] = rewrite{
			name:        rewriteConfig.Name,
			regex:       regex,
			replacement: rewriteConfig.Replacement,
		}
	}

	return &rewriteBody{
		name:     name,
		next:     next,
		rewrites: rewrites,
	}, nil
}

func (r *rewriteBody) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &responseWriter{
		writer:   rw,
		rewrites: r.rewrites,
	}

	r.next.ServeHTTP(wrappedWriter, req)
}

type responseWriter struct {
	writer   http.ResponseWriter
	rewrites []rewrite
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	headers := r.writer.Header()
	req := http.Response{Header: headers}
	cookies := req.Cookies()

	r.writer.Header().Del(setCookieHeader)

	for _, cookie := range cookies {
		for _, rewrite := range r.rewrites {
			if cookie.Name == rewrite.name {
				cookie.Path = rewrite.regex.ReplaceAllString(cookie.Path, rewrite.replacement)
			}
		}
		http.SetCookie(r, cookie)
	}

	r.writer.WriteHeader(statusCode)
}
