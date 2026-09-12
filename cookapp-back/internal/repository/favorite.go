package repository

import (
	"database/sql"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) AddFavorite(userID int, recipeID int) error {
	_, err := r.db.Exec("INSERT INTO favorites (user_id, recipe_id) VALUES ($1, $2)", userID, recipeID)
	return err
}

func (r *FavoriteRepository) DeleteFavorite(userID int, recipeID int) error {
	_, err := r.db.Exec("DELETE FROM favorites WHERE user_id = $1 AND recipe_id = $2", userID, recipeID)
	return err
}

func (r *FavoriteRepository) GetFavoriteCount(recipeID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM favorites WHERE recipe_id = $1", recipeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FavoriteRepository) IsFavorited(userID int, recipeID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = $1 AND recipe_id = $2", userID, recipeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}