package handler

import (
	"context"
	"database/sql"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/stretchr/testify/require"
)

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


