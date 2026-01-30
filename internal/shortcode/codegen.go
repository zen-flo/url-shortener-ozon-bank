package shortcode

import (
	"crypto/rand"
	"math/big"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const codeLength = 10

func GenerateCode() (string, error) {
	code := make([]byte, codeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		code[i] = charSet[num.Int64()]
	}
	return string(code), nil
}
