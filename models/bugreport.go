package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BugReport represents a bug report in the system
type BugReport struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title            string             `json:"title" bson:"title" validate:"required"`
	Description      string             `json:"description" bson:"description" validate:"required"`
	Severity         string             `json:"severity" bson:"severity" validate:"required,oneof=low medium high critical"`
	Status           string             `json:"status" bson:"status" validate:"required,oneof=pending under_review accepted rejected resolved"`
	StepsToReproduce string             `json:"stepsToReproduce" bson:"stepsToReproduce" validate:"required"`
	IsAnonymous      bool               `json:"isAnonymous" bson:"isAnonymous"`
	ReporterID       primitive.ObjectID `json:"reporterId" bson:"reporterId"`
	CompanyID        primitive.ObjectID `json:"companyId" bson:"companyId"`
	CompanyName      string             `json:"companyName" bson:"companyName"`
	SubmittedAt      time.Time          `json:"submittedAt" bson:"submittedAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CreateBugReportRequest represents the request to create a bug report
type CreateBugReportRequest struct {
	Title            string `json:"title" validate:"required"`
	Description      string `json:"description" validate:"required"`
	Severity         string `json:"severity" validate:"required,oneof=low medium high critical"`
	StepsToReproduce string `json:"stepsToReproduce" validate:"required"`
	IsAnonymous      bool   `json:"isAnonymous"`
	CompanyID        string `json:"companyId" validate:"required"`
	CompanyName      string `json:"companyName"`
}
