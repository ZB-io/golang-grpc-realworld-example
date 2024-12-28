package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
)






func TestGetByUsername(t *testing.T) {

	tests := []struct {
		name          string
		username      string
		expectedUser  *model.User
		expectedError error
		setupMock     func(db sqlmock.Sqlmock)
	}{
		{
			name:     "Retrieve User Successfully",
			username: "validuser",
			expectedUser: &model.User{

				Username: "validuser",
			},
			expectedError: nil,
			setupMock: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)$").
					WithArgs("validuser").
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("validuser"))
			},
		},
		{
			name:          "User Not Found",
			username:      "nonexistentuser",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMock: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)$").
					WithArgs("nonexistentuser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:          "Database Connection Error",
			username:      "anyuser",
			expectedUser:  nil,
			expectedError: gorm.ErrInvalidDB,
			setupMock: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)$").
					WithArgs("anyuser").
					WillReturnError(gorm.ErrInvalidDB)
			},
		},
		{
			name:          "Username Case Sensitivity",
			username:      "VALIDUser",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMock: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)$").
					WithArgs("VALIDUser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:          "Empty Username Input",
			username:      "",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMock: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)$").
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			gormDB, err := gorm.Open("postgres", db)
			assert.NoError(t, err)
			defer gormDB.Close()

			userStore := &UserStore{gormDB}
			user, err := userStore.GetByUsername(tt.username)

			t.Logf("Scenario: %s", tt.name)

			assert.Equal(t, tt.expectedUser, user, "User object did not match the expected value")
			assert.Equal(t, tt.expectedError, err, "Error did not match the expected error")

			if user == tt.expectedUser && err == tt.expectedError {
				t.Log("Test passed.")
			} else {
				t.Logf("Test failed. Expected user: %v, got: %v. Expected error: %v, got: %v", tt.expectedUser, user, tt.expectedError, err)
			}
		})
	}
}
