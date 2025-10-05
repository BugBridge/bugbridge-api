package databases

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client represents a MongoDB client
type Client struct {
	client *mongo.Client
}

// NewClient creates a new MongoDB client
func NewClient(config *Config) (*Client, error) {
	clientOptions := options.Client().ApplyURI(config.URL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// Connect connects to the MongoDB database
func (c *Client) Connect() error {
	return c.client.Ping(context.Background(), nil)
}

// DatabaseHelper interface for database operations
type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	GetCollection(name string) *mongo.Collection
}

// MongoDatabase implements DatabaseHelper for MongoDB
type MongoDatabase struct {
	db *mongo.Database
}

// NewDatabase creates a new MongoDatabase instance
func NewDatabase(config *Config, client *Client) *MongoDatabase {
	return &MongoDatabase{db: client.client.Database(config.DatabaseName)}
}

// Collection returns a CollectionHelper for the specified collection name
func (md *MongoDatabase) Collection(name string) CollectionHelper {
	return &MongoCollection{collection: md.db.Collection(name)}
}

// GetCollection returns a MongoDB collection directly
func (md *MongoDatabase) GetCollection(name string) *mongo.Collection {
	return md.db.Collection(name)
}

// CollectionHelper interface for collection operations
type CollectionHelper interface {
	FindOne(ctx context.Context, filter interface{}, result interface{}) error
	Find(ctx context.Context, filter interface{}, results interface{}) error
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}) (int64, error)
}

// MongoCollection implements CollectionHelper for MongoDB collections
type MongoCollection struct {
	collection *mongo.Collection
}

// FindOne finds a single document in the collection
func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}, result interface{}) error {
	return mc.collection.FindOne(ctx, filter).Decode(result)
}

// Find finds multiple documents in the collection
func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, results interface{}) error {
	cursor, err := mc.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

// InsertOne inserts a single document into the collection
func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return mc.collection.InsertOne(ctx, document)
}

// UpdateOne updates a single document in the collection
func (mc *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateOne(ctx, filter, update)
}

// DeleteOne deletes a single document from the collection
func (mc *MongoCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return mc.collection.DeleteOne(ctx, filter)
}

// CountDocuments counts the number of documents in the collection
func (mc *MongoCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return mc.collection.CountDocuments(ctx, filter)
}

// Config represents database configuration
type Config struct {
	URL          string
	DatabaseName string
	BaseURL      string
	Port         string
}

// MockCollectionHelper represents a mock collection helper for testing
type MockCollectionHelper struct{}

// GetCompanies returns mock companies data
func (m *MockCollectionHelper) GetCompanies() ([]interface{}, error) {
	companies := []interface{}{
		map[string]interface{}{
			"id":              "68e23b3d997deadd848a490b",
			"name":            "TechCorp Security",
			"description":     "Leading technology company focused on secure software development",
			"industry":        "Technology",
			"website":         "https://techcorp.com",
			"ownerId":         "68e23b3d997deadd848a490c",
			"bugReportsCount": 12,
			"isActive":        true,
			"createdAt":       time.Date(2025, time.October, 5, 9, 32, 45, 735000000, time.UTC),
			"updatedAt":       time.Date(2025, time.October, 5, 9, 32, 45, 735000000, time.UTC),
		},
		map[string]interface{}{
			"id":              "68e23b3d997deadd848a490d",
			"name":            "StartupXYZ",
			"description":     "Innovative startup building the next generation of mobile applications",
			"industry":        "Technology",
			"website":         "https://startupxyz.com",
			"ownerId":         "68e23b3d997deadd848a490e",
			"bugReportsCount": 8,
			"isActive":        true,
			"createdAt":       time.Date(2025, time.October, 5, 9, 32, 45, 735000000, time.UTC),
			"updatedAt":       time.Date(2025, time.October, 5, 9, 32, 45, 735000000, time.UTC),
		},
		map[string]interface{}{
			"id":              "68e23b3d997deadd848a490f",
			"name":            "FinanceFlow",
			"description":     "Financial services platform for small and medium businesses",
			"industry":        "Finance",
			"website":         "https://financeflow.com",
			"ownerId":         "68e23b3d997deadd848a4910",
			"bugReportsCount": 5,
			"isActive":        true,
			"createdAt":       time.Date(2025, time.October, 5, 9, 32, 45, 735000000, time.UTC),
			"updatedAt":       time.Date(2025, time.October, 5, 9, 32, 45, 735000000, time.UTC),
		},
	}
	return companies, nil
}

// GetReports returns mock bug reports data
func (m *MockCollectionHelper) GetReports() ([]interface{}, error) {
	reports := []interface{}{
		map[string]interface{}{
			"id":          "68e1f040242f027998122205",
			"title":       "this is a very bad bug",
			"description": "",
			"severity":    "medium",
			"status":      "",
			"companyId":   "000000000000000000000000",
			"submittedAt": time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			"updatedAt":   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		map[string]interface{}{
			"id":           "68e23b3d997deadd848a4911",
			"title":        "Authentication Bypass Found",
			"description":  "Discovered a critical authentication bypass that allows unauthorized access to admin panel.",
			"severity":     "high",
			"status":       "pending",
			"companyId":    "68e23b3d997deadd848a490b",
			"companyName":  "TechCorp Security",
			"reporterId":   "68e23b3d997deadd848a4912",
			"reporterName": "Security Researcher",
			"submittedAt":  time.Date(2025, time.October, 5, 2, 32, 45, 805000000, time.FixedZone("", -7*60*60)),
			"updatedAt":    time.Date(2025, time.October, 5, 2, 32, 45, 805000000, time.FixedZone("", -7*60*60)),
			"attachments":  []string{"screenshot1.png", "proof_of_concept.txt"},
		},
		map[string]interface{}{
			"id":           "68e23b3d997deadd848a4913",
			"title":        "SQL Injection in Search",
			"description":  "Found SQL injection vulnerability in the search functionality that can be exploited to extract user data.",
			"severity":     "medium",
			"status":       "under_review",
			"companyId":    "68e23b3d997deadd848a490d",
			"companyName":  "StartupXYZ",
			"reporterId":   "68e23b3d997deadd848a4914",
			"reporterName": "Bug Hunter",
			"submittedAt":  time.Date(2025, time.October, 5, 2, 32, 45, 805000000, time.FixedZone("", -7*60*60)),
			"updatedAt":    time.Date(2025, time.October, 5, 2, 32, 45, 805000000, time.FixedZone("", -7*60*60)),
			"attachments":  []string{"exploit_script.sql"},
		},
		map[string]interface{}{
			"id":           "68e23b3d997deadd848a4915",
			"title":        "XSS in Contact Form",
			"description":  "Stored XSS vulnerability in the contact form that persists in the database.",
			"severity":     "low",
			"status":       "accepted",
			"companyId":    "68e23b3d997deadd848a490f",
			"companyName":  "FinanceFlow",
			"reporterId":   "68e23b3d997deadd848a4916",
			"reporterName": "Security Expert",
			"submittedAt":  time.Date(2025, time.October, 5, 2, 32, 45, 805000000, time.FixedZone("", -7*60*60)),
			"updatedAt":    time.Date(2025, time.October, 5, 2, 32, 45, 805000000, time.FixedZone("", -7*60*60)),
			"attachments":  []string{"xss_payload.html"},
		},
	}
	return reports, nil
}

// GetDashboardStats returns mock dashboard stats
func (m *MockCollectionHelper) GetDashboardStats() (interface{}, error) {
	stats := map[string]interface{}{
		"totalReports":    4,
		"reportsTrend":    0,
		"totalCompanies":  3,
		"companiesTrend":  0,
		"pendingReports":  1,
		"resolvedReports": 1,
	}
	return stats, nil
}
