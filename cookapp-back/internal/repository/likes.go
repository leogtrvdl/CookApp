package repository

import (
	"database/sql"
)

type LikesRepository struct {
	db *sql.DB
}

func NewLikesRepository(db *sql.DB) *LikesRepository {
	return &LikesRepository{db: db}
}

func (r *LikesRepository) AddLike(userID int, recipeID int) error {
	_, err := r.db.Exec("INSERT INTO likes (user_id, recipe_id) VALUES ($1, $2)", userID, recipeID)
	return err
}

func (r *LikesRepository) DeleteLike(userID int, recipeID int) error {
	_, err := r.db.Exec("DELETE FROM likes WHERE user_id = $1 AND recipe_id = $2", userID, recipeID)
	return err
}

func (r *LikesRepository) GetLikeCount(recipeID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE recipe_id = $1", recipeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LikesRepository) IsLiked(userID int, recipeID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND recipe_id = $2", userID, recipeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
