package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)




func TestUnfollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	gormDB, err := gorm.Open("postgres", db)
	assert.NoError(t, err)

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		prepare     func()
		userA       *model.User
		userB       *model.User
		expectedErr error
	}{
		{
			name: "Successfully Unfollowing a User",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 2},
			expectedErr: nil,
		},
		{
			name: "Unfollowing a User Not Followed",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 2},
			expectedErr: nil,
		},
		{
			name: "Database Error on Unfollow",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 2},
			expectedErr: errors.New("database error"),
		},
		{
			name: "Unfollowing with Nil User Parameters",
			prepare: func() {

			},
			userA:       nil,
			userB:       nil,
			expectedErr: gorm.ErrInvalidSQL,
		},
		{
			name: "Unfollowing Same User",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 1},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			err = userStore.Unfollow(tt.userA, tt.userB)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}


