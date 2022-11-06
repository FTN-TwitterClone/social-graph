package repository

import (
	"social-graph/model"
)

type SocialGraphRepository interface {
	SaveFollow(follow *model.Follows) error
	RemoveFollow(follow *model.Follows) error
	Get(username string, query string) ([]model.User, error)
}
