package dalkeeth

type Command struct {
}

func (c *Command) New() *Command {
	return new(Command)
}

func (c *Command) Select() *Command {
	return c
}
