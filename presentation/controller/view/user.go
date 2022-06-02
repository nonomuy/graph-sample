package view

import (
	"eng-1004_nonomuy/domain/repository"
	"log"
	"net/http"
	"text/template"
)

type UserViewController interface {
	chooseUserViewHandler(w http.ResponseWriter, r *http.Request)
	registerUserViewHandler(w http.ResponseWriter, r *http.Request)
	followUserViewHandler(w http.ResponseWriter, r *http.Request)
}

type userViewController struct {
	ur repository.UserRepository
}

func NewUserController(ur repository.UserRepository) UserViewController {
	return &userViewController{ur}
}

type updateUserResponse struct {
	Messages []string
	Users    []userResponse
}

type userResponse struct {
	Id   string
	Name string
}

type followUserModel struct {
	Messages []string
	Users    []userResponse
}

func (u userViewController) chooseUserViewHandler(w http.ResponseWriter, r *http.Request) {
	var userResponses []userResponse
	for _, e := range u.ur.GetUsers() {
		userResponses = append(userResponses, userResponse{e.Id, e.Name})
	}
	model := updateUserResponse{[]string{}, userResponses}
	tmpl, err := template.ParseFiles("web/template/update-user.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
}

func (u userViewController) registerUserViewHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		initializeRegisterUserView(w, updateUserResponse{})
		return
	case "POST":
		var userResponses []userResponse
		for _, e := range u.ur.GetUsers() {
			userResponses = append(userResponses, userResponse{e.Id, e.Name})
		}
		response := updateUserResponse{[]string{}, userResponses}
		err := r.ParseForm()
		if err != nil {
			response.Messages = append(response.Messages, "システムエラーです。")
		}
		name := r.Form.Get("name")
		if len(name) == 0 {
			response.Messages = append(response.Messages, "ユーザ名が入力されていません。")
		}
		if len(response.Messages) > 0 {
			initializeRegisterUserView(w, response)
			return
		}
		u.ur.PostUser(name)
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("location", "./choose-user")
		w.WriteHeader(http.StatusMovedPermanently)
		return
	default:
		w.WriteHeader(405)
		return
	}
}

func initializeRegisterUserView(w http.ResponseWriter, fieldErrors updateUserResponse) {
	tmpl, err := template.ParseFiles("web/template/choose-user.html",
		"web/template/_header.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, fieldErrors)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (u userViewController) followUserViewHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		u.initializeFollowUserViewHandler(w, []string{})
		return
	case "POST":
		messages := make([]string, 0)
		err := r.ParseForm()
		if err != nil {
			messages = append(messages, "システムエラーです。")
		}
		followerId := r.Form.Get("follower")
		if len(followerId) == 0 {
			messages = append(messages, "followerが入力されていません。")
		}
		followeeId := r.Form.Get("followee")
		if len(followeeId) == 0 {
			messages = append(messages, "followeeが入力されていません。")
		}
		if followeeId == followerId {
			messages = append(messages, "followerとfolloweeは異なるユーザを指定してください。")
		}

		if len(messages) > 0 {
			u.initializeFollowUserViewHandler(w, messages)
			return
		}
		u.ur.FollowUser(followerId, followeeId)
		u.initializeFollowUserViewHandler(w, []string{"フォロー登録に成功しました。"})
		return
	default:
		w.WriteHeader(405)
		return
	}
}

func (u userViewController) initializeFollowUserViewHandler(w http.ResponseWriter, messages []string) {
	var userResponses []userResponse
	for _, e := range u.ur.GetUsers() {
		userResponses = append(userResponses, userResponse{e.Id, e.Name})
	}
	tmpl, err := template.ParseFiles("web/template/follow-user.html",
		"web/template/_header.html")
	if err != nil {
		log.Fatal(err)
	}
	model := followUserModel{messages, userResponses}
	err = tmpl.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
	return
}
