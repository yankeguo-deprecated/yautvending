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

// NewCmd create a new Cmd instance
func NewCmd(cli *Cli) *Cmd {
	return &Cmd{Cli: cli}
}

// Exec build the final *exec.Cmd instance
func (c *Cmd) Exec() *exec.Cmd {
	x := exec.Command(c.Cli.Path, c.Args...)
	x.Env = os.Environ()
	if c.Cli.SocketPath != "" {
		x.Env = append(x.Env, fmt.Sprintf("CARDANO_NODE_SOCKET_PATH=%s", c.Cli.SocketPath))
	}
	x.Stdout = os.Stdout
	x.Stderr = os.Stderr
	return x
}

// Run run the command with hooks, most users should use this method
func (c *Cmd) Run(hooks ...Hook) (err error) {
	x := c.Exec()
	for _, hook := range hooks {
		hook.BeforeRun(x)
	}
	err = x.Run()
	for _, hook := range hooks {
		hook.AfterRun(x, &err)
	}
	return
}

// Arg append an argument
func (c *Cmd) Arg(args ...string) *Cmd {
	c.Args = append(c.Args, args...)
	return c
}
