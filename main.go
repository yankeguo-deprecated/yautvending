package main

import (
	"encoding/json"
	"go.guoyk.net/cardanocli"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	optUTXOMaxBatch = 30
	optMinLovelace  = 1000000
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
	if addrDist, err = ReadFile(filepath.Join("addrs", "dist.addr")); err != nil {
		return
	}
	log.Println("AddrDist:", addrDist)

	// output mainnet parameters
	fileProtocol := filepath.Join(dir, "protocol.json")
	if err = cli.Cmd().Query().ProtocolParameters().OptMainnet().OptMaryEra().OptOutFile(fileProtocol).Exec().Run(); err != nil {
		return
	}

	// output utxos with cardano-cli
	fileUTXO := filepath.Join(dir, "utxo.json")
	if err = cli.Cmd().Query().Utxo().OptAddress(addrDist).OptMainnet().OptMaryEra().OptOutFile(fileUTXO).Exec().Run(); err != nil {
		return
	}

	// calculate transactions to handle

	var inputs []string
	{
		mapInputs := map[string]struct{}{}
		type UTXOOutput struct {
			Amount []json.RawMessage `json:"amount"`
		}

		var utxos map[string]UTXOOutput

		if err = ReadJSON(fileUTXO, &utxos); err != nil {
			return
		}

		var totalCount int
		var totalLovelace int64

		for tx, out := range utxos {
			if totalCount >= optUTXOMaxBatch {
				log.Println("Exceeding optUTXOMaxBatch")
				continue
			}
			log.Println("Processing:", tx)
			if len(out.Amount) == 0 {
				log.Println("Invalid number of amount in query utxos output")
				continue
			}
			var lovelace int64
			if lovelace, err = strconv.ParseInt(string(out.Amount[0]), 10, 64); err != nil {
				return
			}
			mapInputs[tx] = struct{}{}
			totalCount += 1
			totalLovelace += lovelace
		}

		if totalLovelace < optMinLovelace {
			log.Println("optMinLovelace not meet")
			return
		}

		for tx := range mapInputs {
			inputs = append(inputs, tx)
		}
		log.Println("Inputs:", "["+strings.Join(inputs, ",")+"]", ", Lovelace =", totalLovelace)
	}

}
