package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/news/v2/list-by-symbol", getArticles)
	http.HandleFunc("/news/v2/save-article", saveArticle)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
