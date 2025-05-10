package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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
func TestUserProtoUser(t *testing.T) {
	type testCase struct {
		name     string
		user     User
		token    string
		expected *pb.User
	}

	var cases = []testCase{
		{
			name: "Convert User to ProtoUser with Valid Token",
			user: User{
				Username: "test_user",
				Email:    "test@example.com",
				Bio:      "A test user bio",
				Image:    "http://example.com/image.png",
			},
			token: "valid-token",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "valid-token",
				Username: "test_user",
				Bio:      "A test user bio",
				Image:    "http://example.com/image.png",
			},
		},
		{
			name: "ProtoUser Token Handling",
			user: User{
				Username: "sample_user",
				Email:    "sample@example.com",
				Bio:      "A sample user bio",
				Image:    "http://example.com/sample.png",
			},
			token: "specific-token",
			expected: &pb.User{
				Email:    "sample@example.com",
				Token:    "specific-token",
				Username: "sample_user",
				Bio:      "A sample user bio",
				Image:    "http://example.com/sample.png",
			},
		},
		{
			name: "Handle Empty User Fields Gracefully",
			user: User{
				Username: "",
				Email:    "",
				Bio:      "",
				Image:    "",
			},
			token: "valid-token",
			expected: &pb.User{
				Email:    "",
				Token:    "valid-token",
				Username: "",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "Large Data Size Handling",
			user: User{
				Username: "longusername" + string(make([]byte, 1000)),
				Email:    "longemail@example.com",
				Bio:      "A long bio" + string(make([]byte, 1000)),
				Image:    "http://example.com/largeimage.png",
			},
			token: "large-data-token",
			expected: &pb.User{
				Email:    "longemail@example.com",
				Token:    "large-data-token",
				Username: "longusername" + string(make([]byte, 1000)),
				Bio:      "A long bio" + string(make([]byte, 1000)),
				Image:    "http://example.com/largeimage.png",
			},
		},
		{
			name: "Validate User Image Field Conversion",
			user: User{
				Email:    "imageuser@example.com",
				Username: "imageuser",
				Bio:      "Bio with image",
				Image:    "http://example.com/imageuser.png",
			},
			token: "image-token",
			expected: &pb.User{
				Email:    "imageuser@example.com",
				Token:    "image-token",
				Username: "imageuser",
				Bio:      "Bio with image",
				Image:    "http://example.com/imageuser.png",
			},
		},
		{
			name: "Null or Nil Token Argument",
			user: User{
				Email:    "email@example.com",
				Username: "emptytokenuser",
				Bio:      "bio",
				Image:    "http://example.com/image.png",
			},
			token: "",
			expected: &pb.User{
				Email:    "email@example.com",
				Token:    "",
				Username: "emptytokenuser",
				Bio:      "bio",
				Image:    "http://example.com/image.png",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.user.ProtoUser(c.token)
			if actual.Email != c.expected.Email || actual.Token != c.expected.Token ||
				actual.Username != c.expected.Username || actual.Bio != c.expected.Bio ||
				actual.Image != c.expected.Image {
				t.Errorf("Test %s failed: expected %+v, got %+v", c.name, c.expected, actual)
			} else {
				t.Logf("Test %s passed", c.name)
			}
		})
	}
}
