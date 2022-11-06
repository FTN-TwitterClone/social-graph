package service

import (
	"social-graph/model"
	"social-graph/repository"
)

type SocialGraphService struct {
	repo repository.SocialGraphRepository
}

func NewSocialGraphService(repo repository.SocialGraphRepository) *SocialGraphService {
	return &SocialGraphService{repo}

}

func (s SocialGraphService) CreateFollow(follow *model.Follows) error {

	err := s.repo.SaveFollow(follow)
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
