package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bookstore/models"
)

var books []models.Book
var nextBookID = 1

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w, r)
	case http.MethodPost:
		createBook(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func BookByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getBookByID(w, r, id)
	case http.MethodPut:
		updateBook(w, r, id)
	case http.MethodDelete:
		deleteBook(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(book.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if book.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	if !authorExists(book.AuthorID) {
		http.Error(w, "author does not exist", http.StatusBadRequest)
		return
	}

	if !categoryExists(book.CategoryID) {
		http.Error(w, "category does not exist", http.StatusBadRequest)
		return
	}

	book.ID = nextBookID
	nextBookID++
	books = append(books, book)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	filteredBooks := books

	categoryName := r.URL.Query().Get("category")
	if categoryName != "" {
		categoryID := findCategoryIDByName(categoryName)
		var result []models.Book
		for _, book := range books {
			if book.CategoryID == categoryID {
				result = append(result, book)
			}
		}
		filteredBooks = result
	}

	page := 1
	limit := 5

	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	start := (page - 1) * limit
	end := start + limit

	if start > len(filteredBooks) {
		filteredBooks = []models.Book{}
	} else {
		if end > len(filteredBooks) {
			end = len(filteredBooks)
		}
		filteredBooks = filteredBooks[start:end]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredBooks)
}

func getBookByID(w http.ResponseWriter, r *http.Request, id int) {
	for _, book := range books {
		if book.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "book not found", http.StatusNotFound)
}

func updateBook(w http.ResponseWriter, r *http.Request, id int) {
	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(updatedBook.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if updatedBook.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	if !authorExists(updatedBook.AuthorID) {
		http.Error(w, "author does not exist", http.StatusBadRequest)
		return
	}

	if !categoryExists(updatedBook.CategoryID) {
		http.Error(w, "category does not exist", http.StatusBadRequest)
		return
	}

	for i, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[i] = updatedBook

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	http.Error(w, "book not found", http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "book deleted successfully",
			})
			return
		}
	}
	http.Error(w, "book not found", http.StatusNotFound)
}

func authorExists(authorID int) bool {
	for _, author := range authors {
		if author.ID == authorID {
			return true
		}
	}
	return false
}

func categoryExists(categoryID int) bool {
	for _, category := range categories {
		if category.ID == categoryID {
			return true
		}
	}
	return false
}

func findCategoryIDByName(name string) int {
	for _, category := range categories {
		if strings.EqualFold(category.Name, name) {
			return category.ID
		}
	}
	return 0
}
