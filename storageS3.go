package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func init() {
	registerStorage("s3", newS3Storage)
}

// S3Storage implements a storage option for Amazon S3
type S3Storage struct {
	bucket string
	path   string
	conn   *s3.S3
}

// NewS3Storage checks config, creates the path and initializes a S3Storage
func newS3Storage(u *url.URL) (storageAdapter, error) {
	return &S3Storage{
		bucket: u.Host,
		path:   u.Path,
		conn:   s3.New(session.New()),
	}, nil
}

// Write store the data of a dataObject into the storage
func (s *S3Storage) Write(identifier string, data io.Reader) error {
	buf := bytes.NewBuffer([]byte{})
	io.Copy(buf, data)

	_, err := s.conn.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Body:   bytes.NewReader(buf.Bytes()),
		Key:    aws.String(path.Join(s.path, identifier)),
	})
	return err
}

// Read reads the data of a dataObject from the storage
func (s *S3Storage) Read(identifier string) (io.Reader, error) {
	out, err := s.conn.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path.Join(s.path, identifier)),
	})

	if err != nil {
		return nil, err
	}

	return out.Body, nil
}

// IsPresent checks for the presence of an userfile identifier
func (s *S3Storage) IsPresent(identifier string) bool {
	out, err := s.conn.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path.Join(s.path, identifier)),
	})

	if err != nil {
		return false
	}

	return *out.ContentLength > 0
}

// Backup creates a backup of the old data
func (s *S3Storage) Backup(identifier string) error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	_, err := s.conn.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		Key:        aws.String(path.Join(s.path, "backup", fmt.Sprintf("%s.%s", identifier, ts))),
		CopySource: aws.String(path.Join(s.bucket, s.path, identifier)),
	})

	return err
}
