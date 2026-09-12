package service

import (
	"cookapp/internal/repository"
)

type FavoriteService struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteService(repo *repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) AddFavorite(userID int, recipeID int) error {
	return s.repo.AddFavorite(userID, recipeID)
}

func (s *FavoriteService) DeleteFavorite(userID int, recipeID int) error {
	return s.repo.DeleteFavorite(userID, recipeID)
}

func (s *FavoriteService) GetFavoriteCount(recipeID int) (int, error) {
	return s.repo.GetFavoriteCount(recipeID)
}

func (s *FavoriteService) IsFavorited(userID int, recipeID int) (bool, error) {
	favoriteCount, err := s.repo.IsFavorited(userID, recipeID)
	return favoriteCount, err
}