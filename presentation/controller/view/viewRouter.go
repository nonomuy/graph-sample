package view

import (
	"net/http"
)

type ViewRouter interface {
	InitializeViewRouting()
}

type viewRouter struct {
	uvc  UserViewController
	tvc  TweetViewController
	tlvc TimeLineViewController
}

func NewRouter(uvc UserViewController, tvc TweetViewController, tlvc TimeLineViewController) ViewRouter {
	return &viewRouter{uvc, tvc, tlvc}
}

func (ro *viewRouter) InitializeViewRouting() {
	// template
	http.HandleFunc("/tweet", ro.tweetViewHandler)

	http.HandleFunc("/follow-user", ro.followUserViewHandler)
	http.HandleFunc("/register-user", ro.registerUserViewHandler)
	http.HandleFunc("/choose-user", ro.chooseUserViewHandler)
	http.HandleFunc("/time-line", ro.timeLineViewHandler)
	http.HandleFunc("/", ro.chooseUserViewHandler)

}

func (ro *viewRouter) timeLineViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ro.tlvc.timeLineViewHandler(w, r)
	default:
		w.WriteHeader(405)
	}
}

func (ro *viewRouter) chooseUserViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ro.uvc.chooseUserViewHandler(w, r)
	default:
		w.WriteHeader(405)
	}
}

func (ro *viewRouter) registerUserViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "POST":
		ro.uvc.registerUserViewHandler(w, r)
	default:
		w.WriteHeader(405)
	}
}

func (ro *viewRouter) followUserViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "POST":
		ro.uvc.followUserViewHandler(w, r)
	default:
		w.WriteHeader(405)
	}
}

func (ro *viewRouter) tweetViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "POST":
		ro.tvc.tweetViewHandler(w, r)
	default:
		w.WriteHeader(405)
	}
}
