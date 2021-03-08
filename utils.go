package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func RandomHex(n int) string {
	buf := make([]byte, n, n)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func ReadFile(f ...string) (c string, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(filepath.Join(f...)); err != nil {
		return
	}
	c = strings.TrimSpace(string(buf))
	return
}

func ReadJSON(f string, data interface{}) (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(f); err != nil {
		return
	}
	if err = json.Unmarshal(buf, data); err != nil {
		return
	}
	return
}
