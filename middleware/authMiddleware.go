package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	identityKey := "email"
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "full stack project NAME TBD",
		Key:         []byte(os.Getenv("JWT_SECRET")),
		Timeout:     time.Hour * 24 * 90,
		MaxRefresh:  time.Hour * 24 * 90,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
					"name":      v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Email: claims[identityKey].(string),
				Name:  claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.SignUpRequest
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Email
			password := loginVals.Password
			var user models.User

			if err := models.DB.Where("email = ?", userID).First(&user).Error; err != nil && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &models.User{
				Email: userID,
				Name:  user.Name,
			}, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	return authMiddleware, err
}
