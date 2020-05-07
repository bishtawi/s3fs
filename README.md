# S3FS

Golang library that implements the [http.FileSystem](https://godoc.org/net/http#FileSystem) interface on top of AWS S3.

![](https://github.com/bishtawi/s3fs/workflows/test/badge.svg)

## Purpose

Why would you use this library instead of directly interacting with the AWS S3 Go library?

Good question! For 99.99% of usecases, you would want to use the official AWS S3 Go library instead of this package.

The only real reason why you would use this package is when you are using some arbitrary Go library that does not have baked-in support for AWS S3 but does provide compatibility with the [http.FileSystem](https://godoc.org/net/http#FileSystem) interface.

The goal of this tool is to implement the `http.FileSystem` interface on top of AWS S3 to allow one to plug in the AWS S3 filesystem into any Go library that supports the `http.FileSystem` interface as input.

Take a look at the [examples](./examples) directory to see how you would use this library with other Go libraries to bring AWS S3 support to a package that might not natively support it.

## Install

```bash
go get github.com/bishtawi/s3fs
```

## Usage

```go
bucketName := "bucket-name"

// Create AWS session object
sess, err := session.NewSession(&aws.Config{...})
if err != nil {
    // Handle error...
}

// Create AWS S3 client
s3client := s3.New(sess)

// Create the s3fs filesystem object
fs, err := s3fs.NewWithS3Client(bucketName, s3client)
if err != nil {
    // Handle error...
}

// The returned object, fs, implements the http.FileSystem interface
// You can then execute all the standard http.FileSystem functions to interact with the AWS S3 bucket

// Open the root directory
dir, err := fs.Open("/")
if err != nil {
    // Handle error...
}
defer dir.Close()

// List all files in the directory
fileInfos, err := dir.Readdir(0)
if err != nil {
    // Handle error...
}

// Open the first file
file, err := fs.Open(fileInfos[0].Name())
if err != nil {
    // Handle error...
}
defer file.Close()

// Read the contents to a buffer
buf := new(bytes.Buffer)
n, err := buf.ReadFrom(file)
if err != nil {
    // Handle error...
}

// Or simply use the fs object as input to a Go library that accepts http.FileSystem objects

// Execute SQL migration files that live in an S3 bucket
// Using rubenv's sql-migrate library: https://github.com/rubenv/sql-migrate
dbconn := "host=localhost port=5432 user=postgres dbname=postgres"
db, err := sql.Open("postgres", dbconn)
if err != nil {
    // Handle error...
}
defer db.Close()

n, err := migrate.Exec(db, "postgres", &migrate.HttpFileSystemMigrationSource{FileSystem: fs}, migrate.Up)
if err != nil {
    // Handle error...
}

// View other example usages in the `examples/` directory
```

## Limitations

All of the limitations below can be addressed, just requires some extra effort to solve.

1. The `Readdir(count int) ([]os.FileInfo, error)` function of `http.File` interface only returns a list of files in the immediate directory and will not return a list of subdirectories due to a limitation of AWS S3's `ListObjects` operation.

1. The `Seek(offset int64, whence int) (int64, error)` function of `http.File` interface is not implemented (stubbed function only) as the underlying AWS S3 objects do not implement the `io.Seeker` interface.

1. The `Size() int64` function of `os.FileInfo` interface only works on files. `0` is returned for directories.

1. The `ModTime() time.Time` function of `os.FileInfo` interface only works on files. Empty struct is returned for directories.

1. The `Mode() os.FileMode` function of `os.FileInfo` interface is another stub function that returns `0` as I don't know if it is worth trying to map S3 permissions to the linux permission model.

## Development

Open to contributions via Pull Requests.

### Testing

Due to this library essentially being a lightweight wrapper around the AWS S3 library, all tests are integration tests and require connecting to an S3 bucket.
To make testing easier, we are using [localstack](https://github.com/localstack/localstack) which allows us to bring up a local S3 bucket instead of directly testing on a live S3 bucket.

In addition to testing against S3, we also test compatibility with other Go libraries that use the `http.FileSystem` interface as an input. These tests are located in the [examples](./examples) directory and require connecting to a running postgres instance.

For simplicity, we use docker-compose to easily bring up localstack and postgres.

```bash
docker-compose up
make test
```

### Linting

- We use [golangci-linter](https://github.com/golangci/golangci-lint) to lint the Go code.
- [Shellcheck](https://github.com/koalaman/shellcheck) to lint bash scripts.
- And [Prettier](https://github.com/prettier/prettier) to format non-code files like json, yaml and markdown files.

```bash
make format
make lint
```
