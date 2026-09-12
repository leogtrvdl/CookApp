package handler

import (
	"cookapp/internal/model"
	"cookapp/internal/service"
	"cookapp/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type RecipeHandler struct {
	service *service.RecipeService
	storage *storage.R2Storage
}

func NewRecipeHandler(service *service.RecipeService, storage *storage.R2Storage) *RecipeHandler {
	return &RecipeHandler{service: service, storage: storage}
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	var ingredients []model.Ingredient
	ingredientsStr := r.FormValue("ingredients")
	if ingredientsStr != "" {
		err = json.Unmarshal([]byte(ingredientsStr), &ingredients)
		if err != nil {
			http.Error(w, "Ingrédients invalides", http.StatusBadRequest)
			return
		}
	}

	var steps []model.Step
	stepsStr := r.FormValue("steps")
	if stepsStr != "" {
		err = json.Unmarshal([]byte(stepsStr), &steps)
		if err != nil {
			http.Error(w, "Étapes invalides", http.StatusBadRequest)
			return
		}
	}

	imagePath := ""
	file, header, fileErr := r.FormFile("image")
	if fileErr == nil {
		defer file.Close()
		url, err := h.storage.UploadFile(file, header.Filename, "recipes")
		if err != nil {
			http.Error(w, "Erreur lors de l'upload de l'image", http.StatusInternalServerError)
			return
		}
		imagePath = url
	}

	err = h.service.CreateRecipe(userID, title, imagePath, ingredients, steps)
	if err != nil {
		http.Error(w, "Erreur lors de la création de la recette", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"recette": "créer",
	})
}

func (h *RecipeHandler) GetRecipeById(w http.ResponseWriter, r *http.Request) {

	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	recipe, err := h.service.GetRecipeByID(recipeID)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	var ingredients []model.Ingredient
	ingredientsStr := r.FormValue("ingredients")
	if ingredientsStr != "" {
		err = json.Unmarshal([]byte(ingredientsStr), &ingredients)
		if err != nil {
			http.Error(w, "Ingrédients invalides", http.StatusBadRequest)
			return
		}
	}

	var steps []model.Step
	stepsStr := r.FormValue("steps")
	if stepsStr != "" {
		err = json.Unmarshal([]byte(stepsStr), &steps)
		if err != nil {
			http.Error(w, "Étapes invalides", http.StatusBadRequest)
			return
		}
	}

	// Récupérer l'image_path actuelle si pas de nouvelle image
	existingRecipe, err := h.service.GetRecipeByID(recipeID)
	if err != nil {
		http.Error(w, "Recette introuvable", http.StatusNotFound)
		return
	}
	imagePath := existingRecipe.ImagePath

	file, header, fileErr := r.FormFile("image")
	if fileErr == nil {
		defer file.Close()
		url, err := h.storage.UploadFile(file, header.Filename, "recipes")
		if err != nil {
			http.Error(w, "Erreur lors de l'upload de l'image", http.StatusInternalServerError)
			return
		}
		h.storage.DeleteFile(existingRecipe.ImagePath)
		imagePath = url
	}

	err = h.service.UpdateRecipe(recipeID, title, imagePath, ingredients, steps, userID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"recipe": "modifié",
	})
}

func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	existingRecipe, err := h.service.GetRecipeByID(recipeID)
	if err != nil {
		http.Error(w, "Recette introuvable", http.StatusNotFound)
		return
	}
	if existingRecipe.ImagePath != "" {
		h.storage.DeleteFile(existingRecipe.ImagePath)
	}

	err = h.service.DeleteRecipe(recipeID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"recipe": "recette supprimé",
	})
}

func (h *RecipeHandler) GetAllRecipesShuffled(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.service.GetAllRecipesShuffled()
	if err != nil {
		http.Error(w, "Recipes not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) GetAllRecipesFromUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipes, err := h.service.GetAllRecipesFromUser(userID)
	if err != nil {
		json.NewEncoder(w).Encode([]any{})
		return
	}
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(model.UserIDKey).(int)
	title := r.URL.Query().Get("title")
	ingredient := r.URL.Query().Get("ingredient")
	sort := r.URL.Query().Get("sort")
	liked := r.URL.Query().Get("liked") == "true"
	favorited := r.URL.Query().Get("favorited") == "true"

	recipes, err := h.service.SearchRecipes(userID, title, ingredient, sort, liked, favorited)

	if err != nil {
		http.Error(w, "Recipes not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(recipes)
}
