package utils

import "math/rand"

func GenerateRandomDigitStringWithLength(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		code += string(rune(rand.Intn(10) + int('0')))
	}
	return code
}
