package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	getCommand map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	getComm, exists := c.getCommand[cmd.name]
	if !exists {
		return fmt.Errorf("command not found")
	}

	err := getComm(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.getCommand[name] = f
}
