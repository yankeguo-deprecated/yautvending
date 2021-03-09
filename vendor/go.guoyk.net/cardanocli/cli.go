package cardanocli

type Cli struct {
	Path       string
	SocketPath string
}

// New create a new instance
func New() *Cli {
	return &Cli{
		Path: "cardano-cli",
	}
}

// Cmd create a command instance
func (c *Cli) Cmd() *Cmd {
	return NewCmd(c)
}
