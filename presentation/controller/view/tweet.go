package view

import (
	"eng-1004_nonomuy/domain/repository"
	"log"
	"net/http"
	"text/template"
)

type TweetViewController interface {
	tweetViewHandler(w http.ResponseWriter, r *http.Request)
}

type tweetViewController struct {
	tr repository.PostRepository
}

func NewTweetViewController(tr repository.PostRepository) TweetViewController {
	return &tweetViewController{tr}
}

type tweetModel struct {
	Messages []string
}

func (t tweetViewController) tweetViewHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		initializeTweetView(w, []string{})
		return
	case "POST":
		messages := make([]string, 0)
		err := r.ParseForm()
		if err != nil {
			messages = append(messages, "システムエラーです。")
		}
		tweetFrom := r.Form.Get("tweetFrom")
		if len(tweetFrom) == 0 {
			messages = append(messages, "システムエラーです。（ユーザID不明）")
		}

		message := r.Form.Get("message")
		if len(message) == 0 {
			messages = append(messages, "メッセージが入力されていません。")
		}

		if len(messages) > 0 {
			initializeTweetView(w, messages)
			return
		}
		t.tr.Tweet(tweetFrom, message)
		initializeTweetView(w, []string{"ツイートを投稿しました。"})
	default:
		w.WriteHeader(405)
		return
	}
}

func initializeTweetView(w http.ResponseWriter, messages []string) {
	tmpl, err := template.ParseFiles("web/template/tweet.html",
		"web/template/_header.html")
	if err != nil {
		log.Fatal(err)
	}

	model := tweetModel{messages}
	err = tmpl.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
	return
}
