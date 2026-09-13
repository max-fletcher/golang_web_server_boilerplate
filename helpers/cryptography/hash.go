package cryptography

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(str string) string {
	hash := sha256.Sum256([]byte(str)) // convert to SHA256 checksum(32 bytes of binary data). Is of type [32]byte
	return hex.EncodeToString(hash[:]) // converts raw bytes into a hexadecimal string. return hexadecimal encoding of SHA256
}

func CompareHash(str1 string, str2 string) bool { // returns true if hash(str1) == str2, else false
	hash := sha256.Sum256([]byte(str1))     // convert to SHA256 checksum
	encoding := hex.EncodeToString(hash[:]) // return hexadecimal encoding of SHA256

	return encoding == str2
}
