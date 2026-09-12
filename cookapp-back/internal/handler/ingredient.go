package handler

import (
	"encoding/json"
	"net/http"
	"cookapp/internal/service"
)

type IngredientHandler struct {
	service *service.IngredientService
}

func NewIngredientHandler(service *service.IngredientService) *IngredientHandler {
	return &IngredientHandler{service: service}
}
func (h *IngredientHandler) GetAllIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.service.GetAllIngredients()
	if err != nil {
		http.Error(w, "Ingredients not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(ingredients)
}
