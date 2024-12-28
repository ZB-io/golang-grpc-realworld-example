package handler

import (
	"context"
	"fmt"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginUser(t *testing.T) {
	type testCase struct {
		desc        string
		req         *pb.LoginUserRequest
		mockSetup   func(h *Handler, dbMock sqlmock.Sqlmock)
		expected    *pb.UserResponse
		expectError bool
		errorCode   codes.Code
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range []testCase{
		{
			desc: "Valid Email and Password",
			req: &pb.LoginUserRequest{
				User: &pb.User{
					Email:    "valid@example.com",
					Password: "correctpassword",
				},
			},
			mockSetup: func(h *Handler, dbMock sqlmock.Sqlmock) {
				user := &model.User{
					ID:       1,
					Email:    "valid@example.com",
					Password: "hashedpassword",
				}
				dbMock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs(user.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
						AddRow(user.ID, user.Email, user.Password))
				h.auth.EXPECT().GenerateToken(user.ID).Return("validtoken", nil)

				model.CheckPassword = func(hashed, password string) bool {
					return password == "correctpassword"
				}
			},
			expected:    &pb.UserResponse{User: &pb.User{Token: "validtoken"}},
			expectError: false,
		},
		{
			desc: "Invalid Email",
			req: &pb.LoginUserRequest{
				User: &pb.User{
					Email:    "nonexistent@example.com",
					Password: "somepassword",
				},
			},
			mockSetup: func(h *Handler, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs("nonexistent@example.com").
					WillReturnError(fmt.Errorf("sql: no rows in result set"))
			},
			expected:    nil,
			expectError: true,
			errorCode:   codes.InvalidArgument,
		},
		{
			desc: "Incorrect Password",
			req: &pb.LoginUserRequest{
				User: &pb.User{
					Email:    "valid@example.com",
					Password: "wrongpassword",
				},
			},
			mockSetup: func(h *Handler, dbMock sqlmock.Sqlmock) {
				user := &model.User{
					ID:    1,
					Email: "valid@example.com",
				}
				dbMock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs(user.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
						AddRow(user.ID, user.Email))

				model.CheckPassword = func(hashed, password string) bool {
					return false
				}
			},
			expected:    nil,
			expectError: true,
			errorCode:   codes.InvalidArgument,
		},
		{
			desc: "Token Generation Failure",
			req: &pb.LoginUserRequest{
				User: &pb.User{
					Email:    "valid@example.com",
					Password: "correctpassword",
				},
			},
			mockSetup: func(h *Handler, dbMock sqlmock.Sqlmock) {
				user := &model.User{
					ID:       1,
					Email:    "valid@example.com",
					Password: "hashedpassword",
				}
				dbMock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs(user.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
						AddRow(user.ID, user.Email, user.Password))
				h.auth.EXPECT().GenerateToken(user.ID).Return("", fmt.Errorf("token generation error"))

				model.CheckPassword = func(hashed, password string) bool {
					return password == "correctpassword"
				}
			},
			expected:    nil,
			expectError: true,
			errorCode:   codes.Aborted,
		},
		{
			desc:        "Nil Request Handling",
			req:         nil,
			mockSetup:   func(h *Handler, dbMock sqlmock.Sqlmock) {},
			expected:    nil,
			expectError: true,
			errorCode:   codes.InvalidArgument,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			handler := &Handler{
				auth:   auth.NewMockAuth(ctrl),
				logger: model.NewMockLogger(ctrl),
			}

			tc.mockSetup(handler, dbMock)

			resp, err := handler.LoginUser(context.Background(), tc.req)

			if tc.expectError {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, tc.errorCode, status.Code(err))
				t.Logf("Expected error for %s: %v", tc.desc, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, resp)
				t.Logf("Successfully passed %s", tc.desc)
			}

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet expectations: %s", err)
			}
		})
	}
}


