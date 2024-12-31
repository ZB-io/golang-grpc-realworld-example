package model

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"github.com/go-ozzo/ozzo-validation/v4"
	pb "your_protobuf_package_path" // Replace with the actual path to your protobuf package
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
				start := time.Now()
				tt.comment.ProtoComment()
				duration := time.Since(start)
				assert.Less(t, duration.Milliseconds(), int64(100), "ProtoComment took too long for large body")
			}
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
				Body: "   \t   ",
			},
			wantErr: true,
		},
		{
			name: "Very Long Body",
			comment: Comment{
				Body: strings.Repeat("a", 10000),
			},
			wantErr: false,
		},
		{
			name: "Special Characters in Body",
			comment: Comment{
				Body: "!@#$%^&*()_+ Special characters are allowed",
			},
			wantErr: false,
		},
		{
			name: "Non-ASCII Characters in Body",
			comment: Comment{
				Body: "こんにちは 你好 🌍🌎🌏",
			},
			wantErr: false,
		},
		{
			name: "Other Fields Populated",
			comment: Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Valid comment with other fields populated",
				UserID:    123,
				ArticleID: 456,
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
