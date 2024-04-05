package controllers

import (
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateFriendInvitation(c *gin.Context) {
	var request models.CreateFriendInvitationRequest
	if !bindRequestToJSON(&request, c) {
		return
	}

	user := models.User{}
	models.DB.Where("email = ?", request.UserEmail).Find(&user)
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeUserNotFound,
			Message: "User not found",
		})
		return
	}

	friend := models.User{}
	models.DB.Where("email = ?", request.FriendEmail).Find(&friend)
	if friend.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeUserNotFound,
			Message: "Friendship not found",
		})
		return
	}

	models.DB.Create(&models.Friendship{
		UserID:   user.ID,
		FriendID: friend.ID,
		Accepted: false,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Friend invitation sent"})
}

func GetAllFriends(c *gin.Context) {
	email := extractEmailFromJWT(c)

	var user models.User
	models.DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeUserNotFound,
			Message: "User not found",
		})
		return
	}

	friendIds := getFriendIDsOfUserById(user.ID)

	// return an empty array if the user has no friends
	if len(friendIds) == 0 {
		c.JSON(http.StatusOK, gin.H{"friends": []models.User{}})
		return
	}

	var friends []*models.User
	// equivalent to SELECT id, email FROM users WHERE id IN (friendIds)
	models.DB.Select("id, email, avatar, name").Find(&friends, friendIds)

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

func GetFriendshipByEmail(c *gin.Context) {
	email := c.Param("email")

	searchedUser := models.User{}
	models.DB.Where("email = ?", email).Find(&searchedUser)
	if searchedUser.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.CustomError{
			Code:    exceptions.CodeUserNotFound,
			Message: "No user found with the given email",
		})
		return
	}

	currentUserEmail := extractEmailFromJWT(c)
	currentUser := models.User{}
	models.DB.Where("email = ?", currentUserEmail).Find(&currentUser)

	var friendShipsOfCurrentUser []models.Friendship
	currentUserId := currentUser.ID

	models.DB.Where("(user_id = ? or friend_id = ?)", currentUserId, currentUserId).Find(&friendShipsOfCurrentUser)
	requestStatus := "unsent"
	searchedUserId := searchedUser.ID
	for _, friendship := range friendShipsOfCurrentUser {
		if friendship.UserID == searchedUserId || friendship.FriendID == searchedUserId {
			if friendship.Accepted {
				requestStatus = "accepted"
			} else {
				requestStatus = "sent"
			}
		}
	}

	if searchedUser == currentUser {
		requestStatus = "yourself"
	}

	searchedUser.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": searchedUser, "request_status": requestStatus})
}
