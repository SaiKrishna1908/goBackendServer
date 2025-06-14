package api

import (
	"encoding/json"
	"goBackendServer/internal/store"
	"goBackendServer/internal/tokens"
	"goBackendServer/internal/utils"
	"log"
	"net/http"
	"time"
)

type TokenHanlder struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHanlder {
	return &TokenHanlder{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHanlder) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)

	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByUsername: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	passwordDoMatch, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("Error PasswordHash.Matches %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordDoMatch {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 25*time.Hour, tokens.ScopeAuth)

	if err != nil {
		h.logger.Printf("ERROR: Creating token %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
