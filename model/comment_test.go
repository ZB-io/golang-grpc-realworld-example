package model

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	pb "your-project-path/proto" // Replace with the actual path to your proto package
)

const ISO8601 = "2006-01-02T15:04:05Z07:00"

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
			name: "Valid Comment",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 15, 11, 45, 0, 0, time.UTC),
				},
				Body: "This is a test comment",
			},
			expected: &pb.Comment{
				Id:        "1",
				Body:      "This is a test comment",
				CreatedAt: "2023-05-15T10:30:00Z",
				UpdatedAt: "2023-05-15T11:45:00Z",
			},
		},
		{
			name: "Zero values",
			comment: Comment{
				Model: gorm.Model{},
				Body:  "",
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Large Body content",
			comment: Comment{
				Model: gorm.Model{
					ID:        999,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body: string(make([]byte, 1024*1024)),
			},
			expected: &pb.Comment{
				Id:        "999",
				Body:      string(make([]byte, 1024*1024)),
				CreatedAt: time.Now().Format(ISO8601),
				UpdatedAt: time.Now().Format(ISO8601),
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
			assert.Nil(t, result.Author)

			if tt.name == "Large Body content" {
				assert.Len(t, result.Body, 1024*1024)
			}
		})
	}
}

func TestProtoCommentPerformance(t *testing.T) {
	largeComment := Comment{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Body: string(make([]byte, 1024*1024)),
	}

	start := time.Now()
	result := largeComment.ProtoComment()
	duration := time.Since(start)

	assert.NotNil(t, result)
	assert.Less(t, duration, 100*time.Millisecond)
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
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Whitespace-only Body",
			comment: Comment{
				Body: "   \t\n",
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Very Long Body",
			comment: Comment{
				Body: strings.Repeat("a", 10000),
			},
			wantErr: false,
		},
		{
			name: "Unicode Character Body",
			comment: Comment{
				Body: "This is a comment with Unicode characters: 你好世界 🌍",
			},
			wantErr: false,
		},
		{
			name:    "Uninitialized Comment",
			comment: Comment{},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Valid Body with Empty Other Fields",
			comment: Comment{
				Body:      "Valid comment",
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
				if err == nil {
					t.Errorf("Comment.Validate() expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Comment.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
