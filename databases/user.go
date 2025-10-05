package databases

import (
	"context"

	"github.com/BugBridge/bugbridge-api/models"
)

const userDBO = "users"

type UserDatabase interface {
	FindOne(ctx context.Context, filter any) (*models.User, error)
	Find(ctx context.Context, filter any) ([]models.User, error)
	InsertOne(ctx context.Context, filter any) (*mongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter, document any) (*mongoUpdateResult, error)
	DeleteOne(ctx context.Context, filter any) (*mongoDeleteOneResult, error)
}

type userDatabase struct {
	db DatabaseHelper
}

func NewUserDatabase(db DatabaseHelper) UserDatabase {
	return &userDatabase{
		db: db,
	}
}

func (u *userDatabase) FindOne(ctx context.Context, filter any) (*models.User, error) {
	user := &models.User{}
	err := u.db.Collection(userDBO).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userDatabase) Find(ctx context.Context, filter any) ([]models.User, error) {
	var users []models.User
	err := u.db.Collection(userDBO).Find(ctx, filter).Decode(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userDatabase) InsertOne(ctx context.Context, document any) (*mongoInsertOneResult, error) {
	result, err := u.db.Collection(userDBO).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *userDatabase) UpdateOne(ctx context.Context, filter, update any) (*mongoUpdateResult, error) {
	result, err := u.db.Collection(userDBO).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *userDatabase) DeleteOne(ctx context.Context, filter any) (*mongoDeleteOneResult, error) {
	result, err := u.db.Collection(userDBO).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
