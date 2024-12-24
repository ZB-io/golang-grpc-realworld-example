
// ********RoostGPT********
/*

roost_feedback [12/24/2024, 12:19:49 PM]:- Add more comments to the test

roost_feedback [12/24/2024, 12:33:56 PM]:- Add more comments to the test
*/

// ********RoostGPT********

package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// TestNewArticleStore checks the functionality of new ArticleStore creation
func TestNewArticleStore(t *testing.T) {
	var cases = []struct {
		name     string
		db       *gorm.DB
		expected *ArticleStore
	}{
		{
			name:     "Valid Database Instance Passed",
			expected: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name:     "Non-Initialized Database Instance Passed",
			expected: &ArticleStore{},
		},
		{
			name:     "Null Database Instance Passed",
			expected: &ArticleStore{},
		},
	}

	// Iterating through different test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			// Initiate a new mock database if db is not nil
			if tc.db != nil {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Unexpected error when opening a stub database connection: %s", err)
				}

				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					t.Fatalf("Unexpected error when opening gorm database connection: %s", err)
				}

				tc.db = gormDB
				defer tc.db.Close() // Ensuring the db is closed post test
			}

			// Testing NewArticleStore function
			got := NewArticleStore(tc.db)

			t.Log("Executed NewArticleStore function with db: ", tc.db)

			// Assert equality for expected and returned values from the function
			assert.Equal(t, tc.expected, got)

			// Verify if all mock expectations are met in test cases with database specified
			if tc.db != nil {
				if assert.NoError(t, mock.ExpectationsWereMet()) {
					t.Log("All mock expectations were fulfilled.")
				} else {
					t.Log("Some mock expectations were not met.")
				}
			}
		})
	}
}
