package model

import "time"

type Post struct {
	Id       string
	Message  string
	PostedAt time.Time
}
