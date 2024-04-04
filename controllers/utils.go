package controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func extractEmailFromJWT(c *gin.Context) string {
	return jwt.ExtractClaims(c)["email"].(string)
}

func getFriendIDsOfUserById(userID uint) []uint {
	var friendships []*models.Friendship
	models.DB.Where("(user_id = ? or friend_id = ? ) AND accepted = ?", userID, userID, true).Find(&friendships)

	var friendIds []uint
	for _, f := range friendships {
		// append whichever side that is not the id of the user
		// (i.e. append the id of their friend)
		if f.UserID != userID {
			friendIds = append(friendIds, f.UserID)
		} else {
			friendIds = append(friendIds, f.FriendID)
		}
	}

	return friendIds
}
