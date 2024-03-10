package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codebuilds-dev/stock-news/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

func main() {
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

	r.Get("/api/v1/articles", getChiArticles)
	r.Post("/api/v1/articles", saveChiArticle)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
