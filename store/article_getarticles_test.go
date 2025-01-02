package store

import (
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)



func TestArticleStoreGetArticles(t *testing.T) {

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

	tests := []struct {
		name          string
		tagName       string
		username      string
		favoritedBy   *model.User
		limit         int64
		offset        int64
		mock          func()
		wantErr       bool
		expectedError error
	}{
		{

			name:     "Normal Operation - Fetching Articles by a Specific Username",
			username: "testUser",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("SELECT").WithArgs("testUser").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{

			name:    "Normal Operation - Fetching Articles by a Specific Tag",
			tagName: "testTag",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("SELECT").WithArgs("testTag").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{

			name: "Error Handling - When the Database Returns an Error",
			mock: func() {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("database error"))
			},
			wantErr:       true,
			expectedError: errors.New("database error"),
		},
		{

			name:   "Edge Case - Fetching Articles with Limit and Offset",
			limit:  10,
			offset: 5,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2).
					AddRow(3).
					AddRow(4).
					AddRow(5).
					AddRow(6).
					AddRow(7).
					AddRow(8).
					AddRow(9).
					AddRow(10)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, but got none")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error: %v, but got error: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
