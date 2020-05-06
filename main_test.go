package s3fs_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/require"

	"github.com/bishtawi/s3fs"
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

	fs, err := s3fs.NewWithS3Client(bucket, s3.New(sess))
	require.NoError(t, err)

	t.Run("root directory", func(t *testing.T) {
		dir, err := fs.Open("/")
		require.NoError(t, err)
		defer dir.Close()

		_, err = dir.Seek(0, io.SeekStart)
		require.Error(t, err)

		buf := new(bytes.Buffer)
		n, err := buf.ReadFrom(dir)
		require.Error(t, err)
		require.Zero(t, n)

		info, err := dir.Stat()
		require.NoError(t, err)
		require.True(t, info.IsDir())
		require.Equal(t, time.Time{}, info.ModTime())
		require.Equal(t, os.ModeDir, info.Mode())
		require.Equal(t, "/", info.Name())
		require.Zero(t, info.Size())
		require.Nil(t, info.Sys())

		fileInfos, err := dir.Readdir(0)
		require.NoError(t, err)
		require.Equal(t, 2, len(fileInfos))

		for _, fileInfo := range fileInfos {
			require.False(t, fileInfo.IsDir())
			require.NotEqual(t, time.Time{}, fileInfo.ModTime())
			require.NotEqual(t, os.ModeDir, fileInfo.Mode())
			require.NotEmpty(t, fileInfo.Name())
			require.NotZero(t, fileInfo.Size())
			require.NotNil(t, fileInfo.Sys())

			file, err := fs.Open(fileInfo.Name())
			require.NoError(t, err)
			defer file.Close()

			_, err = file.Readdir(0)
			require.Error(t, err)

			_, err = file.Seek(0, io.SeekStart)
			require.NoError(t, err)

			_, err = file.Seek(0, io.SeekEnd)
			require.Error(t, err)

			info, err := file.Stat()
			require.NoError(t, err)
			require.Equal(t, fileInfo.IsDir(), info.IsDir())
			require.Equal(t, fileInfo.ModTime().Truncate(time.Second), info.ModTime())
			require.Equal(t, fileInfo.Mode(), info.Mode())
			require.Equal(t, fileInfo.Name(), info.Name())
			require.Equal(t, fileInfo.Size(), info.Size())
			require.NotNil(t, info.Sys())

			buf := new(bytes.Buffer)
			n, err := buf.ReadFrom(file)
			require.NoError(t, err)
			require.Equal(t, fileInfo.Size(), n)

			_, err = file.Seek(0, io.SeekStart)
			require.Error(t, err)

			content := buf.String()

			switch fileInfo.Name() {
			case "one.txt":
				require.Equal(t, "one.txt contents\n", content)
			case "two.md":
				require.Equal(t, "two.md contents\n", content)
			default:
				require.Fail(t, "unexpected file found: "+fileInfo.Name())
			}
		}
	})

	t.Run("subdirectory", func(t *testing.T) {
		dir, err := fs.Open("folder/")
		require.NoError(t, err)
		defer dir.Close()

		_, err = dir.Seek(0, io.SeekStart)
		require.Error(t, err)

		buf := new(bytes.Buffer)
		n, err := buf.ReadFrom(dir)
		require.Error(t, err)
		require.Zero(t, n)

		info, err := dir.Stat()
		require.NoError(t, err)
		require.True(t, info.IsDir())
		require.Equal(t, time.Time{}, info.ModTime())
		require.Equal(t, os.ModeDir, info.Mode())
		require.Equal(t, "folder/", info.Name())
		require.Zero(t, info.Size())
		require.Nil(t, info.Sys())

		fileInfos, err := dir.Readdir(0)
		require.NoError(t, err)
		require.Equal(t, 2, len(fileInfos))

		for _, fileInfo := range fileInfos {
			require.False(t, fileInfo.IsDir())
			require.NotEqual(t, time.Time{}, fileInfo.ModTime())
			require.NotEqual(t, os.ModeDir, fileInfo.Mode())
			require.NotEmpty(t, fileInfo.Name())
			require.NotZero(t, fileInfo.Size())
			require.NotNil(t, fileInfo.Sys())

			file, err := fs.Open(fileInfo.Name())
			require.NoError(t, err)
			defer file.Close()

			_, err = file.Readdir(0)
			require.Error(t, err)

			_, err = file.Seek(0, io.SeekStart)
			require.NoError(t, err)

			_, err = file.Seek(0, io.SeekEnd)
			require.Error(t, err)

			info, err := file.Stat()
			require.NoError(t, err)
			require.Equal(t, fileInfo.IsDir(), info.IsDir())
			require.Equal(t, fileInfo.ModTime().Truncate(time.Second), info.ModTime())
			require.Equal(t, fileInfo.Mode(), info.Mode())
			require.Equal(t, fileInfo.Name(), info.Name())
			require.Equal(t, fileInfo.Size(), info.Size())
			require.NotNil(t, info.Sys())

			buf := new(bytes.Buffer)
			n, err := buf.ReadFrom(file)
			require.NoError(t, err)
			require.Equal(t, fileInfo.Size(), n)

			_, err = file.Seek(0, io.SeekStart)
			require.Error(t, err)

			content := buf.String()

			switch fileInfo.Name() {
			case "folder/a.md":
				require.Equal(t, "folder/a.md contents\n", content)
			case "folder/b.txt":
				require.Equal(t, "folder/b.txt contents\n", content)
			default:
				require.Fail(t, "unexpected file found: "+fileInfo.Name())
			}
		}
	})
}
