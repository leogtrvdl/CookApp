package integration

import (
	"cookapp/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllIngredients(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	createRecipe("pâtes au saumon", "", []model.Ingredient{
		{Name: "saumon", Quantity: "200", Unit: "g"},
		{Name: "pâtes", Quantity: "200", Unit: "g"},
	}, []model.Step{}, 1, &h)

	req := httptest.NewRequest(http.MethodGet, "/ingredients", nil)
	w := httptest.NewRecorder()
	h.Ingredient.GetAllIngredients(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}
