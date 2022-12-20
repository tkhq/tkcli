package display_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkhq/tkcli/internal/display"
)

func TestFormatResponse(t *testing.T) {
	var testCases = []struct {
		code            int
		body            string
		formatter       display.Formatter
		expectedDisplay string
		expectedErr     error
	}{
		{200, "foo", display.FormatterJSON, "foo", nil},
		{200, "foo", display.FormatterPretty, "foo\n", nil},
		{200, `{"foo":"bar", "a": 1}`, display.FormatterJSON, `{
    "foo": "bar",
    "a": 1
}`, nil},
		{200, `{"foo":"bar", "a": 1}`, display.FormatterPretty, trimOutput(`
foo: bar
a: 1
`), nil},
		{200, `{"foo": {"hello":"world","bar":123}, "a": 1}`, display.FormatterJSON, `{
    "foo": {
        "hello": "world",
        "bar": 123
    },
    "a": 1
}`, nil},
		{200, `{"foo": {"hello":"world","bar":123}, "a": 1}`, display.FormatterPretty, trimOutput(`
foo:
    hello: world
    bar: 123
a: 1
`), nil},
		{500, "foo", display.FormatterJSON, `{
    "responseBody": "foo",
    "responseCode": 500
}`, nil},
		{500, "foo", display.FormatterPretty, trimOutput(`
responseBody: foo
responseCode: 500
`), nil},
	}

	for _, testCase := range testCases {
		httpResponse := http.Response{
			StatusCode: testCase.code,
			Body:       io.NopCloser(strings.NewReader(testCase.body)),
		}
		actualDisplay, actualErr := display.FormatResponse(&httpResponse, testCase.formatter)
		assert.Equal(t, testCase.expectedErr, actualErr)
		assert.Equal(t, testCase.expectedDisplay, actualDisplay)
	}
}

func TestFormatStruct(t *testing.T) {
	var testCases = []struct {
		data            interface{}
		formatter       display.Formatter
		expectedDisplay string
		expectedErr     error
	}{
		{map[string]interface{}{
			"foo": "bar",
			"a":   1,
		}, display.FormatterJSON, `{
    "a": 1,
    "foo": "bar"
}`, nil},
		{map[string]interface{}{
			"foo": "bar",
			"a":   1,
		}, display.FormatterPretty, trimOutput(`
a: 1
foo: bar
`), nil},
		{map[string]interface{}{
			"foo": map[string]interface{}{
				"hello": "world",
				"bar":   123,
			},
			"a": 1,
		}, display.FormatterJSON, `{
    "a": 1,
    "foo": {
        "bar": 123,
        "hello": "world"
    }
}`, nil},
		{map[string]interface{}{
			"foo": map[string]interface{}{
				"hello": "world",
				"bar":   123,
			},
			"a": 1,
		}, display.FormatterPretty, trimOutput(`
a: 1
foo:
    bar: 123
    hello: world
`), nil},
	}

	for _, testCase := range testCases {
		result, err := display.FormatStruct(testCase.data.(map[string]interface{}), testCase.formatter)
		assert.Equal(t, testCase.expectedErr, err)
		assert.Equal(t, testCase.expectedDisplay, result)
	}
}

func trimOutput(input string) string {
	return strings.TrimSpace(input) + "\n"
}
