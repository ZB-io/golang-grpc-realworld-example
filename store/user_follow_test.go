package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreFollow(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	assert.NoError(t, err)
	defer gormDB.Close()

	userStore := &UserStore{db: gormDB}

	type testCase struct {
		name        string
		a           *model.User
		b           *model.User
		mockSetup   func()
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Successful Follow Operation",
			a:    &model.User{ID: 1},
			b:    &model.User{ID: 2},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO follows`).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Follow Operation with User already Followed",
			a:    &model.User{ID: 1},
			b:    &model.User{ID: 2},
			mockSetup: func() {

				mock.ExpectExec(`INSERT INTO follows`).WithArgs(1, 2).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedErr: nil,
		},
		{
			name: "Follow Operation with Database Error",
			a:    &model.User{ID: 1},
			b:    &model.User{ID: 2},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO follows`).WithArgs(1, 2).WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedErr: gorm.ErrInvalidSQL,
		},
		{
			name:        "Follow Operation with Null User a",
			a:           nil,
			b:           &model.User{ID: 2},
			mockSetup:   func() {},
			expectedErr: gorm.ErrInvalidSQL,
		},
		{
			name:        "Follow Operation with Null User b",
			a:           &model.User{ID: 1},
			b:           nil,
			mockSetup:   func() {},
			expectedErr: gorm.ErrInvalidSQL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := userStore.Follow(tc.a, tc.b)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
