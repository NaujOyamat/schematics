package command

import "fmt"

type Registry struct {
	commands map[string]Command
}

func NewRegistry(commands ...Command) *Registry {
	reg := &Registry{
		commands: make(map[string]Command),
	}
	for _, c := range commands {
		reg.commands[c.Name()] = c
	}

	return reg
}

func (c *Registry) Get(name string) (Command, bool) {
	com, ok := c.commands[name]
	return com, ok
}

func (c *Registry) Register(command Command) {
	c.commands[command.Name()] = command
}

func (c *Registry) Execute(name string, args ...string) error {
	com, ok := c.Get(name)
	if !ok {
		return fmt.Errorf("command %s dont exists", name)
	}

	return com.Execute(args...)
}
