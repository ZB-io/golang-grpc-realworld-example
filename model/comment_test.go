package github

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
			name: "Valid Comment Conversion",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 5, 1, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 1, 11, 45, 0, 0, time.UTC),
				},
				Body: "This is a valid comment",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "This is a valid comment",
				CreatedAt: "2023-05-01T10:30:00Z",
				UpdatedAt: "2023-05-01T11:45:00Z",
			},
		},
		{
			name:    "Zero Values in Comment",
			comment: Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Large ID Value",
			comment: Comment{
				Model: gorm.Model{
					ID:        math.MaxUint32,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body: "Comment with large ID",
			},
			expected: &pb.Comment{
				Id:   fmt.Sprintf("%d", math.MaxUint32),
				Body: "Comment with large ID",
			},
		},
		{
			name: "Unicode Characters in Body",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body: "Unicode: こんにちは 世界 🌍",
			},
			expected: &pb.Comment{
				Id:   "1",
				Body: "Unicode: こんにちは 世界 🌍",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.ProtoComment()

			assert.Equal(t, tt.expected.Id, result.Id)
			assert.Equal(t, tt.expected.Body, result.Body)

			if tt.name != "Zero Values in Comment" {

				_, err := time.Parse(time.RFC3339, result.CreatedAt)
				assert.NoError(t, err, "CreatedAt should be in ISO8601 format")
				_, err = time.Parse(time.RFC3339, result.UpdatedAt)
				assert.NoError(t, err, "UpdatedAt should be in ISO8601 format")
			} else {
				assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
				assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
			}
		})
	}

	t.Run("Consistency of Multiple Calls", func(t *testing.T) {
		comment := Comment{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Body: "Test comment",
		}

		result1 := comment.ProtoComment()
		result2 := comment.ProtoComment()

		assert.Equal(t, result1, result2, "Multiple calls should produce consistent results")
	})
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
		},
		{
			name: "Comment with Very Long Body",
			comment: Comment{
				Body: string(make([]byte, 10000)),
			},
			wantErr: false,
		},
		{
			name: "Comment with Special Characters in Body",
			comment: Comment{
				Body: "Hello! This is a test: @#$%^&*()",
			},
			wantErr: false,
		},
		{
			name: "Comment with Only Whitespace Characters in Body",
			comment: Comment{
				Body: "   \t\n",
			},
			wantErr: true,
		},
		{
			name: "Validation of Other Comment Fields",
			comment: Comment{
				Body:      "Valid body",
				UserID:    0,
				ArticleID: 0,
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

