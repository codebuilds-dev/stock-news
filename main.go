package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

	r.Get("/api/v1/articles", getArticles)
	r.Post("/api/v1/articles", saveArticle)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}
