package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/BugBridge/bugbridge-api/databases"
	"github.com/BugBridge/bugbridge-api/models"
	"golang.org/x/crypto/bcrypt"
)

// validate is already declared in api.go

// JWT Claims structure
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		ErrorStatus("Invalid request body", http.StatusBadRequest, w, err)
		return
	}

	if err := validate.Struct(loginReq); err != nil {
		ErrorStatus("Validation failed", http.StatusBadRequest, w, err)
		return
	}

	// Get database from context
	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	usersCollection := dbHelper.GetCollection("users")

	// Find user by email
	var user models.User
	err := usersCollection.FindOne(context.Background(), bson.M{"email": loginReq.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ErrorStatus("Invalid credentials", http.StatusUnauthorized, w, err)
			return
		}
		ErrorStatus("Database error", http.StatusInternalServerError, w, err)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		ErrorStatus("Invalid credentials", http.StatusUnauthorized, w, err)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		ErrorStatus("Token generation failed", http.StatusInternalServerError, w, err)
		return
	}

	// Check if user has a company
	companiesCollection := dbHelper.GetCollection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.Background(), bson.M{"ownerId": user.ID}).Decode(&company)

	var companyProfile *models.Company
	if err == nil {
		companyProfile = &company
	}

	// Prepare response
	response := models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Username:  user.Username,
			CompanyID: user.CompanyID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		CompanyProfile: companyProfile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SignupHandler handles user registration
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signupReq models.SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&signupReq); err != nil {
		ErrorStatus("Invalid request body", http.StatusBadRequest, w, err)
		return
	}

	if err := validate.Struct(signupReq); err != nil {
		ErrorStatus("Validation failed", http.StatusBadRequest, w, err)
		return
	}

	// Get database from context
	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	usersCollection := dbHelper.GetCollection("users")

	// Check if user already exists
	var existingUser models.User
	err := usersCollection.FindOne(context.Background(), bson.M{"email": signupReq.Email}).Decode(&existingUser)
	if err == nil {
		ErrorStatus("Email already exists", http.StatusConflict, w, err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		ErrorStatus("Password hashing failed", http.StatusInternalServerError, w, err)
		return
	}

	// Create new user
	user := models.User{
		ID:        primitive.NewObjectID(),
		Email:     signupReq.Email,
		Password:  string(hashedPassword),
		Name:      signupReq.FirstName + " " + signupReq.LastName,
		Username:  signupReq.Email[:len(signupReq.Email)-len("@example.com")], // Simple username generation
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user into database
	_, err = usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		ErrorStatus("Failed to create user", http.StatusInternalServerError, w, err)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		ErrorStatus("Token generation failed", http.StatusInternalServerError, w, err)
		return
	}

	// Prepare response
	response := models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Username:  user.Username,
			CompanyID: user.CompanyID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		CompanyProfile: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUserHandler handles getting current user info
func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID := r.Context().Value("userID").(string)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ErrorStatus("Invalid user ID", http.StatusBadRequest, w, err)
		return
	}

	// Get database from context
	dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
	usersCollection := dbHelper.GetCollection("users")

	// Find user
	var user models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ErrorStatus("User not found", http.StatusNotFound, w, err)
			return
		}
		ErrorStatus("Database error", http.StatusInternalServerError, w, err)
		return
	}

	// Check if user has a company
	companiesCollection := dbHelper.GetCollection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.Background(), bson.M{"ownerId": user.ID}).Decode(&company)

	var companyProfile *models.Company
	if err == nil {
		companyProfile = &company
	}

	// Prepare response
	response := models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Username:  user.Username,
			CompanyID: user.CompanyID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		CompanyProfile: companyProfile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateJWT generates a JWT token
func generateJWT(userID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production" // fallback for development
	}
	return token.SignedString([]byte(jwtSecret))
}
