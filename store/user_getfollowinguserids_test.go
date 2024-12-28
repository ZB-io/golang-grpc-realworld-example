package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)






func TestUserStoreGetFollowingUserIDs(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while opening a stub database connection: %v\n", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("failed to initialize gorm database: %v\n", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	type testData struct {
		description  string
		setup        func()
		input        *model.User
		expectedIDs  []uint
		expectingErr bool
	}

	tests := []testData{
		{
			description: "Successfully Retrieve Following User IDs",
			setup: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)

				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(1).
					WillReturnRows(rows)
			},
			input:        &model.User{ID: 1},
			expectedIDs:  []uint{2, 3, 4},
			expectingErr: false,
		},
		{
			description: "No Following Users",
			setup: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(2).
					WillReturnRows(rows)
			},
			input:        &model.User{ID: 2},
			expectedIDs:  []uint{},
			expectingErr: false,
		},
		{
			description: "Database Connection Error",
			setup: func() {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(3).
					WillReturnError(fmt.Errorf("database connection error"))
			},
			input:        &model.User{ID: 3},
			expectedIDs:  []uint{},
			expectingErr: true,
		},
		{
			description: "User ID is Not Found",
			setup: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(4).
					WillReturnRows(rows)
			},
			input:        &model.User{ID: 4},
			expectedIDs:  []uint{},
			expectingErr: false,
		},
		{
			description: "Invalid User Input (Nil User Object)",
			setup: func() {

			},
			input:        nil,
			expectedIDs:  []uint{},
			expectingErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.setup()

			ids, err := userStore.GetFollowingUserIDs(tc.input)
			if (err != nil) != tc.expectingErr {
				t.Errorf("Got error = %v, expecting error = %v", err, tc.expectingErr)
			}
			if !sameIDs(ids, tc.expectedIDs) {
				t.Errorf("Expected IDs %v, got %v", tc.expectedIDs, ids)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
func sameIDs(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
