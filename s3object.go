package s3fs

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// s3object implements http.File.
type s3object struct {
	s3client       s3iface.S3API
	bucket         string
	key            string
	s3ObjectOutput *s3.GetObjectOutput
	read           bool
	lock           sync.Mutex
}

func (o *s3object) Close() error {
	if o.s3ObjectOutput != nil {
		return o.s3ObjectOutput.Body.Close()
	}

	return nil
}

func (o *s3object) Read(p []byte) (n int, err error) {
	if o.s3ObjectOutput != nil {
		o.lock.Lock()
		defer o.lock.Unlock()

		o.read = true

		return o.s3ObjectOutput.Body.Read(p)
	}

	return 0, os.ErrInvalid // Object is a directory, cannot read it
}

// Seek is technically unsupported due to the fact that the underlying S3 object does not implement the Seeker interface.
// However, we can support a few scenarios where no seeking actually occurs.
// 1. Seeking to the start of a file when no reading has occurred yet.
// 2. Seeking to the current position.
func (o *s3object) Seek(offset int64, whence int) (int64, error) {
	if o.s3ObjectOutput != nil {
		o.lock.Lock()
		defer o.lock.Unlock()

		if (!o.read && offset == 0 && whence == io.SeekStart) || (offset == 0 && whence == io.SeekCurrent) {
			return 0, nil
		}
	}

	return 0, ErrNotSupported
}

func (o *s3object) Readdir(count int) ([]os.FileInfo, error) {
	if o.s3ObjectOutput != nil {
		return nil, os.ErrInvalid // Object is a file, cannot read the directory
	}

	var (
		continuationToken *string
		fileInfos         []os.FileInfo
	)

	for {
		listObjectsV2Output, err := o.s3client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(o.bucket),
			Delimiter:         aws.String("/"),
			Prefix:            aws.String(strings.TrimLeft(o.key, "/")),
			ContinuationToken: continuationToken,
			MaxKeys:           aws.Int64(int64(count)),
		})
		if err != nil {
			return nil, err
		}

		for _, object := range listObjectsV2Output.Contents {
			fileInfos = append(fileInfos, &s3objectInfo{
				key:      aws.StringValue(object.Key),
				s3Object: object,
			})
		}

		continuationToken = listObjectsV2Output.NextContinuationToken

		if !aws.BoolValue(listObjectsV2Output.IsTruncated) || listObjectsV2Output.NextContinuationToken == nil {
			break
		}
	}

	return fileInfos, nil
}

func (o *s3object) Stat() (os.FileInfo, error) {
	return &s3objectInfo{
		key:            o.key,
		s3ObjectOutput: o.s3ObjectOutput,
	}, nil
}
