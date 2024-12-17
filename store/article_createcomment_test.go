package store

import (
	"testing"
	"errors"
	
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestArticleStoreCreateComment(t *testing.T) {
	// Define the test cases for the function
	testCases := []struct {
		Name      string
		Comment   model.Comment
		ShouldErr bool
	}{
		{
			Name: "Test Successful Comment Creation",
			Comment: model.Comment{
				Body:      "Great Article!",
				UserID:    1,
				Author:    model.User{Username: "user1"},
				ArticleID: 1,
				Article:   model.Article{Title: "My Article"},
			},
			ShouldErr: false,
		},
		{
			Name: "Test Failed Comment Creation - Missing Fields",
			Comment: model.Comment{
				Body:    "",
				UserID:  0,
				Article: model.Article{},
			},
			ShouldErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// Initialize the sqlmock database
			sqlDb, mock, _ := sqlmock.New()
			db, _ := gorm.Open("postgres", sqlDb)

			// Initialize the ArticleStore
			store := ArticleStore{db}

			// If we expect an error, mock the gorm Create call to return an error
			if testCase.ShouldErr {
				mock.
					ExpectExec("INSERT").
					WillReturnError(errors.New("failed to connect to database"))
			}

			// Execute the CreateComment function
			err := store.CreateComment(&testCase.Comment)

			// Ensure the function returns an error when it should
			if testCase.ShouldErr && err == nil {
				t.Errorf("Expected an error, but did not get one")
			}

			// Ensure the function does not return an error when it shouldn't
			if !testCase.ShouldErr && err != nil {
				t.Errorf("Received an unexpected error when creating comment: %v", err)
			}
		})
	}

	// TODO Specify a new test case in which the database is disconnected
}
