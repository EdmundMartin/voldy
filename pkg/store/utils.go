package store

import "errors"

var (
	ErrInvalidKey                = errors.New("invalid lookup key")
	ErrObsoleteVersion           = errors.New("obsolete version for provided key")
	ErrUnsupportedForStorageType = errors.New("operation is unsupported on storage type")
)

func IsValidKey(key []byte) bool {
	if len(key) == 0 {
		return false
	}
	return true
}

func BytesToString(key []byte) string {
	return string(key)
}

func StringToBytes(key string) []byte {
	return []byte(key)
}
