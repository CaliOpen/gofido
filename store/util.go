package store

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64(s string) ([]byte, error) {
	for i := 0; i < len(s)%4; i++ {
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func EncodeBase64(buf []byte) string {
	s := base64.URLEncoding.EncodeToString(buf)
	return strings.TrimRight(s, "=")
}
