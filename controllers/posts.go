package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/dekunma/cpsc-519-project-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := extractEmailFromJWT(c)
	var user models.User
	models.DB.Where("email = ?", email).First(&user)

	post.UserID = user.ID
	models.DB.Create(&post)

	c.JSON(http.StatusOK, gin.H{"message": "post created", "post_id": post.ID})
}

func GetAllPostsForUserAndFriends(c *gin.Context) {
	var posts []models.Post
	var user models.User
	email := extractEmailFromJWT(c)
	models.DB.Where("email = ?", email).First(&user)

	friendIDs := getFriendIDsOfUserById(user.ID)
	friendIDs = append(friendIDs, user.ID) // Include the user's own ID to get their posts too

	models.DB.Where("user_id IN (?)", friendIDs).Find(&posts)

	c.JSON(http.StatusOK, posts)
}

func GetPostDetails(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UploadPostImage godoc
//
//	@Summary	Uploads an images
//	@Tags		posts
//	@Accept		multipart/form-data
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/users/upload-post-image [post]
func UploadPostImage(c *gin.Context) {
	postIDStr := c.PostForm("post_id")
	if postIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	// PostID validation
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post ID"})
		return
	}

	awsSession := c.MustGet("awsSession").(*session.Session)
	uploader := s3manager.NewUploader(awsSession)
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	file, header, err := c.Request.FormFile("image")
	if err != nil || header == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file not found or incorrect file format"})
		return
	}

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

	// update the PostImages table with a new entry
	postImage := models.PostImages{
		PostID:   uint(postID),
		FilePath: filepath,
	}
	result := models.DB.Create(&postImage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
	})
}

// GetPostImages godoc
//
//	@Summary	retrieves images for a post
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/posts/get-post-images [get]
func GetPostImages(c *gin.Context) {
	postIDStr := c.Query("post_id")
	if postIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post ID"})
		return
	}

	var postImages []models.PostImages
	result := models.DB.Where("post_id = ?", postID).Find(&postImages)
	if result.Error != nil {
		fmt.Println("Database error:", result.Error)
		c.AbortWithStatusJSON(http.StatusInternalServerError, exceptions.CustomError{
			Code:    exceptions.CodeDatabaseError,
			Message: "Failed to retrieve images",
		})
		return
	}

	if len(postImages) == 0 { // no images found
		c.JSON(http.StatusOK, gin.H{"images": postImages,
			"message": "No images found for this post"})
		return
	}

	// Normal response
	c.JSON(http.StatusOK, gin.H{"images": postImages})
}
