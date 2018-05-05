package s3

import (
	"testing"
)

const (
	// 1x1 pixel PNG in base64 format for testing
	pixelBase64Data = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="
)

func TestNewSession(t *testing.T) {
	_, err := newSession()
	if err != nil {
		t.Error("error creating S3 session:", err.Error())
	}
}

func TestAllTheS3Stuff(t *testing.T) {

	bucketName := RandomBucketName()

	t.Run("test create bucket", func(t *testing.T) {
		if err := CreateBucket(bucketName); err != nil {
			t.Errorf("error creating bucket named %s: %s\n", bucketName, err.Error())
		}
	})

	t.Run("test make bucket public", func(t *testing.T) {
		if err := MakeBucketPublic(bucketName); err != nil {
			t.Errorf("error creating bucket named %s: %s\n", bucketName, err.Error())
		}
	})

	// Also testing that we can put stuff in a subdirectory
	filename := "pixels/pixel.png"

	// This also tests `Upload()` since `UploadBase64()` is just a wrapper
	// around `Upload()`
	t.Run("test upload file", func(t *testing.T) {
		if err := UploadBase64Object(bucketName, filename, pixelBase64Data); err != nil {
			t.Error("error uploading to S3", err.Error())
		}
	})

	t.Run("test delete file", func(t *testing.T) {
		if err := DeleteObject(bucketName, filename); err != nil {
			t.Error("error deleting from S3", err.Error())
			return
		}
	})

	t.Run("test delete bucket", func(t *testing.T) {
		if err := DeleteBucket(bucketName); err != nil {
			t.Errorf("error deleting bucket named %s: %s\n", bucketName, err.Error())
		}
	})
}
