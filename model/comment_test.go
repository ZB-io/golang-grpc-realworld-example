package model

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	pb "path/to/your/proto/package" // Replace with the actual import path
)

/*
ROOST_METHOD_HASH=ProtoComment_f8354e88c8
ROOST_METHOD_SIG_HASH=ProtoComment_ac7368a67c
*/
func TestProtoComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  Comment
		expected *pb.Comment
	}{
		{
			name: "Valid Comment to pb.Comment conversion",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 1, 11, 0, 0, 0, time.UTC),
				},
				Body: "Test comment",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "Test comment",
				CreatedAt: "2023-05-01T10:00:00Z",
				UpdatedAt: "2023-05-01T11:00:00Z",
			},
		},
		{
			name:    "Zero values in Comment struct",
			comment: Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Maximum ID value",
			comment: Comment{
				Model: gorm.Model{
					ID: ^uint(0),
				},
				Body: "Max ID comment",
			},
			expected: &pb.Comment{
				Id:        fmt.Sprintf("%d", ^uint(0)),
				Body:      "Max ID comment",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Future dates",
			comment: Comment{
				Model: gorm.Model{
					ID:        100,
					CreatedAt: time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2050, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				Body: "Future comment",
			},
			expected: &pb.Comment{
				Id:        "100",
				Body:      "Future comment",
				CreatedAt: "2050-01-01T00:00:00Z",
				UpdatedAt: "2050-01-02T00:00:00Z",
			},
		},
		{
			name: "Comment with Author and Article",
			comment: Comment{
				Model: gorm.Model{
					ID:        200,
					CreatedAt: time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 6, 1, 13, 0, 0, 0, time.UTC),
				},
				Body:      "Comment with relations",
				UserID:    1,
				Author:    User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				ArticleID: 1,
				Article:   Article{Model: gorm.Model{ID: 1}, Title: "Test Article"},
			},
			expected: &pb.Comment{
				Id:        "200",
				Body:      "Comment with relations",
				CreatedAt: "2023-06-01T12:00:00Z",
				UpdatedAt: "2023-06-01T13:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.ProtoComment()
			assert.Equal(t, tt.expected.Id, result.Id, "ID mismatch")
			assert.Equal(t, tt.expected.Body, result.Body, "Body mismatch")
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt, "CreatedAt mismatch")
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt, "UpdatedAt mismatch")
			assert.Nil(t, result.Author, "Author should be nil")
		})
	}
}

/*
ROOST_METHOD_HASH=Validate_1df97b5695
ROOST_METHOD_SIG_HASH=Validate_0591f679fe
*/
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid Comment",
			comment: Comment{Body: "This is a valid comment"},
			wantErr: false,
		},
		{
			name:    "Empty Body",
			comment: Comment{Body: ""},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name:    "Whitespace-only Body",
			comment: Comment{Body: "   \t   "},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name:    "Very Long Body",
			comment: Comment{Body: strings.Repeat("a", 10000)},
			wantErr: false,
		},
		{
			name:    "Special Characters in Body",
			comment: Comment{Body: "!@#$%^&*()_+{}[]|:;\"'<>,.?/~`"},
			wantErr: false,
		},
		{
			name:    "Unicode Characters in Body",
			comment: Comment{Body: "こんにちは世界 - Здравствуй, мир - مرحبا بالعالم"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("Comment.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
