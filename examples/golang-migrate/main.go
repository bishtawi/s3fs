package golangmigrateexample

import (
	"database/sql"

	"github.com/bishtawi/s3fs"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq" // required
)

// GolangMigrateExample runs Postgres migrations from an S3 bucket and returns the database connection object.
func GolangMigrateExample(dbconn string, bucket string, s3client s3iface.S3API) (*sql.DB, error) {
	fs, err := s3fs.NewWithS3Client(bucket, s3client)
	if err != nil {
		return nil, err
	}

	source, err := httpfs.New(fs, "/")
	if err != nil {
		return nil, err
	}
	defer source.Close()

	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		return db, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return db, err
	}

	migrater, err := migrate.NewWithInstance("httpfs", source, "postgres", driver)
	if err != nil {
		return db, err
	}

	return db, migrater.Up()
}
