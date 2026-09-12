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

type LikesHandler struct {
	service *service.LikesService
}

func NewLikesHandler(service *service.LikesService) *LikesHandler {
	return &LikesHandler{service: service}
}

func (h *LikesHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.AddLike(userID, recipeID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			http.Error(w, "Like déjà mit", http.StatusConflict)
			return
		}
		http.Error(w, "Erreur lors de l'envoi du like", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "enregistré",
	})
}

func (h *LikesHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteLike(userID, recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "like supprimé",
	})
}

func (h *LikesHandler) GetLikeCount(w http.ResponseWriter, r *http.Request) {
	recipeID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	likeCount, err := h.service.GetLikeCount(recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"like_count": likeCount,
	})
}

func (h *LikesHandler) IsLiked(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	recipeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	isLiked, err := h.service.IsLiked(userID, recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{
		"is_liked": isLiked,
	})
}
