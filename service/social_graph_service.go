package service

import (
	"context"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/tweet"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"social-graph/model"
	"social-graph/repository"
	"social-graph/tls"
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

		conn, err := getgRPCConnection("tweet:9001")
		defer conn.Close()
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		tweetService := tweet.NewTweetServiceClient(conn)
		serviceCtx = metadata.AppendToOutgoingContext(serviceCtx, "authUsername", fromUsername)
		u := tweet.User{
			Username: toUsername,
		}

		_, error := tweetService.UpdateFeed(serviceCtx, &u)
		if error != nil {
			span.SetStatus(codes.Error, err.Error())
			return error
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
	serviceCtx, span := s.tracer.Start(ctx, "SocialGraphService.AcceptRejectFollowRequest")
	defer span.End()
	err := s.repo.AcceptRejectFollowRequest(serviceCtx, from, to, accepted)
	if accepted {
		conn, err := getgRPCConnection("tweet:9001")
		defer conn.Close()
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		tweetService := tweet.NewTweetServiceClient(conn)
		serviceCtx = metadata.AppendToOutgoingContext(serviceCtx, "authUsername", to)
		user := tweet.User{
			Username: from,
		}

		_, error := tweetService.UpdateFeed(serviceCtx, &user)
		if error != nil {
			span.SetStatus(codes.Error, err.Error())
			return error
		}

	}
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
func (s SocialGraphService) GetRecommendationsProfile(ctx context.Context, username string) ([]model.User, error) {
	ctx, span := s.tracer.Start(ctx, "SocialGraphService.GetRecommendationsProfile")
	defer span.End()

	followersCount, err := s.repo.GetFollowers(ctx, username)
	if len(followersCount) == 0 {
		users, err := s.repo.GetAllUsers(ctx, username)
		if err != nil {
			return nil, err
		}
		return users, nil

	}
	users, err := s.repo.GetRecommendationsProfile(ctx, username)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func getgRPCConnection(address string) (*grpc.ClientConn, error) {
	creds := credentials.NewTLS(tls.GetgRPCClientTLSConfig())

	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)

	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}

	return conn, err
}
