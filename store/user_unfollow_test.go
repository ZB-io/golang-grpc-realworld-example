package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreUnfollow(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		setup   func(a *model.User, b *model.User) error
		a       *model.User
		b       *model.User
		wantErr bool
	}{
		{
			name: "Successful Unfollow Operation",
			setup: func(a *model.User, b *model.User) error {
				a.Follows = append(a.Follows, *b)
				return nil
			},
			a:       &model.User{Username: "userA", Email: "userA@gmail.com"},
			b:       &model.User{Username: "userB", Email: "userB@gmail.com"},
			wantErr: false,
		},
		{
			name:    "Unfollow Non-followed User",
			setup:   func(a *model.User, b *model.User) error { return nil },
			a:       &model.User{Username: "userA", Email: "userA@gmail.com"},
			b:       &model.User{Username: "userB", Email: "userB@gmail.com"},
			wantErr: false,
		},
		{
			name:  "Unfollow with Non-existent User",
			setup: func(a *model.User, b *model.User) error { return nil },
			a:     &model.User{Username: "userA", Email: "userA@gmail.com"},
			b:     nil,
			wantErr: true,
		},
		{
			name:  "Unfollow with Nil User",
			setup: func(a *model.User, b *model.User) error { return nil },
			a:     nil,
			b:     &model.User{Username: "userB", Email: "userB@gmail.com"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replace with your actual DB instance
			db := &gorm.DB{}
			s := &UserStore{db: db}

			// Setup initial state
			err := tt.setup(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Failed to setup initial state: %v", err)
			}

			// Act
			err = s.Unfollow(tt.a, tt.b)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				if tt.a != nil && tt.b != nil {
					for _, user := range tt.a.Follows {
						if user.Username == tt.b.Username {
							t.Errorf("User A still follows User B")
						}
					}
				}
			}
		})
	}
}
