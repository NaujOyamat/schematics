package internal

import "fmt"

type NewCommand struct{}

func (n *NewCommand) Name() string {
	return "new"
}

func (n *NewCommand) Execute(args ...string) error {
	lenArgs := len(args)
	if lenArgs < 1 || lenArgs > 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	template := args[0]
	name := "project_name"

	if lenArgs == 2 {
		name = args[1]
	}

	return NewLoader(template).Build(name)
}
