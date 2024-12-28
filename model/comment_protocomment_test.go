package model

import (
	"testing"
	"time"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestCommentProtoComment(t *testing.T) {
	tests := []struct {
		name      string
		comment   model.Comment
		expected  pb.Comment
		wantError bool
	}{
		{
			name: "Valid Comment Conversion to Proto",
			comment: model.Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body: "This is a test comment.",
				Author: model.User{
					Username: "testuser",
				},
			},
			expected: pb.Comment{
				Id:        "1",
				Body:      "This is a test comment.",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
				Author: &pb.Profile{
					Username: "testuser",
				},
			},
			wantError: false,
		},
		{
			name: "Edge Case of Comment with Empty Body",
			comment: model.Comment{
				Model: gorm.Model{
					ID: 2,
				},
				Body: "",
				Author: model.User{
					Username: "testuser",
				},
			},
			expected: pb.Comment{
				Id:   "2",
				Body: "",
				Author: &pb.Profile{
					Username: "testuser",
				},
			},
			wantError: false,
		},
		{
			name: "Valid Comment with Maximum ID Int Limit",
			comment: model.Comment{
				Model: gorm.Model{
					ID:        uint(^uint(0) >> 1),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body: "Max ID comment.",
				Author: model.User{
					Username: "maxuser",
				},
			},
			expected: pb.Comment{
				Id:   "18446744073709551615",
				Body: "Max ID comment.",
				Author: &pb.Profile{
					Username: "maxuser",
				},
			},
			wantError: false,
		},
		{
			name: "Handling Comment with Uninitialized Time Fields",
			comment: model.Comment{
				Model: gorm.Model{
					ID: 3,
				},
				Body: "No timestamps.",
				Author: model.User{
					Username: "notimeuser",
				},
			},
			expected: pb.Comment{
				Id:        "3",
				Body:      "No timestamps.",
				CreatedAt: "",
				UpdatedAt: "",
				Author: &pb.Profile{
					Username: "notimeuser",
				},
			},
			wantError: false,
		},
		{
			name: "Ensures Correct Author Mapping",
			comment: model.Comment{
				Model: gorm.Model{
					ID: 4,
				},
				Body: "Correct author mapping.",
				Author: model.User{
					Username: "correctauthor",
				},
			},
			expected: pb.Comment{
				Id:   "4",
				Body: "Correct author mapping.",
				Author: &pb.Profile{
					Username: "correctauthor",
				},
			},
			wantError: false,
		},
		{
			name: "Conversion with Null Author in Comment",
			comment: model.Comment{
				Model: gorm.Model{
					ID: 5,
				},
				Body: "Null author.",
			},
			expected: pb.Comment{
				Id:     "5",
				Body:   "Null author.",
				Author: nil,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.comment.ProtoComment()
			assert.NotNil(t, got)
			assert.Equal(t, tt.expected.Id, got.Id)
			assert.Equal(t, tt.expected.Body, got.Body)
			assert.Equal(t, tt.expected.CreatedAt, got.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, got.UpdatedAt)

			if tt.expected.Author != nil && got.Author != nil {
				assert.Equal(t, tt.expected.Author.Username, got.Author.Username)
			} else {
				assert.Equal(t, tt.expected.Author, got.Author)
			}

			if tt.wantError {
				t.Errorf("Expecting error, but got none.")
			} else {
				t.Logf("Test case %s passed.", tt.name)
			}
		})
	}
}
