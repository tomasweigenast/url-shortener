package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
)

func RandomId() uint32 {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("unable to read random bytes for generating an id: %s", err)
	}
	return binary.BigEndian.Uint32(bytes)
}

func RandomString(length ...int) string {
	strlen := 32 // 64 characters
	if len(length) == 1 {
		strlen = length[0] / 2
	}

	bytes := make([]byte, strlen)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("unable to read random bytes for generating a string: %s", err)
	}

	return hex.EncodeToString(bytes)
}

func RequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("unable to read random bytes for generating a request id: %s", err)
	}

	return hex.EncodeToString(bytes)
}
