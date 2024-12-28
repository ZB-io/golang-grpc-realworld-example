package store

import (
	"errors"
	"fmt"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type ExpectedBegin struct {
	commonExpectation
	delay time.Duration
}

type ExpectedCommit struct {
	commonExpectation
}

type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

type ExpectedRollback struct {
	commonExpectation
}

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
func TestArticleStoreUpdate(t *testing.T) {
	tests := []struct {
		name           string
		article        model.Article
		updatedArticle model.Article
		setupMock      func(sqlmock.Sqlmock)
		expectError    bool
		errorMessage   string
	}{
		{
			name: "Successful Update of an Existing Article",
			article: model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Original Title",
				Description:    "Original Description",
				Body:           "Original Body",
				UserID:         1,
				FavoritesCount: 10,
			},
			updatedArticle: model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Updated Title",
				Description:    "Updated Description",
				Body:           "Updated Body",
				UserID:         1,
				FavoritesCount: 10,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE .*").WithArgs(
					"Updated Title",
					"Updated Description",
					"Updated Body",
					sqlmock.AnyArg(),
					1,
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError:  false,
			errorMessage: "",
		},
		{
			name: "Fail to Update Non-Existent Article",
			article: model.Article{
				Model:       gorm.Model{ID: 999},
				Title:       "Title",
				Description: "Description",
				Body:        "Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE .*").WithArgs(
					"Title",
					"Description",
					"Body",
					sqlmock.AnyArg(),
					999,
				).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			expectError:  true,
			errorMessage: "record not found",
		},
		{
			name: "Error During Database Transaction",
			article: model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Error Title",
				Description: "Error Description",
				Body:        "Error Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE .*").WithArgs(
					"Error Title",
					"Error Description",
					"Error Body",
					sqlmock.AnyArg(),
					1,
				).WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectError:  true,
			errorMessage: "database error",
		},
		{
			name: "Validate Constraints on Article Fields",
			article: model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "",
				Description: "",
				Body:        "",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError:  true,
			errorMessage: "violating non-null constraints",
		},
		{
			name: "Successful Update Without Changing FavoritesCount",
			article: model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Old Title",
				Description:    "Old Description",
				Body:           "Old Body",
				UserID:         1,
				FavoritesCount: 15,
			},
			updatedArticle: model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "New Title",
				Description:    "Old Description",
				Body:           "Old Body",
				UserID:         1,
				FavoritesCount: 15,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE .*").WithArgs(
					"New Title",
					"Old Description",
					"Old Body",
					sqlmock.AnyArg(),
					1,
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError:  false,
			errorMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error while opening sqlmock database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("error while opening gorm database connection: %v", err)
			}

			store := &ArticleStore{db: gormDB}

			tt.setupMock(mock)

			err = store.Update(&tt.updatedArticle)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMessage {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}

			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unmet expectations: %s", err)
			}
		})
	}
}
