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
	const ISO8601 = "2006-01-02T15:04:05.999Z07:00"

	tests := []struct {
		name     string
		comment  *Comment
		expected *pb.Comment
	}{
		{
			name: "Valid Comment",
			comment: &Comment{
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
				CreatedAt: "2023-05-01T10:00:00.000Z",
				UpdatedAt: "2023-05-01T11:00:00.000Z",
			},
		},
		{
			name:    "Zero values",
			comment: &Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: time.Time{}.Format(ISO8601),
				UpdatedAt: time.Time{}.Format(ISO8601),
			},
		},
		{
			name: "ISO8601 date formatting",
			comment: &Comment{
				Model: gorm.Model{
					CreatedAt: time.Date(2023, 5, 1, 10, 30, 45, 123000000, time.UTC),
					UpdatedAt: time.Date(2023, 5, 1, 11, 45, 30, 456000000, time.UTC),
				},
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "2023-05-01T10:30:45.123Z",
				UpdatedAt: "2023-05-01T11:45:30.456Z",
			},
		},
		{
			name: "Maximum uint64 ID",
			comment: &Comment{
				Model: gorm.Model{
					ID: math.MaxUint64,
				},
			},
			expected: &pb.Comment{
				Id:        fmt.Sprintf("%d", uint64(math.MaxUint64)),
				Body:      "",
				CreatedAt: time.Time{}.Format(ISO8601),
				UpdatedAt: time.Time{}.Format(ISO8601),
			},
		},
		{
			name: "Body with special characters",
			comment: &Comment{
				Body: "Line 1\nLine 2\tTabbed\nSpecial chars: !@#$%^&*()",
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "Line 1\nLine 2\tTabbed\nSpecial chars: !@#$%^&*()",
				CreatedAt: time.Time{}.Format(ISO8601),
				UpdatedAt: time.Time{}.Format(ISO8601),
			},
		},
		{
			name:     "Nil pointer",
			comment:  nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.ProtoComment()
			assert.Equal(t, tt.expected, result)
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
			name: "Valid Comment with Non-Empty Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      "This is a valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid Comment with Empty Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
		},
		{
			name: "Comment with Very Long Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      fmt.Sprintf("%010000d", 0),
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Comment with Only Whitespace in Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      "   \t\n",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
		},
		{
			name: "Comment with Special Characters in Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      "This is a comment with special characters: 🚀 こんにちは",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Validation of Other Fields",
			comment: Comment{
				Model:     gorm.Model{},
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

