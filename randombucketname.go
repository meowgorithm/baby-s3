package babys3

// Adapted from:
// https://www.calhoun.io/creating-random-strings-in-go/

import (
	"math/rand"
	"time"
)

// Charset we'll use for our random string
const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()),
)

// Fetch random chars from our charset, returning a random string
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomBucketName creates a random string valid for AWS bucket names. Note
// that bucket names must be between 3 and 63 characters. Here, we're doing it
// 32 characters in length, because 32 is a cool number.
func RandomBucketName() string {
	// Don't worry about the overhead of the extra function call here; Go will
	// inline it during compilation.
	return stringWithCharset(32, charset)
}
