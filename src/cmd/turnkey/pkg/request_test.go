package pkg

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostFollowsRedirectOnSameHost(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/initial", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/final", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, "test-stamp", r.Header.Get("X-Stamp"))
		assert.Equal(t, []byte(`{"a":1}`), body)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := post(context.Background(), "http", strings.TrimPrefix(server.URL, "http://"), "/initial", []byte(`{"a":1}`), "test-stamp")
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
}

func TestPostDoesNotFollowMethodChangingRedirect(t *testing.T) {
	var finalRequests atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("/initial", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/final", http.StatusFound)
	})
	mux.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		finalRequests.Add(1)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	_, err := post(context.Background(), "http", strings.TrimPrefix(server.URL, "http://"), "/initial", []byte(`{"a":1}`), "test-stamp")
	assert.Error(t, err)
	assert.Equal(t, int32(0), finalRequests.Load())
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

	_, err := post(context.Background(), "http", strings.TrimPrefix(origin.URL, "http://"), "/", []byte(`{"a":1}`), "test-stamp")
	assert.Error(t, err)
	assert.Equal(t, int32(0), otherRequests.Load())
}
