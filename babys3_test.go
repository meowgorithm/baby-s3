package babys3

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
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

func TestBucketOperations(t *testing.T) {
	n := "baby-s3-test-" + RandomBucketName()
	t.Run("create bucket", testCreateBucket(n))
	t.Run("bucket exists", testBucketExists(n))
	t.Run("make bucket public", func(t *testing.T) {
		if err := MakeBucketPublic(n); err != nil {
			t.Errorf("tried to make bucket '%s' public, but could not:", err)
		}
	})
	t.Run("delete bucket", testDeleteBucket(n))
}

func TestObjectOperations(t *testing.T) {
	bucketName := "baby-s3-test-" + RandomBucketName()
	textFilename := "meow.txt"
	textFileBytes := []byte("here, kitty kitty")

	t.Run("create bucket", testCreateBucket(bucketName))

	t.Run("test upload file as bytes", func(t *testing.T) {
		err := UploadObject(bucketName, textFilename, textFileBytes)
		if err != nil {
			t.Error("couldn't upload object:", err)
		}
	})

	// Also testing that we can put stuff in a "directory"
	imageFilename := "pixel/pixel.png"

	t.Run("test upload file as base 64", func(t *testing.T) {
		if err := UploadObjectAsBase64(bucketName, imageFilename, pixelBase64Data); err != nil {
			t.Error("couldn't upload object (from base64)", err)
		}
	})

	t.Run("test object exists", func(t *testing.T) {
		if exists, err := ObjectExists(bucketName, textFilename); err != nil {
			t.Error("could not check if object exists", err)
		} else if !exists {
			t.Errorf("object '%s' doesn't exist in bucket '%s', but it was supposed to be there: %s\n",
				textFilename, bucketName, err)
		}
	})

	t.Run("test object does not exist", func(t *testing.T) {
		if exists, err := ObjectExists(bucketName, "pretend-file.txt"); err != nil {
			t.Error("could not check if object exists", err)
		} else if exists {
			t.Error("a nonexistant object was reported as existing when it should not have")
		}
	})

	t.Run("test head object", func(t *testing.T) {
		var (
			info *s3.HeadObjectOutput
			err  error
		)
		if info, err = HeadObject(bucketName, textFilename); err != nil {
			t.Error("tried to get information about object but we could not:", err)
		}

		bytesReturned := *info.ContentLength
		actualBytes := int64(len(textFileBytes))

		if actualBytes != bytesReturned {
			t.Errorf("curious. S3 says the file is %d bytes, when in fact it's %d. That's an error.", bytesReturned, actualBytes)
		}
	})

	t.Run("test delete files", func(t *testing.T) {
		if err := DeleteObject(bucketName, textFilename); err != nil {
			t.Errorf("error deleting text file '%s' from bucket '%s:' %s\n", textFilename, bucketName, err)
			return
		}
		if err := DeleteObject(bucketName, imageFilename); err != nil {
			t.Errorf("error deleting image '%s' from bucket '%s': %s\n", imageFilename, bucketName, err)
			return
		}
	})

	t.Run("test delete bucket", testDeleteBucket(bucketName))
}

// Test helpers

func testCreateBucket(name string) func(*testing.T) {
	return func(t *testing.T) {
		if err := CreateBucket(name); err != nil {
			t.Errorf("tried to create bucket '%s' but could not: %s", name, err)
		}
	}
}

func testBucketExists(name string) func(*testing.T) {
	return func(t *testing.T) {
		if exists, err := BucketExists(name); err != nil {
			t.Errorf("tried to see if bucket '%s' exists but could not: %s\n", name, err)
		} else if !exists {
			t.Errorf("bucket '%s' does not exist: %s\n", name, err)
		}
	}
}

func testDeleteBucket(name string) func(*testing.T) {
	return func(t *testing.T) {
		if err := DeleteBucket(name); err != nil {
			t.Errorf("tried to delete bucket '%s' but could not: %s", name, err)
		}
	}
}
