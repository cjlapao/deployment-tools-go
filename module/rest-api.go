package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/deployment-tools-go/controllers"
	"github.com/cjlapao/deployment-tools-go/database"
	"github.com/cjlapao/deployment-tools-go/entities"

	"net/http"

	"github.com/cjlapao/deployment-tools-go/repositories"

	"github.com/gorilla/mux"
)

type article = entities.Article

var router mux.Router
var databaseContext database.MongoFactory
var articlesRepo repositories.Repository

func RestApiModuleProcessor() {
	fmt.Println("Testing after")

	logger.Info("this would be the info %v", versionSvc.String())

	logger.Notice("Starting Go Rest API v0.1")
	databaseContext = database.NewFactory()
	articlesRepo = repositories.NewRepo(&databaseContext)
	pushTestData()

	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Homepage!")
	fmt.Println("endpoint Hit: homepage")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/", homePage)
	_ = controllers.NewAPIController(router, articlesRepo)
	logger.Success("Finished Init")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func pushTestData() {
	articles := []article{}
	data, err := ioutil.ReadFile("demo.json")

	helper.CheckError(err)

	json.Unmarshal(data, &articles)
	articlesRepo.UpsertManyArticles(articles)
}
