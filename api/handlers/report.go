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

type Report struct {
	DB databases.ReportDatabase
}

// TODO: add delete and update functionality

// ReportByIDHandler returns a report by a given ID
func (report Report) ReportByObjectIDHandler(w http.ResponseWriter, r *http.Request) {
	reportID := mux.Vars(r)["report_id"]

	rID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := report.DB.FindOne(context.Background(), bson.M{"_id": rID})
	if err != nil {
		config.ErrorStatus("failed to get report by ID", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Create a new report
func (report Report) NewReportHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var details models.ReportDetails // Json data will represent the report details model
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

	newReport := models.Report{
		ID:         primitive.NewObjectID(),
		AuthorID:   details.AuthorID,
		Title:      details.Title,
		Desc:       details.Desc,
		Severity:   -1, // -1 indicates it hasn't been assigned a severity level
		Resolved:   false,
		CommentIDs: []string{},
	}

	result, err := report.DB.InsertOne(ctx, newReport)
	if err != nil {
		config.ErrorStatus("failed to insert project", http.StatusBadRequest, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]interface{}{"result": result},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
