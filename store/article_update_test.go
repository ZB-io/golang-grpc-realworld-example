package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

func TestArticleStoreUpdate(t *testing.T) {
	// create sqlmock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error while creating sqlmock instance: %v", err)
	}
	gormDB, _ := gorm.Open("mysql", db) // TODO: 'mysql' assumed, replace it as per actual DB type

	article := &model.Article{
		Model:       gorm.Model{},
		Title:       "Test Article",
		Description: "This is a test article",
		Body:        "The article body",
		Tags:        nil,
		Author:      model.User{},
		UserID:      0,
		FavoritesCount: 0,
		FavoritedUsers: nil,
		Comments:    nil,
	}

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name   string
		setup  func()
		action func() error
		verify func(error) error
	}{
		{
			name: "Successful Update of Article",
			setup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			action: func() error {
				return store.Update(article)
			},
			verify: func(err error) error {
				if err != nil {
					return err
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Updating Non-Existent Article",
			setup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectRollback()
			},
			action: func() error {
				return store.Update(article)
			},
			verify: func(err error) error {
				if err == nil {
					return errors.New("expected error but got nil")
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Database Connection Failure During Update",
			setup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("db connection failure"))
				mock.ExpectRollback()
			},
			action: func() error {
				return store.Update(article)
			},
			verify: func(err error) error {
				if err == nil {
					return errors.New("expected error but got nil")
				}
				return nil
			},
		},
		{
			name: "Updating Article with Invalid Data",
			setup: func() {
				// Assuming "Title" is required and cannot be blank, setting it blank for this test
				article.Title = ""
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("title should not be blank"))
				mock.ExpectRollback()
			},
			action: func() error {
				return store.Update(article)
			},
			verify: func(err error) error {
				if err == nil || err.Error() != "title should not be blank" {
					return errors.New("expected error about title should not be blank but got nil or different error message")
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			err := test.action()
			if verErr := test.verify(err); verErr != nil {
				t.Error(verErr)
			}
		})
	}
}
