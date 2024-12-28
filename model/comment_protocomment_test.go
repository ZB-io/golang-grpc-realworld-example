package model

import (
	"fmt"
	"testing"
	"time"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/proto"
)





func TestCommentProtoComment(t *testing.T) {

	const ISO8601 = "2006-01-02T15:04:05Z07:00"

	tests := []struct {
		name     string
		comment  Comment
		expected *proto.Comment
	}{
		{
			name: "Typical Comment Conversion",
			comment: Comment{
				ID:        123,
				Body:      "This is a comment.",
				CreatedAt: time.Date(2023, 10, 13, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 14, 10, 0, 0, 0, time.UTC),
			},
			expected: &proto.Comment{
				Id:        "123",
				Body:      "This is a comment.",
				CreatedAt: "2023-10-13T10:00:00Z",
				UpdatedAt: "2023-10-14T10:00:00Z",
			},
		},
		{
			name: "Empty Comment Fields",
			comment: Comment{
				ID:        0,
				Body:      "",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			expected: &proto.Comment{
				Id:        "0",
				Body:      "",
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Large ID Value Conversion",
			comment: Comment{
				ID:        9223372036854775807,
				Body:      "Comment with large ID",
				CreatedAt: time.Date(2023, 10, 13, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 14, 10, 0, 0, time.UTC),
			},
			expected: &proto.Comment{
				Id:        "9223372036854775807",
				Body:      "Comment with large ID",
				CreatedAt: "2023-10-13T10:00:00Z",
				UpdatedAt: "2023-10-14T10:00:00Z",
			},
		},
		{
			name: "Date Formatting Validity",
			comment: Comment{
				ID:        45,
				Body:      "Date formatting test",
				CreatedAt: time.Date(2021, 8, 17, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 9, 17, 12, 0, 0, 0, time.UTC),
			},
			expected: &proto.Comment{
				Id:        "45",
				Body:      "Date formatting test",
				CreatedAt: "2021-08-17T12:00:00Z",
				UpdatedAt: "2021-09-17T12:00:00Z",
			},
		},
		{
			name: "Check Comment Body Integrity",
			comment: Comment{
				ID:        300,
				Body:      "This is a very long comment body that needs to be preserved completely without any truncation during the conversion process.",
				CreatedAt: time.Date(2023, 10, 13, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 14, 10, 0, 0, time.UTC),
			},
			expected: &proto.Comment{
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

			var buffer bytes.Buffer
			fmt.Fprintf(&buffer, "Scenario: %s\nExpected: %+v\nActual: %+v\n", tt.name, tt.expected, actual)
			t.Log(buffer.String())
		})
	}
}



