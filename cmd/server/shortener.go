package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// URL-safe alphabet for shortened URLs
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	// Alphabet length
	alphabetLength = uint64(len(alphabet))
)

// Convert GUID to shortened string (reversible)
func encodeGUID(guid string) string {
	guid = strings.ReplaceAll(guid, "-", "")
	num, _ := strToUint64(guid, 16)
	return encode(num)
}

// Decode shortened string to GUID
func decodeGUID(shortURL string) string {
	num := decode(shortURL)
	return fmt.Sprintf("%032x", num)
}

// Encode integer value to alphabet-based string
func encode(num uint64) string {
	var encodedBuilder strings.Builder
	for num > 0 {
		digit := num % alphabetLength
		encodedBuilder.WriteByte(alphabet[digit])
		num /= alphabetLength
	}
	return encodedBuilder.String()
}

// Decode alphabet-encoded string to integer value
func decode(encoded string) uint64 {
	var num uint64
	for i := range encoded {
		num = num*alphabetLength + uint64(strings.IndexByte(alphabet, encoded[i]))
	}
	return num
}

// Convert string to uint64
func strToUint64(s string, base int) (uint64, error) {
	num, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}
