package integration

import (
	"bytes"
	"context"
	"cookapp/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func TestSendFriendRequest(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	body := map[string]any{"receiver_id": 2}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/friends/request", bytes.NewBuffer(bodyJSON))
	ctx := context.WithValue(req.Context(), model.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.SendFriendRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestSendFriendRequestDuplicate(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)

	body := map[string]any{"receiver_id": 2}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/friends/request", bytes.NewBuffer(bodyJSON))
	ctx := context.WithValue(req.Context(), model.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.SendFriendRequest(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("attendu 409, reçu %d", w.Code)
	}
}

func TestAcceptFriendRequest(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)

	req := httptest.NewRequest(http.MethodPut, "/friends/1/accept", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	// Ajout du userID (user 2 = receiver)
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 2,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.AcceptFriendRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestDeclineFriendRequest(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)

	req := httptest.NewRequest(http.MethodPut, "/friends/1/decline", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	// Ajout du userID (user 2 = receiver)
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 2,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.DeclineFriendRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestDeleteFriend(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)
	acceptFriendRequest("1", 2, &h)

	req := httptest.NewRequest(http.MethodDelete, "/friends/1", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	// Ajout du userID (user 1 = requester, peut aussi supprimer)
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 1,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.DeleteFriend(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestGetFriends(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)
	acceptFriendRequest("1", 2, &h)

	req := httptest.NewRequest(http.MethodGet, "/friends", nil)
	ctx := context.WithValue(req.Context(), model.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.GetFriend(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestGetFriendsEmpty(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	req := httptest.NewRequest(http.MethodGet, "/friends", nil)
	ctx := context.WithValue(req.Context(), model.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.GetFriend(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestAcceptFriendRequestForbidden(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)
	registerUser("test3", "test3@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)

	req := httptest.NewRequest(http.MethodPut, "/friends/1/accept", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 3,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.AcceptFriendRequest(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("attendu 403, reçu %d", w.Code)
	}
}

func TestDeclineFriendRequestForbidden(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)
	registerUser("test3", "test3@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)

	req := httptest.NewRequest(http.MethodPut, "/friends/1/decline", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 3,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.DeclineFriendRequest(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("attendu 403, reçu %d", w.Code)
	}
}

func TestDeleteFriendForbidden(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)
	registerUser("test3", "test3@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)
	acceptFriendRequest("1", 2, &h)

	req := httptest.NewRequest(http.MethodDelete, "/friends/1", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	ctx := context.WithValue(
		context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx),
		model.UserIDKey, 3,
	)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.DeleteFriend(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("attendu 403, reçu %d", w.Code)
	}
}

func TestGetPendingFriendRequests(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)
	registerUser("test2", "test2@test.com", "1234", h.User)

	sendFriendRequest(1, 2, &h)

	req := httptest.NewRequest(http.MethodGet, "/friends/requests", nil)
	ctx := context.WithValue(req.Context(), model.UserIDKey, 2)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.GetPendingFriendRequests(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}

func TestGetPendingFriendRequestsEmpty(t *testing.T) {
	db := setupDB(t)
	h := setupHandlers(db)
	registerUser("test", "test@test.com", "1234", h.User)

	req := httptest.NewRequest(http.MethodGet, "/friends/requests", nil)
	ctx := context.WithValue(req.Context(), model.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Friendship.GetPendingFriendRequests(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("attendu 200, reçu %d", w.Code)
	}
}
