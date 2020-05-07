package golangmigrateexample_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	golangmigrate "github.com/bishtawi/s3fs/examples/golang-migrate"
)

const (
	localstackS3Endpoint = "http://localhost:4572"
	bucket               = "golang-migrate-example"
	dbconn               = "host=localhost port=5433 user=postgres dbname=postgres password=insecure-password sslmode=disable"
)

type migrationRecord struct {
	Version int
	Dirty   bool
}

func Test_golangMigrateExample(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(localstackS3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
	})
	require.NoError(t, err)

	db, err := golangmigrate.GolangMigrateExample(dbconn, bucket, s3.New(sess))

	defer func() {
		if db != nil {
			db.Exec("DROP TABLE schema_migrations; DROP TABLE animal;")
			db.Close()
		}
	}()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	records := []migrationRecord{}
	err = dbx.Select(&records, "SELECT * FROM schema_migrations")
	require.NoError(t, err)
	require.Equal(t, 1, len(records))
	require.False(t, records[0].Dirty)
	require.Equal(t, 2, records[0].Version)
}
