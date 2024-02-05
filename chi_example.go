package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main_chi() {
	r := chi.NewRouter()

	r.Get("/news/v2/list-by-symbol", getArticles)
	r.Post("/save-article", saveArticle)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
