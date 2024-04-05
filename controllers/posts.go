package controllers

import (
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
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

	c.JSON(http.StatusOK, gin.H{"message": "post created"})
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
