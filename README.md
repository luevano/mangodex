# mangodex

Golang API wrapper for MangaDex v5's MVP API.

Full API documentation is found [here](https://api.mangadex.org/docs.html).

This branch contains only essential services, such as manga searching and image downloading.

**Note**: This is a fork of [Bob620/mangodex](https://github.com/Bob620/mangodex). Specially for use in [mangoprovider](https://github.com/luevano/mangoprovider) and [mangal](https://github.com/luevano/mangal).

## Installation

To install, do `go get -u github.com/darylhjd/mangodex@essential`.

## Usage

```golang
package main

import (
	"fmt"
	
	m "github.com/darylhjd/mangodex@simple"
)

func main() {
	// Create new client.
	// Without logging in, you may not be able to access 
	// all API functionality.
	c := m.NewDexClient()

	// Login using your username and password.
	err := c.Auth.Login("user", "password")
	if err != nil {
		fmt.Println("Could not login!")
	}
}
```

## Contributing

Any contributions are welcome.
