package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
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
		expected struct {
			id        string
			body      string
			createdAt string
			updatedAt string
		}
	}{
		{
			name: "Scenario 1: Valid Comment Conversion with All Fields Populated",
			comment: Comment{
				Model: gorm.Model{
					ID:        123,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
				Body: "Test comment content",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "123",
				body:      "Test comment content",
				createdAt: "2023-01-01T12:00:00Z",
				updatedAt: "2023-01-02T12:00:00Z",
			},
		},
		{
			name:    "Scenario 2: Comment Conversion with Zero Values",
			comment: Comment{},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "0",
				body:      "",
				createdAt: "0001-01-01T00:00:00Z",
				updatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Scenario 3: Comment Conversion with Maximum Values",
			comment: Comment{
				Model: gorm.Model{
					ID:        ^uint(0),
					CreatedAt: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
					UpdatedAt: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
				},
				Body: "Very long comment body with lots of content...",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "18446744073709551615",
				body:      "Very long comment body with lots of content...",
				createdAt: "9999-12-31T23:59:59Z",
				updatedAt: "9999-12-31T23:59:59Z",
			},
		},
		{
			name: "Scenario 4: Comment Conversion with Special Characters in Body",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				Body: "Special chars: 你好 👋 !@#$%^&*()",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "1",
				body:      "Special chars: 你好 👋 !@#$%^&*()",
				createdAt: "2023-01-01T12:00:00Z",
				updatedAt: "2023-01-01T12:00:00Z",
			},
		},
		{
			name: "Scenario 5: Comment Conversion Time Format Consistency",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
				},
				Body: "Time zone test",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "1",
				body:      "Time zone test",
				createdAt: "2023-01-01T17:00:00Z",
				updatedAt: "2023-01-01T20:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing:", tt.name)

			result := tt.comment.ProtoComment()

			assert.NotNil(t, result, "ProtoComment result should not be nil")
			assert.Equal(t, tt.expected.id, result.Id, "ID mismatch")
			assert.Equal(t, tt.expected.body, result.Body, "Body mismatch")
			assert.Equal(t, tt.expected.createdAt, result.CreatedAt, "CreatedAt mismatch")
			assert.Equal(t, tt.expected.updatedAt, result.UpdatedAt, "UpdatedAt mismatch")

			t.Log("Test completed successfully")
		})
	}
}

