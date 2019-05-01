package infrastructure

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"test_aws/interfaces"
	"time"
)

type S3Handler struct {
	Conn *s3.S3
	Sess *session.Session
	DbName string
}

type S3File struct {
	File *s3.Object
}

func (handler *S3Handler) Scan() (files []interfaces.RemoteFile, err error) {
	resp, err := handler.Conn.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(handler.DbName)})

	if err != nil {
		fmt.Printf("Unable to list items in bucket %q, %v", handler.DbName, err)
		return nil, err
	}

	files = make([]interfaces.RemoteFile, len(resp.Contents))
	for i, item := range resp.Contents {
		files[i] = &S3File{File: item}
	}

	return files, nil
}

func (handler *S3Handler) Download(dest *os.File, item string) (size int64, err error) {
	downloader := s3manager.NewDownloader(handler.Sess)

	numBytes, err := downloader.Download(dest,
		&s3.GetObjectInput{
			Bucket: aws.String(*aws.String(handler.DbName)),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Println("Unable to download item %q, %v", item, err)
		return 0, err
	}

	return numBytes, nil
}

func (file *S3File) Key() string {
	return *file.File.Key
}

func (file *S3File) LastMod() time.Time {
	return *file.File.LastModified
}

func (file *S3File) Etag() string {
	return *file.File.ETag
}

func (file *S3File) Size() int64 {
	return *file.File.Size
}

func NewS3Handler(appKey string, appSecret string, appRegion string, dbName string) *S3Handler {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(appRegion))},
	)
	if err != nil {
		fmt.Println(err)
		return new(S3Handler)
	}

	s3Handler := new(S3Handler)
	s3Handler.Conn = s3.New(sess)
	s3Handler.Sess = sess
	s3Handler.DbName = dbName

	return s3Handler
}