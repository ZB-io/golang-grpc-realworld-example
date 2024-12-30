package store

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestGetFeedArticles(t *testing.T) {
	tests := []struct {
		name     string
		userIDs  []uint
		limit    int64
		offset   int64
		mockDB   func() (*gorm.DB, sqlmock.Sqlmock)
		expected []model.Article
		wantErr  bool
	}{
		{
			name:    "Successful retrieval of feed articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gdb, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "author_id", "author_username"}).
					AddRow(1, "Article 1", "Desc 1", "Body 1", 1, 1, "user1").
					AddRow(2, "Article 2", "Desc 2", "Body 2", 2, 2, "user2")

				mock.ExpectQuery("SELECT (.+) FROM `articles` (.+) WHERE").
					WithArgs(1, 2).
					WillReturnRows(rows)

				return gdb, mock
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", Description: "Desc 1", Body: "Body 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", Description: "Desc 2", Body: "Body 2", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{999, 1000},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gdb, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "author_id", "author_username"})

				mock.ExpectQuery("SELECT (.+) FROM `articles` (.+) WHERE").
					WithArgs(999, 1000).
					WillReturnRows(rows)

				return gdb, mock
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Database error handling",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gdb, _ := gorm.Open("mysql", db)

				mock.ExpectQuery("SELECT (.+) FROM `articles` (.+) WHERE").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))

				return gdb, mock
			},
			expected: nil,
			wantErr:  true,
		},
		// Add more test cases here as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb, mock := tt.mockDB()
			store := &ArticleStore{
				db: gdb,
			}

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
