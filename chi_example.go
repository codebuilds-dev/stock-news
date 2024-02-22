package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/codebuilds-dev/stock-news/auth"
	"github.com/codebuilds-dev/stock-news/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func apiKeyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-rapidapi-key")

		if !auth.ValidateAPIKey(apiKey) {
			http.Error(w, "Invalid API Key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getChiArticles(w http.ResponseWriter, r *http.Request) {
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

func saveChiArticle(w http.ResponseWriter, r *http.Request) {
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

func main_chi() {
	r := chi.NewRouter()

	//r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(apiKeyValidator)

	r.Get("/api/v1/articles", getChiArticles)
	r.Post("/api/v1/articles", saveChiArticle)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main_chi2() {
	// http.HandleFunc("/news/v2/list-by-symbol", getArticles)
	// http.HandleFunc("/news/v2/save-article", saveArticle)

	// fmt.Println("Server listening on port 8080")
	// err := http.ListenAndServe(":8080", nil)
	// log.Fatal(err)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(apiKeyValidator)

	r.Get("/api/v1/articles", getArticles)
	r.Post("/api/v1/articles", saveArticle)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}
