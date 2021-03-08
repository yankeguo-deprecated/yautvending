package cardanocli

import "os"

type Cli struct {
	Path       string
	SocketPath string
}

func New() *Cli {
	return &Cli{
		Path:       "cardano-cli",
		SocketPath: os.Getenv("CARDANO_NODE_SOCKET_PATH"),
	}
}

func (c *Cli) Cmd() *Cmd {
	return NewCmd(c)
}
