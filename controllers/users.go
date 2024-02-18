package controllers

import (
	"github.com/dekunma/cpsc-519-project-backend/cache"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/dekunma/cpsc-519-project-backend/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func checkUserWithEmailExists(user *models.User, email string) bool {
	models.DB.Where("email = ?", email).Find(&user)
	return user.ID != 0
}

// SendVerificationCode godoc
//
//	@Summary	Send verification code to email
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/send-verification-code [post]
func SendVerificationCode(c *gin.Context) {
	var user models.User
	var request models.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	email := request.Email
	if checkUserWithEmailExists(&user, email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	code := cache.RedisSetVerificationCode(email)

	if err := service.SendVerificationCode(code, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

// CheckVerificationCode godoc
//
//	@Summary	Check verification code
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/check-verification-code [post]
func CheckVerificationCode(c *gin.Context) {
	var request models.CheckVerificationCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	email := request.Email
	code := request.VerificationCode

	if cache.RedisCheckVerificationCode(email, code) {
		c.JSON(http.StatusOK, gin.H{"message": "Verification code is valid"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
	}
}

// SignUp godoc
//
//	@Summary	Sign up
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/sign-up [post]
func SignUp(c *gin.Context) {
	var request models.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	email := request.Email
	code := request.VerificationCode

	if !cache.RedisCheckVerificationCode(email, code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	user := models.User{Email: email, Password: string(password)}

	// check again to prevent from possible attacks
	if checkUserWithEmailExists(&user, email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	models.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
