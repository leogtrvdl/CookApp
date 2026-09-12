package integration

import (
	"bytes"
	"context"
	"cookapp/internal/handler"
	"cookapp/internal/model"
	"cookapp/internal/repository"
	"cookapp/internal/service"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupDB(t *testing.T) *sql.DB {
	godotenv.Load("../../.env")

	dbURL := os.Getenv("TEST_DATABASE_URL")
	
	if dbURL == "" {
		t.Fatal("TEST_DATABASE_URL non définie")
	}
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatal(err)
	}

	repository.InitDB(db)

	_, err = db.Exec(`TRUNCATE users, recipes, ingredients, steps, friendship, likes, favorites RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatal(err)
	}

	repository.InitDB(db)

	return db
}

type Handlers struct {
	User       *handler.UserHandler
	Friendship *handler.FriendshipHandler
	Recipe     *handler.RecipeHandler
	Ingredient *handler.IngredientHandler
}

func setupHandlers(db *sql.DB) Handlers {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	friendshipRepo := repository.NewFriendshipRepository(db)
	friendshipService := service.NewFriendshipService(friendshipRepo)

	recipeRepo := repository.NewRecipeRepository(db)
	likesRepo := repository.NewLikesRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	recipeService := service.NewRecipeService(recipeRepo, likesRepo, favoriteRepo)

	ingredientRepo := repository.NewIngredientRepository(db)
	ingredientService := service.NewIngredientService(ingredientRepo)

	return Handlers{
		User:       handler.NewUserHandler(userService, nil),
		Friendship: handler.NewFriendshipHandler(friendshipService),
		Recipe:     handler.NewRecipeHandler(recipeService, nil),
		Ingredient: handler.NewIngredientHandler(ingredientService),
	}
}

func registerUser(username string, email string, password string, h *handler.UserHandler) {
	body := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)
}

func sendFriendRequest(requesterID int, receiverID int, h *Handlers) {
	body := map[string]any{
		"receiver_id": receiverID,
	}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/friends/request", bytes.NewBuffer(bodyJSON))
	ctx := context.WithValue(req.Context(), model.UserIDKey, requesterID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Friendship.SendFriendRequest(w, req)
}

func acceptFriendRequest(friendshipID string, userID int, h *Handlers) {
	req := httptest.NewRequest(http.MethodPut, "/friends/"+friendshipID+"/accept", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", friendshipID)
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, userID,
	)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Friendship.AcceptFriendRequest(w, req)
}

func createRecipe(title string, image_path string, ingredients []model.Ingredient, steps []model.Step, userID int, h *Handlers) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("title", title)

	ingredientsJSON, _ := json.Marshal(ingredients)
	writer.WriteField("ingredients", string(ingredientsJSON))

	stepsJSON, _ := json.Marshal(steps)
	writer.WriteField("steps", string(stepsJSON))

	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/recipes", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), model.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Recipe.CreateRecipe(w, req)
}
