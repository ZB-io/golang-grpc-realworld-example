package model

import (
	"reflect"
	"strings"
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)






func TestProtoUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		token       string
		expected    *pb.User
		expectPanic bool
	}{
		{
			name: "Convert User Struct to ProtoUser with Valid Data",
			user: &User{
				Email:    "test@example.com",
				Username: "testUser",
				Bio:      "This is a test bio.",
				Image:    "https://example.com/image.jpg",
			},
			token: "validToken123",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "validToken123",
				Username: "testUser",
				Bio:      "This is a test bio.",
				Image:    "https://example.com/image.jpg",
			},
			expectPanic: false,
		},
		{
			name: "Handle Empty Token while Converting User Struct",
			user: &User{
				Email:    "test@example.com",
				Username: "testUser",
				Bio:      "This is a test bio.",
				Image:    "https://example.com/image.jpg",
			},
			token: "",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "",
				Username: "testUser",
				Bio:      "This is a test bio.",
				Image:    "https://example.com/image.jpg",
			},
			expectPanic: false,
		},
		{
			name: "Convert User Struct with Empty User Fields",
			user: &User{
				Email:    "",
				Username: "",
				Bio:      "",
				Image:    "",
			},
			token: "validToken123",
			expected: &pb.User{
				Email:    "",
				Token:    "validToken123",
				Username: "",
				Bio:      "",
				Image:    "",
			},
			expectPanic: false,
		},
		{
			name: "Convert User Struct with Special Characters in Fields",
			user: &User{
				Email:    "test+special@example.com",
				Username: "test!@#$%^&*",
				Bio:      "Bio with special <>?/:;",
				Image:    "https://example.com/image?param=value",
			},
			token: "specialToken@123",
			expected: &pb.User{
				Email:    "test+special@example.com",
				Token:    "specialToken@123",
				Username: "test!@#$%^&*",
				Bio:      "Bio with special <>?/:;",
				Image:    "https://example.com/image?param=value",
			},
			expectPanic: false,
		},
		{
			name: "Convert User Struct with Long String Fields",
			user: &User{
				Email:    strings.Repeat("long", 1000) + "@example.com",
				Username: strings.Repeat("long", 1000),
				Bio:      strings.Repeat("long", 1000),
				Image:    "https://example.com/" + strings.Repeat("long", 995),
			},
			token: "validToken123",
			expected: &pb.User{
				Email:    strings.Repeat("long", 1000) + "@example.com",
				Token:    "validToken123",
				Username: strings.Repeat("long", 1000),
				Bio:      strings.Repeat("long", 1000),
				Image:    "https://example.com/" + strings.Repeat("long", 995),
			},
			expectPanic: false,
		},
		{
			name: "Validate Field Mapping Consistency Across Multiple Calls",
			user: &User{
				Email:    "consistent@example.com",
				Username: "consistentUser",
				Bio:      "Some consistent bio.",
				Image:    "https://example.com/image.jpg",
			},
			token: "consistentToken",
			expected: &pb.User{
				Email:    "consistent@example.com",
				Token:    "consistentToken",
				Username: "consistentUser",
				Bio:      "Some consistent bio.",
				Image:    "https://example.com/image.jpg",
			},
			expectPanic: false,
		},
		{
			name:        "Pass Nil User Struct to ProtoUser",
			user:        nil,
			token:       "validToken123",
			expected:    nil,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected ProtoUser to panic for test case: %s", tt.name)
					}
				}()
			}

			result := tt.user.ProtoUser(tt.token)

			if !tt.expectPanic {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Test case failed: %s. Expected: %+v, Got: %+v", tt.name, tt.expected, result)
				} else {
					t.Logf("Test case passed: %s", tt.name)
				}
			}
		})
	}
}


