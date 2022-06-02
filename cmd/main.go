package main

import (
	"eng-1004_nonomuy/domain/repository"
	"eng-1004_nonomuy/infrastructure/mockRepository"
	"eng-1004_nonomuy/infrastructure/neptuneRepository"
	"eng-1004_nonomuy/presentation/controller/view"
	"os"

	"log"
	"net/http"
)

var ur repository.UserRepository
var tr repository.PostRepository

var uvc view.UserViewController
var tvc view.TweetViewController
var tlvc view.TimeLineViewController
var vro view.ViewRouter

func main() {

	dependencyInjection()

	vro.InitializeViewRouting()

	log.Fatal(http.ListenAndServe(":5000", nil))
}

func dependencyInjection() {
	env := os.Getenv("APP_ENV")
	if env == "MOCK" {
		ur = mockRepository.NewMockUserRepository()
		tr = mockRepository.NewMockTweetRepository()
	} else {
		ur = neptuneRepository.NewNeptuneRepository()
		tr = neptuneRepository.NewNeptuneRepository()
	}

	uvc = view.NewUserController(ur)
	tvc = view.NewTweetViewController(tr)
	tlvc = view.NewTimeLineViewController(ur, tr)
	vro = view.NewRouter(uvc, tvc, tlvc)
}
