// +build appengine

package main

import (
	"os"

	"github.com/Luzifer/rconfig"
)

func loadConfig() *config {
	cfg := &config{}

	// Workaround as GAE supplies more parameters than expected.
	// This removes all CLI flags for parsing and relies only on parsing
	// environment variables.
	os.Args = []string{os.Args[0]}

	rconfig.Parse(cfg)
	return cfg
}
