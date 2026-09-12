package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRegister(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)

	body := map[string]string{
		"username": "test",
		"email":    "test@test.com",
		"password": "1234",
	}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.User.Register(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	//Login de l'user créer précédement
	body := map[string]string{
		"username": "test",
		"password": "1234",
	}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.User.Login(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestRegisterDuplicateUsername(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	//Creation d'un 2eme user avec les mêmes info
	body := map[string]string{
		"username": "test",
		"email":    "test@test.com",
		"password": "1234",
	}
	bodyJSON2, _ := json.Marshal(body)
	req2 := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(bodyJSON2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	h.User.Register(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("attendu 409, reçu %d", w2.Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	//Login de l'user créer précédement avec mauvais password
	body := map[string]string{
		"username": "test",
		"password": "123", //mauvais password
	}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.User.Login(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("attendu 401, reçu %d", w.Code)
	}
}

func TestSearchUsers(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	req := httptest.NewRequest(http.MethodGet, "/users/search?username=test", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.User.SearchUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestSearchUsersNotFound(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	req := httptest.NewRequest(http.MethodGet, "/users/search?username=inexistant", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.User.SearchUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}
