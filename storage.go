package main

import (
	"fmt"
	"io"
	"net/url"
)

var (
	storageAdapters = map[string]storageAdapterInitializer{}
)

type storageAdapter interface {
	Write(identifier string, data io.Reader) error
	Read(identifier string) (io.Reader, error)
	IsPresent(identifier string) bool
	Backup(identifier string) error
}
type storageAdapterInitializer func(*url.URL) (storageAdapter, error)

func getStorageAdapter(cfg *config) (storageAdapter, error) {
	storageURI, _ := cfg.ParsedStorage()

	if sa, ok := storageAdapters[storageURI.Scheme]; ok {
		s, err := sa(storageURI)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return nil, fmt.Errorf("Did not find storage adapter for '%s'", storageURI.Scheme)
}

func registerStorage(scheme string, f storageAdapterInitializer) error {
	if _, ok := storageAdapters[scheme]; ok {
		return fmt.Errorf("Cannot register '%s', is already registered", scheme)
	}

	storageAdapters[scheme] = f
	return nil
}
