package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestArticleStoreDelete function tests the Delete method of the ArticleStore
func TestArticleStoreDelete(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		existing      bool
		invalidID     bool
		expectedError bool
	}{
		{"Successful Deletion of an Article", true, false, false},
		{"Deletion of a Non-Existent Article", false, false, true},
		{"Deletion of an Article with Invalid ID", true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock DB
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			store := &ArticleStore{db: gormDB}

			// Arrange
			article := &model.Article{Title: "Test Article"}

			// If the test requires an existing article, add it to the mock DB
			if tt.existing {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs(article.Title).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				// If the test requires an invalid ID, set the ID to a negative value
				if tt.invalidID {
					article.Model.ID = 0
				}
			}

			// Act
			err = store.Delete(article)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
