package babys3

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// A template for a generally publicly open bucket policy. Note the `%s` which
// should be replaced by the bucket name.
const publicReadBucketPolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "AddPerm",
			"Effect": "Allow",
			"Principal": "*",
			"Action": "s3:GetObject",
			"Resource": "arn:aws:s3:::%s/*"
		}
	]
}`

// AWS makes us do this before we do anything else
func newSession() (svc *s3.S3, err error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, errors.New("environment variable `AWS_REGION` not set")
	}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return s3.New(s), nil
}

// BucketExists determines if a given bucket exists
func BucketExists(name string) (exists bool, err error) {
	svc, err := newSession()
	if err != nil {
		return false, err
	}

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchBucket:
				// Bucket doesn't exist: no error.
				//
				// AWS tells us this is the error to expect, though we
				// personally didn't encounter it in our testing.
				return false, nil
			case "NotFound":
				// Bucket doesn't exist: no error.
				//
				// AWS seems to prefer this error over the former one, which is
				// contratictory to their docs as far as we can tell.
				return false, nil
			default:
				// Error.
				//
				// No way to tell if the bucket doesn't exist because we hit
				// some other error.
				return false, awsErr
			}
		}

		// This should probably never happen, but we should keep an eye out for
		// it anyway.
		return false, err
	}

	// Bucket exists
	return true, nil
}

// CreateBucket creates a new S3 bucket. At the time of writing (2018/05/14)
// creating a bucket that exists and is owned by you will succeed silently.
func CreateBucket(name string) error {
	svc, err := newSession()
	if err != nil {
		return err
	}

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
}

// MakeBucketPublic makes a bucket publicly readable over HTTP/HTTPS
func MakeBucketPublic(name string) error {
	svc, err := newSession()
	if err != nil {
		return err
	}

	_, err = svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(name),
		Policy: aws.String(fmt.Sprintf(publicReadBucketPolicy, name)),
	})

	return err
}

// DeleteBucket does exactly that
func DeleteBucket(name string) error {
	svc, err := newSession()
	if err != nil {
		return err
	}

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

// HeadObject fetches simple information about ab object without actually
// fetching the object.
//
// For information on what exactly is returned see:
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#HeadObjectOutput
func HeadObject(bucket, filename string) (*s3.HeadObjectOutput, error) {

	svc, err := newSession()
	if err != nil {
		return nil, err
	}

	h, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return nil, err
	} else if h == nil {
		return nil, fmt.Errorf("received no error, but also recieved no information about the object '%s' in bucket '%s'",
			filename, bucket)
	}

	return h, err
}

// ObjectExists returns true if the object exists on S3, and false if not
func ObjectExists(bucket, filename string) (exists bool, err error) {
	_, err = HeadObject(bucket, filename)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchKey:
				// Object doesn't exist. No error.
				return false, nil
			case "NotFound":
				// AWS might also return this?
				return false, nil
			default:
				return false, err
			}
		}
		// This should never happen
		return false, err
	}

	// No error: object exists
	return true, nil
}

// UploadObject uploads a given file to an S3 bucket
func UploadObject(bucket, filename string, data []byte) (err error) {
	svc, err := newSession()
	if err != nil {
		return err
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})

	return err
}

// UploadObjectAsBase64 uploads base64 data to a given S3 bucket. It's
// just a convenience function for `UploadObject()`
func UploadObjectAsBase64(bucket, filename, base64Data string) (err error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}
	return UploadObject(bucket, filename, data)
}

// DeleteObject deletes an object from an S3 bucket
func DeleteObject(bucket, filename string) error {
	svc, err := newSession()
	if err != nil {
		return err
	}

	deleteInput := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}
	_, err = svc.DeleteObject(&deleteInput)
	return err
}
