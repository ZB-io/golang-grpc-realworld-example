package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)









func TestArticleStoreDeleteComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	store := &ArticleStore{db: gdb}

	tests := []struct {
		name     string
		comment  *model.Comment
		wantErr  bool
		mockFunc func()
	}{
		{
			name:    "Successful Deletion of a Comment",
			comment: &model.Comment{Model: gorm.Model{ID: 1}},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "Deletion of a non-existent Comment",
			comment: &model.Comment{Model: gorm.Model{ID: 2}},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").WithArgs(2).WillReturnError(errors.New("record not found"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "Deletion of a Comment with a null `Comment` struct",
			comment: nil,
			mockFunc: func() {

			},
			wantErr: true,
		},
		{
			name:    "Deletion of a Comment with invalid foreign keys",
			comment: &model.Comment{Model: gorm.Model{ID: 3}, UserID: 999, ArticleID: 999},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").WithArgs(3).WillReturnError(errors.New("foreign key constraint fails"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := store.DeleteComment(tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
