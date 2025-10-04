package models

type Report struct {
	ID         string   `json:"_id"      bson:"_id"`      //report id
	AuthorID   string   `json:"author"   bson:"Author"`   //ID of author
	Desc       string   `json:"desc"     bson:"Desc"`     //description of report
	Severity   int      `json:"severity" bson:"Severity"` //severity of report
	Resolved   bool     `json:"resolved" bson:"Resolved"` //array of IDs of solutions
	CommentIDs []string `json:"comments" bson:"Comments"` //Array of comments under report
}
