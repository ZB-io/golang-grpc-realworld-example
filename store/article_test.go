package store

import (
	"errors"
	"testing"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {

	tests := []struct {
		name    string
		article *model.Article
		setupFn func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful Article Creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test-tag"}},
			},
			setupFn: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Missing Required Fields",
			article: &model.Article{

				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(errors.New("title cannot be null"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "title cannot be null",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(errors.New("database connection lost"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection lost",
		},
		{
			name: "Article with Related Entities",
			article: &model.Article{
				Title:       "Test Article with Relations",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
				Author:      model.User{},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `article_tags`").
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Maximum Field Length Article",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			if tt.setupFn != nil {
				tt.setupFn(mock)
			}

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Error: %v", tt.name, err)
		})
	}
}

