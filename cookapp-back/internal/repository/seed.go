package repository

import (
	"cookapp/internal/model"
	"database/sql"
	"encoding/json"
	"os"
)

type seedRecipe struct {
	Title       string             `json:"title"`
	ImagePath   string             `json:"image_path"`
	Ingredients []model.Ingredient `json:"ingredients"`
	Steps       []model.Step       `json:"steps"`
}

func SeedDB(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&count)
	if err != nil {
		return
	}
	if count > 0 {
		return
	}

	data, err := os.ReadFile("seeds.json")
	if err != nil {
		return
	}

	var recipes []seedRecipe
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		return
	}

	for _, recipe := range recipes {
		var recipeID int
		err := db.QueryRow("INSERT INTO recipes (user_id, title, image_path) VALUES ($1, $2, $3) RETURNING id", 0, recipe.Title, recipe.ImagePath).Scan(&recipeID)
		if err != nil {
			return
		}

		for _, ing := range recipe.Ingredients {
			db.Exec("INSERT INTO ingredients (recipe_id, name, quantity, unit) VALUES ($1, $2, $3, $4)", recipeID, ing.Name, ing.Quantity, ing.Unit)
		}

		for _, step := range recipe.Steps {
			db.Exec("INSERT INTO steps (recipe_id, step_order, description) VALUES ($1, $2, $3)", recipeID, step.Order, step.Description)
		}
	}
}
