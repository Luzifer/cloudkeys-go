package main

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	storage     storageAdapter
	cookieStore *sessions.CookieStore

	cfg     = loadConfig()
	version = "dev"
)

func initialize() {
	if cfg.VersionAndQuit {
		fmt.Printf("cloudkeys-go %s\n", version)
		os.Exit(0)
	}

	if _, err := cfg.ParsedStorage(); err != nil {
		log.WithError(err).Error("Unable to parse storage URI")
		os.Exit(1)
	}

	if cfg.CookieSigningKey == "" {
		cfg.CookieSigningKey = uuid.NewV4().String()[:32]
		log.Warn("cookie-authkey was set randomly, this will break your sessions!")
	}

	if cfg.CookieEncryptKey == "" {
		cfg.CookieEncryptKey = uuid.NewV4().String()[:32]
		log.Warn("cookie-encryptkey was set randomly, this will break your sessions!")
	}

	cookieStore = sessions.NewCookieStore(
		[]byte(cfg.CookieSigningKey),
		[]byte(cfg.CookieEncryptKey),
	)
}

func initializeStorage() {
	s, err := getStorageAdapter(cfg)
	if err != nil {
		log.WithError(err).Fatal("Could not instanciate storage")
	}
	storage = s
}

func createUserFilename(username string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(cfg.UsernameSalt+username)))
}
