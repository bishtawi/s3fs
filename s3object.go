package s3fs

import "os"

// s3object implements http.File
type s3object struct{}

func (o *s3object) Close() error {
	return errNotYetImplemented
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
	return &s3objectInfo{}, errNotYetImplemented
}
