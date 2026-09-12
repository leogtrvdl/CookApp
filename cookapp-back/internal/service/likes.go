package service

import (
	"cookapp/internal/repository"
)

type LikesService struct {
	repo *repository.LikesRepository
}

func NewLikesService(repo *repository.LikesRepository) *LikesService {
	return &LikesService{repo: repo}
}

func (s *LikesService) AddLike(userID int, recipeID int) error {
	err := s.repo.AddLike(userID, recipeID)
	return err
}

func (s *LikesService) DeleteLike(userID int, recipeID int) error {
	err := s.repo.DeleteLike(userID, recipeID)
	return err
}

func (s *LikesService) GetLikeCount(recipeID int) (int, error) {
	likeCount, err := s.repo.GetLikeCount(recipeID)
	return likeCount, err
}

func (s *LikesService) IsLiked(userID int, recipeID int) (bool, error) {
	likeCount, err := s.repo.IsLiked(userID, recipeID)
	return likeCount, err
}