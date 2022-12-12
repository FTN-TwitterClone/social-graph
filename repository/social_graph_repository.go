package repository

import (
	"context"
	"social-graph/model"
)

type SocialGraphRepository interface {
	CreateNewUser(ctx context.Context, user model.User) error
	SaveApprovedFollow(ctx context.Context, fromUsername string, toUsername string) error
	RemoveApprovedFollow(ctx context.Context, fromUsername string, toUsername string) error
	RemoveFollowRequest(ctx context.Context, fromUsername string, toUsername string) error
	SaveFollowRequest(ctx context.Context, fromUsername string, toUsername string) error
	GetFollowing(ctx context.Context, username string) ([]model.User, error)
	GetFollowers(ctx context.Context, username string) ([]model.User, error)
	CheckIfFollowExists(ctx context.Context, from string, to string) (bool, error)
	AcceptRejectFollowRequest(ctx context.Context, from string, to string, approved bool) error
	GetUser(ctx context.Context, username string) (user *model.User, err error)
	CheckIfFollowRequestExists(ctx context.Context, from string, to string) (bool, error)
	GetAllFollowRequests(ctx context.Context, username string) ([]model.User, error)
	GetAllUsersNotFollowedByUser(ctx context.Context, username string) ([]model.User, error)
	GetRecommendationsProfile(ctx context.Context, username string) ([]model.User, error)
	CanAccessTweetOfAnotherUser(ctx context.Context, usernameFromToken string, usernameForAccess string) (bool, error)
	UpdateUser(ctx context.Context, isPrivate bool, authUsername string) error
	GetTargetGroupUser(ctx context.Context, username string, targetUserGroup model.TargetUserGroup) ([]model.User, error)
}
