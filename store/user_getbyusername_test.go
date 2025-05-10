package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)





type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func TestUserStoreGetByUsername(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		mock        func(mock sqlmock.Sqlmock)
		checkResult func(user *model.User, err error)
	}{
		{
			name:     "Existing User",
			username: "existing_user",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "existing_user", "test@example.com", "password", "bio", "image.jpg")
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\? ORDER BY .+ LIMIT 1$").
					WithArgs("existing_user").
					WillReturnRows(rows)
			},
			checkResult: func(user *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "existing_user", user.Username)
			},
		},
		{
			name:     "Username Not Found",
			username: "non_existing_user",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\? ORDER BY .+ LIMIT 1$").
					WithArgs("non_existing_user").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			checkResult: func(user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			},
		},
		{
			name:     "Database Connection Error",
			username: "any_user",
			mock: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\? ORDER BY .+ LIMIT 1$").
					WithArgs("any_user").
					WillReturnError(gorm.ErrInvalidDataSource)
			},
			checkResult: func(user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, gorm.ErrInvalidDataSource, err)
			},
		},
		{
			name:     "Retrieve User with Associated Data",
			username: "user_with_data",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(2, "user_with_data", "data@example.com", "securepassword", "bio", "image.jpg")
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\? ORDER BY .+ LIMIT 1$").
					WithArgs("user_with_data").
					WillReturnRows(rows)

			},
			checkResult: func(user *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "user_with_data", user.Username)

			},
		},
		{
			name:     "Case Sensitivity in Username",
			username: "MixedCaseUser",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(3, "MixedCaseUser", "mixedcase@example.com", "mypassword", "bio", "image.jpg")
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\? ORDER BY .+ LIMIT 1$").
					WithArgs("MixedCaseUser").
					WillReturnRows(rows)
			},
			checkResult: func(user *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "MixedCaseUser", user.Username)
			},
		},
		{
			name:     "Empty Username Parameter",
			username: "",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\? ORDER BY .+ LIMIT 1$").
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			checkResult: func(user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			assert.NoError(t, err)

			store := UserStore{db: gormDB}

			tt.mock(mock)

			user, err := store.GetByUsername(tt.username)

			tt.checkResult(user, err)

			assert.NoError(t, mock.ExpectationsWereMet(), "there were unmet expectations")
		})
	}
}
func (s *UserStore) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := s.db.Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
