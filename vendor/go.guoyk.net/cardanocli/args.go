//go:generate go run tools/generate-args.go
package cardanocli

func (c *Cmd) OptAddress(arg1 string) *Cmd {
    return c.Append("--address", arg1)
}

func (c *Cmd) OptFee(arg1 string) *Cmd {
    return c.Append("--fee", arg1)
}

func (c *Cmd) OptMainnet() *Cmd {
    return c.Append("--mainnet")
}

func (c *Cmd) OptMaryEra() *Cmd {
    return c.Append("--mary-era")
}

func (c *Cmd) OptMint(arg1 string) *Cmd {
    return c.Append("--mint", arg1)
}

func (c *Cmd) OptOutFile(arg1 string) *Cmd {
    return c.Append("--out-file", arg1)
}

func (c *Cmd) OptPaymentVerificationKeyFile(arg1 string) *Cmd {
    return c.Append("--payment-verification-key-file", arg1)
}

func (c *Cmd) OptProtocolParamsFile(arg1 string) *Cmd {
    return c.Append("--protocol-params-file", arg1)
}

func (c *Cmd) OptScriptFile(arg1 string) *Cmd {
    return c.Append("--script-file", arg1)
}

func (c *Cmd) OptSigningKeyFile(arg1 string) *Cmd {
    return c.Append("--signing-key-file", arg1)
}

func (c *Cmd) OptTestnetMagic(arg1 string) *Cmd {
    return c.Append("--testnet-magic", arg1)
}

func (c *Cmd) OptTxBodyFile(arg1 string) *Cmd {
    return c.Append("--tx-body-file", arg1)
}

func (c *Cmd) OptTxFile(arg1 string) *Cmd {
    return c.Append("--tx-file", arg1)
}

func (c *Cmd) OptTxIn(arg1 string) *Cmd {
    return c.Append("--tx-in", arg1)
}

func (c *Cmd) OptTxInCount(arg1 string) *Cmd {
    return c.Append("--tx-in-count", arg1)
}

func (c *Cmd) OptTxOut(arg1 string) *Cmd {
    return c.Append("--tx-out", arg1)
}

func (c *Cmd) OptTxOutCount(arg1 string) *Cmd {
    return c.Append("--tx-out-count", arg1)
}

func (c *Cmd) OptVerificationKeyFile(arg1 string) *Cmd {
    return c.Append("--verification-key-file", arg1)
}

func (c *Cmd) OptWitnessCount(arg1 string) *Cmd {
    return c.Append("--witness-count", arg1)
}

func (c *Cmd) Address() *Cmd {
    return c.Append("address")
}

func (c *Cmd) Build() *Cmd {
    return c.Append("build")
}

func (c *Cmd) BuildRaw() *Cmd {
    return c.Append("build-raw")
}

func (c *Cmd) CalculateMinFee() *Cmd {
    return c.Append("calculate-min-fee")
}

func (c *Cmd) KeyGen() *Cmd {
    return c.Append("key-gen")
}

func (c *Cmd) KeyHash() *Cmd {
    return c.Append("key-hash")
}

func (c *Cmd) Policyid() *Cmd {
    return c.Append("policyid")
}

func (c *Cmd) ProtocolParameters() *Cmd {
    return c.Append("protocol-parameters")
}

func (c *Cmd) Query() *Cmd {
    return c.Append("query")
}

func (c *Cmd) Sign() *Cmd {
    return c.Append("sign")
}

func (c *Cmd) Submit() *Cmd {
    return c.Append("submit")
}

func (c *Cmd) Transaction() *Cmd {
    return c.Append("transaction")
}

func (c *Cmd) Utxo() *Cmd {
    return c.Append("utxo")
}
