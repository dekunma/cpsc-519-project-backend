package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dekunma/cpsc-519-project-backend/cache"
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/dekunma/cpsc-519-project-backend/service"
	"github.com/dekunma/cpsc-519-project-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
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
	email := extractEmailFromJWT(c)

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
	email := extractEmailFromJWT(c)
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

// UploadAvatar godoc
//
//	@Summary	Upload avatar
//	@Tags		users
//	@Accept		multipart/form-data
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/upload-avatar [post]
func UploadAvatar(c *gin.Context) {
	awsSession := c.MustGet("awsSession").(*session.Session)
	uploader := s3manager.NewUploader(awsSession)
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	file, header, err := c.Request.FormFile("avatar")

	filename := utils.GenerateRandomStringWithLength(10) + header.Filename

	//upload to the s3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		//ACL:    aws.String("public-read"),
		Key:  aws.String(filename),
		Body: file,
	})

	if err != nil {
		fmt.Println("AWS S3 upload error:")
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, exceptions.CustomError{
			Code:    exceptions.CodeUploadFailed,
			Message: "Failed to upload file",
		})

		return
	}
	filepath := "https://" + bucketName + "." + "s3" + ".amazonaws.com/" + filename

	email := extractEmailFromJWT(c)
	models.DB.Model(&models.User{}).Where("email = ?", email).Update("avatar", filepath)

	c.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
	})
}

// ChangePassword godoc
//
//	@Summary	Change password
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/change-password [patch]
func ChangePassword(c *gin.Context) {
	var request models.ChangePasswordRequest
	if !bindRequestToJSON(&request, c) {
		return
	}

	email := extractEmailFromJWT(c)
	var user models.User
	models.DB.Where("email = ?", email).Find(&user)

	password, _ := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	models.DB.Model(&models.User{}).Where("email = ?", email).Update("password", string(password))

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}
