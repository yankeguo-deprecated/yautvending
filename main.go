package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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
	optSocketPath       = "/ipc/node.socket"
	optExplorerAPI      = "http://172.18.0.5:8100"
	optTokenName        = "YAUT"
	optUTXOMaxBatch     = 30
	optMinInputLovelace = 2000000
	optBackLovelace     = 1100000
	optSigner           = "ce0b101709696dbc598485a670972e573c34689a2d3974d8c58337ab"
	optNotAfter         = "181306000"
	optPolicyScript     = `{
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

	var optSubmit bool
	flag.BoolVar(&optSubmit, "submit", false, "submit the signed transaction")
	flag.Parse()

	// create cli
	cli := cardanocli.New()
	cli.SocketPath = optSocketPath

	// ensure temp directory
	dirTemp := filepath.Join("dirTemp", strconv.FormatInt(time.Now().Unix(), 10)+"-"+RandomHex(16))
	if err = os.MkdirAll(dirTemp, 0755); err != nil {
		return
	}
	log.Println("DirTemp:", dirTemp)
	defer os.RemoveAll(dirTemp)

	// generate policy file and policy id
	log.Println("PolicyScript:", optPolicyScript)
	filePolicyScript := filepath.Join(dirTemp, "policy.script")
	if err = WriteFile(filePolicyScript, optPolicyScript); err != nil {
		return
	}

	var policyID string
	{
		x := cli.Cmd().Transaction().Policyid().OptScriptFile(filePolicyScript).Exec()
		out := &bytes.Buffer{}
		x.Stdout = out
		if err = x.Run(); err != nil {
			return
		}
		policyID = strings.TrimSpace(out.String())
		if policyID == "" {
			err = errors.New("invalid policy id")
			return
		}
	}
	log.Println("PolicyID:", policyID)

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
	fileProtocol := filepath.Join(dirTemp, "protocol.json")
	if err = cli.Cmd().Query().ProtocolParameters().OptMainnet().OptMaryEra().OptOutFile(fileProtocol).Exec().Run(); err != nil {
		return
	}

	// output utxos with cardano-cli
	fileUTXO := filepath.Join(dirTemp, "utxo.json")
	if err = cli.Cmd().Query().Utxo().OptAddress(addrDist).OptMainnet().OptMaryEra().OptOutFile(fileUTXO).Exec().Run(); err != nil {
		return
	}

	// calculate transactions to handle

	var utxoIDs []string
	var utxoLovelace int64
	{
		utxoIDMap := map[string]struct{}{}
		type UTXOFileEntry struct {
			Amount []json.RawMessage `json:"amount"`
		}

		var utxoEntries map[string]UTXOFileEntry

		if err = ReadJSON(fileUTXO, &utxoEntries); err != nil {
			return
		}

		var count int
		for utxoID, utxoEntry := range utxoEntries {
			if count >= optUTXOMaxBatch {
				log.Println("Exceeding optUTXOMaxBatch")
				break
			}
			log.Println("Found:", utxoID)
			if len(utxoEntry.Amount) == 0 {
				log.Println("Invalid number of amount in query utxos output")
				continue
			}
			var lovelace int64
			if lovelace, err = strconv.ParseInt(string(utxoEntry.Amount[0]), 10, 64); err != nil {
				return
			}
			utxoIDMap[utxoID] = struct{}{}
			utxoLovelace += lovelace
			count += 1
		}

		if utxoLovelace < optMinInputLovelace {
			log.Println("optMinInputLovelace not meet")
			return
		}

		for tx := range utxoIDMap {
			utxoIDs = append(utxoIDs, tx)
		}
		log.Println("Inputs:", "["+strings.Join(utxoIDs, ",")+"]", ", Lovelace =", utxoLovelace)
	}

	log.Println("Available Lovelace", utxoLovelace)

	// calculate utxo distribution
	type TokenOutput struct {
		Address string
		Amount  int64
	}

	var txTokenMint int64
	var txTokenOutputs []TokenOutput
	{
		txTokenOutputMap := map[string]int64{}

		for _, utxoID := range utxoIDs {
			log.Println("Checking:", utxoID)
			splits := strings.Split(utxoID, "#")
			if len(splits) != 2 {
				log.Println("Invalid TX Input:", utxoID)
			}
			var txID string
			var txIdx int
			txID = splits[0]
			if txIdx, err = strconv.Atoi(splits[1]); err != nil {
				return
			}
			log.Println("Split:", txID, txIdx)
			var tokenOutputs map[string]int64
			if tokenOutputs, err = QueryTransaction(optExplorerAPI, txID, txIdx, addrDist); err != nil {
				return
			}
			for addr, tokenOutput := range tokenOutputs {
				txTokenOutputMap[addr] = txTokenOutputMap[addr] + tokenOutput
			}
		}

		for addr, tokenOutput := range txTokenOutputMap {
			if tokenOutput <= (optMinInputLovelace - 1000) {
				continue
			}
			tokenOutput = tokenOutput - optBackLovelace
			txTokenMint += tokenOutput
			txTokenOutputs = append(txTokenOutputs, TokenOutput{Address: addr, Amount: tokenOutput})
		}
	}

	if len(txTokenOutputs) == 0 {
		log.Println("Nothing to output")
		return
	}

	log.Println("Gringotts:", optAddrGringotts)

	// first build tx
	fileTxRaw := filepath.Join(dirTemp, "tx.raw")
	{
		cmd := cli.Cmd().Transaction().BuildRaw().OptMaryEra().OptFee("0").OptInvalidHereafter(optNotAfter)
		for _, input := range utxoIDs {
			cmd = cmd.OptTxIn(input)
		}
		for _, tokenOutput := range txTokenOutputs {
			cmd = cmd.OptTxOut(fmt.Sprintf("%s+%d+%d %s.%s", tokenOutput.Address, optBackLovelace, tokenOutput.Amount, policyID, optTokenName))
		}
		cmd = cmd.OptTxOut(fmt.Sprintf("%s+%d", optAddrGringotts, utxoLovelace-int64(len(txTokenOutputs))*optBackLovelace))
		cmd = cmd.OptMint(fmt.Sprintf("%d %s.%s", txTokenMint, policyID, optTokenName))
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
		OptTxInCount(strconv.Itoa(len(utxoIDs))).
		OptTxOutCount(strconv.Itoa(len(txTokenOutputs) + 1)).
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

	if fee >= utxoLovelace {
		err = errors.New("fee > available lovelace")
		return
	}

	{
		cmd := cli.Cmd().Transaction().BuildRaw().OptMaryEra().OptFee(strconv.FormatInt(fee, 10)).OptInvalidHereafter(optNotAfter)
		for _, input := range utxoIDs {
			cmd = cmd.OptTxIn(input)
		}
		for _, tokenOutput := range txTokenOutputs {
			cmd = cmd.OptTxOut(fmt.Sprintf("%s+%d+%d %s.%s", tokenOutput.Address, optBackLovelace, tokenOutput.Amount, policyID, optTokenName))
		}
		cmd = cmd.OptTxOut(fmt.Sprintf("%s+%d", optAddrGringotts, utxoLovelace-fee-int64(len(txTokenOutputs))*optBackLovelace))
		cmd = cmd.OptMint(fmt.Sprintf("%d %s.%s", txTokenMint, policyID, optTokenName))
		cmd.OptOutFile(fileTxRaw)
		_ = json.NewEncoder(os.Stdout).Encode(cmd.Args)
		if err = cmd.Exec().Run(); err != nil {
			return
		}
	}

	log.Println("Final TX Built:", fileTxRaw)

	fileTxSigned := filepath.Join(dirTemp, "tx.signed")

	if err = cli.Cmd().Transaction().Sign().
		OptSigningKeyFile(filepath.Join("keys", "dist.skey")).
		OptSigningKeyFile(filepath.Join("keys", "issuer.skey")).
		OptScriptFile(filePolicyScript).
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

	if !optSubmit {
		return
	}

	if err = cli.Cmd().Transaction().Submit().OptTxFile(fileTxSigned).OptMainnet().Exec().Run(); err != nil {
		return
	}

	log.Println("Submitted")
}
