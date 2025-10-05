package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/BugBridge/bugbridge-api/databases"
	"github.com/BugBridge/bugbridge-api/models"
)

// GetBugReportsHandler retrieves all bug reports
func GetBugReportsHandler(w http.ResponseWriter, r *http.Request) {
	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	bugReportsCollection := dbHelper.GetCollection("bug_reports")

	cursor, err := bugReportsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		ErrorStatus("Failed to fetch bug reports", http.StatusInternalServerError, w, err)
		return
	}
	defer cursor.Close(context.Background())

	var bugReports []models.BugReport
	if err = cursor.All(context.Background(), &bugReports); err != nil {
		ErrorStatus("Failed to decode bug reports", http.StatusInternalServerError, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bugReports)
}

// GetUserBugReportsHandler retrieves bug reports for a specific user
func GetUserBugReportsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ErrorStatus("Invalid User ID", http.StatusBadRequest, w, err)
		return
	}

	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	bugReportsCollection := dbHelper.GetCollection("bug_reports")

	cursor, err := bugReportsCollection.Find(context.Background(), bson.M{"reporterId": objID})
	if err != nil {
		ErrorStatus("Failed to fetch user bug reports", http.StatusInternalServerError, w, err)
		return
	}
	defer cursor.Close(context.Background())

	var bugReports []models.BugReport
	if err = cursor.All(context.Background(), &bugReports); err != nil {
		ErrorStatus("Failed to decode bug reports", http.StatusInternalServerError, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bugReports)
}

// CreateBugReportHandler creates a new bug report
func CreateBugReportHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	reporterObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ErrorStatus("Invalid User ID in token", http.StatusUnauthorized, w, err)
		return
	}

	var reportReq struct {
		CompanyID        string `json:"companyId" validate:"required"`
		Title            string `json:"title" validate:"required"`
		Description      string `json:"description" validate:"required"`
		Severity         string `json:"severity" validate:"required,oneof=low medium high critical"`
		StepsToReproduce string `json:"stepsToReproduce"`
		IsAnonymous      bool   `json:"isAnonymous"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reportReq); err != nil {
		ErrorStatus("Invalid request body", http.StatusBadRequest, w, err)
		return
	}

	if err := validate.Struct(reportReq); err != nil {
		ErrorStatus("Validation error", http.StatusBadRequest, w, err)
		return
	}

	companyObjID, err := primitive.ObjectIDFromHex(reportReq.CompanyID)
	if err != nil {
		ErrorStatus("Invalid Company ID", http.StatusBadRequest, w, err)
		return
	}

	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)

	// Get company name
	companiesCollection := dbHelper.GetCollection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.Background(), bson.M{"_id": companyObjID}).Decode(&company)
	if err != nil {
		ErrorStatus("Company not found", http.StatusNotFound, w, err)
		return
	}

	// Get user info for reporter
	usersCollection := dbHelper.GetCollection("users")
	var user models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": reporterObjID}).Decode(&user)
	if err != nil {
		ErrorStatus("User not found", http.StatusNotFound, w, err)
		return
	}

	bugReportsCollection := dbHelper.GetCollection("bug_reports")

	newReport := models.BugReport{
		ID:               primitive.NewObjectID(),
		CompanyID:        companyObjID,
		CompanyName:      company.Name,
		ReporterID:       reporterObjID,
		Title:            reportReq.Title,
		Description:      reportReq.Description,
		Severity:         reportReq.Severity,
		Status:           "pending",
		StepsToReproduce: reportReq.StepsToReproduce,
		IsAnonymous:      reportReq.IsAnonymous,
		SubmittedAt:      time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = bugReportsCollection.InsertOne(context.Background(), newReport)
	if err != nil {
		ErrorStatus("Failed to create bug report", http.StatusInternalServerError, w, err)
		return
	}

	// Update company's bug reports count
	_, err = companiesCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": companyObjID},
		bson.M{"$inc": bson.M{"bugReportsCount": 1}},
	)
	if err != nil {
		// Log error but don't fail the request
		// This is not critical for the main functionality
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newReport)
}

// GetBugReportHandler retrieves a specific bug report by ID
func GetBugReportHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		ErrorStatus("Invalid Report ID", http.StatusBadRequest, w, err)
		return
	}

	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	bugReportsCollection := dbHelper.GetCollection("bug_reports")

	var bugReport models.BugReport
	err = bugReportsCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&bugReport)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ErrorStatus("Bug report not found", http.StatusNotFound, w, err)
			return
		}
		ErrorStatus("Database error", http.StatusInternalServerError, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bugReport)
}

// UpdateBugReportStatusHandler updates the status of a bug report
func UpdateBugReportStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		ErrorStatus("Invalid Report ID", http.StatusBadRequest, w, err)
		return
	}

	var statusReq struct {
		Status string `json:"status" validate:"required,oneof=pending under_review accepted resolved rejected"`
	}

	if err := json.NewDecoder(r.Body).Decode(&statusReq); err != nil {
		ErrorStatus("Invalid request body", http.StatusBadRequest, w, err)
		return
	}

	if err := validate.Struct(statusReq); err != nil {
		ErrorStatus("Validation error", http.StatusBadRequest, w, err)
		return
	}

	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	bugReportsCollection := dbHelper.GetCollection("bug_reports")

	_, err = bugReportsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"status":    statusReq.Status,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		ErrorStatus("Failed to update bug report status", http.StatusInternalServerError, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Status updated successfully"})
}

// AddCommentHandler adds a comment to a bug report
// Note: Comments functionality not implemented in current model
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	ErrorStatus("Comments functionality not implemented", http.StatusNotImplemented, w, nil)
}
