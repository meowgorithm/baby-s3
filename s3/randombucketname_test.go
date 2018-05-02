package s3

import (
	"testing"
)

func TestRandomBucketName(t *testing.T) {
	s1 := RandomBucketName()
	s2 := RandomBucketName()

	if s1 == s2 {
		t.Errorf("generated idential strings when we expected random ones (%s and %s)", s1, s2)
	}
}
