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

type FriendshipHandler struct {
	service *service.FriendshipService
}

func NewFriendshipHandler(service *service.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{service: service}
}

func (h *FriendshipHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	json.NewDecoder(r.Body).Decode(&body)
	requesterID := r.Context().Value(model.UserIDKey).(int)
	receiverID := int(body["receiver_id"].(float64))

	err := h.service.CreateFriendRequest(requesterID, receiverID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			http.Error(w, "Demande d'ami déjà envoyée", http.StatusConflict)
			return
		}
		http.Error(w, "Erreur lors de l'envoi de la demande", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "enregistré",
	})
}

func (h *FriendshipHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {

	receiverID := r.Context().Value(model.UserIDKey).(int)
	friendshipID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.UpdateFriendRequestStatus(friendshipID, "accepted", receiverID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "demande acceptée",
	})
}

func (h *FriendshipHandler) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {

	receiverID := r.Context().Value(model.UserIDKey).(int)
	friendshipID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.UpdateFriendRequestStatus(friendshipID, "declined", receiverID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "demande refusée",
	})
}

func (h *FriendshipHandler) DeleteFriend(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(model.UserIDKey).(int)
	friendshipID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteFriend(friendshipID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ami supprimé",
	})
}

func (h *FriendshipHandler) GetFriend(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(model.UserIDKey).(int)

	friendship, err := h.service.GetFriends(userID)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(friendship)
}

func (h *FriendshipHandler) GetPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)

	friendship, err := h.service.GetPendingFriendRequests(userID)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(friendship)
}