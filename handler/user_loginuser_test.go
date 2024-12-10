package handler

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginUser(t *testing.T) {
	// Mock Logger
	mockLogger := new(MockLogger) // assume MockLogger implements logger interface

	// Table-driven test cases
	testCases := []struct {
		name        string
		req         *pb.LoginUserRequest
		mockSetup   func(*testing.T, *pb.LoginUserRequest, sqlmock.Sqlmock)
		assertFunc  func(*testing.T, *pb.UserResponse, error)
		expectLog   string
	}{
		{
			name: "Successful User Login",
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			mockSetup: func(t *testing.T, req *pb.LoginUserRequest, mock sqlmock.Sqlmock) {
				user := &model.User{
					ID:       1,
					Email:    req.GetUser().GetEmail(),
					PasswordHash: model.HashPassword(req.GetUser().GetPassword()), // Assuming a utility function to hash password
				}
				mock.ExpectQuery("^SELECT\\s+.*\\s+FROM\\s+users\\s+WHERE\\s+email=\\$1$").
					WithArgs(req.GetUser().GetEmail()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
						AddRow(user.ID, user.Email, user.PasswordHash))

				token, err := auth.GenerateToken(user.ID)
				if err != nil {
					t.Fatalf("Failed to generate token during setup: %v", err)
				}
				user.Token = token
			},
			assertFunc: func(t *testing.T, res *pb.UserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.User.Email, "test@example.com")
				assert.NotEmpty(t, res.User.Token)
			},
			expectLog: "login user",
		},
		{
			name: "Invalid Email Login Attempt",
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "invalid@example.com",
					Password: "validpassword",
				},
			},
			mockSetup: func(t *testing.T, req *pb.LoginUserRequest, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT\\s+.*\\s+FROM\\s+users\\s+WHERE\\s+email=\\$1$").
					WithArgs(req.GetUser().GetEmail()).
					WillReturnError(fmt.Errorf("email not found"))
			},
			assertFunc: func(t *testing.T, res *pb.UserResponse, err error) {
				assert.Nil(t, res)
				assert.EqualError(t, err, status.Error(codes.InvalidArgument, "invalid email or password").Error())
			},
			expectLog: "failed to login due to wrong email",
		},
		{
			name: "Incorrect Password Attempt",
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
			},
			mockSetup: func(t *testing.T, req *pb.LoginUserRequest, mock sqlmock.Sqlmock) {
				user := &model.User{
					ID:           1,
					Email:        req.GetUser().GetEmail(),
					PasswordHash: model.HashPassword("correctpassword"), // Correct password hash for comparison
				}
				mock.ExpectQuery("^SELECT\\s+.*\\s+FROM\\s+users\\s+WHERE\\s+email=\\$1$").
					WithArgs(req.GetUser().GetEmail()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
						AddRow(user.ID, user.Email, user.PasswordHash))
			},
			assertFunc: func(t *testing.T, res *pb.UserResponse, err error) {
				assert.Nil(t, res)
				assert.EqualError(t, err, status.Error(codes.InvalidArgument, "invalid email or password").Error())
			},
			expectLog: "failed to login due to receive wrong password",
		},
		{
			name: "Token Generation Failure",
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			mockSetup: func(t *testing.T, req *pb.LoginUserRequest, mock sqlmock.Sqlmock) {
				user := &model.User{
					ID:           1,
					Email:        req.GetUser().GetEmail(),
					PasswordHash: model.HashPassword(req.GetUser().GetPassword()),
				}
				mock.ExpectQuery("^SELECT\\s+.*\\s+FROM\\s+users\\s+WHERE\\s+email=\\$1$").
					WithArgs(req.GetUser().GetEmail()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
						AddRow(user.ID, user.Email, user.PasswordHash))

				// Override GenerateToken to simulate failure
				auth.GenerateToken = func(id uint) (string, error) {
					return "", fmt.Errorf("token generation error")
				}
			},
			assertFunc: func(t *testing.T, res *pb.UserResponse, err error) {
				assert.Nil(t, res)
				assert.EqualError(t, err, status.Error(codes.Aborted, "internal server error").Error())
			},
			expectLog: "Failed to create token",
		},
		{
			name: "Logging on User Login Attempts",
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "correctpassword",
				},
			},
			mockSetup: func(t *testing.T, req *pb.LoginUserRequest, mock sqlmock.Sqlmock) {
				user := &model.User{
					ID:           1,
					Email:        req.GetUser().GetEmail(),
					PasswordHash: model.HashPassword(req.GetUser().GetPassword()),
				}
				mock.ExpectQuery("^SELECT\\s+.*\\s+FROM\\s+users\\s+WHERE\\s+email=\\$1$").
					WithArgs(req.GetUser().GetEmail()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
						AddRow(user.ID, user.Email, user.PasswordHash))

				token, err := auth.GenerateToken(user.ID)
				if err != nil {
					t.Fatalf("Failed to generate token during setup: %v", err)
				}
				user.Token = token
			},
			assertFunc: func(t *testing.T, res *pb.UserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.User.Email, "test@example.com")
				assert.NotEmpty(t, res.User.Token)
			},
			expectLog: "login user", // assert that logging occurred
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Start sqlmock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %s", err)
			}
			defer db.Close()

			// Setup mock environment
			tc.mockSetup(t, tc.req, mock)

			// Redirect and buffer log output to capture it
			var logOutput strings.Builder
			mockLogger.OutputFunc = func(calldepth int, s string) error {
				_, err := logOutput.WriteString(s)
				return err
			}

			// Initialize handler
			handler := &Handler{
				us:     model.NewUserStore(db), // assuming NewUserStore is a valid constructor
				logger: mockLogger,
			}

			// Act
			res, err := handler.LoginUser(context.Background(), tc.req)

			// Assert
			tc.assertFunc(t, res, err)
			if tc.expectLog != "" {
				assert.Contains(t, logOutput.String(), tc.expectLog)
			}
		})
	}
}
