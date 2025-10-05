package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const TokenName = "bugbridge"

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Pass  string `json:"-"`
}

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*UserInfo, error)
}

type AuthHandler struct {
	Auth *AuthService
	Repo UserRepo
}

func NewAuthHandler(auth *AuthService, repo UserRepo) *AuthHandler {
	return &AuthHandler{Auth: auth, Repo: repo}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid")
		return
	}
	// Lookup user
	user, err := h.Repo.GetByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	//Check password
	if bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Sign JWT with sign func
	token, err := h.Auth.Sign(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to sign token")
		return
	}

	// return token
	resp := loginResponse{Token: token}
	resp.User.ID = user.ID
	resp.User.Email = user.Email

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Logout clears the JWT cookie or
// instructs client to drop the token.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
	w.WriteHeader(http.StatusNoContent)
}

// writeError returns error if there is one.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
