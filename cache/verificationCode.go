package cache

import (
	"context"
	"github.com/dekunma/cpsc-519-project-backend/utils"
	"time"
)

var CODE_EXPIRATION_MINUTES = 10
var ctx = context.Background()

func RedisSetVerificationCode(email string) string {
	code := utils.GenerateRandomDigitStringWithLength(4)
	RDB.Set(ctx, email, code, time.Duration(CODE_EXPIRATION_MINUTES)*time.Minute)
	return code
}

func RedisCheckVerificationCode(email, code string) bool {
	val, err := RDB.Get(ctx, email).Result()
	return err == nil && val == code
}
