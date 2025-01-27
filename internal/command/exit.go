package command

import (
	"fmt"
	"os"
)

type ExitCommandHandler struct {
}

func NewExitCommand() *ExitCommandHandler {
	return &ExitCommandHandler{}
}

func (c *ExitCommandHandler) Execute(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
