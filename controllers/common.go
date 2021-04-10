package controllers

import (
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/security"
	"github.com/cjlapao/deployment-tools-go/repositories"

	"github.com/gorilla/mux"
)

var logger = log.Get()

// Controllers Controller structure
type Controllers struct {
	Router     *mux.Router
	Repository *repositories.Repository
}

// NewAPIController  Creates a new controller
func NewAPIController(router *mux.Router, repo repositories.Repository) Controllers {
	controller := Controllers{
		Router:     router,
		Repository: &repo,
	}

	controller.Router.Handle("/article", security.AuthenticateMiddleware(controller.GetAllArticles)).Methods("GET")
	controller.Router.Handle("/article", security.AuthenticateMiddleware(controller.PostArticle)).Methods("POST")
	controller.Router.Handle("/article/{id}", security.AuthenticateMiddleware(controller.GetArticle)).Methods("GET")
	controller.Router.Handle("/article/{id}", security.AuthenticateMiddleware(controller.PutArticle)).Methods("PUT")
	controller.Router.Handle("/article/{id}", security.AuthenticateMiddleware(controller.DeleteArticle)).Methods("DELETE")

	controller.Router.HandleFunc("/login", controller.Login).Methods("POST")
	controller.Router.HandleFunc("/validate", controller.Validate).Methods("GET")

	controller.Router.HandleFunc("/test", controller.Test).Methods("GET")

	return controller
}
