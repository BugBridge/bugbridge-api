package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID        primitive.ObjectID `json:"_id"        bson:"_id"`       // report id
	AuthorID  string             `json:"authorId"     bson:"author"`  // ID of author
	ProjectID string             `json:"projectId"  bson:"projectId"` // Project ID report is submitted to
	Title     string             `json:"title"      bson:"title"`     // title of id
	Des       string             `json:"des"        bson:"des"`       // description of report
	Severity  int                `json:"severity"   bson:"severity"`  // severity of report
	Resolved  bool               `json:"resolved"   bson:"resolved"`  // array of IDs of solutions
}

// Data structure of the json object received in POST to create report
type ReportDetails struct {
	AuthorID  string `json:"authorId"  validate:"required"`          // ID of author
	ProjectID string `json:"projectId" validate:"required"`          // Project ID report is submitted to
	Title     string `json:"title"     validate:"required,max=50"`   // title of the report
	Des       string `json:"des"       validate:"required,max=1000"` // description of report
}

// Data structure of the json object received in PATCH to update report
type ReportUpdateDetails struct {
	AuthorID  string `json:"authorId"`                      // ID of author
	ProjectID string `json:"projectId"`                     // Project ID report is submitted to
	Title     string `json:"title"     validate:"max=50"`   // title of the report
	Des       string `json:"des"       validate:"max=1000"` // description of report
}
