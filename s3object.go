package s3fs

import (
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// s3object implements http.File
type s3object struct {
	s3client       s3iface.S3API
	bucket         string
	key            string
	s3ObjectOutput *s3.GetObjectOutput
	read           bool
}

func (o *s3object) Close() error {
	if o.s3ObjectOutput != nil {
		return o.s3ObjectOutput.Body.Close()
	}

	return nil
}

func (o *s3object) Read(p []byte) (n int, err error) {
	if o.s3ObjectOutput != nil {
		o.read = true
		return o.s3ObjectOutput.Body.Read(p)
	}

	return 0, os.ErrInvalid // Object is a directory, cannot read it
}

func (o *s3object) Seek(offset int64, whence int) (int64, error) {
	if (!o.read && offset == 0 && whence == io.SeekStart) || (offset == 0 && whence == io.SeekCurrent) {
		return 0, nil
	}

	return 0, ErrNotSupported
}

func (o *s3object) Readdir(count int) ([]os.FileInfo, error) {
	if o.s3ObjectOutput != nil {
		return nil, os.ErrInvalid // Object is a file, cannot read the directory
	}

	var marker *string

	var fileInfos []os.FileInfo

	for {
		listObjectsOutput, err := o.s3client.ListObjects(&s3.ListObjectsInput{
			Bucket:    aws.String(o.bucket),
			Delimiter: aws.String("/"),
			Prefix:    aws.String(strings.Trim(o.key, "/")),
			Marker:    marker,
			MaxKeys:   aws.Int64(int64(count)),
		})
		if err != nil {
			return nil, err
		}

		for _, object := range listObjectsOutput.Contents {
			fileInfos = append(fileInfos, &s3objectInfo{
				key:      aws.StringValue(object.Key),
				s3Object: object,
			})
		}

		if aws.BoolValue(listObjectsOutput.IsTruncated) {
			if listObjectsOutput.NextMarker != nil {
				marker = listObjectsOutput.NextMarker
			} else {
				marker = listObjectsOutput.Contents[len(listObjectsOutput.Contents)-1].Key
			}
		} else {
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
