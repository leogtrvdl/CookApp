package service

import (
	"cookapp/internal/model"
	"cookapp/internal/repository"
	"errors"
)

type RecipeService struct {
	repo         *repository.RecipeRepository
	likesRepo    *repository.LikesRepository
	favoriteRepo *repository.FavoriteRepository
}

func NewRecipeService(repo *repository.RecipeRepository, likesRepo *repository.LikesRepository, favoriteRepo *repository.FavoriteRepository) *RecipeService {
	return &RecipeService{repo: repo, likesRepo: likesRepo, favoriteRepo: favoriteRepo}
}

func (s *RecipeService) CreateRecipe(userID int, title string, imagePath string, ingredients []model.Ingredient, steps []model.Step) error {

	recipeID, err := s.repo.CreateRecipe(userID, title, imagePath)
	if err != nil {
		return err
	}

	for _, ing := range ingredients {
		err = s.repo.AddIngredient(recipeID, ing.Name, ing.Quantity, ing.Unit)
		if err != nil {
			return err
		}
	}

	for _, step := range steps {
		err = s.repo.AddStep(recipeID, step.Order, step.Description)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *RecipeService) GetRecipeByID(recipeID int) (*model.Recipe, error) {
	recipe, err := s.repo.GetRecipeByID(recipeID)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients, err = s.repo.GetIngredientsByRecipeID(recipeID)
	if err != nil {
		return nil, err
	}
	recipe.Steps, err = s.repo.GetStepsByRecipeID(recipeID)
	if err != nil {
		return nil, err
	}
	recipe.LikeCount, err = s.likesRepo.GetLikeCount(recipeID)
	if err != nil {
		return nil, err
	}
	recipe.FavoriteCount, err = s.favoriteRepo.GetFavoriteCount(recipeID)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func (s *RecipeService) UpdateRecipe(recipeID int, title string, imagePath string, ingredients []model.Ingredient, steps []model.Step, userID int) error {

	//Vérification pour savoir si l'user peu modifier la recette
	recipe, err := s.repo.GetRecipeByID(recipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return errors.New("forbidden")
	}

	//Update recette
	err = s.repo.UpdateRecipe(recipeID, title, imagePath)
	if err != nil {
		return err
	}

	//Delete ingredients et steps
	err = s.repo.DeleteIngredientsByRecipeID(recipeID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteStepsByRecipeID(recipeID)
	if err != nil {
		return err
	}

	//Re création ingredients et steps
	for _, ing := range ingredients {
		err = s.repo.AddIngredient(recipeID, ing.Name, ing.Quantity, ing.Unit)
		if err != nil {
			return err
		}
	}

	for _, step := range steps {
		err = s.repo.AddStep(recipeID, step.Order, step.Description)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *RecipeService) DeleteRecipe(recipeID int, userID int) error {

	//Vérification pour savoir si l'user peu delete la recette
	recipe, err := s.repo.GetRecipeByID(recipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return errors.New("forbidden")
	}
	err = s.repo.DeleteIngredientsByRecipeID(recipeID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteStepsByRecipeID(recipeID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteRecipe(recipeID)
	if err != nil {
		return err
	}

	return err
}

func (s *RecipeService) GetAllRecipesShuffled() ([]*model.Recipe, error) {
	recipes, err := s.repo.GetAllRecipesShuffled()
	return recipes, err
}

func (s *RecipeService) GetAllRecipesFromUser(userID int) ([]*model.Recipe, error) {
	recipes, err := s.repo.GetAllRecipesFromUser(userID)
	return recipes, err
}

func (s *RecipeService) SearchRecipes(userID int, title string, ingredient string, sort string, liked bool, favorited bool) ([]*model.Recipe, error) {
	recipes, err := s.repo.SearchRecipes(userID, title, ingredient, sort, liked, favorited)
	return recipes, err
}