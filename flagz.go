// Flagz adds automatic detection of environment variables to the flag
// package. It's a convenience for when you would like to arbitrarily set
// configuration via command-line arguments or environment variables.
//
// Note: environment variables are the strings.ToUpper() version of the flag
// name by convention.
//
// To use, call flagz.Parse() in place of flag.Parse()
//
// As an example, you can create a new file (main.go)
//
//  package main
//
//  import (
//  	"flag"
//  	"log"
//
//  	"github.com/danott/flagz"
//  )
//
//  func main() {
//  	var i int
//  	flag.IntVar(&i, "port", 2012, "This is the port we'll run on")
//  	flagz.Parse()
//  	log.Printf("port: %v", i)
//  }
//
// Run your example to see the precedence in action:
//
//  go run main.go
//  go run main.go --port=2013
//  PORT=2014 go run main.go
//  PORT=2014 go run main.go --port=2013
package flagz

import (
	"flag"
	"os"
	"strings"
)

// Setup your flags with the flag package as normal. Just call flagz.Parse
// instead of flag.Parse to get the bonus of reading environment variables too.
func Parse() {
	flag.Parse()

	for _, name := range defaultedFlags() {
		if value, ok := getenv(strings.ToUpper(name)); ok {
			flag.Set(name, value)
		}
	}
}

// Flags that were not set via command-line arguments, and have defaulted. It's
// smart enough to respect a flag that was set to the default value via
// command-line arguments.
func defaultedFlags() []string {
	m := make(map[string]bool)

	flag.VisitAll(func(f *flag.Flag) {
		m[f.Name] = true
	})

	flag.Visit(func(f *flag.Flag) {
		delete(m, f.Name)
	})

	s := make([]string, 0)

	for name, _ := range m {
		s = append(s, name)
	}

	return s
}

// Just like os.Getenv, but with a second return value; a boolean specifying
// if name was actually set in the environment.
func getenv(name string) (s string, ok bool) {
	m := make(map[string]bool)

	for _, keyVal := range os.Environ() {
		split := strings.Split(keyVal, "=")
		m[split[0]] = true
	}

	if _, ok = m[name]; ok {
		s = os.Getenv(name)
	}

	return
}
