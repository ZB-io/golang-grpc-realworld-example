package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)






func TestProtoProfile(t *testing.T) {
	type testCase struct {
		description    string
		user           User
		following      bool
		expectedResult *pb.Profile
	}

	testCases := []testCase{
		{
			description: "Successful Profile Generation",
			user: User{
				Username: "john_doe",
				Bio:      "Software Developer",
				Image:    "http://example.com/image.jpg",
			},
			following: true,
			expectedResult: &pb.Profile{
				Username:  "john_doe",
				Bio:       "Software Developer",
				Image:     "http://example.com/image.jpg",
				Following: true,
			},
		},
		{
			description: "Empty User Profile",
			user: User{
				Username: "",
				Bio:      "",
				Image:    "",
			},
			following: false,
			expectedResult: &pb.Profile{
				Username:  "",
				Bio:       "",
				Image:     "",
				Following: false,
			},
		},

		{
			description: "Following Status Toggle",
			user: User{
				Username: "user",
				Bio:      "A regular user",
				Image:    "http://example.com/user.jpg",
			},
			following: true,
			expectedResult: &pb.Profile{
				Username:  "user",
				Bio:       "A regular user",
				Image:     "http://example.com/user.jpg",
				Following: true,
			},
		},
		{
			description: "Null Image Field",
			user: User{
				Username: "no_image_user",
				Bio:      "No image URL provided",
				Image:    "",
			},
			following: false,
			expectedResult: &pb.Profile{
				Username:  "no_image_user",
				Bio:       "No image URL provided",
				Image:     "",
				Following: false,
			},
		},
	}

	for i, tc := range testCases {
		t.Logf("Running test case %d: %s", i+1, tc.description)

		result := tc.user.ProtoProfile(tc.following)

		if result.Username != tc.expectedResult.Username {
			t.Errorf("Test case %d failed, expected Username: %v, got: %v", i+1, tc.expectedResult.Username, result.Username)
		}
		if result.Bio != tc.expectedResult.Bio {
			t.Errorf("Test case %d failed, expected Bio: %v, got: %v", i+1, tc.expectedResult.Bio, result.Bio)
		}
		if result.Image != tc.expectedResult.Image {
			t.Errorf("Test case %d failed, expected Image: %v, got: %v", i+1, tc.expectedResult.Image, result.Image)
		}
		if result.Following != tc.expectedResult.Following {
			t.Errorf("Test case %d failed, expected Following: %v, got: %v", i+1, tc.expectedResult.Following, result.Following)
		}

		t.Logf("Test case %d succeeded", i+1)
	}
}
