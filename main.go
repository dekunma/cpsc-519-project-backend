package main

import (
	"fmt"
	"github.com/dekunma/cpsc-519-project-backend/cache"
	"github.com/dekunma/cpsc-519-project-backend/controllers"
	_ "github.com/dekunma/cpsc-519-project-backend/docs"
	"github.com/dekunma/cpsc-519-project-backend/exceptions"
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
	envMode := os.Getenv("GIN_ENV_MODE")

	if envMode == "" {
		envMode = "local" // Set a default if not specified
	}

	// Load the appropriate .env file
	err := godotenv.Load(fmt.Sprintf(".env.%s", envMode))
	if err != nil {
		fmt.Println("Error loading .env file "+envMode+": ", err)
	}

	// Load the .env file
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
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

	r := gin.Default()
	r.Use(exceptions.ErrorHandler)

	awsSession := middleware.ConnectAws()
	r.Use(func(c *gin.Context) {
		c.Set("awsSession", awsSession)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Hello World!"})
	})

	v1 := r.Group("/v1")

	// Swagger ui
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// As an example from https://blog.logrocket.com/rest-api-golang-gin-gorm/
	books := v1.Group("/books")
	{
		books.GET("/", controllers.FindBooks)
		books.POST("/", controllers.CreateBook)
		books.GET("/:id", controllers.GetBookById)
		books.PATCH("/:id", controllers.UpdateBookById)
	}

	middleware.SetupMiddleware()
	// users
	users := v1.Group("/users")
	{
		users.POST("send-verification-code", controllers.SendVerificationCode)
		users.POST("check-verification-code", controllers.CheckVerificationCode)
		users.POST("sign-up", controllers.SignUp)
		users.POST("log-in", middleware.Auth.LoginHandler)
		users.GET("refresh-token", middleware.Auth.RefreshHandler)
	}

	// authenticated-only routes
	users.Use(middleware.Auth.MiddlewareFunc())
	{
		users.PATCH("update-profile", controllers.UpdateProfile)
		users.GET("profile", controllers.GetOwnProfile)
		users.POST("upload-avatar", controllers.UploadAvatar)
	}

	// friendships
	friendships := v1.Group("/friendships")
	friendships.Use(middleware.Auth.MiddlewareFunc())
	{
		friendships.POST("create-friend-invitation", controllers.CreateFriendInvitation)
		friendships.GET("get-all-friends", controllers.GetAllFriends)
	}

	PORT := os.Getenv("PORT")
	_ = r.Run(":" + PORT)
}
