package repository

import "eng-1004_nonomuy/domain/model"

type PostRepository interface {
	Tweet(tweetFrom string, message string)
	GetPosts(userId string) []model.Post
}
