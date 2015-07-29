package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

func init() {
	registerStorage("local", newLocalStorage)
}

// LocalStorage implements a storage option for local file storage
type LocalStorage struct {
	path string
}

// NewLocalStorage checks config, creates the path and initializes a LocalStorage
func newLocalStorage(u *url.URL) (storageAdapter, error) {
	p := u.Path[1:]

	if len(p) == 0 {
		return nil, fmt.Errorf("Path not present.")
	}

	if err := os.MkdirAll(path.Join(p, "backup"), 0755); err != nil {
		return nil, fmt.Errorf("Unable to create path '%s'", p)
	}

	return &LocalStorage{
		path: p,
	}, nil
}

// Write store the data of a dataObject into the storage
func (l *LocalStorage) Write(identifier string, data io.Reader) error {
	f, err := os.Create(path.Join(l.path, identifier))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	return err
}

// Read reads the data of a dataObject from the storage
func (l *LocalStorage) Read(identifier string) (io.Reader, error) {
	return os.Open(path.Join(l.path, identifier))
}

// IsPresent checks for the presence of an userfile identifier
func (l *LocalStorage) IsPresent(identifier string) bool {
	_, err := os.Stat(path.Join(l.path, identifier))
	return err == nil
}

// Backup creates a backup of the old data
func (l *LocalStorage) Backup(identifier string) error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	o, err := os.Open(path.Join(l.path, identifier))
	if err != nil {
		return err
	}
	n, err := os.Create(path.Join(l.path, "backup", fmt.Sprintf("%s.%s", identifier, ts)))
	if err != nil {
		return err
	}

	defer o.Close()
	defer n.Close()

	_, err = io.Copy(n, o)
	return err
}
