package mockRepository

import (
	"eng-1004_nonomuy/domain/model"
	"fmt"
	"time"
)

var posts = make(map[string][]model.Post)

type mockTweetRepository struct {
}

func NewMockTweetRepository() *mockTweetRepository {
	return &mockTweetRepository{}
}

func (m mockTweetRepository) Tweet(tweetFrom string, message string) {
	newTweet := model.Post{
		Id: tweetFrom, Message: message, PostedAt: time.Now()}
	posts[tweetFrom] = append(posts[tweetFrom], newTweet)
}

func (m mockTweetRepository) GetPosts(userId string) []model.Post {
	fmt.Println("mockTweetRepository : head of read")

	if posts[userId] == nil {
		return []model.Post{}
	}
	return posts[userId]
}
