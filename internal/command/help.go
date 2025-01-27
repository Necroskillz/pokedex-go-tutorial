package command

import "fmt"

type HelpCommandHandler struct {
	commands map[string]Command
}

func NewHelpCommand(commands map[string]Command) *HelpCommandHandler {
	return &HelpCommandHandler{commands}
}

func (c *HelpCommandHandler) Execute(args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range c.commands {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}
