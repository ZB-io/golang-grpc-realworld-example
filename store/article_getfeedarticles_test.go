package store

import (
	"errors"
	"reflect"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)







func TestGetFeedArticles(t *testing.T) {

	type testCase struct {
		description string
		userIDs     []uint
		limit       int64
		offset      int64
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    []model.Article
		expectErr   bool
	}

	tests := []testCase{
		{
			description: "Normal Operation with Valid Inputs",
			userIDs:     []uint{1, 2},
			limit:       2,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM articles WHERE user_id in \\(\\$1, \\$2\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "user_id"}).
						AddRow(1, "Article 1", 1).
						AddRow(2, "Article 2", 2))
			},
			expected: []model.Article{
				{ID: 1, Title: "Article 1", UserID: 1},
				{ID: 2, Title: "Article 2", UserID: 2},
			},
			expectErr: false,
		},
		{
			description: "Edge Case with Empty User ID List",
			userIDs:     []uint{},
			limit:       2,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM articles WHERE user_id in \\(\\)").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "user_id"}))
			},
			expected:  []model.Article{},
			expectErr: false,
		},
		{
			description: "Edge Case with Limit Exceeding Total Articles",
			userIDs:     []uint{1},
			limit:       5,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM articles WHERE user_id in \\(\\$1\\)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "user_id"}).
						AddRow(1, "Article 1", 1).
						AddRow(2, "Article 2", 1))
			},
			expected: []model.Article{
				{ID: 1, Title: "Article 1", UserID: 1},
				{ID: 2, Title: "Article 2", UserID: 1},
			},
			expectErr: false,
		},
		{
			description: "Handling of Database Errors",
			userIDs:     []uint{1},
			limit:       2,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM articles WHERE user_id in \\(\\$1\\)").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			expected:  []model.Article{},
			expectErr: true,
		},
		{
			description: "Test with Offset Beyond Total Articles",
			userIDs:     []uint{1},
			limit:       2,
			offset:      5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM articles WHERE user_id in \\(\\$1\\)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "user_id"}))
			},
			expected:  []model.Article{},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm db: %s", err)
			}
			defer gormDB.Close()

			articleStore := &ArticleStore{db: gormDB}
			tc.mockSetup(mock)

			articles, err := articleStore.GetFeedArticles(tc.userIDs, tc.limit, tc.offset)

			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}

			if !reflect.DeepEqual(articles, tc.expected) {
				t.Errorf("Expected articles: %+v, got: %+v", tc.expected, articles)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %s", err)
			}
		})
	}
}


