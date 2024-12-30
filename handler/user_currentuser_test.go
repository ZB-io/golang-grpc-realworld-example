package handler

import (
	"context"
	"fmt"
	"os"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var generateToken = auth.GenerateToken
var getUserID = auth.GetUserIDtype ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

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

type UserStore struct {
	db *gorm.DB
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
func TestCurrentUser(t *testing.T) {
	mockUserID := uint(1)
	mockToken := "mockedToken"
	mockUser := &model.User{
		ID:       mockUserID,
		Email:    "john.doe@example.com",
		Username: "johndoe",
		Bio:      "Bio of John Doe",
		Image:    "image.png",
	}

	getUserID = func(ctx context.Context) (uint, error) {
		return mockUserID, nil
	}
	generateToken = func(id uint) (string, error) {
		return mockToken, nil
	}

	tests := []struct {
		name         string
		setupMocks   func(*Handler, sqlmock.Sqlmock)
		expectedResp *pb.UserResponse
		expectedErr  error
	}{
		{
			name: "Scenario 1: Successfully Retrieve Current User",
			setupMocks: func(h *Handler, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM users WHERE id=?").
					WithArgs(mockUserID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "bio", "image"}).
						AddRow(mockUser.ID, mockUser.Email, mockUser.Username, mockUser.Bio, mockUser.Image))
			},
			expectedResp: &pb.UserResponse{User: mockUser.ProtoUser(mockToken)},
			expectedErr:  nil,
		},
		{
			name: "Scenario 2: Unauthenticated User Request",
			setupMocks: func(h *Handler, mock sqlmock.Sqlmock) {
				getUserID = func(ctx context.Context) (uint, error) {
					return 0, fmt.Errorf("unauthenticated")
				}
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Scenario 3: User Not Found",
			setupMocks: func(h *Handler, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM users WHERE id=?").
					WithArgs(mockUserID).
					WillReturnError(fmt.Errorf("not found"))

				getUserID = func(ctx context.Context) (uint, error) {
					return mockUserID, nil
				}
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Scenario 4: Token Generation Failure",
			setupMocks: func(h *Handler, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM users WHERE id=?").
					WithArgs(mockUserID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "bio", "image"}).
						AddRow(mockUser.ID, mockUser.Email, mockUser.Username, mockUser.Bio, mockUser.Image))

				generateToken = func(id uint) (string, error) {
					return "", fmt.Errorf("token generation failed")
				}
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error occurred while opening stub database: %v", err)
			}
			defer db.Close()

			userStore := store.NewUserStore(db)
			logger := zerolog.New(os.Stdout)
			handler := &Handler{
				logger: &logger,
				us:     userStore,
			}

			tt.setupMocks(handler, mock)

			resp, err := handler.CurrentUser(context.Background(), &pb.Empty{})

			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)

			if err != nil {
				t.Logf("expected error: %v", err)
			} else {
				t.Log("test passed without error")
			}
		})
	}
}
