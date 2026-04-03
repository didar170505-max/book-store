package main

import (
	"fmt"
	"net/http"

	"bookstore/handlers"
)

func main() {
	http.HandleFunc("/authors", handlers.AuthorsHandler)
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/books", handlers.BooksHandler)
	http.HandleFunc("/books/", handlers.BookByIDHandler)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
