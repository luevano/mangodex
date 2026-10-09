package mangodex

import "fmt"

// MangaDex requires a User-Agent that isn't spoofed (e.g. a browser one).
const defaultUserAgent = "mangodex (+https://github.com/luevano/mangodex)"

type Options struct {
	UserAgent string
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
