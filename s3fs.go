package s3fs

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// s3fs implements http.FileSystem
type s3fs struct {
	s3client s3iface.S3API
}

func (s *s3fs) Open(name string) (http.File, error) {
	return &s3object{}, errNotYetImplemented
}
