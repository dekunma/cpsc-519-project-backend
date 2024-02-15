// controllers/books.go

package controllers

import (
	"errors"
	"net/http"

	"github.com/dekunma/cpsc-519-project-backend/models"
	"github.com/gin-gonic/gin"
)

// FindBooks             godoc
//
//	@Summary		Get books array
//	@Description	Responds with the list of all books as JSON.
//	@Tags			books
//	@Produce		json
//	@Success		200	{array}	models.Book
//	@Router			/books [get]
func FindBooks(c *gin.Context) {
	var books []models.Book
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// CreateBook             godoc
//
//	@Summary		Create a single book
//	@Description	Responds with the book created
//	@Tags			books
//	@Produce		json
//	@Success		200	{object}	models.Book
//	@Router			/books [post]
func CreateBook(c *gin.Context) {
	var request models.CreatBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create book
	book := models.Book{Title: request.Title, Author: request.Author}
	models.DB.Create(&book)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func getBookByIdFromDB(book *models.Book, id string) error {
	if err := models.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return errors.New("No book with id: " + id)
	}

	return nil
}

// GetBookById             godoc
//
//	@Summary		Get a single book by its id
//	@Description	Responds with the book
//	@Tags			books
//	@Produce		json
//	@Param			id	path		string	true	"search book by its id"
//	@Success		200	{object}	models.Book
//	@Router			/books/{id} [get]
func GetBookById(c *gin.Context) {
	var book models.Book

	id := c.Param("id")
	if err := getBookByIdFromDB(&book, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// UpdateBookById             godoc
//
//	@Summary		Update a single book by its id
//	@Description	Responds with the updated book
//	@Tags			books
//	@Produce		json
//	@Param			id	path		string	true	"search book by its id"
//	@Success		200	{object}	models.Book
//	@Router			/books/{id} [patch]
func UpdateBookById(c *gin.Context) {
	var book models.Book

	id := c.Param("id")

	if err := getBookByIdFromDB(&book, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request models.UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	models.DB.Model(&book).Updates(request)

	c.JSON(http.StatusOK, gin.H{"data": book})
}
