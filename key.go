package cosmos

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type Key []byte

func ParseKey(key string) (Key, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, ErrInvalidKey
	}

	return Key(bytes), nil
}

func (k Key) Sign(data []byte) string {
	h := hmac.New(sha256.New, k)

	if _, err := h.Write(data); err != nil {
		return ""
	}

	b := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}
