package repository

import (
	"cookapp/internal/model"
	"database/sql"
)

type IngredientRepository struct {
	db *sql.DB
}

func NewIngredientRepository(db *sql.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

func (r *IngredientRepository) GetAllIngredients() ([]*model.Ingredient, error) {
	rows, err := r.db.Query("SELECT DISTINCT name FROM ingredients WHERE name NOT IN ('eau', 'sel', 'poivre') ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []*model.Ingredient
	for rows.Next() {
		i := &model.Ingredient{}
		err := rows.Scan(&i.Name)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	if ingredients == nil {
		ingredients = []*model.Ingredient{}
	}

	return ingredients, nil
}
