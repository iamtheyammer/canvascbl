package util

import (
	"crypto/rand"
	"fmt"
	"strings"
)

func GenerateGiftCardClaimCodes(count int) []string {
	var codes []string
	for i := 0; i < count; i++ {
		base := fmt.Sprintf(
			"%s-%s-%s",
			GenerateRandomString(2),
			GenerateRandomString(2),
			GenerateRandomString(2),
		)

		code := strings.ToUpper(base)

		codes = append(codes, code)
	}

	return codes
}

func GenerateRandomString(n uint64) string {
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
