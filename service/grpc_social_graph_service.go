package service

import (
	"context"
	social_graph "github.com/FTN-TwitterClone/grpc-stubs/social-graph"
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

func (s gRPCSocialGraphService) RegisterUser(ctx context.Context, user *social_graph.User) (*empty.Empty, error) {
	//serviceCtx
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterUser")
	defer span.End()

	return new(empty.Empty), nil
}

func (s gRPCSocialGraphService) RegisterBusinessUser(ctx context.Context, user *social_graph.BusinessUser) (*empty.Empty, error) {
	//serviceCtx
	_, span := s.tracer.Start(ctx, "gRPCSocialGraphService.RegisterBusinessUser")
	defer span.End()

	return new(empty.Empty), nil
}
