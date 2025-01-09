package model

import (
	"math"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)








/*
ROOST_METHOD_HASH=ProtoComment_f8354e88c8
ROOST_METHOD_SIG_HASH=ProtoComment_ac7368a67c

FUNCTION_DEF=func (c *Comment) ProtoComment() *pb.Comment 

 */
func TestCommentProtoComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  Comment
		expected *pb.Comment
	}{
		{
			name: "Successful Conversion",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 5, 1, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 1, 11, 45, 0, 0, time.UTC),
				},
				Body: "Test comment",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "Test comment",
				CreatedAt: "2023-05-01T10:30:00Z",
				UpdatedAt: "2023-05-01T11:45:00Z",
			},
		},
		{
			name:    "Zero Values",
			comment: Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Correct Timestamp Formatting",
			comment: Comment{
				Model: gorm.Model{
					CreatedAt: time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "2023-12-31T23:59:59Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			},
		},
		{
			name: "Large ID Value",
			comment: Comment{
				Model: gorm.Model{
					ID: math.MaxUint32,
				},
			},
			expected: &pb.Comment{
				Id:        "4294967295",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Body Content Preservation",
			comment: Comment{
				Body: "This is a test comment with special characters: !@#$%^&*()_+\nNew line and Unicode: 你好",
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "This is a test comment with special characters: !@#$%^&*()_+\nNew line and Unicode: 你好",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Deleted Comment",
			comment: Comment{
				Model: gorm.Model{
					DeletedAt: func() *time.Time {
						t := time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC)
						return &t
					}(),
				},
				Body: "Deleted comment",
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "Deleted comment",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.ProtoComment()
			assert.Equal(t, tt.expected.Id, result.Id)
			assert.Equal(t, tt.expected.Body, result.Body)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}


/*
ROOST_METHOD_HASH=Validate_1df97b5695
ROOST_METHOD_SIG_HASH=Validate_0591f679fe

FUNCTION_DEF=func (c Comment) Validate() error 

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
				Body: "This is a valid comment",
			},
			wantErr: false,
		},
		{
			name: "Invalid Comment with Empty Body",
			comment: Comment{
				Body: "",
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Comment with Very Long Body",
			comment: Comment{
				Body: string(make([]rune, 10000)),
			},
			wantErr: false,
		},
		{
			name: "Comment with Special Characters in Body",
			comment: Comment{
				Body: "!@#$%^&*()_+",
			},
			wantErr: false,
		},
		{
			name: "Comment with Only Whitespace in Body",
			comment: Comment{
				Body: "   \t\n",
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Validation of Other Fields",
			comment: Comment{
				Body:      "Valid body",
				UserID:    0,
				ArticleID: 0,
			},
			wantErr: false,
		},
		{
			name: "Comment with Unicode Characters in Body",
			comment: Comment{
				Body: "こんにちは世界 👋🌍",
			},
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
			if tt.wantErr {
				if err == nil {
					t.Errorf("Comment.Validate() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Comment.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

