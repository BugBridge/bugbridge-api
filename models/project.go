package models

type Project struct {
	ID         string   `json:"_id"     bson:"_id"`          //Id of a project
	Name       string   `json:"name"    bson:"Name"`         //Name of project
	Desc       string   `json:"desc"    bson:"Desc"`         // Public description of the project
	Template   string   `json:"template" bson:"Template"`    // Template that bug reports should be submitted
	AdminsIDs  []string `json:"adminIds"  bson:"AdminsIDs"`  //Array of IDs for users as admins
	MembersIDs []string `json:"memberIds" bson:"MembersIDs"` //Array of IDs of users
	ReportIDs  []string `json:"reportIds" bson:"ReportsIDs"` //Array of IDs of reports
}
