package store

import (
	"errors"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
)






func TestGetComments(t *testing.T) {
	tests := []struct {
		name          string
		article       *model.Article
		mockSetupFunc func(sqlmock.Sqlmock)
		expectedErr   error
		expectedCount int
	}{
		{
			name: "Successfully Retrieve Comments",
			article: &model.Article{
				ID: 1,
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \?\)`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "content"}).
						AddRow(1, 1, "Great article").
						AddRow(2, 1, "Very informative"))
			},
			expectedErr:   nil,
			expectedCount: 2,
		},
		{
			name: "Return Empty List When No Comments Exist",
			article: &model.Article{
				ID: 2,
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \?\)`).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "content"}))
			},
			expectedErr:   nil,
			expectedCount: 0,
		},
		{
			name:    "Handle Database Errors Gracefully",
			article: &model.Article{ID: 3},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \?\)`).WithArgs(3).
					WillReturnError(errors.New("DB error"))
			},
			expectedErr:   errors.New("DB error"),
			expectedCount: 0,
		},
		{
			name:          "Correct Handling of Nil Article Input",
			article:       nil,
			mockSetupFunc: nil,
			expectedErr:   errors.New("article is nil"),
			expectedCount: 0,
		},
		{
			name: "Ensure Preloading Author Data with Comments",
			article: &model.Article{
				ID: 4,
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \?\)`).
					WithArgs(4).
					WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "content"}).
						AddRow(3, 4, "Nice read"))
				mock.ExpectQuery(`SELECT \* FROM "authors" WHERE \(id = \?\)`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "John Doe"))
			},
			expectedErr:   nil,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer db.Close()

			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(mock)
			}

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			articleStore := &ArticleStore{db: gormDB}

			comments, err := articleStore.GetComments(tt.article)

			if err != nil && tt.expectedErr == nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectedErr != nil {
				t.Errorf("expected error but got nil")
			}
			if tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error() {
				t.Errorf("expected error '%v' but got '%v'", tt.expectedErr.Error(), err.Error())
			}
			if len(comments) != tt.expectedCount {
				t.Errorf("expected comments count %d, got %d", tt.expectedCount, len(comments))
			}
		})
	}
}

