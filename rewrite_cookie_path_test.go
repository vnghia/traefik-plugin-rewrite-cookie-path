package traefik_plugin_rewrite_cookie_path // nolint

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc          string
		rewrites      []Rewrite
		reqHeader     http.Header
		expRespHeader http.Header
	}{
		{
			desc: "should replace foo by bar",
			rewrites: []Rewrite{
				{
					Name:        "someName",
					Regex:       "foo",
					Replacement: "bar",
				},
			},
			reqHeader: map[string][]string{
				"set-cookie":    {"someName=someValue; Path=/foo"},
				"anotherHeader": {"Path=/"},
			},
			expRespHeader: map[string][]string{
				"set-cookie":    {"someName=someValue; Path=/bar"},
				"anotherHeader": {"Path=/"},
			},
		},
		{
			desc: "should replace the same name",
			rewrites: []Rewrite{
				{
					Name:        "someName1",
					Regex:       "foo",
					Replacement: "bar",
				},
			},
			reqHeader: map[string][]string{
				"set-cookie":    {"someName2=someValue; Path=foo"},
				"anotherHeader": {"Path=/"},
			},
			expRespHeader: map[string][]string{
				"set-cookie":    {"someName2=someValue; Path=foo"},
				"anotherHeader": {"Path=/"},
			},
		},
		{
			desc: "should replace http by https",
			rewrites: []Rewrite{
				{
					Name:        "someName",
					Regex:       "^http://(.+)$",
					Replacement: "https://$1",
				},
			},
			reqHeader: map[string][]string{
				"set-cookie":    {"someName=someValue; Path=http://test:1000"},
				"anotherHeader": {"Path=/"},
			},
			expRespHeader: map[string][]string{
				"set-cookie":    {"someName=someValue; Path=https://test:1000"},
				"anotherHeader": {"Path=/"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := &Config{
				Rewrites: test.rewrites,
			}

			next := func(rw http.ResponseWriter, req *http.Request) {
				for k, v := range test.reqHeader {
					for _, h := range v {
						rw.Header().Add(k, h)
					}
				}
				rw.WriteHeader(http.StatusOK)
			}

			rewriteBody, err := New(context.Background(), http.HandlerFunc(next), config, "rewriteCookie")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			rewriteBody.ServeHTTP(recorder, req)
			for k, expected := range test.expRespHeader {
				values := recorder.Header().Values(k)

				if !testEq(values, expected) {
					t.Errorf("Slice arent equals: expect: %+v, result: %+v", expected, values)
				}
			}
		})
	}
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
