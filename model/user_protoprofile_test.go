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
func TestUserProtoProfile(t *testing.T) {

	type testCase struct {
		name      string
		user      User
		following bool
		expected  *pb.Profile
	}

	tests := []testCase{
		{
			name: "Scenario 1: Basic Profile Data",
			user: User{
				Username: "testuser",
				Bio:      "Simple bio",
				Image:    "testimage.jpg",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "testuser",
				Bio:       "Simple bio",
				Image:     "testimage.jpg",
				Following: false,
			},
		},
		{
			name: "Scenario 2: Following Status Check",
			user: User{
				Username: "testuser",
				Bio:      "Simple bio",
				Image:    "testimage.jpg",
			},
			following: true,
			expected: &pb.Profile{
				Username:  "testuser",
				Bio:       "Simple bio",
				Image:     "testimage.jpg",
				Following: true,
			},
		},
		{
			name: "Scenario 3: Handling Empty Fields",
			user: User{
				Username: "",
				Bio:      "",
				Image:    "",
			},
			following: true,
			expected: &pb.Profile{
				Username:  "",
				Bio:       "",
				Image:     "",
				Following: true,
			},
		},
		{
			name: "Scenario 4: Large Data Handling",
			user: User{
				Username: "a very long username that would be typically unrealistic",
				Bio:      "a very long bio that goes on and on and on",
				Image:    "a very long image path that doesn't make sense but is for testing",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "a very long username that would be typically unrealistic",
				Bio:       "a very long bio that goes on and on and on",
				Image:     "a very long image path that doesn't make sense but is for testing",
				Following: false,
			},
		},
		{
			name: "Scenario 5: Stress Test with Multiple Invocations",
			user: User{
				Username: "stressuser",
				Bio:      "stress bio",
				Image:    "stress.jpg",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "stressuser",
				Bio:       "stress bio",
				Image:     "stress.jpg",
				Following: false,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actualProfile := tc.user.ProtoProfile(tc.following)

			if actualProfile.GetUsername() != tc.expected.GetUsername() {
				t.Errorf("expected username: %v, got: %v", tc.expected.GetUsername(), actualProfile.GetUsername())
			}
			if actualProfile.GetBio() != tc.expected.GetBio() {
				t.Errorf("expected bio: %v, got: %v", tc.expected.GetBio(), actualProfile.GetBio())
			}
			if actualProfile.GetImage() != tc.expected.GetImage() {
				t.Errorf("expected image: %v, got: %v", tc.expected.GetImage(), actualProfile.GetImage())
			}
			if actualProfile.GetFollowing() != tc.expected.GetFollowing() {
				t.Errorf("expected following: %v, got: %v", tc.expected.GetFollowing(), actualProfile.GetFollowing())
			}

			t.Logf("Test executed: %s", tc.name)
		})
	}
}
