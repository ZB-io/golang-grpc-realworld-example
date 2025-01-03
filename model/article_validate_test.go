package model

import (
	"testing"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type Comment struct {
	Body string
}
type Tag struct {
	Name string
}
type User struct {
	Username string `gorm:"unique_index;not null"`
	Email    string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Bio      string `gorm:"not null"`
	Image    string `gorm:"not null"`
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestArticleValidate(t *testing.T) {
	type testCase struct {
		name          string
		article       Article
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successful Validation with All Required Fields Present",
			article: Article{
				Title: "Sample Title",
				Body:  "Sample Body",
				Tags:  []Tag{{Name: "Go"}, {Name: "Programming"}},
			},
			expectedError: nil,
		},
		{
			name: "Validation Failure Due to Missing Title",
			article: Article{
				Body: "Sample Body",
				Tags: []Tag{{Name: "Go"}, {Name: "Programming"}},
			},
			expectedError: validation.Errors{"Title": validation.ErrRequired},
		},
		{
			name: "Validation Failure Due to Missing Body",
			article: Article{
				Title: "Sample Title",
				Tags:  []Tag{{Name: "Go"}, {Name: "Programming"}},
			},
			expectedError: validation.Errors{"Body": validation.ErrRequired},
		},
		{
			name: "Validation Failure Due to Missing Tags",
			article: Article{
				Title: "Sample Title",
				Body:  "Sample Body",
			},
			expectedError: validation.Errors{"Tags": validation.ErrRequired},
		},
		{
			name: "Validation with Edge Case of Minimum Field Content",
			article: Article{
				Title: "A",
				Body:  "B",
				Tags:  []Tag{{Name: "C"}},
			},
			expectedError: nil,
		},
		{
			name: "Validation with All Fields Filled Beyond Requirements",
			article: Article{
				Title:          "Sample Title",
				Body:           "Sample Body",
				Tags:           []Tag{{Name: "Go"}},
				Author:         User{Username: "AuthorName"},
				UserID:         1,
				FavoritesCount: 10,
				FavoritedUsers: []User{{Username: "User1"}, {Username: "User2"}},
				Comments:       []Comment{{Body: "Nice Article!"}},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if err != nil && tt.expectedError == nil {
				t.Errorf("Expected nil error, got %v", err)
				t.Logf("Test Case Failed: %s - Article: %+v", tt.name, tt.article)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error, got nil")
				t.Logf("Test Case Failed: %s - Article: %+v", tt.name, tt.article)
			}
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				t.Logf("Test Case Result: expectedError: %v, actualError: %v", tt.expectedError, err)
			}
			if err == nil {
				t.Logf("Test Case Success: %s", tt.name)
			}
		})
	}
}

