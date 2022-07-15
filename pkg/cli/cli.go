package cli

import (
	"flag"
)

type Command struct {
	ptr          any
	Value        any
	Name         string
	DefaultValue any
	Usage        string
}

type CLI struct {
	flag.FlagSet
}

func (c *CLI) Add(name string, defval string, usage string) {
	flag.Set(name, defval)
}

func (c *CLI) Execute() error {

	return nil
}
