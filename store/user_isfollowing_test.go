package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)






func TestUserStoreIsFollowing(t *testing.T) {
	testCases := []struct {
		name           string
		userA          *model.User
		userB          *model.User
		mockDBFunc     func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User)
		expectedResult bool
		expectError    bool
	}{
		{
			name:  "Successful following check between two valid users",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDBFunc: func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User) {
				mock.ExpectQuery("SELECT count(*) FROM `follows` WHERE `from_user_id` = \\? AND `to_user_id` = \\?").
					WithArgs(userA.ID, userB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedResult: true,
			expectError:    false,
		},
		{
			name:  "Successful following check between two valid users where A is not following B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDBFunc: func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User) {
				mock.ExpectQuery("SELECT count(*) FROM `follows` WHERE `from_user_id` = \\? AND `to_user_id` = \\?").
					WithArgs(userA.ID, userB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name:  "Null user A provided to the function",
			userA: nil,
			userB: &model.User{Model: gorm.Model{ID: 1}},
			mockDBFunc: func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User) {

			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name:  "Null user B provided to the function",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: nil,
			mockDBFunc: func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User) {

			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name:  "Database error during the operation",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDBFunc: func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User) {
				mock.ExpectQuery("SELECT count(*) FROM `follows` WHERE `from_user_id` = \\? AND `to_user_id` = \\?").
					WithArgs(userA.ID, userB.ID).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			gormDB, err := gorm.Open("postgres", db)
			assert.NoError(t, err)

			tc.mockDBFunc(mock, tc.userA, tc.userB)

			store := &UserStore{db: gormDB}

			result, err := store.IsFollowing(tc.userA, tc.userB)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
