package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.db != nil {
				var mock sqlmock.Sqlmock
				var err error
				var db *sql.DB
				db, mock, err = sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening gorm database", err)
				}

				tc.db = gormDB
				defer tc.db.Close()

			}

			got := NewArticleStore(tc.db)
			t.Log("Executed function NewArticleStore with db - ", tc.db)

			assert.Equal(t, tc.expected, got)

			if tc.db != nil {
				if assert.NoError(t, mock.ExpectationsWereMet()) {
					t.Log("All mock expectations met")
				} else {
					t.Log("Mock expectations failed")
				}
			}
		})
	}
}
