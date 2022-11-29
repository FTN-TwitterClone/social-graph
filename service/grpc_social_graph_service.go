package service

import (
	"context"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/social_graph"
	"github.com/golang/protobuf/ptypes/empty"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
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
	serviceCtx, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterUser")
	defer span.End()
	err := s.repo.CreateNewUser(serviceCtx, user.Username, true)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (s gRPCSocialGraphService) RegisterBusinessUser(ctx context.Context, user *social_graph.SocialGraphBusinessUser) (*empty.Empty, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterBusinessUser")
	defer span.End()
	err := s.repo.CreateNewUser(serviceCtx, user.Username, false)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (s gRPCSocialGraphService) CheckVisibility(ctx context.Context, gRPCUsername *social_graph.SocialGraphUsername) (*social_graph.SocialGraphVisibilityUserResponse, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCSocialGraphService.CheckVisibility")
	defer span.End()

	md, _ := metadata.FromIncomingContext(ctx)
	authUsername := md.Get("authUsername")[0]

	visible, _ := s.repo.CanAccessTweetOfAnotherUser(serviceCtx, authUsername, gRPCUsername.Username)

	return &social_graph.SocialGraphVisibilityUserResponse{Visibility: visible}, nil
}

func (s gRPCSocialGraphService) GetMyFollowers(ctx context.Context, empty *emptypb.Empty) (*social_graph.SocialGraphFollowers, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCSocialGraphService.GetMyFollowers")
	defer span.End()

	md, _ := metadata.FromIncomingContext(ctx)
	authUsername := md.Get("authUsername")[0]

	users, _ := s.repo.GetFollowers(serviceCtx, authUsername)
	usersUsername := []*social_graph.SocialGraphUsername{}
	for _, user := range users {
		usersUsername = append(usersUsername, &social_graph.SocialGraphUsername{Username: user.Username})
	}
	return &social_graph.SocialGraphFollowers{Usernames: usersUsername}, nil
}

func (s gRPCSocialGraphService) SocialGraphUpdateUser(ctx context.Context, user *social_graph.SocialGraphUpdatedUser) (*empty.Empty, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCSocialGraphService.SocialGraphUpdatedUser")
	defer span.End()

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return new(empty.Empty), nil
	}

	authUsername := md.Get("authUsername")[0]

	err := s.repo.UpdateUser(serviceCtx, user.Private, authUsername)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}
