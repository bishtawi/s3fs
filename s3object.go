package s3fs

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
)

// s3object implements http.File
type s3object struct {
	key            string
	s3ObjectOutput *s3.GetObjectOutput
}

func (o *s3object) Close() error {
	return o.s3ObjectOutput.Body.Close()
}

func (o *s3object) Read(p []byte) (n int, err error) {
	return -1, errNotYetImplemented
}

func (o *s3object) Seek(offset int64, whence int) (int64, error) {
	return -1, errNotYetImplemented
}

func (o *s3object) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errNotYetImplemented
}

func (o *s3object) Stat() (os.FileInfo, error) {
	return &s3objectInfo{o.key, o.s3ObjectOutput}, nil
}
