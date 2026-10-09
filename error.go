package mangodex

import (
	"fmt"
	"strings"
)

// ErrorResponse: Typical response for errored requests.
type ErrorResponse struct {
	Result string  `json:"result"`
	Errors []Error `json:"errors"`
}

// GetErrors: Get the errors for this particular request, in one line.
func (er *ErrorResponse) GetErrors() string {
	var errors []string
	for _, err := range er.Errors {
		if err.Detail == "" {
			errors = append(errors, err.Title)
			continue
		}
		errors = append(errors, fmt.Sprintf("%s: %s", err.Title, err.Detail))
	}
	return strings.Join(errors, "; ")
}

// Error: Struct containing details of an error.
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
