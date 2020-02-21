package s3fs

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var errNotYetImplemented = errors.New("not yet implemented")

// New returns a new S3 http.FileSystem implementation
func New(s3client s3iface.S3API) http.FileSystem {
	return &s3fs{s3client}
}
