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
	Name     string       `json:"name"`
	Des      string       `json:"des"`
	OwnerID  string       `json:"adminId"`
	Template TemplateData `json:"template"`
}

type TemplateData struct {
	Title          string `json:"title"     bson:"title"`
	Des            string `json:"des"       bson:"des"`
	Steps          string `json:"steps"     bson:"steps"`
	Behaviour      string `json:"behaviour" bson:"behaviour"`
	AdditionalInfo string `json:"addInfo"   bson:"additionalInfo"`
}
