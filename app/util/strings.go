package util

import (
	"encoding/hex"

	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

// Keccak256 returns a keccak256(string)
func Keccak256(content string) string {
	hash := solsha3.SoliditySHA3(
		[]string{"string"},
		[]interface{}{
			content,
		},
	)
	return hex.EncodeToString(hash)
}
