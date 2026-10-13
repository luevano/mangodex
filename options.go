package mangodex

import (
	"fmt"
	"net/http"
	"time"
)

// MangaDex requires a User-Agent that isn't spoofed (e.g. a browser one).
const defaultUserAgent = "mangodex (+https://github.com/luevano/mangodex)"

// defaultTimeout: Timeout of the client used when Options.HTTPClient is nil.
const defaultTimeout = time.Minute

type Options struct {
	UserAgent string

	// HTTPClient: Client used for every request.
	//
	// Optional, a client with a 1 minute timeout is used when nil.
	HTTPClient *http.Client
}

func (o Options) validate() error {
	if o.UserAgent == "" {
		return fmt.Errorf("UserAgent is empty")
	}
	return nil
}

func DefaultOptions() Options {
	return Options{
		UserAgent: defaultUserAgent,
	}
}
