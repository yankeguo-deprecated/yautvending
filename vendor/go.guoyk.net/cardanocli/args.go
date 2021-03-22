//go:generate go run tools/generate-args.go
package cardanocli

// OptAddress append argument '--address'
func (c *Cmd) OptAddress(arg1 string) *Cmd {
    return c.Arg("--address", arg1)
}

// OptFee append argument '--fee'
func (c *Cmd) OptFee(arg1 string) *Cmd {
    return c.Arg("--fee", arg1)
}

// OptInvalidBefore append argument '--invalid-before'
func (c *Cmd) OptInvalidBefore(arg1 string) *Cmd {
    return c.Arg("--invalid-before", arg1)
}

// OptInvalidHereafter append argument '--invalid-hereafter'
func (c *Cmd) OptInvalidHereafter(arg1 string) *Cmd {
    return c.Arg("--invalid-hereafter", arg1)
}

// OptMainnet append argument '--mainnet'
func (c *Cmd) OptMainnet() *Cmd {
    return c.Arg("--mainnet")
}

// OptMaryEra append argument '--mary-era'
func (c *Cmd) OptMaryEra() *Cmd {
    return c.Arg("--mary-era")
}

// OptMint append argument '--mint'
func (c *Cmd) OptMint(arg1 string) *Cmd {
    return c.Arg("--mint", arg1)
}

// OptOutFile append argument '--out-file'
func (c *Cmd) OptOutFile(arg1 string) *Cmd {
    return c.Arg("--out-file", arg1)
}

// OptPaymentVerificationKeyFile append argument '--payment-verification-key-file'
func (c *Cmd) OptPaymentVerificationKeyFile(arg1 string) *Cmd {
    return c.Arg("--payment-verification-key-file", arg1)
}

// OptProtocolParamsFile append argument '--protocol-params-file'
func (c *Cmd) OptProtocolParamsFile(arg1 string) *Cmd {
    return c.Arg("--protocol-params-file", arg1)
}

// OptScriptFile append argument '--script-file'
func (c *Cmd) OptScriptFile(arg1 string) *Cmd {
    return c.Arg("--script-file", arg1)
}

// OptSigningKeyFile append argument '--signing-key-file'
func (c *Cmd) OptSigningKeyFile(arg1 string) *Cmd {
    return c.Arg("--signing-key-file", arg1)
}

// OptTestnetMagic append argument '--testnet-magic'
func (c *Cmd) OptTestnetMagic(arg1 string) *Cmd {
    return c.Arg("--testnet-magic", arg1)
}

// OptTxBodyFile append argument '--tx-body-file'
func (c *Cmd) OptTxBodyFile(arg1 string) *Cmd {
    return c.Arg("--tx-body-file", arg1)
}

// OptTxFile append argument '--tx-file'
func (c *Cmd) OptTxFile(arg1 string) *Cmd {
    return c.Arg("--tx-file", arg1)
}

// OptTxIn append argument '--tx-in'
func (c *Cmd) OptTxIn(arg1 string) *Cmd {
    return c.Arg("--tx-in", arg1)
}

// OptTxInCount append argument '--tx-in-count'
func (c *Cmd) OptTxInCount(arg1 string) *Cmd {
    return c.Arg("--tx-in-count", arg1)
}

// OptTxOut append argument '--tx-out'
func (c *Cmd) OptTxOut(arg1 string) *Cmd {
    return c.Arg("--tx-out", arg1)
}

// OptTxOutCount append argument '--tx-out-count'
func (c *Cmd) OptTxOutCount(arg1 string) *Cmd {
    return c.Arg("--tx-out-count", arg1)
}

// OptVerificationKeyFile append argument '--verification-key-file'
func (c *Cmd) OptVerificationKeyFile(arg1 string) *Cmd {
    return c.Arg("--verification-key-file", arg1)
}

// OptWitnessCount append argument '--witness-count'
func (c *Cmd) OptWitnessCount(arg1 string) *Cmd {
    return c.Arg("--witness-count", arg1)
}

// Address append argument 'address'
func (c *Cmd) Address() *Cmd {
    return c.Arg("address")
}

// Build append argument 'build'
func (c *Cmd) Build() *Cmd {
    return c.Arg("build")
}

// BuildRaw append argument 'build-raw'
func (c *Cmd) BuildRaw() *Cmd {
    return c.Arg("build-raw")
}

// CalculateMinFee append argument 'calculate-min-fee'
func (c *Cmd) CalculateMinFee() *Cmd {
    return c.Arg("calculate-min-fee")
}

// KeyGen append argument 'key-gen'
func (c *Cmd) KeyGen() *Cmd {
    return c.Arg("key-gen")
}

// KeyHash append argument 'key-hash'
func (c *Cmd) KeyHash() *Cmd {
    return c.Arg("key-hash")
}

// Policyid append argument 'policyid'
func (c *Cmd) Policyid() *Cmd {
    return c.Arg("policyid")
}

// ProtocolParameters append argument 'protocol-parameters'
func (c *Cmd) ProtocolParameters() *Cmd {
    return c.Arg("protocol-parameters")
}

// Query append argument 'query'
func (c *Cmd) Query() *Cmd {
    return c.Arg("query")
}

// Sign append argument 'sign'
func (c *Cmd) Sign() *Cmd {
    return c.Arg("sign")
}

// Submit append argument 'submit'
func (c *Cmd) Submit() *Cmd {
    return c.Arg("submit")
}

// Tip append argument 'tip'
func (c *Cmd) Tip() *Cmd {
    return c.Arg("tip")
}

// Transaction append argument 'transaction'
func (c *Cmd) Transaction() *Cmd {
    return c.Arg("transaction")
}

// Utxo append argument 'utxo'
func (c *Cmd) Utxo() *Cmd {
    return c.Arg("utxo")
}
