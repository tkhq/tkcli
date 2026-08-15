package pkg

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckRedirectLimit(t *testing.T) {
	origin, err := url.Parse("https://example.com")
	require.NoError(t, err)

	redirect := &http.Request{
		URL:      origin,
		Response: &http.Response{StatusCode: http.StatusTemporaryRedirect},
	}
	previous := make([]*http.Request, 11)
	previous[0] = &http.Request{URL: origin}

	require.NoError(t, checkRedirect(redirect, previous[:10]))
	assert.EqualError(t, checkRedirect(redirect, previous), "stopped after 10 redirects")
}

func TestEffectiveOrigin(t *testing.T) {
	for rawURL, expected := range map[string]string{
		"https://EXAMPLE.com/path":        "https://example.com:443",
		"http://example.com:80":           "http://example.com:80",
		"http://[0:0:0:0:0:0:0:1]:8080/x": "http://[::1]:8080",
		"https://BÜCHER.example":          "https://xn--bcher-kva.example:443",
	} {
		parsed, err := url.Parse(rawURL)
		require.NoError(t, err)

		origin, err := effectiveOrigin(parsed)
		require.NoError(t, err)
		assert.Equal(t, expected, origin, rawURL)
	}
}
