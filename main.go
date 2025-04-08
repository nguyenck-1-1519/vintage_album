package main

import (
	"fmt"

	bookService "example.com/book_service"
	auth "example.com/my_authentication"
	"github.com/gin-gonic/gin"
)

func main() {

	token, err := auth.GETJWTTokenString()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Auth for this section: ", token)

	router := gin.Default()
	authMiddleware := auth.AuthMiddleware()

	router.GET("/books", bookService.GetBooksWithPagination)
	router.GET("/books/:id", bookService.GetBookWithID)

	bookRouter := router.Group("/books")
	bookRouter.Use(authMiddleware) // Áp dụng middleware cho group này
	{
		bookRouter.POST("", bookService.ImportNewBook)
		bookRouter.PUT("/:id", bookService.UpdateBookInfoWithID)
		bookRouter.DELETE("/:id", bookService.DeleteABookWithID)
	}

	router.Run("localhost:8080")
}
