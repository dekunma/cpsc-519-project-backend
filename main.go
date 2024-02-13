package main

import (
	"fmt"
	"github.com/dekunma/cpsc-519-project-backend/controllers"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func loadEnvFile() {
	env_mode := os.Getenv("GIN_ENV_MODE")

	if env_mode == "" {
		env_mode = "local" // Set a default if not specified
	}

	// Load the appropriate .env file
	err := godotenv.Load(fmt.Sprintf(".env.%s", env_mode))
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	loadEnvFile()

	r := gin.Default()

	models.ConnectDatabase()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Hello World!"})
	})

	// As an example from https://blog.logrocket.com/rest-api-golang-gin-gorm/
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	r.GET("books/:id", controllers.GetBookById)
	r.PATCH("books/:id", controllers.UpdateBookById)

	PORT := os.Getenv("PORT")
	_ = r.Run(":" + PORT)
}
