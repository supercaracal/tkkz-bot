package command

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestOKHandler struct{}

func (h *TestOKHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{'D', 'u', 'd', 'e'})
}

type TestNGHandler struct{}

func (h *TestNGHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte{'O', 'o', 'p', 's'})
}

func TestGetDefaultReply(t *testing.T) {
	cases := []struct {
		desc    string
		handler http.Handler
		want    string
	}{
		{
			"200",
			&TestOKHandler{},
			"Dude（ﾎﾞﾛﾝ",
		},
		{
			"302",
			http.RedirectHandler("https://www.example.com", http.StatusFound),
			// "brain api replied 302: <a href=\"https://www.example.com\">Found</a>.\n\n",
			"brain api replied 302: ",
		},
		{
			"404",
			http.NotFoundHandler(),
			"brain api replied 404: 404 page not found\n",
		},
		{
			"500",
			&TestNGHandler{},
			"brain api replied 500: Oops",
		},
		{
			"timeout",
			http.TimeoutHandler(&TestOKHandler{}, (0 * time.Second), "zzz"),
			"brain api replied 503: zzz",
		},
	}

	for i, c := range cases {
		i := i
		c := c
		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(c.handler)
			defer ts.Close()

			if got := GetDefaultReply(ts.URL, "Hi"); got != c.want {
				t.Errorf("%d: want: %s, got: %s", i, c.want, got)
			}
		})
	}
}

func TestTrimReply(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"foobar", "foobar"},
		{"<@U123ABC>foobar", "foobar"},
		{"foo<@U123ABC>bar", "foobar"},
		{"foobar<@U123ABC>", "foobar"},
		{":+1:foobar", "foobar"},
		{"foobar:+1:", "foobar"},
		{"foo:+1:bar", "foobar"},
		{"<http://example.com>foobar", "foobar"},
		{"foobar<http://example.com>", "foobar"},
		{"foo<http://example.com>bar", "foobar"},
		{" 　\t:+1: < foo <123ABC>bar\r\n", "foo bar"},
	}

	for i, c := range cases {
		if got := trimReply(c.text); got != c.want {
			t.Errorf("%d: want: %s, got: %s", i, c.want, got)
		}
	}
}

func TestPersonalizeReply(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"foobar", "foobar（ﾎﾞﾛﾝ"},
	}

	for i, c := range cases {
		if got := personalizeReply(c.text); got != c.want {
			t.Errorf("%d: want: %s, got: %s", i, c.want, got)
		}
	}
}
