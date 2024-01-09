# mangodex

Golang API wrapper for MangaDex v5's API.

Full API documentation is found [here](https://api.mangadex.org/docs.html).

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
    c := mangodex.NewDexClient()

    // Create search params.
    params := url.Values{}
    params.Set("title", "tengoku daimakyou")
    params.Set("translatedLanguage[]", "en")

    // Get list of mangas by search query.
    mangaList, err := c.Manga.List(params)
    if err != nil {
        panic(err)
    }
    for _, manga := range mangaList {
        // Do something
        fmt.Println(manga.GetTitle("en"))
    }
}
```

## Contributing

Any contributions are welcome.
