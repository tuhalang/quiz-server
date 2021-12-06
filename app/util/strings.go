package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

const (
	prefix        = "0x"
	addressLength = 40
	dataLength    = 64
	zeroChar      = 48

	// StatusDraft is a status when quiz is unconfirmed
	StatusDraft = -1
	// StatusDone is a status when quiz is confirmed
	StatusDone = 1
	// AddressZero is an address(0)
	AddressZero = "0x0000000000000000000000000000000000000000"
)

func HexTo32Bytes(id string) (*[32]byte, error) {
	var rs [32]byte
	src, err := hex.DecodeString(id[2:])
	if err != nil {
		return nil, err
	}
	copy(rs[:], src[:])
	return &rs, nil
}

func ByteToString(data [32]byte) string {
	var dst = make([]byte, 32)
	copy(dst[:], data[:])
	text := hex.EncodeToString(dst)
	return prefix + text
}

// Keccak256 returns a keccak256(string)
func Keccak256(content string) string {
	sum := sha256.Sum256([]byte(content))
	content = fmt.Sprintf("%x", sum)
	hash := solsha3.SoliditySHA3(
		[]string{"string"},
		[]interface{}{
			content,
		},
	)
	return prefix + hex.EncodeToString(hash)
}

// FormatAddress returns a valid address
func FormatAddress(address string) string {
	length := len(address)
	if length > addressLength {
		return prefix + address[length-addressLength:]
	}
	return address
}

// SplitData returns a array string hex data
func SplitData(data string) []string {
	length := len(data)
	var result []string
	for i := 0; i < length; i += dataLength {
		result = append(result, prefix+data[i:i+dataLength])
	}
	return result
}

// FormatHexNumber returns a valid hex number
func FormatHexNumber(data string, hasPrefix bool) string {
	length := len(data)
	i := 2
	for i < length {
		if data[i] != zeroChar {
			break
		}
		i++
	}
	if hasPrefix {
		return prefix + data[i:]
	}
	return data[i:]
}
