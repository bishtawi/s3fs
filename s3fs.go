package s3fs

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// s3fs implements http.FileSystem
type s3fs struct {
	bucket   string
	s3client s3iface.S3API
}

func (s *s3fs) Open(name string) (http.File, error) {
	getObjectOutput, err := s.s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	return &s3object{name, getObjectOutput}, nil
}
