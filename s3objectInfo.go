package s3fs

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// s3objectInfo implements os.FileInfo
type s3objectInfo struct {
	key            string
	s3ObjectOutput *s3.GetObjectOutput
}

func (i *s3objectInfo) Name() string {
	return i.key
}

func (i *s3objectInfo) Size() int64 {
	return aws.Int64Value(i.s3ObjectOutput.ContentLength)
}

func (i *s3objectInfo) Mode() os.FileMode {
	return os.ModeIrregular
}

func (i *s3objectInfo) ModTime() time.Time {
	return aws.TimeValue(i.s3ObjectOutput.LastModified)
}

func (i *s3objectInfo) IsDir() bool {
	return false
}

func (i *s3objectInfo) Sys() interface{} {
	return i.s3ObjectOutput
}
