package rubenvsqlmigrateexample_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"

	rubenvSQLMigrate "github.com/bishtawi/s3fs/examples/rubenv-sql-migrate"
)

const (
	localstackS3Endpoint = "http://localhost:4572"
	bucket               = "rubenv-sql-migrate-example"
	dbconn               = "host=localhost port=5433 user=postgres dbname=postgres password=insecure-password sslmode=disable"
)

func Test_rubenvSQLMigrateExample(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(localstackS3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
	})
	require.NoError(t, err)

	db, n, err := rubenvSQLMigrate.RubenvSQLMigrateExample(dbconn, bucket, s3.New(sess))

	defer func() {
		if db != nil {
			db.Exec("DROP TABLE gorp_migrations; DROP TABLE people;")
			db.Close()
		}
	}()

	require.NoError(t, err)
	require.Equal(t, 2, n)

	dbx := sqlx.NewDb(db, "postgres")
	records := []migrate.MigrationRecord{}
	err = dbx.Select(&records, "SELECT * FROM gorp_migrations")
	require.NoError(t, err)
	require.Equal(t, n, len(records))
	require.Equal(t, "0001_create_people_table.sql", records[0].Id)
	require.Equal(t, "0002_add_name_column.sql", records[1].Id)
}
