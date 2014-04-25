# package envflag

Set flags via environment variables.

## Usage

Flags are still defined using the stdlib `package "flag"`. The only change to your code is calling `envflag.Parse()` in place of `flag.Parse()`.

If your flag wasn't set via command-line argument, an equivalent environment variable will be used.

Precedence is: command-line agrument, environment variable, default.

## Example

```go
package main

import (
	"flag"
	"log"

	"github.com/danott/envflag"
)

func main() {
	var i int
	flag.IntVar(&i, "port", 2112, "Run on this port.")
	envflag.Parse()
	log.Printf("port: %v", i)
}
```

Run your example to see the precedence in action:

```bash
go run main.go
go run main.go --port=2113
PORT=2114 go run main.go
PORT=2114 go run main.go --port=2113
```
