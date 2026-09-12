package service

import (
	"cookapp/internal/model"
	"cookapp/internal/repository"
)

type IngredientService struct {
	repo *repository.IngredientRepository
}

func NewIngredientService(repo *repository.IngredientRepository) *IngredientService {
	return &IngredientService{repo: repo}
}

func (s *IngredientService) GetAllIngredients() ([]*model.Ingredient, error) {
	ingredients, err := s.repo.GetAllIngredients()
	return ingredients, err
}
