package mangodex

import (
	"fmt"
	"strings"
)

// ErrorResponse : Typical response for errored requests.
type ErrorResponse struct {
	Result string  `json:"result"`
	Errors []Error `json:"errors"`
}

func (er *ErrorResponse) GetResult() string {
	return er.Result
}

// GetErrors : Get the errors for this particular request.
func (er *ErrorResponse) GetErrors() string {
	var errors strings.Builder
	for _, err := range er.Errors {
		errors.WriteString(fmt.Sprintf("%s: %s\n", err.Title, err.Detail))
	}
	return errors.String()
}

// Error : Struct containing details of an error.
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
