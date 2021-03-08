package main

import (
	"encoding/json"
	"go.guoyk.net/cardanocli"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Println("exited with error:", (*err).Error())
			os.Exit(1)
		} else {
			log.Println("exited")
		}
	}(&err)

	// create cli
	cli := cardanocli.New()
	cli.SocketPath = "/ipc/node.socket"

	// ensure dir
	dir := filepath.Join("tmp", strconv.FormatInt(time.Now().Unix(), 10)+"-"+RandomHex(16))
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}
	log.Println("TempDir:", dir)
	defer os.RemoveAll(dir)

	// load dist.addr
	var addrDist string
	if addrDist, err = ReadFile("addrs", "dist.addr"); err != nil {
		return
	}
	log.Println("AddrDist:", addrDist)

	// output utxos with cardano-cli
	utxoFile := filepath.Join(dir, "utxo.json")

	if err = cli.Cmd().Query().Utxo().OptAddress(addrDist).OptMainnet().OptMaryEra().OptOutFile(utxoFile).Exec().Run(); err != nil {
		return
	}

	type UTXOOutput struct {
		Amount []json.RawMessage `json:"amount"`
	}

	var utxos map[string]UTXOOutput

	if err = ReadJSON(utxoFile, &utxos); err != nil {
		return
	}

	for tx, out := range utxos {
		if len(out.Amount) == 0 {
			continue
		}
		log.Println("UTXO:", tx, string(out.Amount[0]))
	}

}
