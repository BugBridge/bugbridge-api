package databases

import (
	"context"

	"github.com/BugBridge/bugbridge-api/models"
)

const reportDBO = "reports"

type ReportDatabase interface {
	FindOne(ctx context.Context, filter interface{}) (*models.Report, error)
	Find(ctx context.Context, filter interface{}) ([]models.Report, error)
	InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter, document interface{}) (*mongoUpdateResult, error)
	DeleteOne(document any, condition bool) (*mongoDeleteOneResult, error)
}

type reportDatabase struct {
	db DatabaseHelper
}

func NewReportDatabase(db DatabaseHelper) ReportDatabase {
	return &reportDatabase{
		db: db,
	}
}

func (u *reportDatabase) FindOne(ctx context.Context, filter interface{}) (*models.Report, error) {
	report := &models.Report{}
	err := u.db.Collection(reportDBO).FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (u *reportDatabase) Find(ctx context.Context, filter interface{}) ([]models.Report, error) {
	var reports []models.Report
	err := u.db.Collection(reportDBO).Find(ctx, filter).Decode(&reports)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (u *reportDatabase) InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error) {
	result, err := u.db.Collection(reportDBO).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *reportDatabase) UpdateOne(ctx context.Context, filter, update interface{}) (*mongoUpdateResult, error) {
	result, err := u.db.Collection(reportDBO).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *reportDatabase) DeleteOne(document any, condition bool) (*mongoDeleteOneResult, error) {
	result, err := u.db.Collection(reportDBO).DeleteOne(document, conditions)
	if err != nil {
		return nil, err
	}
	return &result, nil
}