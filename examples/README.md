## Examples

Some example usages of the s3fs library working with other libraries that accept objects that implement the `http.Filesystem` interface.

To run the accompanying tests, you will need to bring the test environment via docker-compose, or modify the tests to connect to your own AWS S3 bucket and Postgres instance.

### golang-migrate

An example which uses the [golang-migrate](https://github.com/golang-migrate/migrate) library to run postgres migrations where the migration files are stored on AWS S3.

Golang-migrate library already supports AWS S3 natively but this example shows how it can be done via `http.Filesystem` allowing the Golang-migrate codebase to be simplified.

### rubenv-sql-migrate

Another example which uses [sql-migrate](https://github.com/rubenv/sql-migrate) library to run postgres migrations where the migration files are stored on AWS S3.

rubenv's sql-migrate library does not have built in support for AWS S3, thus demonstrating the entire purpose of s3fs.
