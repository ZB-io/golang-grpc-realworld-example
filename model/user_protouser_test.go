package model

import (
	"testing"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestProtoUser(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		token    string
		expected *pb.User
	}{
		{
			name: "Successful conversion",
			user: User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			token: "test_token",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "test_token",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
		},
		{
			name: "Conversion with empty fields",
			user: User{
				Email:    "empty@example.com",
				Username: "emptyuser",
			},
			token: "empty_token",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "empty_token",
				Username: "emptyuser",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "Conversion with maximum length values",
			user: User{
				Email:    "verylongemail@verylongdomain.com",
				Username: "verylongusernamewithalotofcharacters",
				Bio:      "This is a very long bio that contains a lot of characters to test the maximum length handling of the ProtoUser method. It should be able to handle long strings without truncation.",
				Image:    "https://example.com/very/long/image/path/with/many/subdirectories/and/a/long/filename.jpg",
			},
			token: "very_long_token_string_for_testing_maximum_length_handling",
			expected: &pb.User{
				Email:    "verylongemail@verylongdomain.com",
				Token:    "very_long_token_string_for_testing_maximum_length_handling",
				Username: "verylongusernamewithalotofcharacters",
				Bio:      "This is a very long bio that contains a lot of characters to test the maximum length handling of the ProtoUser method. It should be able to handle long strings without truncation.",
				Image:    "https://example.com/very/long/image/path/with/many/subdirectories/and/a/long/filename.jpg",
			},
		},
		{
			name: "Conversion with special characters",
			user: User{
				Email:    "special!#$%&'*+-/=?^_`{|}~@example.com",
				Username: "user@name",
				Bio:      "Bio with unicode: こんにちは 你好 😊 and HTML: <script>alert('test')</script>",
				Image:    "https://example.com/image?param1=value1&param2=value2",
			},
			token: "token!@#$%^&*()",
			expected: &pb.User{
				Email:    "special!#$%&'*+-/=?^_`{|}~@example.com",
				Token:    "token!@#$%^&*()",
				Username: "user@name",
				Bio:      "Bio with unicode: こんにちは 你好 😊 and HTML: <script>alert('test')</script>",
				Image:    "https://example.com/image?param1=value1&param2=value2",
			},
		},
		{
			name: "Token handling",
			user: User{
				Email:    "token@example.com",
				Username: "tokenuser",
			},
			token: "unique_test_token_123",
			expected: &pb.User{
				Email:    "token@example.com",
				Token:    "unique_test_token_123",
				Username: "tokenuser",
				Bio:      "",
				Image:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.ProtoUser(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}
