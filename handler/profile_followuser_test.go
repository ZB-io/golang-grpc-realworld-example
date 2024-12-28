package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}



type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestHandlerFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := store.NewMockUserStore(ctrl)
	logger := zerolog.New(os.Stdout)

	h := &Handler{logger: &logger, us: mockUserStore}

	t.Run("Scenario 1: User is Unauthenticated", func(t *testing.T) {
		mockCtx := context.TODO()
		mockRequest := &pb.FollowRequest{Username: "targetUser"}

		auth.GetUserID = func(ctx context.Context) (uint, error) {
			return 0, errors.New("unauthenticated")
		}

		_, err := h.FollowUser(mockCtx, mockRequest)

		if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
			t.Errorf("expected Unaunthenticated error, got %v", err)
		}
	})

	t.Run("Scenario 2: Current User Not Found", func(t *testing.T) {
		mockCtx := context.TODO()
		mockRequest := &pb.FollowRequest{Username: "targetUser"}

		auth.GetUserID = func(ctx context.Context) (uint, error) {
			return 1, nil
		}
		mockUserStore.EXPECT().GetByID(uint(1)).Return(nil, errors.New("not found"))

		_, err := h.FollowUser(mockCtx, mockRequest)

		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound || s.Message() != "user not found" {
			t.Errorf("expected NotFound error with message 'user not found', got %v", err)
		}
	})

	t.Run("Scenario 3: User Attempts to Follow Themselves", func(t *testing.T) {
		mockCtx := context.TODO()
		mockRequest := &pb.FollowRequest{Username: "currentUser"}

		auth.GetUserID = func(ctx context.Context) (uint, error) {
			return 2, nil
		}
		mockUserStore.EXPECT().GetByID(uint(2)).Return(&model.User{Username: "currentUser"}, nil)

		_, err := h.FollowUser(mockCtx, mockRequest)

		if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
			t.Errorf("expected InvalidArgument error, got %v", err)
		}
	})

	t.Run("Scenario 4: Target User Not Found", func(t *testing.T) {
		mockCtx := context.TODO()
		mockRequest := &pb.FollowRequest{Username: "targetUser"}

		auth.GetUserID = func(ctx context.Context) (uint, error) {
			return 3, nil
		}
		mockUserStore.EXPECT().GetByID(uint(3)).Return(&model.User{Username: "currentUser"}, nil)
		mockUserStore.EXPECT().GetByUsername("targetUser").Return(nil, errors.New("not found"))

		_, err := h.FollowUser(mockCtx, mockRequest)

		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
			t.Errorf("expected NotFound error, got %v", err)
		}
	})

	t.Run("Scenario 5: Successful Follow Operation", func(t *testing.T) {
		mockCtx := context.TODO()
		mockRequest := &pb.FollowRequest{Username: "targetUser"}

		auth.GetUserID = func(ctx context.Context) (uint, error) {
			return 4, nil
		}
		currentUser := &model.User{ID: 4, Username: "currentUser"}
		targetUser := &model.User{ID: 5, Username: "targetUser"}

		mockUserStore.EXPECT().GetByID(uint(4)).Return(currentUser, nil)
		mockUserStore.EXPECT().GetByUsername("targetUser").Return(targetUser, nil)
		mockUserStore.EXPECT().Follow(currentUser, targetUser).Return(nil)

		resp, err := h.FollowUser(mockCtx, mockRequest)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp == nil || resp.Profile.Username != "targetUser" || !resp.Profile.Following {
			t.Errorf("unexpected response: %v", resp)
		}
	})

	t.Log("All test scenarios executed.")
}
