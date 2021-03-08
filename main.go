package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.guoyk.net/cardanocli"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	optSocketPath   = "/ipc/node.socket"
	optExplorerAPI  = "http://172.18.0.5:8100"
	optTokenName    = "YAUT"
	optUTXOMaxBatch = 30
	optMinLovelace  = 1000000
	optSigner       = "ce0b101709696dbc598485a670972e573c34689a2d3974d8c58337ab"
	optNotAfter     = "181306000"
	optPolicyScript = `{
  "type": "all",
  "scripts": [
    {
      "type": "sig",
      "keyHash": "` + optSigner + `"
    },
    {
      "type": "before",
      "slot": ` + optNotAfter + `
    }
  ]
}`
	optAddrGringotts = "addr1q9aemmfl4qr8sjp2xj5zupzvuamuw36z5awv865qt0lsl3pj72alpak07tadfuusgl5guq3ndtr3r2aknt4c3tgny7eqna8kkj"
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
	cli.SocketPath = optSocketPath

	// ensure tmp
	tmp := filepath.Join("tmp", strconv.FormatInt(time.Now().Unix(), 10)+"-"+RandomHex(16))
	if err = os.MkdirAll(tmp, 0755); err != nil {
		return
	}
	log.Println("TempDir:", tmp)
	defer os.RemoveAll(tmp)

	// generate policy file and policy id
	log.Println("Policy:", optPolicyScript)
	filePolicyScript := filepath.Join(tmp, "policy.script")
	if err = WriteFile(filePolicyScript, optPolicyScript); err != nil {
		return
	}

	var policyId string
	{
		x := cli.Cmd().Transaction().Policyid().OptScriptFile(filePolicyScript).Exec()
		out := &bytes.Buffer{}
		x.Stdout = out
		if err = x.Run(); err != nil {
			return
		}
		policyId = strings.TrimSpace(out.String())
		if policyId == "" {
			err = errors.New("invalid policy id")
			return
		}
	}
	log.Println("PolicyID:", policyId)

	// load dist.addr
	var addrDist string
	{
		x := cli.Cmd().Address().Build().OptPaymentVerificationKeyFile(filepath.Join("keys", "dist.vkey")).OptMainnet().Exec()
		out := &bytes.Buffer{}
		x.Stdout = out
		if err = x.Run(); err != nil {
			return
		}
		addrDist = strings.TrimSpace(out.String())
		if addrDist == "" {
			err = errors.New("invalid dist addr")
			return
		}
	}
	log.Println("AddrDist:", addrDist)

	// output mainnet parameters
	fileProtocol := filepath.Join(tmp, "protocol.json")
	if err = cli.Cmd().Query().ProtocolParameters().OptMainnet().OptMaryEra().OptOutFile(fileProtocol).Exec().Run(); err != nil {
		return
	}

	// output utxos with cardano-cli
	fileUTXO := filepath.Join(tmp, "utxo.json")
	if err = cli.Cmd().Query().Utxo().OptAddress(addrDist).OptMainnet().OptMaryEra().OptOutFile(fileUTXO).Exec().Run(); err != nil {
		return
	}

	// calculate transactions to handle

	var availableInputs []string
	var availableLovelace int64
	{
		mapInputs := map[string]struct{}{}
		type UTXOOutput struct {
			Amount []json.RawMessage `json:"amount"`
		}

		var utxos map[string]UTXOOutput

		if err = ReadJSON(fileUTXO, &utxos); err != nil {
			return
		}

		var count int
		for tx, out := range utxos {
			if count >= optUTXOMaxBatch {
				log.Println("Exceeding optUTXOMaxBatch")
				break
			}
			log.Println("Found:", tx)
			if len(out.Amount) == 0 {
				log.Println("Invalid number of amount in query utxos output")
				continue
			}
			var lovelace int64
			if lovelace, err = strconv.ParseInt(string(out.Amount[0]), 10, 64); err != nil {
				return
			}
			mapInputs[tx] = struct{}{}
			count += 1
			availableLovelace += lovelace
		}

		if availableLovelace < optMinLovelace {
			log.Println("optMinLovelace not meet")
			return
		}

		for tx := range mapInputs {
			availableInputs = append(availableInputs, tx)
		}
		log.Println("Inputs:", "["+strings.Join(availableInputs, ",")+"]", ", Lovelace =", availableLovelace)
	}

	log.Println("Available Lovelace", availableLovelace)

	tokenOuts := map[string]int64{}

	// calculate utxo distribution
	for _, input := range availableInputs {
		log.Println("Checking:", input)
		splits := strings.Split(input, "#")
		if len(splits) != 2 {
			log.Println("Invalid TX Input:", input)
		}
		var txid string
		var txidx int
		txid = splits[0]
		if txidx, err = strconv.Atoi(splits[1]); err != nil {
			return
		}
		log.Println("Split:", txid, txidx)
		var contributes map[string]int64
		if contributes, err = QueryTransaction(optExplorerAPI, txid, txidx); err != nil {
			return
		}
		for addr, contrib := range contributes {
			if contrib <= 0 {
				continue
			}
			log.Println("Contrib:", addr, contrib)
			tokenOuts[addr] = tokenOuts[addr] + contrib
		}
	}

	log.Println("Gringotts:", optAddrGringotts)

	// first build tx
	fileTxRaw := filepath.Join(tmp, "tx.raw")
	var countOut int
	var countMint int64
	{
		cmd := cli.Cmd().Transaction().BuildRaw().OptMaryEra().OptFee("0").OptInvalidHereafter(optNotAfter)
		for _, input := range availableInputs {
			cmd = cmd.OptTxIn(input)
		}
		for addr, tokenCount := range tokenOuts {
			countOut++
			countMint += tokenCount
			cmd = cmd.OptTxOut(fmt.Sprintf("%s+0+%d %s.%s", addr, tokenCount, policyId, optTokenName))
		}
		cmd = cmd.OptTxOut(fmt.Sprintf("%s+%d", optAddrGringotts, availableLovelace))
		countOut++
		cmd = cmd.OptMint(fmt.Sprintf("%d %s.%s", countMint, policyId, optTokenName))
		cmd.OptOutFile(fileTxRaw)
		if err = cmd.Exec().Run(); err != nil {
			return
		}
	}

	// calculate fee
	var fee int64
	out := &bytes.Buffer{}
	x := cli.Cmd().Transaction().CalculateMinFee().
		OptTxBodyFile(fileTxRaw).
		OptTxInCount(strconv.Itoa(len(availableInputs))).
		OptTxOutCount(strconv.Itoa(countOut)).
		OptWitnessCount("2").
		OptMainnet().
		OptProtocolParamsFile(fileProtocol).Exec()
	x.Stdout = out
	if err = x.Run(); err != nil {
		return
	}
	feeSplits := strings.Split(strings.TrimSpace(out.String()), " ")
	if len(feeSplits) != 2 {
		err = fmt.Errorf("invalid fee output: %s", out.String())
		return
	}
	if feeSplits[1] != "Lovelace" {
		err = fmt.Errorf("invalid fee unit: %s", feeSplits[1])
		return
	}
	if fee, err = strconv.ParseInt(feeSplits[0], 10, 64); err != nil {
		return
	}

	log.Println("Fee Calculated:", fee)

	if fee >= availableLovelace {
		err = errors.New("fee > available lovelace")
		return
	}

	{
		cmd := cli.Cmd().Transaction().BuildRaw().OptMaryEra().OptFee(strconv.FormatInt(fee, 10)).OptInvalidHereafter(optNotAfter)
		for _, input := range availableInputs {
			cmd = cmd.OptTxIn(input)
		}
		for addr, tokenCount := range tokenOuts {
			cmd = cmd.OptTxOut(fmt.Sprintf("%s+0+%d %s.%s", addr, tokenCount, policyId, optTokenName))
		}
		cmd = cmd.OptTxOut(fmt.Sprintf("%s+%d", optAddrGringotts, availableLovelace-fee))
		cmd = cmd.OptMint(fmt.Sprintf("%d %s.%s", countMint, policyId, optTokenName))
		cmd.OptOutFile(fileTxRaw)
		_ = json.NewEncoder(os.Stdout).Encode(cmd.Args)
		if err = cmd.Exec().Run(); err != nil {
			return
		}
	}

	log.Println("Final TX Built:", fileTxRaw)

	fileTxSigned := filepath.Join(tmp, "tx.signed")

	if err = cli.Cmd().Transaction().Sign().
		OptSigningKeyFile(filepath.Join("keys", "dist.skey")).
		OptSigningKeyFile(filepath.Join("keys", "issuer.skey")).
		OptScriptFile(filepath.Join("policy", "policy.script")).
		OptMainnet().
		OptTxBodyFile(fileTxRaw).
		OptOutFile(fileTxSigned).Exec().Run(); err != nil {
		return
	}

	var signed string
	if signed, err = ReadFile(fileTxSigned); err != nil {
		return
	}

	log.Println("Signed:", signed)

	if err = cli.Cmd().Transaction().Submit().OptTxFile(fileTxSigned).OptMainnet().Exec().Run(); err != nil {
		return
	}

	log.Println("Submitted")
}
