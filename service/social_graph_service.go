package service

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"social-graph/model"
	"social-graph/repository"
)

type SocialGraphService struct {
	repo   repository.SocialGraphRepository
	tracer trace.Tracer
}

func NewSocialGraphService(repo repository.SocialGraphRepository, tracer trace.Tracer) *SocialGraphService {
	return &SocialGraphService{
		repo,
		tracer,
	}
}

func (s SocialGraphService) CreateFollow(ctx context.Context, fromUsername string, toUsername string) error {
	serviceCtx, span := s.tracer.Start(ctx, "SocialGraphService.CreateFollow")
	defer span.End()
	user, err2 := s.repo.GetUser(serviceCtx, toUsername)
	if err2 != nil {
		return err2
	}
	if user.IsPrivate {
		err := s.repo.SaveFollowRequest(serviceCtx, fromUsername, toUsername)
		if err != nil {
			return err
		}

	} else {
		err := s.repo.SaveApprovedFollow(serviceCtx, fromUsername, toUsername)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s SocialGraphService) RemoveFollow(ctx context.Context, fromUsername string, toUsername string) error {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.RemoveFollow")
	defer span.End()
	err := s.repo.RemoveApprovedFollow(ctx, fromUsername, toUsername)
	if err != nil {
		return err
	}

	return nil
}
func (s SocialGraphService) GetFollowing(ctx context.Context, username string) ([]model.User, error) {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.GetFollowing")
	defer span.End()
	users, err := s.repo.GetFollowing(ctx, username)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s SocialGraphService) GetFollowers(ctx context.Context, username string) ([]model.User, error) {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.GetFollowers")
	defer span.End()
	users, err := s.repo.GetFollowers(ctx, username)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s SocialGraphService) CheckIfFollowExists(ctx context.Context, from string, to string) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.CheckIfFollowExists")
	defer span.End()
	exists, err := s.repo.CheckIfFollowExists(ctx, from, to)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s SocialGraphService) CheckIfFollowRequestExists(ctx context.Context, from string, to string) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.CheckIfFollowRequestExists")
	defer span.End()
	exists, err := s.repo.CheckIfFollowRequestExists(ctx, from, to)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (s SocialGraphService) AcceptRejectFollowRequest(ctx context.Context, from string, to string, accepted bool) error {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.AcceptRejectFollowRequest")
	defer span.End()
	err := s.repo.AcceptRejectFollowRequest(ctx, from, to, accepted)
	if err != nil {
		return err
	}
	return nil
}

func (s SocialGraphService) GetAllFollowRequests(ctx context.Context, username string) ([]model.User, error) {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.GetAllFollowRequests")
	defer span.End()
	users, err := s.repo.GetAllFollowRequests(ctx, username)
	if err != nil {
		return nil, err
	}

	return users, nil
}
