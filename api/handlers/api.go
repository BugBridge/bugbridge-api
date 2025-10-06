package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/BugBridge/bugbridge-api/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/BugBridge/bugbridge-api/config"
	"github.com/BugBridge/bugbridge-api/databases"
)

var validate = validator.New()

// App stores the router and db connection so it can be reused
type App struct {
	Router   *mux.Router
	DB       databases.CollectionHelper
	Config   config.Config
	dbHelper databases.DatabaseHelper
}

// New creates a new mux router and all the routes
func (a *App) New() *mux.Router {

	// System initialization (removed newSystem for now)

	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, req)
		})
	})

	// Add app context middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), "app", a)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	})

	// create database handlers like this
	// cow := Cow{DB: databases.NewCowDatabase(a.dbHelper)}

	// healthcheck
	r.HandleFunc("/health", healthCheckHandler)

	apiCreate := r.PathPrefix("/api/v1").Subrouter()

	// Authentication endpoints
	apiCreate.HandleFunc("/auth/login", LoginHandler).Methods("POST", "OPTIONS")
	apiCreate.HandleFunc("/auth/signup", SignupHandler).Methods("POST", "OPTIONS")
	apiCreate.Handle("/auth/me", Middleware(http.HandlerFunc(GetCurrentUserHandler))).Methods("GET", "OPTIONS")

	// Company endpoints
	apiCreate.HandleFunc("/companies", GetCompaniesHandler).Methods("GET", "OPTIONS")
	apiCreate.HandleFunc("/companies", CreateCompanyHandler).Methods("POST", "OPTIONS")
	apiCreate.Handle("/companies/join", Middleware(http.HandlerFunc(JoinCompanyHandler))).Methods("POST", "OPTIONS")
	apiCreate.HandleFunc("/companies/{companyId}/reports", GetCompanyBugReportsHandler).Methods("GET", "OPTIONS")

	// Bug Reports endpoints
	apiCreate.HandleFunc("/bug-reports", GetBugReportsHandler).Methods("GET", "OPTIONS")
	apiCreate.HandleFunc("/bug-reports", CreateBugReportHandler).Methods("POST", "OPTIONS")
	apiCreate.Handle("/bug-reports/{id}", Middleware(http.HandlerFunc(GetBugReportHandler))).Methods("GET", "OPTIONS")
	apiCreate.Handle("/bug-reports/{id}/status", Middleware(http.HandlerFunc(UpdateBugReportStatusHandler))).Methods("PUT", "OPTIONS")
	apiCreate.Handle("/bug-reports/{id}/comments", Middleware(http.HandlerFunc(AddCommentHandler))).Methods("POST", "OPTIONS")

	// User-specific endpoints
	apiCreate.Handle("/users/{userId}/reports", Middleware(http.HandlerFunc(GetUserBugReportsHandler))).Methods("GET", "OPTIONS")

	return r
}

func (a *App) Initialize() error {
	// Convert config.Config to databases.Config
	dbConfig := &databases.Config{
		URL:          a.Config.URL,
		DatabaseName: a.Config.DatabaseName,
		BaseURL:      a.Config.BaseURL,
		Port:         a.Config.Port,
	}

	client, err := databases.NewClient(dbConfig)
	if err != nil {
		// if we fail to create a new database client, the kill the pod
		zap.S().With(err).Error("failed to create new client")
		return err
	}

	a.dbHelper = databases.NewDatabase(dbConfig, client)
	// Client is already connected from NewClient, no need to connect again
	zap.S().Info("BugBridge API has connected to the database")

	// initialize api router
	a.initializeRoutes()
	return nil

}

func (a *App) initializeRoutes() {
	a.Router = a.New()
}

// GetDBHelper returns the database helper
func (a *App) GetDBHelper() databases.DatabaseHelper {
	return a.dbHelper
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(models.HealthCheckResponse{
		Alive: true,
	})
	_, _ = io.WriteString(w, string(b))
}
