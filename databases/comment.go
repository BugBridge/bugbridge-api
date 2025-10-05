package databases

import (
	"context"

	"github.com/BugBridge/bugbridge-api/models"
)

const commentDBO = "comments"

type CommentDatabase interface {
	FindOne(ctx context.Context, filter interface{}) (*models.Comment, error)
	Find(ctx context.Context, filter interface{}) ([]models.Comment, error)
	InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter, document interface{}) (*mongoUpdateResult, error)
	// DeleteOne(ctx context.Context, filter interface{}) (*mongoDeleteOneResult, error)
}

type commentDatabase struct {
	db DatabaseHelper
}

func NewCommentDatabase(db DatabaseHelper) CommentDatabase {
	return &commentDatabase{
		db: db,
	}
}

func (u *commentDatabase) FindOne(ctx context.Context, filter interface{}) (*models.Comment, error) {
	comment := &models.Comment{}
	err := u.db.Collection(commentDBO).FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (u *commentDatabase) Find(ctx context.Context, filter interface{}) ([]models.Comment, error) {
	var comments []models.Comment
	err := u.db.Collection(commentDBO).Find(ctx, filter).Decode(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (u *commentDatabase) InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error) {
	result, err := u.db.Collection(commentDBO).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *commentDatabase) UpdateOne(ctx context.Context, filter, update interface{}) (*mongoUpdateResult, error) {
	result, err := u.db.Collection(commentDBO).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// func (u *commentDatabase) DeleteOne(ctx context.Context, filter interface{}) (*mongoDeleteOneResult, error) {
// 	result, err := u.db.Collection(commentDBO).DeleteOne(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }
