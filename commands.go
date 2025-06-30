package main

import (
	"fmt"

	"github.com/FlamestarRS/blogaggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handlerFunc, ok := c.cmds[cmd.name]
	if !ok {
		fmt.Println("Command does not exist")
		return nil
	}
	return handlerFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
