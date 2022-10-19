package main

import (
	"flag"
	"fmt"

	"github.com/naujoyamat/schematics/internal"
	"github.com/naujoyamat/schematics/internal/command"
)

func main() {
	flag.Parse()

	regCommands := command.NewRegistry(&internal.NewCommand{}, &internal.TemplateCommand{})

	cmd := ""
	if len(flag.Args()) > 0 {
		cmd = flag.Arg(0)
	}

	if cmd == "" {
		fmt.Printf("usage:\n\tschematics [command] [args]\n\nexample:\n\tschematics new template_name project_name")
		return
	}

	args := flag.Args()[1:]

	err := regCommands.Execute(cmd, args...)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
