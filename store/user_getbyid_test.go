package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// UserStore struct
type UserStore struct {
	db *gorm.DB
}

// GetByID function
func (s *UserStore) GetByID(id uint) (*model.User, error) {
	var m model.User
	if err := s.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func TestUserStoreGetByID(t *testing.T) {
	testCases := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, id uint)
		id            uint
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Successful Retrieval of User by ID",
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"ID", "Username", "Email", "Password", "Bio", "Image"}).
					AddRow(id, "TestUser", "test@example.com", "password", "Bio", "image.png")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = \\?$").
					WithArgs(id).
					WillReturnRows(rows)
			},
			id:            1,
			expectedUser:  &model.User{Model: gorm.Model{ID: 1}, Username: "TestUser", Email: "test@example.com", Password: "password", Bio: "Bio", Image: "image.png"},
			expectedError: nil,
		},
		{
			name: "User ID Does Not Exist",
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = \\?$").
					WithArgs(id).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			id:            2,
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = \\?$").
					WithArgs(id).
					WillReturnError(errors.New("database connection error"))
			},
			id:            1,
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name: "User ID is Zero",
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = \\?$").
					WithArgs(id).
					WillReturnError(errors.New("invalid ID"))
			},
			id:            0,
			expectedUser:  nil,
			expectedError: errors.New("invalid ID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlmock", db)
			if err != nil {
				t.Fatalf("failed to open gorm database: %v", err)
			}

			tc.setupMock(mock, tc.id)

			store := UserStore{db: gormDB}

			user, err := store.GetByID(tc.id)

			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
