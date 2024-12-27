package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"strings"
)

/*
ROOST_METHOD_HASH=ProtoComment_f8354e88c8
ROOST_METHOD_SIG_HASH=ProtoComment_ac7368a67c
*/
func TestCommentProtoComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  Comment
		expected *pb.Comment
	}{
		{
			name: "Typical Comment Conversion",
			comment: Comment{
				Model: gorm.Model{ID: 123},
				Body:  "This is a comment.",
				CreatedAt: time.Date(2023, 10, 13, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 14, 10, 0, 0, 0, time.UTC),
			},
			expected: &pb.Comment{
				Id:        "123",
				Body:      "This is a comment.",
				CreatedAt: "2023-10-13T10:00:00Z",
				UpdatedAt: "2023-10-14T10:00:00Z",
			},
		},
		{
			name: "Empty Comment Fields",
			comment: Comment{
				Model:    gorm.Model{ID: 0},
				Body:     "",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			expected: &pb.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Large ID Value Conversion",
			comment: Comment{
				Model: gorm.Model{ID: 9223372036854775807},
				Body:  "Comment with large ID",
				CreatedAt: time.Date(2023, 10, 13, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 14, 10, 0, 0, time.UTC),
			},
			expected: &pb.Comment{
				Id:        "9223372036854775807",
				Body:      "Comment with large ID",
				CreatedAt: "2023-10-13T10:00:00Z",
				UpdatedAt: "2023-10-14T10:00:00Z",
			},
		},
		{
			name: "Date Formatting Validity",
			comment: Comment{
				Model: gorm.Model{ID: 45},
				Body:  "Date formatting test",
				CreatedAt: time.Date(2021, 8, 17, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 9, 17, 12, 0, 0, 0, time.UTC),
			},
			expected: &pb.Comment{
				Id:        "45",
				Body:      "Date formatting test",
				CreatedAt: "2021-08-17T12:00:00Z",
				UpdatedAt: "2021-09-17T12:00:00Z",
			},
		},
		{
			name: "Check Comment Body Integrity",
			comment: Comment{
				Model: gorm.Model{ID: 300},
				Body:  "This is a very long comment body that needs to be preserved completely without any truncation during the conversion process.",
				CreatedAt: time.Date(2023, 10, 13, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 14, 10, 0, 0, time.UTC),
			},
			expected: &pb.Comment{
				Id:        "300",
				Body:      "This is a very long comment body that needs to be preserved completely without any truncation during the conversion process.",
				CreatedAt: "2023-10-13T10:00:00Z",
				UpdatedAt: "2023-10-14T10:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.comment.ProtoComment()
			
			assert.Equal(t, tt.expected, actual)
		})
	}
}

/*
ROOST_METHOD_HASH=Validate_1df97b5695
ROOST_METHOD_SIG_HASH=Validate_0591f679fe
*/
func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		comment   Comment
		expectErr bool
	}{
		{
			name: "Valid Comment Body",
			comment: Comment{
				Body: "This is a valid comment",
			},
			expectErr: false,
		},
		{
			name: "Empty Comment Body",
			comment: Comment{
				Body: "",
			},
			expectErr: true,
		},
		{
			name: "Whitespace-only Comment Body",
			comment: Comment{
				Body: "    ",
			},
			expectErr: true,
		},
		{
			name: "Large Comment Body",
			comment: Comment{
				Body: strings.Repeat("a", 10*1024),
			},
			expectErr: false,
		},
		{
			name: "Comment Body with Special Characters",
			comment: Comment{
				Body: "!@#%&*()",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()

			if tt.expectErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
		})
	}
}
