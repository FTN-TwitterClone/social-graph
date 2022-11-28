package service

import (
	"context"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/social_graph"
	"github.com/golang/protobuf/ptypes/empty"
	"go.opentelemetry.io/otel/trace"
	"social-graph/model"
	"social-graph/repository"
)

type gRPCSocialGraphService struct {
	social_graph.UnimplementedSocialGraphServiceServer
	tracer trace.Tracer
	repo   repository.SocialGraphRepository
}

func NewgRPCSocialGraphService(tracer trace.Tracer, repo repository.SocialGraphRepository) *gRPCSocialGraphService {
	return &gRPCSocialGraphService{
		tracer: tracer,
		repo:   repo,
	}
}

func (s gRPCSocialGraphService) RegisterUser(ctx context.Context, user *social_graph.SocialGraphUser) (*empty.Empty, error) {
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterUser")
	defer span.End()
	err := s.repo.CreateNewUser(ctx, user.Username, true)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (s gRPCSocialGraphService) RegisterBusinessUser(ctx context.Context, user *social_graph.SocialGraphBusinessUser) (*empty.Empty, error) {
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterBusinessUser")
	defer span.End()
	err := s.repo.CreateNewUser(ctx, user.Username, false)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (s gRPCSocialGraphService) CheckVisibility(ctx context.Context, username string) (bool, error) {
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.CheckVisibility")
	defer span.End()
	authUser := ctx.Value("authUser").(model.AuthUser)
	visible, _ := s.repo.CanAccessTweetOfAnotherUser(ctx, authUser.Username, username)
	return visible, nil
}

func (s gRPCSocialGraphService) GetMyFollowers(ctx context.Context) ([]string, error) {
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.GetMyFollowers")
	defer span.End()
	authUser := ctx.Value("authUser").(model.AuthUser)
	users, _ := s.repo.GetFollowers(ctx, authUser.Username)
	usersUsername := []string{}
	for _, user := range users {
		usersUsername = append(usersUsername, user.Username)
	}
	return usersUsername, nil
}

func (s gRPCSocialGraphService) SocialGraphUpdateUser(ctx context.Context, isPrivate bool) (*empty.Empty, error) {
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.SocialGraphUpdatedUser")
	defer span.End()

	err := s.repo.UpdateUser(ctx, isPrivate)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}
