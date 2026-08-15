package pkg

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostFollowsRedirectOnSameHost(t *testing.T) {
	var gotStamp string
	var gotBody []byte

	mux := http.NewServeMux()
	mux.HandleFunc("/initial", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/final", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		gotStamp = r.Header.Get("X-Stamp")
		gotBody, _ = io.ReadAll(r.Body)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	response, err := post(context.Background(), "http", serverURL.Host, "/initial", []byte(`{"a":1}`), "test-stamp")
	require.NoError(t, err)

	defer response.Body.Close() //nolint: errcheck

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "test-stamp", gotStamp)
	assert.Equal(t, []byte(`{"a":1}`), gotBody)
}

func TestPostDoesNotFollowRedirectToOtherHost(t *testing.T) {
	var otherRequests atomic.Int32

	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otherRequests.Add(1)
	}))
	defer other.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL, http.StatusTemporaryRedirect)
	}))
	defer origin.Close()

	originURL, err := url.Parse(origin.URL)
	require.NoError(t, err)

	response, err := post(context.Background(), "http", originURL.Host, "/", []byte(`{"a":1}`), "test-stamp")
	if response != nil {
		response.Body.Close() //nolint: errcheck
	}

	assert.Error(t, err)
	assert.Equal(t, int32(0), otherRequests.Load())
}
