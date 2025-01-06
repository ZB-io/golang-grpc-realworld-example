package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// UserStore type
type UserStore struct {
	db *gorm.DB
}

// GetByUsername method
func (s *UserStore) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := s.db.Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// TestUserStoreGetByUsername is a test function for the UserStore's GetByUsername
func TestUserStoreGetByUsername(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		username string
		setup    func(mock sqlmock.Sqlmock)
		check    func(user *model.User, err error)
	}{
		{
			name:     "Valid User Retrieval",
			username: "testuser",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ID", "Username", "Email", "Password", "Bio", "Image"}).
					AddRow(1, "testuser", "testuser@example.com", "testpassword", "testbio", "testimage")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("testuser").WillReturnRows(rows)
			},
			check: func(user *model.User, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if user.Username != "testuser" {
					t.Errorf("Expected username testuser but got %s", user.Username)
				}
			},
		},
		{
			name:     "User Not Found",
			username: "unknownuser",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("unknownuser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			check: func(user *model.User, err error) {
				if err != gorm.ErrRecordNotFound {
					t.Errorf("Expected error %v but got %v", gorm.ErrRecordNotFound, err)
				}
			},
		},
		{
			name:     "Database Error",
			username: "dberroruser",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("dberroruser").
					WillReturnError(errors.New("database error"))
			},
			check: func(user *model.User, err error) {
				if err == nil || err.Error() != "database error" {
					t.Errorf("Expected database error but got %v", err)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm DB: %v", err)
			}

			tc.setup(mock)

			store := &UserStore{db: gormDB}
			user, err := store.GetByUsername(tc.username)
			tc.check(user, err)
		})
	}
}
