package command

type Command struct {
	Name        string
	Description string
	Handler     CommandHandler
}

type CommandHandler interface {
	Execute(args []string) error
}
