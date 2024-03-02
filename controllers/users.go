package controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dekunma/cpsc-519-project-backend/cache"
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/dekunma/cpsc-519-project-backend/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func abortIfUserWithEmailExists(email string, c *gin.Context) bool {
	user := &models.User{}
	models.DB.Where("email = ?", email).Find(&user)
	if user.ID != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeEmailAlreadyExists,
			Message: "User already exists",
		})
		return true
	}
	return false
}

func bindRequestToJSON(request any, c *gin.Context) bool {
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeParamInvalid,
			Message: err.Error(),
		})
		return false
	}
	return true
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
	var request models.SendVerificationCodeRequest
	if !bindRequestToJSON(&request, c) {
		return
	}

	email := request.Email
	if abortIfUserWithEmailExists(email, c) {
		return
	}

	code := cache.RedisSetVerificationCode(email)

	if err := service.SendVerificationCode(code, email); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, exceptions.CustomError{
			Code:    exceptions.CodeSendEmailFailed,
			Message: "Failed to send verification code",
		})
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
	if !bindRequestToJSON(&request, c) {
		return
	}

	email := request.Email
	code := request.VerificationCode

	if cache.RedisCheckVerificationCode(email, code) {
		c.JSON(http.StatusOK, gin.H{"message": "Verification code is valid"})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeVerificationCodeInvalid,
			Message: "Invalid verification code",
		})
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
	if !bindRequestToJSON(&request, c) {
		return
	}

	email := request.Email
	code := request.VerificationCode

	if !cache.RedisCheckVerificationCode(email, code) {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeVerificationCodeInvalid,
			Message: "Invalid verification code",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	// check again to prevent from possible attacks
	if abortIfUserWithEmailExists(email, c) {
		return
	}

	user := models.User{Email: email, Password: string(password)}
	models.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// UpdateProfile godoc
//
//	@Summary	Update a user's name
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/update-name [patch]
func UpdateProfile(c *gin.Context) {
	var request models.UpdateProfileRequest
	if !bindRequestToJSON(&request, c) {
		return
	}

	var user models.User
	email := jwt.ExtractClaims(c)["email"].(string)

	models.DB.Where("email = ?", email).Find(&user)
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeUserNotFound,
			Message: "User not found",
		})
		return
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Avatar != "" {
		user.Avatar = request.Avatar
	}

	models.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// GetOwnProfile godoc
//
//	@Summary	Get user's own profile
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Router		/users/profile [get]
func GetOwnProfile(c *gin.Context) {
	email := jwt.ExtractClaims(c)["email"].(string)
	var user models.User
	models.DB.Where("email = ?", email).Find(&user)
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeUserNotFound,
			Message: "User not found",
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
