// controllers/books.go

package controllers

import (
	"errors"
	"net/http"

	"github.com/dekunma/gin-books-test/models"
	"github.com/gin-gonic/gin"
)

// GET /books
// Get all books
func FindBooks(c *gin.Context) {
	var books []models.Book
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

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

func GetBookById(c *gin.Context) {
	var book models.Book

	id := c.Param("id")
	if err := getBookByIdFromDB(&book, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

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
