package handler

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/logger"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		logger: logger.New(),
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


