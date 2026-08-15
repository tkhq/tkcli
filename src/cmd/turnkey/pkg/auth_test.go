package pkg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkhq/go-sdk/pkg/api/client/sessions"
)

func TestAPIClientDoesNotFollowRedirectToOtherHost(t *testing.T) {
	var otherRequests atomic.Int32

	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otherRequests.Add(1)
	}))
	defer other.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL, http.StatusTemporaryRedirect)
	}))
	defer origin.Close()

	apiClient := newAPIClient("http", strings.TrimPrefix(origin.URL, "http://"))

	_, err := apiClient.Sessions.GetWhoami(&sessions.GetWhoamiParams{Context: context.Background()}, nil)
	assert.Error(t, err)
	assert.Equal(t, int32(0), otherRequests.Load())
}
