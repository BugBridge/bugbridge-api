package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/BugBridge/bugbridge-api/api/auth"
	"github.com/BugBridge/bugbridge-api/config"
	"github.com/BugBridge/bugbridge-api/databases"
	"github.com/BugBridge/bugbridge-api/models"
	"github.com/BugBridge/bugbridge-api/util"
)

type User struct {
	DB   databases.UserDatabase
	Auth *auth.AuthService
}

// temp
const TokenName = "bugbridge"

func (user User) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// defer and context

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.ErrorStatus("Invalid login request", http.StatusBadRequest, w, err)
		return
	}

	// Lookup user
	dbResp, err := user.DB.FindOne(context.Background(), bson.M{"email": req.Email})
	if err != nil {
		config.ErrorStatus("failed to get user by email", http.StatusNotFound, w, err)
		return
	}

	if dbResp == nil {
		config.ErrorStatus("user not found", http.StatusNotFound, w, err)
		return
	}

	// Check password with bcrypt
	// if bcrypt.CompareHashAndPassword([]byte(dbResp.Password), []byte(req.Password)) != nil {
	if dbResp.Password != req.Password {
		config.ErrorStatus("Incorrect password", http.StatusUnauthorized, w, err)
		return
	}

	// Sign JWT with sign func
	token, err := user.Auth.Sign(dbResp.ID.Hex())
	if err != nil {
		config.ErrorStatus("Failed to sign token", http.StatusInternalServerError, w, err)
		return
	}

	// return token
	resp := models.LoginResponse{Token: token}
	resp.User = *dbResp

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Logout clears the JWT cookie or
// instructs client to drop the token.
func (user User) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // change to true for https
	})

	w.WriteHeader(http.StatusNoContent)
}

// UserByIDHandler returns a user by a given ID
func (user User) UserByObjectIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := user.DB.FindOne(context.Background(), bson.M{"_id": uID})
	if err != nil {
		config.ErrorStatus("failed to get user by ID", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Create new user
func (user User) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var details models.UserDetails // Json data will represent the user details model
	defer cancel()

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&details); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	// TODO: validate inputs, ensure username is available, email, password

	newUser := models.User{
		ID:         primitive.NewObjectID(),
		ProjectIDs: []string{},
		Username:   details.Username,
		Email:      details.Email,
		Password:   details.Password,
	}

	result, err := user.DB.InsertOne(ctx, newUser)
	if err != nil {
		config.ErrorStatus("failed to insert user", http.StatusBadRequest, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]any{"result": result},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

// UpdateUserHandler updates the attributes of a user
func (user User) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var newDetails models.UserUpdateDetails
	defer cancel()

	userID := mux.Vars(r)["user_id"]

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&newDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&newDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	update := util.BuildUpdate(newDetails)

	dbResp, err := user.DB.UpdateOne(
		ctx,
		bson.M{"_id": uID},
		bson.M{"$set": update},
	)

	if err != nil {
		config.ErrorStatus("the user could not be updated", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (user User) DeleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := mux.Vars(r)["user_id"]

	uID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := user.DB.DeleteOne(ctx, bson.M{"_id": uID})
	if err != nil {
		config.ErrorStatus("failed to delete user", http.StatusNotFound, w, err)
		return
	}

	if dbResp.Dr.DeletedCount == 0 {
		config.ErrorStatus("User not found", http.StatusNotFound, w, nil)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
