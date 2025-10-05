package models

type Template struct {
	Title          string `json:"title"     bson:"title"`
	Desc           string `json:"desc"      bson:"desc"`
	Steps          string `json:"steps"     bson:"steps"`
	Behaviour      string `json:"behaviour" bson:"behaviour"`
	AdditionalInfo string `json:"addInfo"   bson:"additionalInfo"`
}
