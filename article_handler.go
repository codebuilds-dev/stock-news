package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"time"
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

type article struct {
	Symbol    string
	CreatedAt time.Time
	Headline  string
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	numberStr := r.URL.Query().Get("number")
	id := r.URL.Query().Get("id")

	// Default values for size and number if not provided
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 5 // default size
	}
	number, _ := strconv.Atoi(numberStr)
	if number == 0 {
		number = 1 // default number
	}

	file, err := os.ReadFile("articles.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var allArticles []article
	if err := json.Unmarshal(file, &allArticles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter articles by ID if provided
	var filteredArticles []article
	for _, art := range allArticles {
		if id == "" || art.Symbol == id {
			filteredArticles = append(filteredArticles, art)
		}
	}

	// Pagination logic (simplified)
	start := (number - 1) * size
	end := start + size
	if end > len(filteredArticles) {
		end = len(filteredArticles)
	}
	if start >= len(filteredArticles) {
		start = len(filteredArticles)
	}

	selectedArticles := filteredArticles[start:end]

	var articleRes articleResponse

	for _, art := range selectedArticles {
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

	var art article
	err := json.NewDecoder(r.Body).Decode(&art)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var articles []article
	file, err := os.ReadFile("articles.json")
	if err == nil {
		json.Unmarshal(file, &articles)
	}

	articles = append(articles, art)
	data, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	os.WriteFile("articles.json", data, fs.FileMode(0644))
	fmt.Fprintf(w, "Article saved with ID: %s", art.Symbol)
}
