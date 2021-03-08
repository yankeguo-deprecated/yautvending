package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.guoyk.net/requo"
	"log"
	"math"
	"math/big"
	"strconv"
)

type GetCoin struct {
	Value int64
}

func (g *GetCoin) UnmarshalJSON(buf []byte) (err error) {
	var s struct {
		V string `json:"getCoin"`
	}
	if err = json.Unmarshal(buf, &s); err != nil {
		return
	}
	if g.Value, err = strconv.ParseInt(s.V, 10, 64); err != nil {
		return
	}
	return
}

type QueryTransactionResponse struct {
	Right struct {
		TotalInput  GetCoin `json:"ctsTotalInput"`
		TotalOutput GetCoin `json:"ctsTotalOutput"`
		Inputs      []struct {
			Address string  `json:"ctaAddress"`
			Amount  GetCoin `json:"ctaAmount"`
			TxIndex *int    `json:"ctaTxIndex"`
		} `json:"ctsInputs"`
		Outputs []struct {
			Address string  `json:"ctaAddress"`
			Amount  GetCoin `json:"ctaAmount"`
			TxIndex *int    `json:"ctaTxIndex"`
		} `json:"ctsOutputs"`
	} `json:"Right"`
}

func (resp QueryTransactionResponse) Check() (err error) {
	{
		for _, ip := range resp.Right.Inputs {
			if ip.Address == "" {
				err = errors.New("input missing address")
				return
			}
			if ip.Amount.Value == 0 {
				err = errors.New("input missing amount")
				return
			}
			if ip.TxIndex == nil {
				err = errors.New("input missing ctaTxIndex")
				return
			}
		}
	}
	{
		for _, ip := range resp.Right.Outputs {
			if ip.Address == "" {
				err = errors.New("output missing address")
				return
			}
			if ip.Amount.Value == 0 {
				err = errors.New("output missing amount")
				return
			}
			if ip.TxIndex == nil {
				err = errors.New("output missing ctaTxIndex")
				return
			}
		}
	}

	{
		var ti int64
		for _, ip := range resp.Right.Inputs {
			ti += ip.Amount.Value
		}
		if ti != resp.Right.TotalInput.Value {
			err = errors.New("total input mismatch")
			return
		}
	}

	{
		var to int64
		for _, ip := range resp.Right.Outputs {
			to += ip.Amount.Value
		}
		if to != resp.Right.TotalOutput.Value {
			err = errors.New("total output mismatch")
			return
		}
	}
	return
}

func QueryTransaction(endpoint string, txid string, txoidx int, selfAddress string) (tokenOutputs map[string]int64, err error) {
	tokenOutputs = make(map[string]int64)
	var resp QueryTransactionResponse
	if err = requo.JSONGet(context.Background(), endpoint+"/api/txs/summary/"+txid, &resp); err != nil {
		return
	}
	if err = resp.Check(); err != nil {
		return
	}
	var ok bool
	var toMeOutput int64
	for _, out := range resp.Right.Outputs {
		if *out.TxIndex == txoidx {
			if out.Address != selfAddress {
				err = errors.New("not my address:" + out.Address)
				return
			}
			toMeOutput = out.Amount.Value
			ok = true
			break
		}
	}
	if toMeOutput == 0 || !ok {
		err = fmt.Errorf("missing ok output")
		return
	}
	toMeOutputF := big.NewFloat(0).SetInt64(toMeOutput)
	totalOutF := big.NewFloat(0).SetInt64(resp.Right.TotalOutput.Value)
	ratioFee := big.NewFloat(0).Quo(big.NewFloat(0).SetInt64(resp.Right.TotalOutput.Value), big.NewFloat(0).SetInt64(resp.Right.TotalInput.Value))
	ratioToMe := big.NewFloat(0).Quo(toMeOutputF, totalOutF)
	ratio := big.NewFloat(0).Mul(ratioToMe, ratioFee)

	for _, input := range resp.Right.Inputs {
		if input.Address == selfAddress {
			continue
		}
		log.Printf("From: %s %d", input.Address, input.Amount.Value)
		value := big.NewFloat(0).SetInt64(input.Amount.Value)
		lovelace, _ := big.NewFloat(0).Mul(ratio, value).Int64()
		lovelace = int64(math.RoundToEven(float64(lovelace)/10) * 10)
		if lovelace < (optMinInputLovelace - 1000) {
			continue
		}
		tokenOutputs[input.Address] = tokenOutputs[input.Address] + lovelace
		log.Printf("Received: %s %d", input.Address, lovelace)
	}

	return
}
