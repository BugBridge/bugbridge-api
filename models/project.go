package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	ID        primitive.ObjectID `json:"_id"       bson:"_id"`      // ID of a project
	Name      string             `json:"name"      bson:"name"`     // Name of project
	Des       string             `json:"des"       bson:"des"`      // Public description of the project
	OwnerID   string             `json:"ownerId"   bson:"ownerId"`  // Owner ID of the project
	AdminsIDs []string           `json:"adminIds"  bson:"adminIds"` // Array of IDs for users with admin privilages
	Template  TemplateData       `json:"template"  bson:"template"` // Template that bug reports should be submitted
}

// Data structure of the json object received in POST to create project
type ProjectDetails struct {
	Name     string       `json:"name"      validate:"required,min=3,max=50"`
	Des      string       `json:"des"       validate:"required,max=500"`
	OwnerID  string       `json:"ownerId"   validate:"required"`
	Template TemplateData `json:"template"  validate:"required"`
}

type TemplateData struct {
	Title          string `json:"title"     bson:"title"     validate:"required,min=3,max=50"`
	Des            string `json:"des"       bson:"des"       validate:"required,max=500"`
	Steps          string `json:"steps"     bson:"steps"     validate:"required,max=1000"`
	Behaviour      string `json:"behaviour" bson:"behaviour" validate:"required,max=1000"`
	AdditionalInfo string `json:"addInfo"   bson:"addInfo"   validate:"max=1000"`
}

// Data structure of the json object received in PATCH to update project
type ProjectUpdateDetails struct {
	Name     string             `json:"name"      validate:"max=50"`
	Des      string             `json:"des"       validate:"max=500"`
	Template TemplateUpdateData `json:"template"`
}

type TemplateUpdateData struct {
	Title          string `json:"title"     validate:"max=50"`
	Des            string `json:"des"       validate:"max=500"`
	Steps          string `json:"steps"     validate:"max=1000"`
	Behaviour      string `json:"behaviour" validate:"max=1000"`
	AdditionalInfo string `json:"addInfo"   validate:"max=1000"`
}
