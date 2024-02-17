package cache

import (
	"context"
	"math/rand"
	"time"
)

var CODE_EXPIRATION_MINUTES = 10
var ctx = context.Background()

func generateRandomDigitStringWithLength(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		code += string(rune(rand.Intn(10) + int('0')))
	}
	return code
}

func RedisSetVerificationCode(email string) string {
	code := generateRandomDigitStringWithLength(4)
	RDB.Set(ctx, email, code, time.Duration(CODE_EXPIRATION_MINUTES)*time.Minute)
	return code
}

func RedisCheckVerificationCode(email, code string) bool {
	val, err := RDB.Get(ctx, email).Result()
	return err == nil && val == code
}
