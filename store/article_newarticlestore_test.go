package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewArticleStore(t *testing.T) {
	// Define test scenarios
	scenarios := []struct {
		name     string
		db       *gorm.DB
		expected *gorm.DB
	}{
		{name: "Successful creation of a new ArticleStore", db: &gorm.DB{}, expected: &gorm.DB{}},
		{name: "Passing a nil gorm.DB instance", db: nil, expected: nil},
		{name: "Passing an uninitialized gorm.DB instance", db: &gorm.DB{}, expected: &gorm.DB{}},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// Act
			articleStore := NewArticleStore(s.db)

			// Assert
			assert.NotNil(t, articleStore, "Expected ArticleStore instance to be not nil")
			assert.Equal(t, s.expected, articleStore.db, "Expected gorm.DB instance to be equal to the provided one")

			if s.db == nil {
				t.Log("The gorm.DB instance is nil")
			} else if s.db == &gorm.DB{} {
				t.Log("The gorm.DB instance is uninitialized")
			} else {
				t.Log("The gorm.DB instance is initialized")
			}
		})
	}
}

// Function NewArticleStore
func NewArticleStore(db *gorm.DB) *ArticleStore {
	return &ArticleStore{
		db: db,
	}
}

// Type ArticleStore
type ArticleStore struct {
	db *gorm.DB
}
