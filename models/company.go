package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Company represents a company in the system
type Company struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name" validate:"required"`
	Industry          string             `json:"industry" bson:"industry"`
	Description       string             `json:"description" bson:"description"`
	Website           string             `json:"website" bson:"website"`
	OwnerID           primitive.ObjectID `json:"ownerId" bson:"ownerId"`
	AcceptingReports  bool               `json:"acceptingReports" bson:"acceptingReports"`
	BugReportTemplate interface{}        `json:"bugReportTemplate" bson:"bugReportTemplate"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CreateCompanyRequest represents the request to create a company
type CreateCompanyRequest struct {
	Name              string `json:"name" validate:"required"`
	Industry          string `json:"industry"`
	Description       string `json:"description"`
	Website           string `json:"website"`
	BugReportTemplate string `json:"bugReportTemplate"`
}
