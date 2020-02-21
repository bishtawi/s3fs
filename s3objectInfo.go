package s3fs

import (
	"os"
	"time"
)

// s3objectInfo implements os.FileInfo
type s3objectInfo struct{}

func (i *s3objectInfo) Name() string {
	return "TODO"
}

func (i *s3objectInfo) Size() int64 {
	return -1
}

func (i *s3objectInfo) Mode() os.FileMode {
	return 0
}

func (i *s3objectInfo) ModTime() time.Time {
	return time.Time{}
}

func (i *s3objectInfo) IsDir() bool {
	return false
}

func (i *s3objectInfo) Sys() interface{} {
	return nil
}
