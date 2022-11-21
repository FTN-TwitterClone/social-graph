package repository

import (
	"context"
	"social-graph/model"
)

type SocialGraphRepository interface {
	SaveFollow(ctx context.Context, follow *model.Follows) error
	RemoveFollow(follow *model.Follows) error
	GetFollowing(username string) ([]model.User, error)
	GetFollowers(username string) ([]model.User, error)
}
