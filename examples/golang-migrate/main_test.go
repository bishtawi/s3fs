package golangmigrateexample

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/require"
)

const (
	localstackS3Endpoint = "http://localhost:4572"
	bucket               = "golang-migrate-example"
	dbconn               = "host=localhost port=5433 user=postgres dbname=postgres password=insecure-password sslmode=disable"
)

func Test_golangMigrateExample(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(localstackS3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
	})
	require.NoError(t, err)

	db, err := golangMigrateExample(dbconn, bucket, s3.New(sess))
	defer func() {
		if db != nil {
			db.Exec("DROP TABLE animal; DROP TABLE schema_migrations;")
			db.Close()
		}
	}()
	require.NoError(t, err)
}
