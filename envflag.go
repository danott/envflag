// Set flags via environment variables.
//
// Flags are still defined using the stdlib package flag. The only change to
// your code is calling envflag.Parse() in place of flag.Parse().
//
// If your flag wasn't set via command-line argument, an equivalent environment
// variable will be used.
//
// Precedence is: command-line agrument, environment variable, default.
//
// As an example, you can create a new file (main.go)
//
//  package main
//
//  import (
//  	"flag"
//  	"log"
//
//  	"github.com/danott/envflag"
//  )
//
//  func main() {
//  	var i int
//  	flag.IntVar(&i, "port", 2112, "Run on this port.")
//  	envflag.Parse()
//  	log.Printf("port: %v", i)
//  }
//
// Run your example to see the precedence in action:
//
//  go run main.go
//  go run main.go --port=2113
//  PORT=2114 go run main.go
//  PORT=2114 go run main.go --port=2113
package envflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	EnvfmtFlag    = "%[1]s"       // Envfmt for flag name
	EnvfmtProgram = "%[2]s_%[1]s" // Envfmt for program name and flag name
)

// Configure how flag names are translated to environment variable names.
// Accepts a string to be interpolated using Sprintf.
//
//  "%[1]s" - the flag name
//  "%[2]s" - the program name
var Envfmt = EnvfmtFlag

// Define your flags with package flag. Call envflag.Parse() in place of
// flag.Parse() to set flags via environment variables (if they weren't set via
// command-line arguments).
func Parse() {
	if !flag.Parsed() {
		flag.Parse()
	}

	for _, name := range defaultedFlags() {
		if value, ok := getenv(name); ok {
			flag.Set(name, value)
		}
	}
}

// Identical to os.Environ, but limited to the environment variable equivalents
// for the flags your program cares about.
func Environ() []string {
	s := make([]string, 0)

	flag.VisitAll(func(f *flag.Flag) {
		s = append(s, flagAsEnv(f.Name)+"="+os.Getenv(flagAsEnv(f.Name)))
	})

	return s
}

// Flags that were not set via command-line arguments, and have defaulted. It's
// smart enough to respect a flag that was set to the default value via
// command-line arguments. Must be called after flag.Parse()
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

	name = flagAsEnv(name)
	if _, ok = m[name]; ok {
		s = os.Getenv(name)
	}

	return
}

// To be unix'y, we translate flagnames to their uppercase equivalents.
func flagAsEnv(name string) string {
	name = strings.ToUpper(fmt.Sprintf(Envfmt, name, programName()))
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	return name
}

// The name of the currently running program
func programName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s)-1]
}
