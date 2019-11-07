// +build !appengine

package main

import "github.com/Luzifer/rconfig/v2"

func loadConfig() *config {
	cfg := &config{}
	rconfig.Parse(cfg)
	return cfg
}
