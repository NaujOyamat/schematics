package internal

import (
	"fmt"
)

type TemplateCommand struct{}

func (t *TemplateCommand) Name() string {
	return "template"
}

func (t *TemplateCommand) Execute(args ...string) error {
	lenArgs := len(args)
	if lenArgs != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	name := args[0]

	return NewLoader(name).CreateTemplate(name)
}
