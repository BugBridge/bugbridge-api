package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/BugBridge/bugbridge-api/databases"
	"github.com/BugBridge/bugbridge-api/models"
)

// GetCompaniesHandler retrieves all companies
func GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	// Get database helper from context
	dbHelperInterface := r.Context().Value("dbHelper")
	if dbHelperInterface == nil {
		ErrorStatus("Database connection not available", http.StatusInternalServerError, w, nil)
		return
	}

	dbHelper := dbHelperInterface.(databases.DatabaseHelper)
	companiesCollection := dbHelper.GetCollection("companies")

	cursor, err := companiesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		ErrorStatus("Failed to fetch companies", http.StatusInternalServerError, w, err)
		return
	}
	defer cursor.Close(context.Background())

	var companies []models.Company
	if err = cursor.All(context.Background(), &companies); err != nil {
		ErrorStatus("Failed to decode companies", http.StatusInternalServerError, w, err)
		return
	}

	// If no companies found, return empty array instead of error
	if companies == nil {
		companies = []models.Company{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

// CreateCompanyHandler creates a new company
func CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ErrorStatus("Invalid User ID in token", http.StatusUnauthorized, w, err)
		return
	}

	var companyReq struct {
		Name              string      `json:"name" validate:"required"`
		Industry          string      `json:"industry"`
		Description       string      `json:"description"`
		Website           string      `json:"website"`
		BugReportTemplate interface{} `json:"bugReportTemplate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&companyReq); err != nil {
		ErrorStatus("Invalid request body", http.StatusBadRequest, w, err)
		return
	}

	if err := validate.Struct(companyReq); err != nil {
		ErrorStatus("Validation error", http.StatusBadRequest, w, err)
		return
	}

	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	companiesCollection := dbHelper.GetCollection("companies")

	// Check if user already has a company
	var existingCompany models.Company
	err = companiesCollection.FindOne(context.Background(), bson.M{"ownerId": objID}).Decode(&existingCompany)
	if err == nil {
		ErrorStatus("User already has a company", http.StatusConflict, w, err)
		return
	}
	if err != mongo.ErrNoDocuments {
		ErrorStatus("Database error", http.StatusInternalServerError, w, err)
		return
	}

	newCompany := models.Company{
		ID:                primitive.NewObjectID(),
		OwnerID:           objID,
		Name:              companyReq.Name,
		Industry:          companyReq.Industry,
		Description:       companyReq.Description,
		Website:           companyReq.Website,
		AcceptingReports:  true,
		BugReportTemplate: companyReq.BugReportTemplate,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	_, err = companiesCollection.InsertOne(context.Background(), newCompany)
	if err != nil {
		ErrorStatus("Failed to create company", http.StatusInternalServerError, w, err)
		return
	}

	// Link user to the company
	usersCollection := dbHelper.GetCollection("users")
	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"companyId": newCompany.ID,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		// Log error but don't fail the request
		// The company was created successfully
		zap.S().With(err).Error("Failed to link user to company")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCompany)
}

// JoinCompanyHandler allows a user to join an existing company
func JoinCompanyHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ErrorStatus("Invalid User ID in token", http.StatusUnauthorized, w, err)
		return
	}

	var joinReq struct {
		CompanyID string `json:"companyId" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&joinReq); err != nil {
		ErrorStatus("Invalid request body", http.StatusBadRequest, w, err)
		return
	}

	if err := validate.Struct(joinReq); err != nil {
		ErrorStatus("Validation error", http.StatusBadRequest, w, err)
		return
	}

	companyObjID, err := primitive.ObjectIDFromHex(joinReq.CompanyID)
	if err != nil {
		ErrorStatus("Invalid Company ID", http.StatusBadRequest, w, err)
		return
	}

	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)

	// Check if company exists
	companiesCollection := dbHelper.GetCollection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.Background(), bson.M{"_id": companyObjID}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ErrorStatus("Company not found", http.StatusNotFound, w, err)
			return
		}
		ErrorStatus("Database error", http.StatusInternalServerError, w, err)
		return
	}

	// Check if user is already in a company
	usersCollection := dbHelper.GetCollection("users")
	var user models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		ErrorStatus("User not found", http.StatusNotFound, w, err)
		return
	}

	if user.CompanyID != nil {
		ErrorStatus("User is already part of a company", http.StatusConflict, w, nil)
		return
	}

	// Link user to the company
	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"companyId": companyObjID,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		ErrorStatus("Failed to join company", http.StatusInternalServerError, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Successfully joined company",
		"companyId":   companyObjID.Hex(),
		"companyName": company.Name,
	})
}
