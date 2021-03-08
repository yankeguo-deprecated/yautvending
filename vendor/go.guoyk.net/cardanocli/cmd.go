package cardanocli

import (
	"fmt"
	"os"
	"os/exec"
)

type Cmd struct {
	*Cli
	Args []string
}

func NewCmd(cli *Cli) *Cmd {
	return &Cmd{Cli: cli}
}

func (c *Cmd) Exec() *exec.Cmd {
	x := exec.Command(c.Cli.Path, c.Args...)
	if c.Cli.SocketPath != "" {
		x.Env = os.Environ()
		x.Env = append(x.Env, fmt.Sprintf("CARDANO_NODE_SOCKET_PATH=%s", c.Cli.SocketPath))
	}
	return x
}

func (c *Cmd) Append(args ...string) *Cmd {
	c.Args = append(c.Args, args...)
	return c
}
