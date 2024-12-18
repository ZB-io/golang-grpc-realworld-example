package store_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name          string
		mockFn        func(mock sqlmock.Sqlmock)
		inputArticle  *model.Article
		expectedError error
	}{
		{
			name: "Successfully creating a new article",
            mockFn: func(mock sqlmock.Sqlmock) {
            	//Async query and result mock
            },
			inputArticle: &model.Article{Title: "Test Title", Description: "Test Description", Body: "Test Body", 
				UserID: 1},
			expectedError: nil,
		},
		{
			name: "Creating an article without necessary fields",
            mockFn: func(mock sqlmock.Sqlmock) {
            	//Async query and result mock
            },
			inputArticle: &model.Article{Title: "", Description: "", Body: "", UserID: 0},
			expectedError: errors.New("Incomplete article fields"),
		},
		{
			name: "Database connectivity issues when creating an article",
            mockFn: func(mock sqlmock.Sqlmock) {
            	//Async error scenario mock
            },
			inputArticle: &model.Article{Title: "Test Title", Description: "Test Description", Body: "Test Body", 
				UserID: 1},
			expectedError: errors.New("DB connection error"),
		},
		{
			name: "Null input to the create function",
            mockFn: func(mock sqlmock.Sqlmock) {},
			inputArticle: nil,
			expectedError: errors.New("Null input"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error init mock DB: %v", err)
			}
			defer db.Close()

			tt.mockFn(mock)

			gormDB, _ := gorm.Open("mysql", db)
			articleStore := &store.ArticleStore{Db: gormDB}
			err = articleStore.Create(tt.inputArticle)

			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
