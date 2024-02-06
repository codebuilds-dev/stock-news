package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/codebuilds-dev/stock-news/auth"
	"github.com/codebuilds-dev/stock-news/db"
)

type attributes struct {
	PublishOn time.Time `json:"publishOn"`
	Title     string    `json:"title"`
}

type dataItem struct {
	Attributes attributes `json:"attributes"`
}

type articleResponse struct {
	Data []dataItem `json:"data"`
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("x-rapidapi-key")

	if !auth.ValidateAPIKey(apiKey) {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	sizeStr := r.URL.Query().Get("size")

	// Default values for size and number if not provided
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 10 // default size
	}

	articles, err := db.GetArticles(id, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var articleRes articleResponse

	for _, art := range articles {
		articleRes.Data = append(articleRes.Data, dataItem{
			Attributes: attributes{
				PublishOn: art.CreatedAt,
				Title:     art.Headline,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articleRes)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("x-rapidapi-key")

	if !auth.ValidateAPIKey(apiKey) {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var art db.Article
	err := json.NewDecoder(r.Body).Decode(&art)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.SaveArticle(art)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(art)
}
