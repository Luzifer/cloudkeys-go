package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
)

type authorizedAccounts []authorizedAccount
type authorizedAccount struct {
	Name     string
	UserFile string
}

func init() {
	gob.Register(authorizedAccounts{})
}

type dataObject struct {
	MetaData struct {
		Version  string `json:"version"`
		Password string `json:"password"`
	} `json:"metadata"`
	Data string `json:"data"`
}

func readDataObject(in io.Reader) (*dataObject, error) {
	t := &dataObject{}
	return t, json.NewDecoder(in).Decode(t)
}

func (d *dataObject) GetData() (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	return buf, json.NewEncoder(buf).Encode(d)
}
