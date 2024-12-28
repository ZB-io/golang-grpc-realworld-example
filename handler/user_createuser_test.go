package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type MockUserService struct {
	mockedUser *model.User
	ctrl       *gomock.Controller
}


func (m *MockUserService) Create(user *model.User) error {

	return nil
}
func (m *MockUserService) HashPassword() error {

	return nil
}
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	return &MockUserService{
		ctrl: ctrl,
	}
}
func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUserService := NewMockUserService(ctrl)

	testLogger := model.TestLogger{}

	h := handler.Handler{
		Us:     mockedUserService,
		Logger: testLogger,
	}

	tests := []struct {
		name       string
		req        *pb.CreateUserRequest
		setupMocks func()
		assert     func(*pb.UserResponse, error)
	}{
		{
			name: "Successful User Creation with Valid Input",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "validuser",
					Email:    "valid@example.com",
					Password: "validpassword",
				},
			},
			setupMocks: func() {
				mockedUser := model.User{
					Username: "validuser",
					Email:    "valid@example.com",
					Password: "validpassword",
				}
				mockedUserService.EXPECT().Validate().Return(nil)
				mockedUserService.EXPECT().HashPassword().Return(nil)
				mockedUserService.EXPECT().Create(gomock.Any()).Return(nil)
				auth.EXPECT().GenerateToken(mockedUser.ID).Return("sometoken", nil)
			},
			assert: func(resp *pb.UserResponse, err error) {
				assert.NotNil(t, resp)
				assert.NoError(t, err)
				assert.Equal(t, "sometoken", resp.User.Token)
			},
		},
		{
			name: "Validation Failure on User Input",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "",
					Email:    "invalidemail.com",
					Password: "pass",
				},
			},
			setupMocks: func() {
				mockedUserService.EXPECT().Validate().Return(errors.New("invalid input"))
			},
			assert: func(resp *pb.UserResponse, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Password Hashing Failure",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "validuser",
					Email:    "valid@example.com",
					Password: "validpassword",
				},
			},
			setupMocks: func() {
				mockedUserService.EXPECT().Validate().Return(nil)
				mockedUserService.EXPECT().HashPassword().Return(errors.New("hashing error"))
			},
			assert: func(resp *pb.UserResponse, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, codes.Aborted, status.Code(err))
			},
		},
		{
			name: "User Creation Failure at Repository Level",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "validuser",
					Email:    "valid@example.com",
					Password: "validpassword",
				},
			},
			setupMocks: func() {
				mockedUserService.EXPECT().Validate().Return(nil)
				mockedUserService.EXPECT().HashPassword().Return(nil)
				mockedUserService.EXPECT().Create(gomock.Any()).Return(errors.New("creation error"))
			},
			assert: func(resp *pb.UserResponse, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, codes.Canceled, status.Code(err))
			},
		},
		{
			name: "Token Generation Failure",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "validuser",
					Email:    "valid@example.com",
					Password: "validpassword",
				},
			},
			setupMocks: func() {
				mockedUserService.EXPECT().Validate().Return(nil)
				mockedUserService.EXPECT().HashPassword().Return(nil)
				mockedUserService.EXPECT().Create(gomock.Any()).Return(nil)
				auth.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("token error"))
			},
			assert: func(resp *pb.UserResponse, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, codes.Aborted, status.Code(err))
			},
		},
		{
			name: "Missing or Nil Request",
			req:  nil,
			setupMocks: func() {

			},
			assert: func(resp *pb.UserResponse, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "nil")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			resp, err := h.CreateUser(context.Background(), tc.req)
			tc.assert(resp, err)
		})
	}
}
func (m *MockUserService) Validate() error {

	return nil
}


