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

func (s SocialGraphService) CreateFollow(ctx context.Context, follow *model.Follows) error {
	serviceCtx, span := s.tracer.Start(ctx, "SocialGraphService.CreateFollow")
	defer span.End()

	err := s.repo.SaveFollow(serviceCtx, follow)
	if err != nil {
		return err
	}

	return nil
}
func (s SocialGraphService) RemoveFollow(follow *model.Follows) error {

	err := s.repo.RemoveFollow(follow)
	if err != nil {
		return err
	}

	return nil
}
func (s SocialGraphService) GetFollowing(username string) ([]model.User, error) {

	users, err := s.repo.GetFollowing(username)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s SocialGraphService) GetFollowers(username string) ([]model.User, error) {

	users, err := s.repo.GetFollowers(username)
	if err != nil {
		return nil, err
	}

	return users, nil
}
