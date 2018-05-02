package s3

import (
	"bytes"
	"encoding/base64"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWS makes us do this before we do anything else
func newSession() (svc *session.Session, err error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, errors.New("environment variable `AWS_REGION` not set")
	}

	return session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
}

// CreateBucket does exactly that
func CreateBucket(name string) error {
	s, err := newSession()
	if err != nil {
		return err
	}

	svc := s3.New(s)

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})

	return err
}

// DeleteBucket does exactly that
func DeleteBucket(name string) error {
	s, err := newSession()
	if err != nil {
		return err
	}

	svc := s3.New(s)

	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
}

// UploadObject uploads a given file to an S3 bucket
func UploadObject(bucket, filename string, data []byte) (err error) {
	s, err := newSession()
	if err != nil {
		return err
	}

	svc := s3.New(s)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})

	return err
}

// UploadBase64Object uploads base64 data to a given S3 bucket. It's just a
// convenience function for `UploadObject()`
func UploadBase64Object(bucket, filename, base64Data string) (err error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}
	return UploadObject(bucket, filename, data)
}

// DeleteObject deletes an object from an S3 bucket
func DeleteObject(bucket, filename string) error {
	s, err := newSession()
	if err != nil {
		return err
	}

	svc := s3.New(s)

	deleteInput := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}
	_, err = svc.DeleteObject(&deleteInput)
	return err
}