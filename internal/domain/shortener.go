package domain

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// URL-safe alphabet for shortened URLs
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	// Alphabet length
	alphabetLength = uint64(len(alphabet))
)

// ADD TO Test
// var guid = "123e4567-e89b-12d3-a456-426655440000"
// 	var shortURL = encodeGUID(guid)
// 	fmt.Println("GUID:", guid)
// 	fmt.Println("Short URL:", shortURL)

// 	decodedGUID := decodeGUID(shortURL)
// 	fmt.Println("Decoded GUID:", decodedGUID)

// Convert GUID to shortened string (reversible)
func EncodeGUID(guid string) string {
	guid = strings.ReplaceAll(guid, "-", "")
	num1, _ := strToUint64(guid[:16], 16)
	num2, _ := strToUint64(guid[16:], 16)
	return encode(num1) + "-" + encode(num2)
}

// Decode shortened string to GUID
func DecodeGUID(shortURL string) (string, error) {
	var num1, num2 uint64

	urls := strings.Split(shortURL, "-")
	if len(urls) != 2 {
		return "", fmt.Errorf("DecodeGUID: wrong input data: %s", shortURL)
	}
	num1 = decode(urls[0])
	num2 = decode(urls[1])
	uuid := uint64ToUUID(num1, num2)
	return uuid, nil
}

// Encode integer value to alphabet-based string
func encode(num uint64) string {
	var encodedBuilder strings.Builder
	for num > 0 {
		digit := num % alphabetLength
		encodedBuilder.WriteByte(alphabet[digit])
		num /= alphabetLength
	}
	encoded := encodedBuilder.String()
	return reverse(encoded)
}

// Decode alphabet-encoded string to integer value
func decode(encoded string) uint64 {
	var num uint64
	for i := range encoded {
		num = num*alphabetLength + uint64(strings.IndexByte(alphabet, encoded[i]))
	}
	return num
}

// Reverse string
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Convert string to uint64
func strToUint64(s string, base int) (uint64, error) {
	num, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// func uint64ToUUID(uint1 uint64, uint2 uint64, setFlags bool) string {
// 	bytes := make([]byte, 16)

// 	// Fill first 8 bytes from uint1
// 	for i := 0; i < 8; i++ {
// 		bytes[7-i] = byte(uint1 >> (8 * i))
// 	}

// 	// Fill last 8 bytes from uint2
// 	for i := 8; i < 16; i++ {
// 		bytes[15-i] = byte(uint2 >> (8 * (i - 8)))
// 	}

// 	if setFlags {
// 		// Set UUID version (v4 - random)
// 		bytes[6] = (bytes[6] & 0x0f) | 0x40

// 		// Set variant flag (RFC 4122)
// 		bytes[8] = (bytes[8] & 0x3f) | 0x80
// 	}

// 	// Form UUID string
// 	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
// }

func uint64ToUUID(uint1 uint64, uint2 uint64) string {
	uuid := fmt.Sprintf("%16x%16x", uint1, uint2)
	uuidWithDelimiter := fmt.Sprintf("%s-%s-%s-%s-%s",
		uuid[0:8],
		uuid[8:12],
		uuid[12:16],
		uuid[16:20],
		uuid[20:],
	)
	return uuidWithDelimiter
}
