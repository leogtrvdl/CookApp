package handler

import (
	"cookapp/internal/model"
	"cookapp/internal/service"
	"cookapp/internal/storage"
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
	service *service.UserService
	storage *storage.R2Storage
}

func NewUserHandler(service *service.UserService, storage *storage.R2Storage) *UserHandler {
	return &UserHandler{service: service, storage: storage}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)

	err := h.service.Register(body["username"], body["email"], body["password"])
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			http.Error(w, "Username ou email déjà utilisé", http.StatusConflict)
			return
		}
		http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "enregistré",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)

	tokenString, err := h.service.Login(body["username"], body["password"])
	if err != nil {
		http.Error(w, "Identifiants incorrects", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)

	userID := r.Context().Value(model.UserIDKey).(int)
	err := h.service.UpdateUser(userID, body["username"], body["email"], body["bio"])
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(body)
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)

	existingUser, err := h.service.GetUserByID(userID)
	if err == nil && existingUser.AvatarPath != "" {
		h.storage.DeleteFile(existingUser.AvatarPath)
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "compte supprimé",
	})
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDKey).(int)
	file, header, err := r.FormFile("avatar")

	if err != nil {
		http.Error(w, "Fichier invalide", http.StatusBadRequest)
		return
	}
	defer file.Close()
	url, err := h.storage.UploadFile(file, header.Filename, "avatars")
	if err != nil {
		http.Error(w, "Erreur lors de l'upload de l'image", http.StatusInternalServerError)
		return
	}
	existingUser, err := h.service.GetUserByID(userID)
	if err == nil && existingUser.AvatarPath != "" {
		h.storage.DeleteFile(existingUser.AvatarPath)
	}
	avatarPath := url

	err = h.service.UpdateAvatar(userID, avatarPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status": "avatar ajouté",
	})
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	user, err := h.service.SearchUsers(username)
	if err != nil {
		json.NewEncoder(w).Encode([]any{})
		return
	}
	json.NewEncoder(w).Encode(user)
}
