package s3fs

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// s3objectInfo implements os.FileInfo.
type s3objectInfo struct {
	key            string
	s3ObjectOutput *s3.GetObjectOutput
	s3Object       *s3.Object
}

func (i *s3objectInfo) Name() string {
	return i.key
}

func (i *s3objectInfo) Size() int64 {
	if i.s3ObjectOutput != nil {
		return aws.Int64Value(i.s3ObjectOutput.ContentLength)
	}

	if i.s3Object != nil {
		return aws.Int64Value(i.s3Object.Size)
	}

	return 0
}

func (i *s3objectInfo) Mode() os.FileMode {
	if i.s3ObjectOutput != nil || i.s3Object != nil {
		return 0 // Is it worth trying to map S3 file permissions to the linux permissions?
	}

	return os.ModeDir
}

func (i *s3objectInfo) ModTime() time.Time {
	if i.s3ObjectOutput != nil {
		return aws.TimeValue(i.s3ObjectOutput.LastModified)
	}

	if i.s3Object != nil {
		return aws.TimeValue(i.s3Object.LastModified)
	}

	return time.Time{}
}

func (i *s3objectInfo) IsDir() bool {
	return i.s3ObjectOutput == nil && i.s3Object == nil
}

func (i *s3objectInfo) Sys() interface{} {
	if i.s3ObjectOutput != nil {
		return i.s3ObjectOutput
	}

	return i.s3Object
}
