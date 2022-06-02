package repository

import (
	"eng-1004_nonomuy/domain/model"
)

type UserRepository interface {
	GetUsers() []model.User
	PostUser(name string)
	FollowUser(followerId string, followeeId string)
	GetFollowees(followerId string) []model.User
}
