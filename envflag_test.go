package envflag

import (
	"flag"
	"os"
	"testing"
)

var (
	name = ""
)

func setup() {
	os.Clearenv()
	EnvPrefix = ""
	FlagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	FlagSet.StringVar(&name, "name", "Michael", "This is a member of the bluth family")
}

func TestParse(t *testing.T) {
	setup()

	os.Setenv("NAME", "Gob")
	Parse()

	if name != "Gob" {
		t.Errorf("Selfish %s.", name)
	}
}

func TestParseWithCustomEnvfmt(t *testing.T) {
	setup()

	EnvPrefix = "CUSTOM_"
	os.Setenv("NAME", "Gob")
	os.Setenv("CUSTOM_NAME", "Tobias")
	Parse()

	if name != "Tobias" {
		t.Errorf("%s, you blow hard", name)
	}
}

func TestParsePrecedence(t *testing.T) {
	setup()
	os.Setenv("NAME", "George Michael")

	// You'd thing this needs to be -name="George Oscar Bluth", but it doesn't.
	// Just trust me.
	parse([]string{"-name=George Oscar Bluth"})

	if name != "George Oscar Bluth" {
		t.Errorf("Lemons! %s", name)
	}

	setup()
	os.Setenv("NAME", "George Michael")
	Parse()

	if name != "George Michael" {
		t.Errorf("Lemons! %s", name)
	}
}

func TestEnvironEmpty(t *testing.T) {
	setup()

	if len(Environ()) != 0 {
		t.Errorf("There shouldn't be anything here. %v", Environ())
	}
}

func TestEnvironSet(t *testing.T) {
	setup()

	os.Setenv("NAME", "Maeby")
	os.Setenv("MONEY_IN_THE_BANANA_STAND", "1")

	if len(Environ()) != 1 {
		t.Errorf("There shouldn only be one thing herone thing here %v", Environ())
	}
}
