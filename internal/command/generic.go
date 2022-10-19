package command

import "fmt"

type GenericCommand struct {
	name     string
	function func(...string) error
}

func New(name string, function func(...string) error) *GenericCommand {
	return &GenericCommand{name, function}
}

func (g *GenericCommand) Name() string {
	return g.name
}

func (g *GenericCommand) Execute(args ...string) error {
	if g.function == nil {
		return fmt.Errorf("function to execute is nil")
	}

	return g.function(args...)
}
