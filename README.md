# mangodex

Golang API wrapper for MangaDex v5's API (5.13.1).

Full API documentation is found [here](https://api.mangadex.org/docs/).

> **Warning**
> 
> The API implementation is not stable and may change at any time.

**Note**: This is a fork of [Bob620/mangodex](https://github.com/Bob620/mangodex). Specially for use in [mangoprovider](https://github.com/luevano/mangoprovider) and [mangal](https://github.com/luevano/mangal).

## Installation

To install, do `go get -u github.com/luevano/mangodex@latest`.

## Usage

Basic usage example.

```go
package main

import (
    "fmt"
    "net/url"

    "github.com/luevano/mangodex"
)

func main() {
    // Create new client.
    c, err := mangodex.NewDexClient(mangodex.DefaultOptions())
    if err != nil {
        panic(err)
    }

    // Create search params.
    params := url.Values{}
    params.Set("title", "tengoku daimakyou")
    params.Set("availableTranslatedLanguage[]", "en")

    // Get list of mangas by search query.
    mangaList, err := c.Manga.List(params)
    if err != nil {
        panic(err)
    }
    for _, manga := range mangaList {
        // Do something
        fmt.Println(manga.GetTitle("en", true))
    }
}
```

Every request method has a `...Context` version that takes a `context.Context` (e.g.
`c.Manga.ListContext(ctx, params)`), the plain one uses `context.Background()`.

Errors can be checked with `errors.Is` against `ErrNotFound`, `ErrRateLimited` and `ErrMaintenance`, or
inspected with `errors.As` as `*HTTPError` (non-200 status) or `*RateLimitError` (429, with `RetryAfter`).

Follow MangaDex's [usage rules](https://api.mangadex.org/docs/2-limitations/): don't spoof the
User-Agent (the default one identifies this library), stay around 5 requests per second.

## Contributing

Any contributions are welcome.
