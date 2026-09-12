package repository

import (
	"cookapp/internal/model"
	"database/sql"
	"fmt"
	"strings"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) CreateRecipe(userID int, title string, imagePath string) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO recipes (user_id, title, image_path) VALUES ($1, $2, $3) RETURNING id", userID, title, imagePath).Scan(&id)
	return int(id), err
}

func (r *RecipeRepository) AddIngredient(recipeID int, name string, quantity string, unit string) error {
	_, err := r.db.Exec("INSERT INTO ingredients (recipe_id, name, quantity, unit) VALUES ($1, $2, $3, $4)", recipeID, name, quantity, unit)
	return err
}

func (r *RecipeRepository) AddStep(recipeID int, order int, description string) error {
	_, err := r.db.Exec("INSERT INTO steps (recipe_id, step_order, description) VALUES ($1, $2, $3)", recipeID, order, description)
	return err
}

func (r *RecipeRepository) GetRecipeByID(recipeID int) (*model.Recipe, error) {
	recipe := &model.Recipe{}
	err := r.db.QueryRow("SELECT id, user_id, title, image_path, created_at, updated_at FROM recipes WHERE id = $1", recipeID).Scan(
		&recipe.ID,
		&recipe.UserID,
		&recipe.Title,
		&recipe.ImagePath,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func (r *RecipeRepository) GetIngredientsByRecipeID(recipeID int) ([]*model.Ingredient, error) {
	rows, err := r.db.Query("SELECT id, recipe_id, name, quantity, unit FROM ingredients WHERE recipe_id = $1", recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []*model.Ingredient
	for rows.Next() {
		i := &model.Ingredient{}
		err := rows.Scan(&i.ID, &i.RecipeID, &i.Name, &i.Quantity, &i.Unit)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *RecipeRepository) GetStepsByRecipeID(recipeID int) ([]*model.Step, error) {
	rows, err := r.db.Query("SELECT id, recipe_id, step_order, description FROM steps WHERE recipe_id = $1", recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*model.Step
	for rows.Next() {
		s := &model.Step{}
		err := rows.Scan(&s.ID, &s.RecipeID, &s.Order, &s.Description)
		if err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *RecipeRepository) UpdateRecipe(recipeID int, title string, imagePath string) error {
	_, err := r.db.Exec("UPDATE recipes SET title = $1, image_path = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3", title, imagePath, recipeID)
	return err
}

func (r *RecipeRepository) DeleteRecipe(recipeID int) error {
	_, err := r.db.Exec("DELETE FROM recipes WHERE id = $1", recipeID)
	return err
}

func (r *RecipeRepository) DeleteIngredientsByRecipeID(recipeID int) error {
	_, err := r.db.Exec("DELETE FROM ingredients WHERE recipe_id = $1", recipeID)
	return err
}

func (r *RecipeRepository) DeleteStepsByRecipeID(recipeID int) error {
	_, err := r.db.Exec("DELETE FROM steps WHERE recipe_id = $1", recipeID)
	return err
}

func (r *RecipeRepository) GetAllRecipesShuffled() ([]*model.Recipe, error) {
	rows, err := r.db.Query(`SELECT id, user_id, title, image_path,
	(SELECT COUNT(*) FROM likes WHERE recipe_id = recipes.id) as like_count,
	(SELECT COUNT(*) FROM favorites WHERE recipe_id = recipes.id) as favorite_count,
	created_at, updated_at
	FROM recipes ORDER BY RANDOM()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*model.Recipe
	for rows.Next() {
		r := &model.Recipe{}
		err := rows.Scan(&r.ID, &r.UserID, &r.Title, &r.ImagePath, &r.LikeCount, &r.FavoriteCount, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	if recipes == nil {
		recipes = []*model.Recipe{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *RecipeRepository) GetAllRecipesFromUser(userID int) ([]*model.Recipe, error) {
	rows, err := r.db.Query(`SELECT id, user_id, title, image_path,
	(SELECT COUNT(*) FROM likes WHERE recipe_id = recipes.id) as like_count,
	(SELECT COUNT(*) FROM favorites WHERE recipe_id = recipes.id) as favorite_count,
	created_at, updated_at
	FROM recipes WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*model.Recipe
	for rows.Next() {
		r := &model.Recipe{}
		err := rows.Scan(&r.ID, &r.UserID, &r.Title, &r.ImagePath, &r.LikeCount, &r.FavoriteCount, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	if recipes == nil {
		recipes = []*model.Recipe{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *RecipeRepository) SearchRecipes(userID int, title string, ingredient string, sort string, liked bool, favorited bool) ([]*model.Recipe, error) {
	query := `SELECT DISTINCT recipes.id, recipes.user_id, recipes.title, recipes.image_path,
		(SELECT COUNT(*) FROM likes WHERE recipe_id = recipes.id) as like_count,
		(SELECT COUNT(*) FROM favorites WHERE recipe_id = recipes.id) as favorite_count,
		recipes.created_at, recipes.updated_at, RANDOM() as rand_order
		FROM recipes`

	conditions := []string{}
	args := []any{}

	if ingredient != "" {
		query += " JOIN ingredients ON ingredients.recipe_id = recipes.id"
		conditions = append(conditions, fmt.Sprintf("LOWER(ingredients.name) LIKE $%d", len(args)+1))
		args = append(args, "%"+strings.ToLower(ingredient)+"%")
	}

	if title != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(recipes.title) LIKE $%d", len(args)+1))
		args = append(args, "%"+strings.ToLower(title)+"%")
	}

	if liked {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM likes WHERE likes.recipe_id = recipes.id AND likes.user_id = $%d)", len(args)+1))
		args = append(args, userID)
	}

	if favorited {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM favorites WHERE favorites.recipe_id = recipes.id AND favorites.user_id = $%d)", len(args)+1))
		args = append(args, userID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	switch sort {
	case "most_liked":
		query += " ORDER BY like_count DESC"
	case "most_favorited":
		query += " ORDER BY favorite_count DESC"
	case "newest":
		query += " ORDER BY recipes.created_at DESC"
	default:
		query += " ORDER BY rand_order"
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*model.Recipe
	var randOrder float64
	for rows.Next() {
		rec := &model.Recipe{}
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.Title, &rec.ImagePath, &rec.LikeCount, &rec.FavoriteCount, &rec.CreatedAt, &rec.UpdatedAt, &randOrder)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if recipes == nil {
		recipes = []*model.Recipe{}
	}

	return recipes, nil
}
