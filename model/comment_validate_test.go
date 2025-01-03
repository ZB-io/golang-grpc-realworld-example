package model

import (
	"fmt"
	"testing"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	"strings"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}



type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestCommentValidate(t *testing.T) {
	tests := []struct {
		name        string
		comment     Comment
		expectedErr error
	}{
		{
			name: "Scenario 1: Validate With Populated Comment Body",
			comment: Comment{
				Body:      "This is a valid comment.",
				UserID:    1,
				ArticleID: 1,
			},
			expectedErr: nil,
		},
		{
			name: "Scenario 2: Validate With Empty Comment Body",
			comment: Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			expectedErr: fmt.Errorf("cannot be empty"),
		},
		{
			name: "Scenario 3: Validate With All Mandatory Fields Populated",
			comment: Comment{
				Body:      "This comment has all required fields.",
				UserID:    2,
				ArticleID: 3,
			},
			expectedErr: nil,
		},
		{
			name: "Scenario 4: Validate With Special Characters in Comment Body",
			comment: Comment{
				Body:      "Special characters !@#$%^&*()_+{}|:\"<>?",
				UserID:    1,
				ArticleID: 1,
			},
			expectedErr: nil,
		},
		{
			name: "Scenario 5: Validate with Very Long Comment Body",
			comment: Comment{
				Body:      strings.Repeat("a", 10000),
				UserID:    1,
				ArticleID: 2,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()

			if (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) ||
				(err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr == nil) {
				t.Errorf("Validate() error = %v, expectedErr %v", err, tt.expectedErr)
			} else {
				t.Logf("Passed: %s", tt.name)
			}
		})
	}
}

