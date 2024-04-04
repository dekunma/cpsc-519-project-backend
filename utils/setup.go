package utils

import (
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

func GenerateRandomDigitStringWithLength(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		code += string(rune(rand.Intn(10) + int('0')))
	}
	return code
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenerateRandomStringWithLength(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func BindRequestToJSON(request any, c *gin.Context) bool {
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeParamInvalid,
			Message: err.Error(),
		})
		return false
	}
	return true
}
