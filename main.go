package s3fs

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// ErrNotSupported is the error returned when a feature is not yet supported
var ErrNotSupported = errors.New("feature not supported")

// New returns a http.FileSystem object representing an AWS S3 bucket
func New(bucket string) (http.FileSystem, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return NewWithS3Client(bucket, s3.New(sess))
}

// NewWithS3Client returns a http.FileSystem object representing an AWS S3 bucket given an s3client
func NewWithS3Client(bucket string, s3client s3iface.S3API) (http.FileSystem, error) {
	return &s3fs{
		s3client: s3client,
		bucket:   bucket,
	}, nil
}
