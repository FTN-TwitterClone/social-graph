package repository

import (
	"social-graph/model"
)

type SocialGraphRepository interface {
	SaveFollow(follow *model.Follows) error
	RemoveFollow(follow *model.Follows) error
	GetFollowing(username string) ([]model.User, error)
	GetFollowers(username string) ([]model.User, error)
}
