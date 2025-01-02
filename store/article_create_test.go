package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)









func TestArticleStoreCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql DB: %v", err)
	}
	defer db.Close()

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	s := &ArticleStore{
		db: gdb,
	}

	tests := []struct {
		name     string
		article  *model.Article
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Successful Article Creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Article Creation with Missing Mandatory Fields",
			article: &model.Article{
				Title: "Test Article",
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("missing mandatory fields"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database Connection Error during Article Creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
