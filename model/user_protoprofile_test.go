package model

import (
	"testing"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestProtoProfile(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		following bool
		want      *pb.Profile
	}{
		{
			name: "Basic Profile Creation",
			user: User{
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "testuser",
				Bio:       "Test bio",
				Image:     "https://example.com/image.jpg",
				Following: true,
			},
		},
		{
			name: "Profile Creation with Empty Fields",
			user: User{
				Username: "",
				Bio:      "",
				Image:    "",
			},
			following: false,
			want: &pb.Profile{
				Username:  "",
				Bio:       "",
				Image:     "",
				Following: false,
			},
		},
		{
			name: "Following Status True",
			user: User{
				Username: "followeduser",
				Bio:      "Followed user bio",
				Image:    "https://example.com/followed.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "followeduser",
				Bio:       "Followed user bio",
				Image:     "https://example.com/followed.jpg",
				Following: true,
			},
		},
		{
			name: "Following Status False",
			user: User{
				Username: "unfolloweduser",
				Bio:      "Unfollowed user bio",
				Image:    "https://example.com/unfollowed.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "unfolloweduser",
				Bio:       "Unfollowed user bio",
				Image:     "https://example.com/unfollowed.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Creation with Maximum Length Fields",
			user: User{
				Username: "maxlengthusername",
				Bio:      "This is a very long bio that reaches the maximum allowed length for testing purposes.",
				Image:    "https://example.com/very/long/image/url/that/reaches/maximum/allowed/length/for/testing/purposes.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "maxlengthusername",
				Bio:       "This is a very long bio that reaches the maximum allowed length for testing purposes.",
				Image:     "https://example.com/very/long/image/url/that/reaches/maximum/allowed/length/for/testing/purposes.jpg",
				Following: true,
			},
		},
		{
			name: "Profile Creation with Special Characters",
			user: User{
				Username: "special_user_🚀",
				Bio:      "Bio with special chars: ©®™°C²♥★☆☂☀☁☔",
				Image:    "https://example.com/image_🖼️.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "special_user_🚀",
				Bio:       "Bio with special chars: ©®™°C²♥★☆☂☀☁☔",
				Image:     "https://example.com/image_🖼️.jpg",
				Following: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.ProtoProfile(tt.following)
			assert.Equal(t, tt.want, got)
		})
	}
}
