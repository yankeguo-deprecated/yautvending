# cardanocli

a Golang wrapper for executing cardano-cli commands

I know it's stupid to invoke `cardano-cli` other than communicate with `cardano-node` via `node.socket` directly, but
this is maybe the best solution for Golang right now.

## Usage

```golang
package main

import (
	"go.guoyk.net/cardanocli"
)

func main() {
    cli := cardanocli.New()
    cli.SocketPath = "/path/to/cardano-node/node.socket"

    // example: get policy id from policy script
    {
    	var policyID string
        // cardano-cli transaction policyid --script-file /path/to/policy.script
        err := cli.Cmd().Transaction().
        	Policyid().
        	OptScriptFile("/path/to/policy.script").
        	Run(cardanocli.CollectStdout(&policyID))
        if err != nil {
            panic(err)
        }
    }
}
```

## Credits

Guo Y.K.， MIT License
