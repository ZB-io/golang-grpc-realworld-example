package model

import (
	"strings"
	"testing"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

/*
ROOST_METHOD_HASH=Validate_1df97b5695
ROOST_METHOD_SIG_HASH=Validate_0591f679fe


 */
func TestCommentValidate(t *testing.T) {

	tests := []struct {
		name    string
		comment Comment
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Comment with Non-Empty Body",
			comment: Comment{
				Body:      "This is a valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Empty Body Field",
			comment: Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
			errMsg:  "cannot be blank",
		},
		{
			name: "Body Field with Only Whitespace",
			comment: Comment{
				Body:      "    \t\n",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
			errMsg:  "cannot be blank",
		},
		{
			name: "Very Long Body Content",
			comment: Comment{
				Body:      strings.Repeat("a", 10000),
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Body with Special Characters",
			comment: Comment{
				Body:      "Comment with emoji 😀 and специальные символы",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "Zero Value Comment",
			comment: Comment{},
			wantErr: true,
			errMsg:  "cannot be blank",
		},
		{
			name: "Valid Required Fields but Empty Body",
			comment: Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
				Author:    User{Username: "testuser"},
				Article:   Article{Title: "Test Article"},
			},
			wantErr: true,
			errMsg:  "cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.name)

			err := tt.comment.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			}

			if err != nil {
				t.Logf("Validation failed as expected: %v", err)
			} else {
				t.Logf("Validation passed successfully")
			}
		})
	}
}

