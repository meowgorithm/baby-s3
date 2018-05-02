package s3

import (
	"testing"
)

func TestRandomBucketName(t *testing.T) {
	s1 := RandomBucketName(32)
	s2 := RandomBucketName(32)

	if s1 == s2 {
		t.Errorf("generated idential strings when we expected random ones (%s and %s)", s1, s2)
	}
}
