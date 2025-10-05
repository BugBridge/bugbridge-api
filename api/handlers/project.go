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

type Project struct {
	DB databases.ProjectDatabase
}

// TODO: add delete and update functionality

// ProjectByIDHandler returns a project by a given ID
func (project Project) ProjectByObjectIDHandler(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	pID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := project.DB.FindOne(context.Background(), bson.M{"_id": pID})
	if err != nil {
		config.ErrorStatus("failed to get project by ID", http.StatusNotFound, w, err)
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

// Create a new project
func (project Project) NewProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var details models.ProjectDetails // Json data will represent the report details model
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

	// TODO: add validation to title / description length

	newProject := models.Project{
		ID:        primitive.NewObjectID(),
		Name:      details.Name,
		Des:       details.Des,
		Template:  details.Template,
		OwnerID:   details.OwnerID,
		AdminsIDs: []string{},
	}

	result, err := project.DB.InsertOne(ctx, newProject)
	if err != nil {
		config.ErrorStatus("failed to insert project", http.StatusBadRequest, w, err)
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

func (project Project) DeleteProjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	uID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := project.DB.DeleteOne(context.Background(), bson.M{"_id": uID})
	if err != nil {
		config.ErrorStatus("failed to delete project", http.StatusNotFound, w, err)
		return
	}

	if dbResp.Dr.DeletedCount == 0 {
		config.ErrorStatus("Project not found", http.StatusNotFound, w, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
