package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.guoyk.net/requo"
	"log"
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

func QueryTransaction(endpoint string, txid string, txoidx int) (contributes map[string]int64, err error) {
	contributes = make(map[string]int64)
	var resp QueryTransactionResponse
	if err = requo.JSONGet(context.Background(), endpoint+"/api/txs/summary/"+txid, &resp); err != nil {
		return
	}
	if err = resp.Check(); err != nil {
		return
	}
	var ok bool
	var lovelaceReceived int64
	for _, out := range resp.Right.Outputs {
		if *out.TxIndex == txoidx {
			lovelaceReceived = out.Amount.Value
			ok = true
			break
		}
	}
	if lovelaceReceived == 0 || !ok {
		err = fmt.Errorf("missing ok output")
		return
	}
	lovelaceRF := big.NewFloat(0).SetInt64(lovelaceReceived)
	totalOF := big.NewFloat(0).SetInt64(resp.Right.TotalOutput.Value)
	ratio := big.NewFloat(0).Quo(lovelaceRF, totalOF)

	for _, in := range resp.Right.Inputs {
		log.Printf("From: %s %d", in.Address, in.Amount.Value)
		f := big.NewFloat(0).SetInt64(in.Amount.Value)
		contrib, _ := big.NewFloat(0).Mul(ratio, f).Int64()
		contributes[in.Address] = contributes[in.Address] + contrib
		log.Printf("Contrib: %s %d", in.Address, contrib)
	}

	return
}
