package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cjlapao/deployment-tools-go/entities"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

// GetArticle Gets an article by it's id from the database
func (c *Controllers) GetArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Articles Get Endpoint Hit")
	vars := mux.Vars(r)
	key := vars["id"]

	article := c.Repository.GetArticleByID(key)
	json.NewEncoder(w).Encode(article)
}

// GetAllArticles Gets all the articles from the database
func (c *Controllers) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Articles GetAll Endpoint Hit")

	articles := c.Repository.GetAllArticles()

	json.NewEncoder(w).Encode(articles)
}

// PostArticle Post new article into database
func (c *Controllers) PostArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Articles PostArticle Endpoint Hit")

	reqBody, _ := ioutil.ReadAll(r.Body)
	article := entities.Article{}
	json.Unmarshal(reqBody, &article)
	article.ID = xid.New().String()

	result := c.Repository.UpsertArticle(article)
	fmt.Print(result)
	if result.UpsertedCount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// PutArticle Updates an article in the database
func (c *Controllers) PutArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Articles PutArticle Endpoint Hit")
	vars := mux.Vars(r)
	key := vars["id"]

	if len(key) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	article := entities.Article{}
	json.Unmarshal(reqBody, &article)
	article.ID = key

	result := c.Repository.UpdateArticle(article)

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// DeleteArticle Deletes an article from the database
func (c *Controllers) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Articles DeleteArticle Endpoint Hit")
	vars := mux.Vars(r)
	key := vars["id"]

	if len(key) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := c.Repository.DeleteArticle(key)

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
