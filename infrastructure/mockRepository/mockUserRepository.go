package mockRepository

import (
	"eng-1004_nonomuy/domain/model"
	"strconv"
)

var users = map[string]model.User{
	"1": {"1", "foo"},
	"2": {"2", "bar"},
}

var follows = make(map[string][]string)

type mockUserRepository struct {
}

func (m mockUserRepository) GetUsers() []model.User {
	var result []model.User
	for _, user := range users {
		result = append(result, user)
	}
	return result
}

func (m mockUserRepository) PostUser(name string) {
	id := strconv.Itoa(len(users) + 1)

	newUser := model.User{Id: id, Name: name}
	users[id] = newUser
}

func (m mockUserRepository) GetFollowees(followerId string) []model.User {
	var result []model.User
	var empty model.User
	followeeIds := follows[followerId]
	for _, id := range followeeIds {

		if users[id] == empty {
			continue
		}
		result = append(result, users[id])
	}
	return result
}

func (m mockUserRepository) FollowUser(followerId string, followeeId string) {
	followees := follows[followerId]
	for _, id := range followees {
		if id == followeeId {
			return
		}
	}
	// HACK: 排他制御してないがmockなので…。
	follows[followerId] = append(followees, followeeId)
}

func NewMockUserRepository() *mockUserRepository {
	return &mockUserRepository{}
}
