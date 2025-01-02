package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)









func TestUserStoreFollow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)

	userStore := &UserStore{db: gormDB}

	userA := &model.User{Username: "userA", Email: "userA@example.com", Password: "passwordA", Bio: "bioA", Image: "imageA"}
	userB := &model.User{Username: "userB", Email: "userB@example.com", Password: "passwordB", Bio: "bioB", Image: "imageB"}

	testCases := []struct {
		name          string
		userA         *model.User
		userB         *model.User
		mockBehaviour func()
		expectedError error
	}{
		{
			name:  "Successful Follow",
			userA: userA,
			userB: userB,
			mockBehaviour: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO follows").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:  "User A already follows User B",
			userA: userA,
			userB: userB,
			mockBehaviour: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO follows").WillReturnError(errors.New("user A already follows user B"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("user A already follows user B"),
		},
		{
			name:          "User A tries to follow themselves",
			userA:         userA,
			userB:         userA,
			mockBehaviour: func() {},
			expectedError: errors.New("a user cannot follow themselves"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehaviour()

			err := userStore.Follow(tc.userA, tc.userB)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
