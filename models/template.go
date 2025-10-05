package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	Title				string	`json:"title" bson:"Title"`
	Desc				string	`json:"desc" bson:"Desc"`
	Steps				string	`json:"steps" bson:"Steps"`
	Behaviour			string	`json:"behaviour" bson:"Behaviour"`
	AdditionalInfo		string	`json:"addInfo" bson:"AdditionalInfo"`
}