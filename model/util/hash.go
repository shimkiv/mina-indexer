package util

import (
	"crypto/sha1"
	"fmt"
)

// SHA1 returns a SHA1 digest of a given string in hex format
func SHA1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil))
}
