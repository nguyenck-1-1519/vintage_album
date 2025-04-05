package book_service

import (
	"errors"
	"net/http"
	"slices"
	"strconv"

	response "example.com/base_response"
	messages "example.com/messages"
	"github.com/gin-gonic/gin"
)

var Books = []Book{
	{ID: 1, Title: "A", Author: "X", Price: 1},
	{ID: 2, Title: "B", Author: "X", Price: 20},
	{ID: 3, Title: "B", Author: "X", Price: 300},
	{ID: 4, Title: "B", Author: "X", Price: 4000},
	{ID: 5, Title: "B", Author: "X", Price: 50000},
	{ID: 6, Title: "B", Author: "Y", Price: 600000},
	{ID: 7, Title: "B", Author: "Y", Price: 7000000, Stock: 1},
}

/*
[GET] /books?page=<int>&limit=<int>
Handle request GET with pagination
- page int <optional>: number of page, start from 0
- limit int <optinoal>: number of items per page, greater than 0, default value is 10
*/
func GetBooksWithPagination(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 0
	limit := 10

	// Get & assign default value of page
	if pageStr != "" {
		pageConvert, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			throwBadRequestWithError(&err, c, messages.InvalidParameter)
			return
		}
		page = pageConvert
	}

	// Get & assign default value of limit
	if limitStr != "" {
		limitConvert, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			throwBadRequestWithError(&err, c, messages.InvalidParameter)
			return
		}
		limit = limitConvert
	}

	books, pageInfo, err := GetBooksFromDB(page, limit)
	if err != nil {
		throwBadRequestWithError(&err, c, "")
		return
	}

	startIndex := page * limit
	if startIndex >= len(books) {
		c.IndentedJSON(http.StatusOK, response.BaseResponse{
			Status:  response.StatusOK,
			Data:    []Book{},
			Message: messages.OK,
			Page:    pageInfo,
		})
		return
	}

	endIndex := min(startIndex+limit, len(books))
	resultBooks := books[startIndex:endIndex]

	c.IndentedJSON(http.StatusOK, response.BaseResponse{
		Status:  response.StatusOK,
		Data:    resultBooks,
		Message: messages.OK,
		Page:    pageInfo,
	})
}

/*
[GET] /books/:id
Handle request GET book with designate id
- id int: id of book
*/
func GetBookWithID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	// id == 0 default value when parse string to int
	if err != nil || id == 0 {
		throwBadRequestWithError(&err, c, "")
		return
	}

	book, err := GetBookInfoFromDB(id)
	if err != nil {
		throwBadRequestWithError(&err, c, messages.ResultNotFound)
	}
	c.IndentedJSON(http.StatusAccepted, response.BaseResponse{
		Status:  response.StatusOK,
		Message: messages.OK,
		Data:    book,
	})
}

/*
[POST] /books/
Handle request POST[append] a book to current Books
- data will be decompressed from BODY
*/
func ImportNewBook(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil || !checkBindingConditionForNewBook(newBook) {
		throwBadRequestWithError(&err, c, "")
		return
	}

	err := InsertBookToDB(newBook)
	if err != nil {
		throwBadRequestWithError(&err, c, "")
		return
	}
	c.IndentedJSON(http.StatusOK, response.BaseResponse{
		Status:  response.StatusOK,
		Message: messages.OK,
		Data:    newBook,
	})
}

func checkBindingConditionForNewBook(b Book) bool {
	// ID will be manage by BE side, client don't need to pass ID
	if b.ID != 0 {
		return false
	}
	if b.Title == "" || b.Author == "" {
		return false
	}
	if b.Price <= 0 || b.Stock < 0 {
		return false
	}
	return true
}

/*
[PUT] /books/:id
Handle request PUT[update] a book info
- data will be decompressed from BODY
*/
func UpdateBookInfoWithID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	// id == 0 default value when parse string to int
	if err != nil || id == 0 {
		throwBadRequestWithError(&err, c, "")
		return
	}

	var adjustingBook Book
	if err := c.BindJSON(&adjustingBook); err != nil || !checkBindingConditionForNewBook(adjustingBook) {
		throwBadRequestWithError(&err, c, "")
		return
	}

	err = UpdateBookInfoToDB(adjustingBook, id)

	if err != nil {
		throwBadRequestWithError(&err, c, "")
		return
	}

	c.IndentedJSON(http.StatusOK, response.BaseResponse{
		Status:  response.StatusOK,
		Message: messages.OK,
	})
}

func getBookWithID(id int) (*Book, int, error) {
	for i, book := range Books {
		if book.ID == id {
			return &Books[i], i, nil
		}
	}
	return nil, -1, errors.New("not found book")
}

/*
[DELETE] /books/:id
Handle request DELETE a book with ID
- data will be decompressed from BODY
*/

func DeleteABookWithID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	// id == 0 default value when parse string to int
	if err != nil || id == 0 {
		c.IndentedJSON(http.StatusNoContent, response.BaseResponse{
			Status:  response.StatusError,
			Message: messages.ResultNotFound,
			Error:   &err,
		})
		return
	}

	_, i, err := getBookWithID(id)
	if err != nil || i < 0 {
		c.IndentedJSON(http.StatusNotFound, response.BaseResponse{
			Status:  response.StatusError,
			Message: messages.ResultNotFound,
			Error:   &err,
		})
		return
	}

	Books = slices.Delete(Books, i, i+1)
	c.IndentedJSON(http.StatusOK, response.BaseResponse{
		Status:  response.StatusOK,
		Message: messages.OK,
	})
}

func throwBadRequestWithError(err *error, c *gin.Context, m string) {
	message := m
	if m == "" {
		message = messages.BadRequest
	}

	c.IndentedJSON(http.StatusBadRequest, response.BaseResponse{
		Status:  response.StatusError,
		Message: message,
		Error:   err,
	})
}
