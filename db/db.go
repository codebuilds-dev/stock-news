package db

import (
	"encoding/json"
	"os"
	"time"
)

type Article struct {
	Symbol    string
	CreatedAt time.Time
	Headline  string
}

func GetArticles(id string, size int) ([]Article, error) {

	file, err := os.ReadFile("db/articles.json")
	if err != nil {
		return nil, err
	}

	var allArticles []Article
	if err := json.Unmarshal(file, &allArticles); err != nil {
		return nil, err
	}

	// Filter articles by ID if provided
	var filteredArticles []Article
	for _, art := range allArticles {
		if id == "" || art.Symbol == id {
			filteredArticles = append(filteredArticles, art)
		}
	}

	artCount := len(filteredArticles)

	if artCount > 0 {
		if artCount < size {
			size = artCount
		}
		filteredArticles = filteredArticles[0:size]
	}

	return filteredArticles, nil
}

func SaveArticle(art Article) error {
	file, err := os.ReadFile("db/articles.json")
	if err != nil {
		return err
	}

	var articles []Article
	json.Unmarshal(file, &articles)

	articles = append(articles, art)
	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}

	os.WriteFile("db/articles.json", data, os.FileMode(0644))
	return nil
}
