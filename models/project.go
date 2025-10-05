package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	ID         primitive.ObjectID `json:"_id"       bson:"_id"`      //Id of a project
	Name       string             `json:"name"      bson:"name"`     //Name of project
	Desc       string             `json:"desc"      bson:"desc"`     // Public description of the project
	Template   TemplateData       `json:"template"  bson:"template"` // Template that bug reports should be submitted
	OwnerID    string             `json:"ownerId"   bson:"ownerId"`
	AdminsIDs  []string           `json:"adminIds"  bson:"adminIds"`  //Array of IDs for users as admins
	MembersIDs []string           `json:"memberIds" bson:"memberIds"` //Array of IDs of users
	ReportIDs  []string           `json:"reportIds" bson:"reportIds"` //Array of IDs of reports
}

// Data structure of the json object received in POST to create project
type ProjectDetails struct {
	Name     string       `json:"name"`
	Desc     string       `json:"desc"`
	Template TemplateData `json:"template"`
	OwnerID  string       `json:"adminID"`
}

type TemplateData struct {
	Title          string `json:"title"     bson:"title"`
	Desc           string `json:"desc"      bson:"desc"`
	Steps          string `json:"steps"     bson:"steps"`
	Behaviour      string `json:"behaviour" bson:"behaviour"`
	AdditionalInfo string `json:"addInfo"   bson:"additionalInfo"`
}
