package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/BugBridge/bugbridge-api/api"
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

	r := mux.NewRouter()

	// create database handlers like this
	users := User{DB: databases.NewUserDatabase(a.dbHelper)}
	projects := Project{DB: databases.NewProjectDatabase(a.dbHelper)}
	reports := Report{DB: databases.NewReportDatabase(a.dbHelper)}
	comments := Comment{DB: databases.NewCommentDatabase(a.dbHelper)}

	// healthcheck
	r.HandleFunc("/health", healthCheckHandler)

	apiCreate := r.PathPrefix("/api").Subrouter()

	// API endpoints
	apiCreate.Handle("/user/{user_id}", api.Middleware(http.HandlerFunc(users.UserByObjectIDHandler))).Methods("GET")
	apiCreate.Handle("/user/create", api.Middleware(http.HandlerFunc(users.NewUserHandler))).Methods("POST")

	apiCreate.Handle("/report/{report_id}", api.Middleware(http.HandlerFunc(reports.ReportByObjectIDHandler))).Methods("GET")
	apiCreate.Handle("/report/create", api.Middleware(http.HandlerFunc(reports.NewReportHandler))).Methods("POST")

	apiCreate.Handle("/projects/{project_id}", api.Middleware(http.HandlerFunc(projects.ProjectByObjectIDHandler))).Methods("GET")
	apiCreate.Handle("/projects/create", api.Middleware(http.HandlerFunc(projects.ProjectByObjectIDHandler))).Methods("POST")

	apiCreate.Handle("/comment/{comment_id}", api.Middleware(http.HandlerFunc(comments.CommentByObjectIDHandler))).Methods("GET")
	apiCreate.Handle("/comment/report/{report_id}", api.Middleware(http.HandlerFunc(comments.CommentsByReportIDHandler))).Methods("GET")
	apiCreate.Handle("/comment/create", api.Middleware(http.HandlerFunc(comments.NewCommentHandler))).Methods("POST")

	return r
}

func (a *App) Initialize() error {
	client, err := databases.NewClient(&a.Config)
	if err != nil {
		// if we fail to create a new database client, the kill the pod
		zap.S().With(err).Error("failed to create new client")
		return err
	}

	a.dbHelper = databases.NewDatabase(&a.Config, client)
	err = client.Connect()
	if err != nil {
		// if we fail to connect to the database, the kill the pod
		zap.S().With(err).Error("failed to connect to database")
		return err
	}
	zap.S().Info("DeviceBookingAPI has connected to the database")

	// initialize api router
	a.initializeRoutes()
	return nil

}

func (a *App) initializeRoutes() {
	a.Router = a.New()
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(models.HealthCheckResponse{
		Alive: true,
	})
	_, _ = io.WriteString(w, string(b))
}
