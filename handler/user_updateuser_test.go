package handler

import (
	"context"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type UserStore struct {
	db *gorm.DB
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestHandlerUpdateUser(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(us *store.UserStore, ctx context.Context)
		req          *pb.UpdateUserRequest
		wantResponse *pb.UserResponse
		wantErr      error
	}{
		{
			name: "Successful User Update",
			setupMocks: func(us *store.UserStore, ctx context.Context) {

				ctx = context.WithValue(ctx, contextKey("userID"), uint(1))

				us.On("GetByID", mock.Anything).Return(&model.User{
					ID:       1,
					Username: "old_user",
					Email:    "old_email@example.com",
					Password: "old_password",
				}, nil)

				us.On("Update", mock.Anything).Return(nil)

				auth.GenerateToken = func(id uint) (string, error) {
					return "newToken", nil
				}
			},
			req: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "new_user",
					Email:    "new_email@example.com",
					Password: "new_password",
					Bio:      "new bio",
					Image:    "new_image_url",
				},
			},
			wantResponse: &pb.UserResponse{
				User: &pb.User{
					Username: "new_user",
					Email:    "new_email@example.com",
					Token:    "newToken",
					Bio:      "new bio",
					Image:    "new_image_url",
				},
			},
			wantErr: nil,
		},
		{
			name: "Unauthenticated User Request",
			setupMocks: func(us *store.UserStore, ctx context.Context) {

				ctx = context.WithValue(ctx, contextKey("userID"), uint(0))
			},
			req: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{},
			},
			wantResponse: nil,
			wantErr:      status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "User Not Found in Database",
			setupMocks: func(us *store.UserStore, ctx context.Context) {

				ctx = context.WithValue(ctx, contextKey("userID"), uint(1))

				us.On("GetByID", mock.Anything).Return(nil, status.Error(codes.NotFound, "not user found"))
			},
			req: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{},
			},
			wantResponse: nil,
			wantErr:      status.Error(codes.NotFound, "not user found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			us := store.UserStore{db: mock}
			logger := zerolog.New(os.Stdout)
			handler := &Handler{logger: &logger, us: &us}

			ctx := context.Background()
			tt.setupMocks(&us, ctx)

			resp, err := handler.UpdateUser(ctx, tt.req)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.Equal(t, tt.wantResponse.User.Username, resp.User.Username)
				assert.Equal(t, tt.wantResponse.User.Email, resp.User.Email)
			}
		})
	}
}
