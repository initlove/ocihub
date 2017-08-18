package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

var MaxSize int64

func init() {
	MaxSize = 1024 * 1024 * 1024
}

// TODO: lots of todo, need to verify the digest
func Snap(digestFull string) (string, string) {
	var digest string
	strs := strings.Split(digestFull, ":")
	if len(strs) < 2 {
		digest = strs[0]
	} else {
		digest = strs[1]
	}

	if len(digest) < 2 {
		panic("Invalid digest")
	}

	return digest[:2], digest
}

func GetDigest(alg string, data []byte) string {
	if alg == "sha256" {
		sum := sha256.Sum256(data)
		return fmt.Sprintf("%x", sum)
	}

	return ""
}
