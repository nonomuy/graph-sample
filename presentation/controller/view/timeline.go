package view

import (
	"eng-1004_nonomuy/domain/repository"
	"log"
	"net/http"
	"sort"
	"text/template"
	"time"
)

type TimeLineViewController interface {
	timeLineViewHandler(w http.ResponseWriter, r *http.Request)
}

type timeLineViewController struct {
	ur repository.UserRepository
	tr repository.PostRepository
}

func NewTimeLineViewController(ur repository.UserRepository, tr repository.PostRepository) TimeLineViewController {
	return &timeLineViewController{ur, tr}
}

type timeLineModel struct {
	Posts []postResponse
}

type postResponse struct {
	Id       string
	Message  string
	PostedAt string
}

func (u timeLineViewController) timeLineViewHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := query["userId"][0]

	posts := u.tr.GetPosts(userId)
	for _, user := range u.ur.GetFollowees(userId) {
		log.Println(user)
		posts = append(posts, u.tr.GetPosts(user.Id)...)
	}

	var postResponses []postResponse
	for _, p := range posts {
		log.Println(p)
		postResponses = append(postResponses, postResponse{Id: p.Id, Message: p.Message, PostedAt: p.PostedAt.Format(time.RFC822)})
	}
	sort.SliceStable(postResponses, func(i, j int) bool { return postResponses[i].PostedAt > postResponses[j].PostedAt })

	model := timeLineModel{postResponses}
	tmpl, err := template.ParseFiles("web/template/time-line.html",
		"web/template/_header.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
}
