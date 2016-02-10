package models

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

var DB LibStore

// LibStore intarace has Get and Put operations
type LibStore interface {
	Get(key string) (interface{}, error)
	Put(key string, val interface{}) error
}

// Format key for user login
func FormatUserLoginKey(userID string) string {
	return fmt.Sprintf("%s:usrid", userID)
}

// hash function to hash password
func Hash(msg string) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(msg))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

// Initialize
func init() {
	DB, _ = NewLibStore()
}
