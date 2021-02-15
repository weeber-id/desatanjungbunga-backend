package tools

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"strings"
)

// PasswordHashing generator
// 10 times SHA512 with salt
// for every iteration:
// 		prev_hash + identity + iteration(integer to string)
func PasswordHashing(password string, identity string) string {
	hashi := password

	for idx := 0; idx < 10; idx++ {

		input := hashi + identity + strconv.Itoa(idx)
		result, err := encodeSHA512(input)

		if err != nil {
			log.Fatalf("error when hashing password: %v", err)
		}

		hashi = result
	}

	return hashi
}

func encodeSHA512(input string) (string, error) {
	inputReader := strings.NewReader(input)

	hash := sha512.New()
	if _, err := io.Copy(hash, inputReader); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	return hex.EncodeToString(sum), nil
}
