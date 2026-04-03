package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"bookstore/models"
)

var authors []models.Author
var nextAuthorID = 1

func AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authors)

	case http.MethodPost:
		var author models.Author
		if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(author.Name) == "" {
			http.Error(w, "author name is required", http.StatusBadRequest)
			return
		}

		author.ID = nextAuthorID
		nextAuthorID++
		authors = append(authors, author)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(author)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
