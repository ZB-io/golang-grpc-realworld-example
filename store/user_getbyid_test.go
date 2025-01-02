package store

import (
	"testing"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)







func TestUserStoreGetByID(t *testing.T) {

	testCases := []struct {
		desc          string
		userID        uint
		mockDBHandler func(mock sqlmock.Sqlmock, id uint)
		expectedError error
	}{
		{
			desc:   "Positive test with valid user ID",
			userID: 1,
			mockDBHandler: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)\\?$").
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedError: nil,
		},
		{
			desc:   "Negative test with invalid user ID",
			userID: 2,
			mockDBHandler: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)\\?$").
					WithArgs(id).
					WillReturnError(errors.New("record not found"))
			},
			expectedError: errors.New("record not found"),
		},
		{
			desc:   "Edge case with user ID of zero",
			userID: 0,
			mockDBHandler: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(0, "testuser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)\\?$").
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedError: nil,
		},
		{
			desc:   "Negative test with database connection error",
			userID: 3,
			mockDBHandler: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)\\?$").
					WithArgs(id).
					WillReturnError(errors.New("connection error"))
			},
			expectedError: errors.New("connection error"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			db, mock, _ := sqlmock.New()
			gdb, _ := gorm.Open("postgres", db)
			defer db.Close()

			store := &UserStore{db: gdb}
			tC.mockDBHandler(mock, tC.userID)

			user, err := store.GetByID(tC.userID)

			if tC.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, tC.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tC.userID, user.ID)
			}
		})
	}
}
