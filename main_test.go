package s3fs

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/require"
)

const (
	localstackS3Endpoint = "http://localhost:4572"
	bucket               = "integration-test"
)

func TestNewWithS3Client(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(localstackS3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
	})
	require.NoError(t, err)

	fs, err := NewWithS3Client(bucket, s3.New(sess))
	require.NoError(t, err)

	dir, err := fs.Open("/")
	require.NoError(t, err)

	files, err := dir.Readdir(0)
	require.NoError(t, err)
	require.Equal(t, 2, len(files))

	for _, file := range files {
		require.False(t, file.IsDir())
		require.NotEqual(t, time.Time{}, file.ModTime())
		require.NotEqual(t, os.ModeDir, file.Mode())
		require.NotEmpty(t, file.Name())
		require.NotZero(t, file.Size())
		require.NotNil(t, file.Sys())

		f, err := fs.Open(file.Name())
		require.NoError(t, err)

		defer f.Close()

		info, err := f.Stat()
		require.NoError(t, err)

		require.Equal(t, file.IsDir(), info.IsDir())
		require.Equal(t, file.Mode(), info.Mode())
		require.Equal(t, file.Name(), info.Name())
		require.Equal(t, file.Size(), info.Size())

		buf := new(bytes.Buffer)
		n, err := buf.ReadFrom(f)
		require.NoError(t, err)
		require.Equal(t, file.Size(), n)

		content := buf.String()
		require.True(t, content == "one.txt contents\n" || content == "two.md contents\n")
	}
}
