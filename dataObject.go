package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const currentSchemaVersion = 1

type authorizedAccounts []authorizedAccount
type authorizedAccount struct {
	Name     string
	UserFile string
}

func init() {
	gob.Register(authorizedAccounts{})
}

type dataObject struct {
	SchemaVersion int `json:"schema_version"`
	MetaData      struct {
		Version   string `json:"version"`
		Password  string `json:"password"`
		MFASecret string `json:"mfa_secret"`
	} `json:"metadata"`
	Data string `json:"data"`
}

func newDataObject() *dataObject {
	return &dataObject{
		SchemaVersion: currentSchemaVersion,
	}
}

func dataObjectFromStorage(ctx context.Context, storage storageAdapter, filename string) (*dataObject, error) {
	userFileRaw, err := storage.Read(ctx, filename)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read data file from storage")
	}
	t := &dataObject{}
	return t, json.NewDecoder(userFileRaw).Decode(t)
}

// FIXME (kahlers): remove
func (d *dataObject) GetData() (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	return buf, json.NewEncoder(buf).Encode(d)
}

func (d *dataObject) migrate(ctx context.Context, storage storageAdapter, filename, password string) error {
	needsMigrate := true

	for needsMigrate {
		switch d.SchemaVersion {

		case 0: // Initial data file created before v2.0.0
			// Ensure a bcrypt hashed password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return errors.Wrap(err, "Unable to generate bcrypt hash")
			}
			d.MetaData.Password = string(hashedPassword)

		default: // No migration for this schema version defined, everything fine
			needsMigrate = false

		}

		// Increase schema version, see if there are more migrates
		d.SchemaVersion++
	}

	return d.writeToStorage(ctx, storage, filename)
}

func (d dataObject) writeToStorage(ctx context.Context, storage storageAdapter, filename string) error {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(d); err != nil {
		return errors.Wrap(err, "Unable to marshal data object")
	}
	return errors.Wrap(storage.Write(ctx, filename, buf), "Unable to write data file to storage")
}
