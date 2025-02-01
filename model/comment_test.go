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
	now := time.Now()
	iso8601Format := "2006-01-02T15:04:05.999Z07:00"

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
					CreatedAt: now,
					UpdatedAt: now,
				},
				Body: "Test comment",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "Test comment",
				CreatedAt: now.Format(iso8601Format),
				UpdatedAt: now.Format(iso8601Format),
			},
		},
		{
			name:    "Zero values",
			comment: &Comment{},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: time.Time{}.Format(iso8601Format),
				UpdatedAt: time.Time{}.Format(iso8601Format),
			},
		},
		{
			name: "Large ID",
			comment: &Comment{
				Model: gorm.Model{
					ID:        math.MaxUint32,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Body: "Large ID comment",
			},
			expected: &pb.Comment{
				Id:        fmt.Sprintf("%d", math.MaxUint32),
				Body:      "Large ID comment",
				CreatedAt: now.Format(iso8601Format),
				UpdatedAt: now.Format(iso8601Format),
			},
		},
		{
			name:     "Nil Comment",
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
			name: "Comment with Whitespace-Only Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      "   ",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
		},
		{
			name: "Comment with Very Long Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      string(make([]rune, 10000)),
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Comment with Special Characters in Body",
			comment: Comment{
				Model:     gorm.Model{},
				Body:      "This is a comment with special characters: !@#$%^&*()_+ and emojis: 😀🎉👍",
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
		{
			name: "Comment with Null Body",
			comment: Comment{
				Model:     gorm.Model{},
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
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

