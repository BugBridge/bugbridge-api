package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/BugBridge/bugbridge-api/config"
	"github.com/BugBridge/bugbridge-api/databases"
	"github.com/BugBridge/bugbridge-api/models"
)

type User struct {
	DB databases.UserDatabase
}

// TODO: add delete and update functionality

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
