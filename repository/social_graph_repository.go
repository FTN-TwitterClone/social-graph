package repository

import (
	"social-graph/model"
)

type SocialGraphRepository interface {
	SaveFollow(follow *model.Follows, isPrivate bool) error
	RemoveFollow(follow *model.Follows) error
	GetFollowing(username string) ([]model.User, error)
	GetFollowers(username string) ([]model.User, error)
	CheckIfFollowExists(from string, to string, includeAll bool) (bool, error)
	AcceptRejectFollowRequest(from string, to string, approved bool) error
}
