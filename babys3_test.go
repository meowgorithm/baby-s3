package babys3

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

	t.Run("eh", func(t *testing.T) {
		if _, err := HeadObject(bucketName, "eh"); err != nil {
			t.Errorf("wut %s\n", err)
		}
	})

	t.Run("test create bucket", func(t *testing.T) {
		if err := CreateBucket(bucketName); err != nil {
			t.Errorf("error creating bucket named %s: %s\n", bucketName, err)
		}
	})

	t.Run("test make bucket public", func(t *testing.T) {
		if err := MakeBucketPublic(bucketName); err != nil {
			t.Errorf("error creating bucket named %s: %s\n", bucketName, err)
		}
	})

	textFilename := "meow.txt"

	t.Run("test upload file as bytes", func(t *testing.T) {
		err := UploadObject(bucketName, textFilename, []byte("here, kitty kitty"))
		if err != nil {
			t.Error("couldn't upload bytes", err)
		}
	})

	// Also testing that we can put stuff in a "directory"
	imageFilename := "pixels/pixel.png"

	t.Run("test upload file as base 64", func(t *testing.T) {
		if err := UploadObjectAsBase64(bucketName, imageFilename, pixelBase64Data); err != nil {
			t.Error("error uploading to S3", err)
		}
	})

	t.Run("test delete files", func(t *testing.T) {
		if err := DeleteObject(bucketName, textFilename); err != nil {
			t.Error("error deleting text file from S3", err)
			return
		}
		if err := DeleteObject(bucketName, imageFilename); err != nil {
			t.Error("error deleting image from S3", err)
			return
		}
	})

	t.Run("test delete bucket", func(t *testing.T) {
		if err := DeleteBucket(bucketName); err != nil {
			t.Errorf("error deleting bucket named %s: %s\n", bucketName, err)
		}
	})
}
