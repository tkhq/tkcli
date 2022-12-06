package display_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkhq/tkcli/cmd/tk/internal/display"
)

func TestDisplayResponse(t *testing.T) {
	var testCases = []struct {
		code            int
		body            string
		expectedDisplay string
		expectedErr     error
	}{
		{200, "foo", "foo", nil},
		{200, `{"foo":"bar"}`, `{
    "foo": "bar"
}`, nil},
		{500, "foo", `{
    "responseBody": "foo",
    "responseCode": 500
}`, nil},
	}

	for _, testCase := range testCases {
		httpResponse := http.Response{
			StatusCode: testCase.code,
			Body:       io.NopCloser(strings.NewReader(testCase.body)),
		}
		actualDisplay, actualErr := display.DisplayResponse(&httpResponse)
		assert.Equal(t, testCase.expectedErr, actualErr)
		assert.Equal(t, testCase.expectedDisplay, actualDisplay)
	}
}
