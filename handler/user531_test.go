package handler

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)


/*
ROOST_METHOD_HASH=CurrentUser_e3fa631d55
ROOST_METHOD_SIG_HASH=CurrentUser_29413339e9
*/

func TestHandlerCurrentUser(t *testing.T) {
	t.Run("Scenario 1: Valid User Retrieval", func(t *testing.T) {
		h, mockAuth, _, mockDB := setupTestHandler(t)
		defer tearDownTestHandler(mockDB)

		mockAuth.On("GetUserID", mock.Anything).Return("valid-user-id", nil)

		expectedUser := &model.User{ID: "valid-user-id", Name: "John Doe"}
		mockDB.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
			WithArgs("valid-user-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUser.ID, expectedUser.Name))

		expectedToken := "valid-token"
		mockAuth.On("GenerateToken", expectedUser.ID).Return(expectedToken, nil)

		ctx := context.WithValue(context.Background(), auth.UserContextKey, "valid-user-id")
		resp, err := h.CurrentUser(ctx, &pb.Empty{})

		require.NoError(t, err)
		require.Equal(t, expectedUser.ID, resp.User.ID)
		require.Equal(t, expectedUser.Name, resp.User.Name)
		t.Log("Success: User retrieved and token generated correctly.")
	})

	t.Run("Scenario 2: User ID Not Found in Context", func(t *testing.T) {
		h, mockAuth, _, mockDB := setupTestHandler(t)
		defer tearDownTestHandler(mockDB)

		mockAuth.On("GetUserID", mock.Anything).Return("", status.Errorf(codes.Unauthenticated, "unauthenticated"))

		ctx := context.Background()
		_, err := h.CurrentUser(ctx, &pb.Empty{})

		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, status.Code(err))
		t.Log("Success: Correctly identified unauthenticated user.")
	})

	t.Run("Scenario 3: User Not Found in Database", func(t *testing.T) {
		h, mockAuth, _, mockDB := setupTestHandler(t)
		defer tearDownTestHandler(mockDB)

		mockAuth.On("GetUserID", mock.Anything).Return("valid-user-id", nil)

		mockDB.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
			WithArgs("valid-user-id").
			WillReturnError(sql.ErrNoRows)

		ctx := context.WithValue(context.Background(), auth.UserContextKey, "valid-user-id")
		_, err := h.CurrentUser(ctx, &pb.Empty{})

		require.Error(t, err)
		require.Equal(t, codes.NotFound, status.Code(err))
		t.Log("Success: Correctly handled user not found in database.")
	})

	t.Run("Scenario 4: Token Generation Fails", func(t *testing.T) {
		h, mockAuth, _, mockDB := setupTestHandler(t)
		defer tearDownTestHandler(mockDB)

		mockAuth.On("GetUserID", mock.Anything).Return("valid-user-id", nil)

		expectedUser := &model.User{ID: "valid-user-id", Name: "John Doe"}
		mockDB.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
			WithArgs("valid-user-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUser.ID, expectedUser.Name))

		mockAuth.On("GenerateToken", expectedUser.ID).Return("", status.Error(codes.Aborted, "internal error"))

		ctx := context.WithValue(context.Background(), auth.UserContextKey, "valid-user-id")
		_, err := h.CurrentUser(ctx, &pb.Empty{})

		require.Error(t, err)
		require.Equal(t, codes.Aborted, status.Code(err))
		t.Log("Success: Correctly handled token generation failure.")
	})

	t.Run("Scenario 5: Nil Request Handling", func(t *testing.T) {
		h, mockAuth, _, mockDB := setupTestHandler(t)
		defer tearDownTestHandler(mockDB)

		mockAuth.On("GetUserID", mock.Anything).Return("valid-user-id", nil)

		expectedUser := &model.User{ID: "valid-user-id", Name: "John Doe"}
		mockDB.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
			WithArgs("valid-user-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUser.ID, expectedUser.Name))

		expectedToken := "valid-token"
		mockAuth.On("GenerateToken", expectedUser.ID).Return(expectedToken, nil)

		ctx := context.WithValue(context.Background(), auth.UserContextKey, "valid-user-id")
		resp, err := h.CurrentUser(ctx, nil)

		require.NoError(t, err)
		require.Equal(t, expectedUser.ID, resp.User.ID)
		require.Equal(t, expectedUser.Name, resp.User.Name)
		t.Log("Success: Function handled nil request safely.")
	})
}


/*
ROOST_METHOD_HASH=LoginUser_079a321a92
ROOST_METHOD_SIG_HASH=LoginUser_e7df23a6bd
*/


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

	authMock := auth.NewMockAuth(ctrl)

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
				authMock.EXPECT().GenerateToken(user.ID).Return("validtoken", nil)

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
					WillReturnError(sql.ErrNoRows)
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
				authMock.EXPECT().GenerateToken(user.ID).Return("", status.Error(codes.Aborted, "token generation error"))

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

			h := &Handler{
				logger: model.NewMockLogger(ctrl),
				auth:   authMock,
			}

			tc.mockSetup(h, dbMock)

			resp, err := h.LoginUser(context.Background(), tc.req)

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


/*
ROOST_METHOD_HASH=CreateUser_f2f8a1c84a
ROOST_METHOD_SIG_HASH=CreateUser_a3af3934da
*/


func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUserService := NewMockUserService(ctrl)

	h := &Handler{
		us:     mockedUserService,
		logger: model.NewMockLogger(ctrl),
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
				mockedUserService.EXPECT().Create(gomock.Any()).Return(nil)
				auth.EXPECT().GenerateToken(gomock.Any()).Return("sometoken", nil)
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


/*
ROOST_METHOD_HASH=UpdateUser_6fa4ecf979
ROOST_METHOD_SIG_HASH=UpdateUser_883937d25b
*/

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := auth.NewMockAuth(ctrl)

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, %s", err)
	}
	defer mockDB.Close()

	h := &Handler{
		logger: model.NewMockLogger(ctrl),
		us:     model.NewSqlUserStore(mockDB),
	}

	type args struct {
		ctx context.Context
		req *pb.UpdateUserRequest
	}

	tests := []struct {
		name      string
		prepare   func()
		args      args
		wantErr   codes.Code
		wantToken string
	}{
		{
			name: "Successful User Update with All Fields Provided",
			prepare: func() {

				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil).AnyTimes()

				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).AddRow(1, "testuser", "test@test.com", "bio", "image"))

				mock.ExpectExec("UPDATE users SET (.+) WHERE id=\\?").
					WithArgs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mockAuth.EXPECT().GenerateToken(gomock.Any()).Return("new-token", nil)
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Username: "newname",
						Email:    "newemail@test.com",
						Password: "newpass",
						Image:    "newimage",
						Bio:      "newbio",
					},
				},
			},
			wantErr:   codes.OK,
			wantToken: "new-token",
		},
		{
			name: "Unauthenticated User",
			prepare: func() {

				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(0, errors.New("unauthenticated")).AnyTimes()
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{},
			},
			wantErr: codes.Unauthenticated,
		},
		{
			name: "User Not Found After Authentication",
			prepare: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\?").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{},
			},
			wantErr: codes.NotFound,
		},
		{
			name: "Validation Error on User Fields",
			prepare: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).AddRow(1, "testuser", "test@test.com", "bio", "image"))
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Email: "invalid-email",
					},
				},
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "Password Hashing Failure",
			prepare: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).AddRow(1, "testuser", "test@test.com", "bio", "image"))

				h.us.EXPECT().HashPassword(gomock.Any()).Return(errors.New("hash error"))
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Password: "newpass",
					},
				},
			},
			wantErr: codes.Aborted,
		},
		{
			name: "Internal Server Error During Update",
			prepare: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).AddRow(1, "testuser", "test@test.com", "bio", "image"))

				mock.ExpectExec("UPDATE users SET (.+) WHERE id=\\?").
					WithArgs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					WillReturnError(errors.New("db error"))
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Username: "newname",
					},
				},
			},
			wantErr: codes.Internal,
		},
		{
			name: "Token Generation Failure",
			prepare: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).AddRow(1, "testuser", "test@test.com", "bio", "image"))

				mock.ExpectExec("UPDATE users SET (.+) WHERE id=\\?").
					WithArgs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mockAuth.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("token generation error"))
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Username: "newname",
					},
				},
			},
			wantErr: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			resp, err := h.UpdateUser(tt.args.ctx, tt.args.req)

			if tt.wantErr == codes.OK {
				assert.Equal(t, tt.wantToken, resp.User.Token)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErr, st.Code())
			}
		})
	}
}
