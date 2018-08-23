package cmq_go

import (
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"crypto/hmac"
	"encoding/base64"
)

const (
	SIGN_ALGORITHM_SHA1   = "sha1"
)

func Sign(src, key, method string) string{
	var mac hash.Hash
	if method == SIGN_ALGORITHM_SHA1{
		mac = hmac.New(sha1.New, []byte(key))
	}else{
		mac = hmac.New(sha256.New, []byte(key))
	}

	mac.Write([]byte(src))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}