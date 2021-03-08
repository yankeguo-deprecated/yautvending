package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.guoyk.net/requo"
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
			TxIndex int     `json:"ctaTxIndex"`
		} `json:"ctxInputs"`
		Outputs []struct {
			Address string  `json:"ctaAddress"`
			Amount  GetCoin `json:"ctaAmount"`
			TxIndex int     `json:"ctaTxIndex"`
		} `json:"ctxOutputs"`
	} `json:"Right"`
}

func (resp QueryTransactionResponse) Check() (err error) {
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

/*
{
    "Right": {
        "ctsId": "2d814167828c27f43148adb8afd09efeac99ccb6d4a2868d5961c0178d86d205",
        "ctsTxTimeIssued": 1614564156,
        "ctsBlockTimeIssued": 1614564156,
        "ctsBlockHeight": 5403265,
        "ctsBlockEpoch": 250,
        "ctsBlockSlot": 361065,
        "ctsBlockHash": "f6a0a7065b31ae99918120a3e3f35fc7aae3508cfa4b2484f410963f4a8f0be5",
        "ctsRelayedBy": null,
        "ctsTotalInput": {
            "getCoin": "40983191022"
        },
        "ctsTotalOutput": {
            "getCoin": "40983018305"
        },
        "ctsFees": {
            "getCoin": "172717"
        },
        "ctsInputs": [
            {
                "ctaAddress": "DdzFFzCqrht6wUoxsseKN3bwqxaHrbQsUKzwMYA4KpMHoEbtyBRhnSXU71q4mj4KRJnssqAkHRw3V8kKSmTKzJdc1HJfjXJ94BsUdcyR",
                "ctaAmount": {
                    "getCoin": "40983191022"
                },
                "ctaTxHash": "95928b53aae97be6e8f39981ce6003ad7516136d8549e9091540d417870a872d",
                "ctaTxIndex": 1
            }
        ],
        "ctsOutputs": [
            {
                "ctaAddress": "Ae2tdPwUPEZHHcpw5Kz63Mh7SFapUfjNhjv9JXUZxukjc2AR1tUH92UMqpT",
                "ctaAmount": {
                    "getCoin": "17162323"
                },
                "ctaTxHash": "2d814167828c27f43148adb8afd09efeac99ccb6d4a2868d5961c0178d86d205",
                "ctaTxIndex": 1
            },
            {
                "ctaAddress": "DdzFFzCqrhseocHFPNhED2KXwgPMFkR28AWQTdTKYR95LitEVE1FSbT9hBh3592SU5WSENqt8bRYmW17qUCxCyBykKgrNKfoWMouGzy2",
                "ctaAmount": {
                    "getCoin": "40965855982"
                },
                "ctaTxHash": "2d814167828c27f43148adb8afd09efeac99ccb6d4a2868d5961c0178d86d205",
                "ctaTxIndex": 0
            }
        ]
    }
}
*/

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
		if out.TxIndex == txoidx {
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
		f := big.NewFloat(0).SetInt64(in.Amount.Value)
		contrib, _ := big.NewFloat(0).Mul(ratio, f).Int64()
		contributes[in.Address] = contributes[in.Address] + contrib
	}

	return
}
