// +build !appengine

package main

import "github.com/Luzifer/rconfig"

func loadConfig() *config {
	cfg := &config{}
	rconfig.Parse(cfg)
	return cfg
}
