package controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
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
