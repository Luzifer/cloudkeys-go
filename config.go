package main

import (
	"net/url"

	"github.com/Luzifer/rconfig"
)

type config struct {
	// General Config
	PasswordSalt string `env:"passwordSalt" flag:"password-salt" description:"A random unique salt for encrypting the passwords"`
	UsernameSalt string `env:"usernameSalt" flag:"username-salt" description:"A random unique salt for encrypting the usernames"`
	Storage      string `env:"storage" flag:"storage" default:"local:///./data" description:"Configuration for storage adapter (see README.md)"`
	Listen       string `flag:"listen" default:":3000" description:"IP and port to listen on"`

	CookieSigningKey string `flag:"cookie-authkey" description:"Key used to authenticate the session"`
	CookieEncryptKey string `flag:"cookie-encryptkey" description:"Key used to encrypt the session"`

	VersionAndQuit bool `flag:"version" default:"false" description:"Show version and quit"`
}

func (c config) ParsedStorage() (*url.URL, error) {
	return url.Parse(c.Storage)
}

func loadConfig() *config {
	cfg := &config{}
	rconfig.Parse(cfg)
	return cfg
}
