package s3fs

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// s3fs implements http.FileSystem.
type s3fs struct {
	s3client s3iface.S3API
	bucket   string
}

func (s *s3fs) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, "/") {
		return &s3object{
			s3client: s.s3client,
			bucket:   s.bucket,
			key:      name,
		}, nil
	}

	getObjectOutput, err := s.s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(strings.TrimLeft(name, "/")),
	})
	if err != nil {
		return nil, err
	}

	return &s3object{
		s3client:       s.s3client,
		bucket:         s.bucket,
		key:            name,
		s3ObjectOutput: getObjectOutput,
	}, nil
}
