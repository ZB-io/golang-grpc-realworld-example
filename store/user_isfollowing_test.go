package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





func TestIsFollowing(t *testing.T) {
	type testCase struct {
		name      string
		userA     *model.User
		userB     *model.User
		mock      func(sqlmock.Sqlmock)
		expected  bool
		expectErr bool
	}

	tests := []testCase{
		{
			name:  "A follows B",
			userA: &model.User{ID: 1},
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM follows`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:  "A does not follow B",
			userA: &model.User{ID: 1},
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM follows`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "A is nil, B is valid",
			userA: nil,
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {

			},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "A is valid, B is nil",
			userA: &model.User{ID: 1},
			userB: nil,
			mock: func(mock sqlmock.Sqlmock) {

			},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "Both A and B nil",
			userA: nil,
			userB: nil,
			mock: func(mock sqlmock.Sqlmock) {

			},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "Database error occurs",
			userA: &model.User{ID: 1},
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM follows`).
					WithArgs(1, 2).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expected:  false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}

			userStore := &UserStore{db: gormDB}
			tt.mock(mock)

			got, err := userStore.IsFollowing(tt.userA, tt.userB)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if got != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
