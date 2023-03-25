package display

import (
	"fmt"
)

// ErrorResponse is a structured format to display an HTTP error response.
type ErrorResponse struct {
	Code int    `json:"responseCode"`
	Text string `json:"responseBody"`
}

func (r *ErrorResponse) Error() string {
   return fmt.Sprintf("%d: %s", r.Code, r.Text)
}
