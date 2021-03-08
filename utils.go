package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"
)

func RandomHex(n int) string {
	buf := make([]byte, n, n)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func ReadFile(f string) (c string, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(f); err != nil {
		return
	}
	c = strings.TrimSpace(string(buf))
	return
}

func WriteFile(f string, content string) (err error) {
	err = ioutil.WriteFile(f, []byte(strings.TrimSpace(content)), 0640)
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
