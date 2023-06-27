package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	uploader *s3manager.Uploader
	s3       *s3.S3
}

func new(uploader *s3manager.Uploader, s3 *s3.S3) *Storage {
	return &Storage{uploader: uploader, s3: s3}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Connect() *Storage {

	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	})
	uploader := s3manager.NewUploader(session)
	s3 := s3.New(session)

	if err != nil {
		panic(err)
	}

	credentials, _ := session.Config.Credentials.Get()
	log.Println("AWS Session Started: ", credentials)
	log.Println("S3 client started: ", uploader)

	return new(uploader, s3)
}

func (s *Storage) Upload(fileName string, file io.ReadSeeker) {

	AWS_VIDEO_UPLOAD_BUCKET := os.Getenv("AWS_VIDEO_UPLOAD_BUCKET")
	file.Seek(0, io.SeekStart)

	//write input to s3
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_VIDEO_UPLOAD_BUCKET),
		Body:   file,
		Key:    aws.String("inputvideos/" + fileName),
	})

	if err != nil {
		log.Println(err.Error())
		exitErrorf("Unable to upload: ", fileName, " to bucket: ", AWS_VIDEO_UPLOAD_BUCKET)
	}

	log.Println("Successfully uploaded: ", fileName, " to bucket: ", AWS_VIDEO_UPLOAD_BUCKET)
}

func (s *Storage) GetFileLink(fileName string) string {

	AWS_VIDEO_UPLOAD_BUCKET := os.Getenv("AWS_VIDEO_UPLOAD_BUCKET")

	request, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(AWS_VIDEO_UPLOAD_BUCKET),
		Key:    aws.String("inputvideos/" + fileName),
	})

	urlStr, err := request.Presign(30 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	return urlStr
}
