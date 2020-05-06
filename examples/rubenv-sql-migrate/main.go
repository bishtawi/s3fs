package rubenvsqlmigrateexample

import (
	"database/sql"

	"github.com/bishtawi/s3fs"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	_ "github.com/lib/pq" // required
	migrate "github.com/rubenv/sql-migrate"
)

// RubenvSQLMigrateExample runs Postgres SQL migrations from an S3 bucket and returns the database connection object.
func RubenvSQLMigrateExample(dbconn string, bucket string, s3client s3iface.S3API) (*sql.DB, int, error) {
	fs, err := s3fs.NewWithS3Client(bucket, s3client)
	if err != nil {
		return nil, 0, err
	}

	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		return db, 0, err
	}

	n, err := migrate.Exec(db, "postgres", &migrate.HttpFileSystemMigrationSource{FileSystem: fs}, migrate.Up)
	return db, n, err
}
