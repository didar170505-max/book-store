package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"bookstore/models"
)

var categories []models.Category
var nextCategoryID = 1

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)

	case http.MethodPost:
		var category models.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(category.Name) == "" {
			http.Error(w, "category name is required", http.StatusBadRequest)
			return
		}

		category.ID = nextCategoryID
		nextCategoryID++
		categories = append(categories, category)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(category)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
