package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xuyu/goredis"
)

func init() {
	registerStorage("redis+tcp", newRedisStorage)
	registerStorage("redis+udp", newRedisStorage)
}

// RedisStorage implements a storage option for redis server
type RedisStorage struct {
	conn   *goredis.Redis
	prefix string
}

// NewRedisStorage checks config, creates the path and initializes a RedisStorage
func newRedisStorage(u *url.URL) (storageAdapter, error) {
	client, err := goredis.DialURL(strings.Replace(u.String(), "redis+", "", -1))
	if err != nil {
		return nil, err
	}

	return &RedisStorage{
		conn:   client,
		prefix: u.Query().Get("prefix"),
	}, nil
}

// Write store the data of a dataObject into the storage
func (r *RedisStorage) Write(ctx context.Context, identifier string, data io.Reader) error {
	d, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	return r.conn.Set(r.prefix+identifier, string(d), 0, 0, false, false)
}

// Read reads the data of a dataObject from the storage
func (r *RedisStorage) Read(ctx context.Context, identifier string) (io.Reader, error) {
	content, err := r.conn.Get(r.prefix + identifier)
	return bytes.NewReader(content), err
}

// IsPresent checks for the presence of an userfile identifier
func (r *RedisStorage) IsPresent(ctx context.Context, identifier string) bool {
	e, err := r.conn.Exists(r.prefix + identifier)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"storagedriver": "redis",
			"identifier":    identifier,
		}).Error("Unable to check key existence")
	}
	return e && err == nil
}

// Backup creates a backup of the old data
func (r *RedisStorage) Backup(ctx context.Context, identifier string) error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	data, err := r.Read(ctx, identifier)
	if err != nil {
		return err
	}

	return r.Write(ctx, identifier+":backup:"+ts, data)
}
