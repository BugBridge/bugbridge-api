package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
    ID      string      `json:"_id" bson:"_id"`         //Id of a project
    Name    string      `json:"name" bson:"Name"`       //Name of project
    Admins  []string    `json:"admins" bson:"Admins"`   //Array of IDs for users as admins
    Members []string    `json:"members" bson:"Members"` //Array of IDs of users
    Reports []string    `json:"reports" bson:"Reports"` //Array of IDs of reports
}