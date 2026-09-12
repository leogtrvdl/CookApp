package handler

import (
	"cookapp/internal/model"
	"cookapp/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(service *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

func (h *FavoriteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.AddFavorite(userID, recipeID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			http.Error(w, "Déjà en favoris", http.StatusConflict)
			return
		}
		http.Error(w, "Erreur lors de l'ajout en favoris", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ajouté aux favoris",
	})
}

func (h *FavoriteHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteFavorite(userID, recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "retiré des favoris",
	})
}

func (h *FavoriteHandler) GetFavoriteCount(w http.ResponseWriter, r *http.Request) {
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	count, err := h.service.GetFavoriteCount(recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"favorite_count": count,
	})
}

func (h *FavoriteHandler) IsFavorited(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	IsFavorited, err := h.service.IsFavorited(userID, recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{
		"is_favorited": IsFavorited,
	})
}
