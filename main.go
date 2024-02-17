package main

import (
	"fmt"
	"github.com/dekunma/cpsc-519-project-backend/cache"
	"github.com/dekunma/cpsc-519-project-backend/controllers"
	_ "github.com/dekunma/cpsc-519-project-backend/docs"
	"github.com/dekunma/cpsc-519-project-backend/middleware"
	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

//	@title		API for CPSC 519 Project Group 6
//	@version	1.0
//
// @BasePath	/v1
func main() {
	loadEnvFile()
	models.ConnectDatabase()
	cache.ConnectRedis()
	middleware.SetupMiddleware()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Hello World!"})
	})

	v1 := r.Group("/v1")

	// Swagger ui
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// As an example from https://blog.logrocket.com/rest-api-golang-gin-gorm/
	v1.GET("/books", controllers.FindBooks)
	v1.POST("/books", controllers.CreateBook)
	v1.GET("books/:id", controllers.GetBookById)
	v1.PATCH("books/:id", controllers.UpdateBookById)

	// users
	v1.POST("/users/send-verification-code", controllers.SendVerificationCode)
	v1.POST("/users/sign-up", controllers.SignUp)
	v1.POST("/users/log-in", middleware.Auth.LoginHandler)

	PORT := os.Getenv("PORT")
	_ = r.Run(":" + PORT)
}
