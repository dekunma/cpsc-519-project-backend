package middleware

import jwt "github.com/appleboy/gin-jwt/v2"

var Auth *jwt.GinJWTMiddleware

func SetupMiddleware() {
	var err error
	if Auth, err = GetAuthMiddleware(); err != nil {
		panic(err)
	}
}
