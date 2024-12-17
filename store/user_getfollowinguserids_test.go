package store

import (
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	// create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// wrap the mock database with gorm
	gdb, err := gorm.Open("postgres", db)
	require.NoError(t, err)

	// creating an instance of our store
	st := &UserStore{
		db: gdb,
	}

	testCases := []struct {
		name      string
		user      model.User
		mock      func()
		expectIDs []uint
		expectErr bool
	}{
		{
			name: "Positive test for the GetFollowingUserIDs function",
			user: model.User{Model: gorm.Model{ID: 1}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(1).
					AddRow(2).
					AddRow(3)
				mock.ExpectQuery(`SELECT (.+) FROM "follows" WHERE (.+)$`).WillReturnRows(rows)
			},
			expectIDs: []uint{1, 2, 3},
			expectErr: false,
		},
		{
			name: "Negative test for the GetFollowingUserIDs function with zero followers",
			user: model.User{Model: gorm.Model{ID: 1}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery(`SELECT (.+) FROM "follows" WHERE (.+)$`).WillReturnRows(rows)
			},
			expectIDs: nil,
			expectErr: false,
		},
		{
			name: "Negative test for the GetFollowingUserIDs function upon database read failure",
			user: model.User{Model: gorm.Model{ID: 1}},
			mock: func() {
				mock.ExpectQuery(`SELECT (.+) FROM "follows" WHERE (.+)$`).WillReturnError(fmt.Errorf("some error"))
			},
			expectIDs: nil,
			expectErr: true,
		},
		{
			name: "Negative test for the GetFollowingUserIDs function with a non-existing user",
			user: model.User{Model: gorm.Model{ID: 1000}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery(`SELECT (.+) FROM "follows" WHERE (.+)$`).WillReturnRows(rows)
			},
			expectIDs: nil,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the mock instructions
			tc.mock()

			// Call the actual function
			ids, err := st.GetFollowingUserIDs(&tc.user)

			// Validate function results
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, ids, len(tc.expectIDs))
				for i, id := range ids {
					assert.Equal(t, tc.expectIDs[i], id)
				}
			}

			// Assert all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unexecuted expectations: %s", err)
			}
		})
	}
}
