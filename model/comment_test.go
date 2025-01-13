package model

import (
	"fmt"
	"math"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	validation "github.com/go-ozzo/ozzo-validation"
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
			name: "Successfully Convert Comment to Proto Comment",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 15, 11, 45, 0, 0, time.UTC),
				},
				Body: "Test comment body",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "Test comment body",
				CreatedAt: "2023-05-15T10:30:00Z",
				UpdatedAt: "2023-05-15T11:45:00Z",
			},
		},
		{
			name:    "Handle Zero Values in Comment Struct",
			comment: Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Verify ISO8601 Timestamp Formatting",
			comment: Comment{
				Model: gorm.Model{
					CreatedAt: time.Date(2023, 5, 15, 10, 30, 45, 123456789, time.UTC),
					UpdatedAt: time.Date(2023, 5, 15, 11, 45, 30, 987654321, time.UTC),
				},
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "2023-05-15T10:30:45Z",
				UpdatedAt: "2023-05-15T11:45:30Z",
			},
		},
		{
			name: "Handle Maximum Values for ID Field",
			comment: Comment{
				Model: gorm.Model{
					ID: math.MaxUint32,
				},
			},
			expected: &pb.Comment{
				Id:        fmt.Sprintf("%d", uint(math.MaxUint32)),
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Verify Unchanged Body Field",
			comment: Comment{
				Body: "This is a test comment.\nIt contains special characters: !@#$%^&*()_+\nAnd multiple lines.",
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "This is a test comment.\nIt contains special characters: !@#$%^&*()_+\nAnd multiple lines.",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Check Omission of Unexported Fields",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body:      "Test comment",
				UserID:    2,
				ArticleID: 3,
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "Test comment",
				CreatedAt: "",
				UpdatedAt: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.ProtoComment()

			assert.Equal(t, tt.expected.Id, result.Id)
			assert.Equal(t, tt.expected.Body, result.Body)

			if tt.name == "Check Omission of Unexported Fields" {

				assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`, result.CreatedAt)
				assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`, result.UpdatedAt)
			} else {
				assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
				assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
			}

			assert.Nil(t, result.Author)
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
	}{
		{
			name: "Valid Comment",
			comment: Comment{
				Body: "This is a valid comment",
			},
			wantErr: false,
		},
		{
			name: "Empty Body",
			comment: Comment{
				Body: "",
			},
			wantErr: true,
		},
		{
			name: "Whitespace-only Body",
			comment: Comment{
				Body: "   \t\n",
			},
			wantErr: true,
		},
		{
			name: "Very Long Body",
			comment: Comment{
				Body: string(make([]rune, 10000)),
			},
			wantErr: false,
		},
		{
			name: "Special Characters in Body",
			comment: Comment{
				Body: "This is a comment with special characters: 你好, world! 🌍",
			},
			wantErr: false,
		},
		{
			name: "Unset Fields",
			comment: Comment{
				Body: "Valid comment with other fields unset",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if _, ok := err.(validation.Errors); !ok {
					t.Errorf("Expected validation.Errors, got %T", err)
				}
			}
		})
	}
}

