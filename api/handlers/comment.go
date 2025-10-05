package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/BugBridge/bugbridge-api/config"
	"github.com/BugBridge/bugbridge-api/databases"
	"github.com/BugBridge/bugbridge-api/models"
)

type Comment struct {
	DB databases.CommentDatabase
}

// TODO: add delete and update functionality

// CommentByObjectIDHandler returns a comment by a given ID
func (comment Comment) CommentByObjectIDHandler(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["comment_id"]

	cID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := comment.DB.FindOne(context.Background(), bson.M{"_id": cID})
	if err != nil {
		config.ErrorStatus("failed to get comment by ID", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// CommentsByReportIDHandler returns a comment by a given ID
func (comment Comment) CommentsByReportIDHandler(w http.ResponseWriter, r *http.Request) {
	reportID := mux.Vars(r)["report_id"]

	dbResp, err := comment.DB.Find(context.Background(), bson.M{"reportId": reportID})
	if err != nil {
		config.ErrorStatus("failed to get comment by ID", http.StatusNotFound, w, err)
		return
	}

	if len(dbResp) == 0 {
		dbResp = []models.Comment{}
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Create a new comment
func (comment Comment) NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var details models.CommentDetails // Json data will represent the report details model
	defer cancel()

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&details); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	// TODO: add validation comment attributes

	newComment := models.Comment{
		ID:       primitive.NewObjectID(),
		AuthorID: details.AuthorID,
		ReportID: details.ReportID,
		Content:  details.Content,
	}

	result, err := comment.DB.InsertOne(ctx, newComment)
	if err != nil {
		config.ErrorStatus("failed to insert comment", http.StatusBadRequest, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]any{"result": result},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

// UpdateCommentHandler updates the content for an existing comment
func (comment Comment) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var newDetails models.CommentUpdateDetails
	defer cancel()

	commentID := mux.Vars(r)["comment_id"]

	cID, err := primitive.ObjectIDFromHex(commentID)

	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&newDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&newDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	dbResp, err := comment.DB.UpdateOne(
		ctx,
		bson.M{"_id": cID},
		bson.M{"$set": bson.M{"content": newDetails.Content}},
	)

	if err != nil {
		config.ErrorStatus("the comment could not be updated", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (comment Comment) DeleteCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commentID := mux.Vars(r)["comment_id"]

	uID, err := primitive.ObjectIDFromHex(commentID)

	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := comment.DB.DeleteOne(ctx, bson.M{"_id": uID})
	if err != nil {
		config.ErrorStatus("failed to delete comment", http.StatusNotFound, w, err)
		return
	}

	if dbResp.Dr.DeletedCount == 0 {
		config.ErrorStatus("Comment not found", http.StatusNotFound, w, nil)
		return
	}

	b, err := json.Marshal(
		models.DataResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]any{"result": dbResp},
		},
	)

	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
