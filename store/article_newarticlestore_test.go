package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewArticleStore(t *testing.T) {
   // Table-driven tests for robustness
   scenarios := []struct {
		name         string
		setupMock    func() *gorm.DB
		wantNil      bool
	}{
		{
			name: "Test if function returns correct ArticleStore.",
			setupMock: func() *gorm.DB {
				sqlDB, _, _ := sqlmock.New() // Create a mock DB
				db, _ := gorm.Open("postgres", sqlDB)
				return db
			},
			wantNil: false,
		},
		{
			name: "Check if function can handle nil DB input.",
			setupMock: func() *gorm.DB {
				return nil
			},
			wantNil: true,
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// TODO: Arrange for each scenario
			mockDB := s.setupMock()

			// TODO: Act: Invoke the NewArticleStore function in each scenario.
			res := NewArticleStore(mockDB)

			// Tests the scenarios where DB is expected to be nil
			if s.wantNil {
				// TODO: Assert: Check if returned ArticleStore's DB is nil
				assert.Nil(t, res.db, "The ArticleStore DB is not nil when it should be nil.")
			} else {
				// TODO: Assert: Check if returned ArticleStore's DB is the same as the mock DB.
				assert.Equal(t, mockDB, res.db, "The returned ArticleStore's DB is not the same as the mock DB.")
			}
		})
	}
}

// TODO: Implement the scenario "Test if function can handle multiple consecutive calls" using similar structure as the previous scenarios.

// TODO: Implement the scenario "Test if function can handle a DB with data" using similar structure as the previous scenarios.
