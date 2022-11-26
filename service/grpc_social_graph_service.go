package service

import (
	"context"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/social_graph"
	"github.com/golang/protobuf/ptypes/empty"
	"go.opentelemetry.io/otel/trace"
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
	//serviceCtx
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterUser")
	defer span.End()
	err := s.repo.CreateNewUser(ctx, user.Username, true)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (s gRPCSocialGraphService) RegisterBusinessUser(ctx context.Context, user *social_graph.SocialGraphBusinessUser) (*empty.Empty, error) {
	//serviceCtx
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterBusinessUser")
	defer span.End()
	err := s.repo.CreateNewUser(ctx, user.Username, false)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}
