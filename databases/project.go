package databases

import (
	"context"

	"github.com/BugBridge/bugbridge-api/models"
)

const projectDBO = "projects"

type ProjectDatabase interface {
	FindOne(ctx context.Context, filter interface{}) (*models.Project, error)
	Find(ctx context.Context, filter interface{}) ([]models.Project, error)
	InsertOne(ctx context.Context, filter interface{}) (*mongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}) (*mongoUpdateOneResult, error)
	DeleteOne(document any, condition bool) (*mongoDeleteOneResult, error)
}

type projectDatabase struct {
	db DatabaseHelper
}

func NewProjectDatabase(db DatabaseHelper) ProjectDatabase {
	return &projectDatabase{
		db: db,
	}
}

func (u *projectDatabase) FindOne(ctx context.Context, filter interface{}) (*models.Project, error) {
	project := &models.Project{}
	err := u.db.Collection(projectDBO).FindOne(ctx, filter).Decode(&project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (u *projectDatabase) Find(ctx context.Context, filter interface{}) ([]models.Project, error) {
	var projects []models.Project
	err := u.db.Collection(projectDBO).Find(ctx, filter).Decode(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (u *projectDatabase) InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error) {
	result, err := u.db.Collection(projectDBO).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *projectDatabase) UpdateOne(ctx context.Context, filter, update interface{}) (*mongoUpdateResult, error) {
	result, err := u.db.Collection(projectDBO).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *projectDatabase) DeleteOne(document any, condition bool) (*mongoDeleteOneResult, error) {
	result, err := u.db.Collection(projectDBO).DeleteOne(document, condition)
	if err != nil {
		return nil, err
	}
	return &result, nil
}