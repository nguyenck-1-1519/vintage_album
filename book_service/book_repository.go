package book_service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	response "example.com/base_response"
	_ "github.com/go-sql-driver/mysql"
)

var (
	QueryGetBooksWithPagination = "SELECT id, title, author, price, stock FROM books LIMIT ? OFFSET ?"
	QueryGetTotalItemCount      = "SELECT COUNT(*) FROM books"
	QueryGetBookInfo            = "SELECT id, title, author, price, stock FROM books WHERE id = ?"
	QueryInsertBook             = "INSERT INTO books (title, author, price, stock) VALUES (?, ?, ?, ?)"
	QueryUpdateBookInfo         = "UPDATE books SET title = ?, author = ?, price = ?, stock = ? WHERE id = ?"
)

func GetBooksFromDB(page int, limit int) ([]Book, response.PaginationData, error) {
	// open connection to db
	db, err := sql.Open("mysql", "root:Aa@123456@tcp(localhost:3306)/book_management")
	if err != nil {
		log.Fatal("open db failed")
		return nil, response.PaginationData{}, errors.New("open db failed")
	}
	defer db.Close()

	// Check connection before query
	if err := db.Ping(); err != nil {
		log.Fatal("keep connection to db failed")
		return nil, response.PaginationData{}, errors.New("keep connection to db failed")
	}

	// Calculate offset
	offset := max(page*limit, 0)

	// Query
	rows, err := db.Query(QueryGetBooksWithPagination, limit, offset)
	if err != nil {
		log.Fatal("query items from db failed")
		return nil, response.PaginationData{}, errors.New("query items from db failed")
	}
	defer rows.Close()

	// access rows & write to books
	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Stock); err != nil {
			log.Fatal("scan items db failed")
			return nil, response.PaginationData{}, errors.New("scan items db failed")
		}
		books = append(books, book)
	}

	// Query get total item count
	var totalItems int
	err = db.QueryRow(QueryGetTotalItemCount).Scan(&totalItems)
	if err != nil {
		log.Fatal("query total count of db failed")
		return nil, response.PaginationData{}, errors.New("query total count of db failed")
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	pageInfo := response.PaginationData{
		TotalItems:  totalItems,
		CurrentPage: page,
		PageSize:    limit,
		TotalPages:  totalPages,
	}

	return books, pageInfo, nil
}

func GetBookInfoFromDB(id int) (Book, error) {
	// open connection to db
	db, err := sql.Open("mysql", "root:Aa@123456@tcp(localhost:3306)/book_management")
	if err != nil {
		return Book{}, errors.New("open db failed")
	}
	defer db.Close()

	// Check connection before query
	if err := db.Ping(); err != nil {
		return Book{}, errors.New("keep connection to db failed")
	}

	// Query get Book Info
	var book Book
	err = db.QueryRow(QueryGetBookInfo, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Stock)
	if err != nil {
		log.Fatal("query items from db failed ", err)
		return Book{}, errors.New("query items from db failed")
	}

	return book, nil
}

func InsertBookToDB(book Book) error {
	// open connection to db
	db, err := sql.Open("mysql", "root:Aa@123456@tcp(localhost:3306)/book_management")
	if err != nil {
		return errors.New("open db failed")
	}
	defer db.Close()

	// Check connection before query
	if err := db.Ping(); err != nil {
		return errors.New("keep connection to db failed")
	}

	// Query get Book Info
	result, err := db.Exec(QueryInsertBook, book.Title, book.Author, book.Price, book.Stock)
	if err != nil {
		errmessages := fmt.Sprintf("query items from db failed %v", err)
		log.Fatal(errmessages)
		return errors.New("query items from db failed")
	}
	if result != nil {
		return nil
	}
	return errors.New("unknown error")
}

func UpdateBookInfoToDB(book Book, id int) error {
	// open connection to db
	db, err := sql.Open("mysql", "root:Aa@123456@tcp(localhost:3306)/book_management")
	if err != nil {
		return errors.New("open db failed")
	}
	defer db.Close()

	// Check connection before query
	if err := db.Ping(); err != nil {
		return errors.New("keep connection to db failed")
	}

	//Query update book
	result, err := db.Exec(QueryUpdateBookInfo, book.Title, book.Author, book.Price, book.Stock, id)

	if err != nil || result == nil {
		log.Fatal("query items from db failed ", err)
		return errors.New("query items from db failed")
	}
	return nil
}
