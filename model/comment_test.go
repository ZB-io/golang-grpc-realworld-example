package model

import (
	"math"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"github.com/go-ozzo/ozzo-validation"
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
			name: "Valid Comment",
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
			name:    "Zero values",
			comment: Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Verify time format",
			comment: Comment{
				Model: gorm.Model{
					CreatedAt: time.Date(2023, 5, 1, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 1, 11, 45, 0, 0, time.UTC),
				},
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "2023-05-01T10:30:00Z",
				UpdatedAt: "2023-05-01T11:45:00Z",
			},
		},
		{
			name: "Maximum ID value",
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
			name: "Empty Body",
			comment: Comment{
				Model: gorm.Model{
					ID: 1,
				},
				Body: "",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.ProtoComment()
			assert.Equal(t, tt.expected, result)

			if tt.name == "Valid Comment" {
				result2 := tt.comment.ProtoComment()
				assert.Equal(t, result, result2)
			}
		})
	}

	t.Run("Author field not included", func(t *testing.T) {
		comment := Comment{
			Model: gorm.Model{ID: 1},
			Body:  "Test",
			Author: User{
				Model:    gorm.Model{ID: 2},
				Username: "testuser",
			},
		}
		result := comment.ProtoComment()
		assert.NotContains(t, result, "Author")
		assert.NotContains(t, result, "testuser")
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
			name: "Valid Comment",
			comment: Comment{
				Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Body:  "This is a valid comment",
			},
			wantErr: false,
		},
		{
			name: "Empty Body",
			comment: Comment{
				Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Body:  "",
			},
			wantErr: true,
		},
		{
			name: "Whitespace-only Body",
			comment: Comment{
				Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Body:  "   \t\n",
			},
			wantErr: true,
		},
		{
			name: "Very Long Body",
			comment: Comment{
				Model: gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Body:  string(make([]rune, 10000)),
			},
			wantErr: false,
		},
		{
			name: "Special Characters in Body",
			comment: Comment{
				Model: gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Body:  "Special chars: ©®™£€¥§±×÷ and emojis: 😀🎉🌈",
			},
			wantErr: false,
		},
		{
			name: "Unset Fields",
			comment: Comment{
				Body: "Only Body is set",
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
				if _, ok := err.(validation.Errors); !ok {
					t.Errorf("Expected validation.Errors, got %T", err)
				}
			}
		})
	}
}

