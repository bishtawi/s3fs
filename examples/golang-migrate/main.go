package golangmigrateexample

import (
	"database/sql"

	"github.com/bishtawi/s3fs"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/aws_s3" // Required import
	awss3 "github.com/golang-migrate/migrate/v4/source/aws_s3"
	_ "github.com/lib/pq" // required
)

// TODO: Switch to use filesystem instead of native aws

func golangMigrateExample(dbconn string, bucket string, s3client s3iface.S3API) (*sql.DB, error) {
	_, err := s3fs.NewWithS3Client(bucket, s3client)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		return db, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return db, err
	}

	source, err := awss3.WithInstance(s3client, &awss3.Config{Bucket: bucket})
	if err != nil {
		return db, err
	}

	defer source.Close()

	seeder, err := migrate.NewWithInstance("s3", source, "postgres", driver)
	if err != nil {
		return db, err
	}

	if err := seeder.Up(); err != nil && err != migrate.ErrNoChange {
		return db, err
	}

	return db, nil
}
