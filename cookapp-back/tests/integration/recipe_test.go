package integration

import (
	"bytes"
	"context"
	"cookapp/internal/model"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func TestCreateRecipe(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Ma recette")
	writer.WriteField("ingredients", `[{"name":"farine","quantity":"200","unit":"g"}]`)
	writer.WriteField("steps", `[{"order":1,"description":"Mélanger"}]`)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/recipes", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), model.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Recipe.CreateRecipe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestGetRecipeByID(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("recette", "image.jpg", []model.Ingredient{}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()
	h.Recipe.GetRecipeById(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestUpdateRecipe(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("recette", "", []model.Ingredient{}, []model.Step{}, 1, &h)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "recette modifiée")
	writer.WriteField("ingredients", `[{"name":"farine","quantity":"200","unit":"g"}]`)
	writer.WriteField("steps", `[{"order":1,"description":"Mélanger"}]`)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPut, "/recipes/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 1,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Recipe.UpdateRecipe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestUpdateRecipeForbidden(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)
	createRecipe("recette", "", []model.Ingredient{}, []model.Step{}, 1, &h)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "recette modifiée")
	writer.WriteField("ingredients", `[]`)
	writer.WriteField("steps", `[]`)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPut, "/recipes/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 2,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Recipe.UpdateRecipe(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("attendu 403, reçu %d", w.Code)
	}
}

func TestDeleteRecipe(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)
	createRecipe("recette", "image.jpg", []model.Ingredient{}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodDelete, "/recipes/1", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 1,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Recipe.DeleteRecipe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestDeleteRecipeForbidden(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("recette", "image.jpg", []model.Ingredient{}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodDelete, "/recipes/1", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 2,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Recipe.DeleteRecipe(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("attendu 403, reçu %d", w.Code)
	}
}

func TestSearchRecipesNoFilter(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("pâtes au saumon", "", []model.Ingredient{
		{Name: "pâtes", Quantity: "200", Unit: "g"},
	}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search", nil)
	w := httptest.NewRecorder()
	h.Recipe.SearchRecipes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestSearchRecipesByTitle(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("pâtes au saumon", "", []model.Ingredient{}, []model.Step{}, 1, &h)
	createRecipe("pizza margherita", "", []model.Ingredient{}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?title=pâtes", nil)
	w := httptest.NewRecorder()
	h.Recipe.SearchRecipes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}

	var recipes []*model.Recipe
	json.NewDecoder(w.Body).Decode(&recipes)
	if len(recipes) != 1 {
		t.Errorf("attendu 1 recette, reçu %d", len(recipes))
	}
}

func TestSearchRecipesByIngredient(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("pâtes au saumon", "", []model.Ingredient{
		{Name: "saumon", Quantity: "200", Unit: "g"},
	}, []model.Step{}, 1, &h)
	createRecipe("pizza margherita", "", []model.Ingredient{
		{Name: "mozzarella", Quantity: "100", Unit: "g"},
	}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?ingredient=saumon", nil)
	w := httptest.NewRecorder()
	h.Recipe.SearchRecipes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}

	var recipes []*model.Recipe
	json.NewDecoder(w.Body).Decode(&recipes)
	if len(recipes) != 1 {
		t.Errorf("attendu 1 recette, reçu %d", len(recipes))
	}
}

func TestSearchRecipesSortByMostLiked(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("pâtes au saumon", "", []model.Ingredient{}, []model.Step{}, 1, &h)
	createRecipe("pizza margherita", "", []model.Ingredient{}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?sort=most_liked", nil)
	w := httptest.NewRecorder()
	h.Recipe.SearchRecipes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}
