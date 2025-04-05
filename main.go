package main

import (
	bookService "example.com/book_service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/books", bookService.GetBooksWithPagination)
	router.GET("/books/:id", bookService.GetBookWithID)
	router.POST("/books", bookService.ImportNewBook)
	router.PUT("/books/:id", bookService.UpdateBookInfoWithID)
	router.DELETE("/books/:id", bookService.DeleteABookWithID)

	router.Run("localhost:8080")
}
