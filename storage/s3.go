package storage

import (
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	client *s3.S3
	bucket string
	root   string
}

func NewS3Storage(client *s3.S3, bucket, root string) *S3Storage {
	return &S3Storage{
		client: client,
		bucket: bucket,
		root:   root,
	}
}

func (s *S3Storage) Upload(path string, body io.ReadSeeker) error {
	key := filepath.Join(s.root, path)
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (s *S3Storage) Download(w io.Writer, path string) (err error) {
	key := filepath.Join(s.root, path)
	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return
	}

	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	return
}
