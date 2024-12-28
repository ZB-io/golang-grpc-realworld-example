package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/jinzhu/gorm"
	"log"
)

func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening the gorm with mock", err)
	}

	userStore := &UserStore{db: gormDB}
	testUser := &model.User{Model: gorm.Model{ID: 1}}

	tests := []struct {
		name           string
		mockSetup      func()
		expectedResult []uint
		expectedError  error
	}{
		{
			name: "Scenario 1: User with One Followed User",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT to_user_id FROM follows WHERE from_user_id = ?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}).AddRow(2))
			},
			expectedResult: []uint{2},
			expectedError:  nil,
		},
		{
			name: "Scenario 2: User with No Followed Users",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT to_user_id FROM follows WHERE from_user_id = ?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}))
			},
			expectedResult: []uint{},
			expectedError:  nil,
		},
		{
			name: "Scenario 3: User Following Multiple Users",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT to_user_id FROM follows WHERE from_user_id = ?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}).AddRow(2).AddRow(3).AddRow(4))
			},
			expectedResult: []uint{2, 3, 4},
			expectedError:  nil,
		},
		{
			name: "Scenario 4: Non-existing User Input",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT to_user_id FROM follows WHERE from_user_id = ?`).
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedResult: []uint{},
			expectedError:  gorm.ErrRecordNotFound,
		},
		{
			name: "Scenario 5: Database Error Occurrence",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT to_user_id FROM follows WHERE from_user_id = ?`).
					WithArgs(1).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedResult: []uint{},
			expectedError:  gorm.ErrInvalidSQL,
		},
		{
			name: "Scenario 6: Mixed State of Following",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT to_user_id FROM follows WHERE from_user_id = ?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}).AddRow(2).AddRow(3))
			},
			expectedResult: []uint{2, 3},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result, err := userStore.GetFollowingUserIDs(testUser)
			if (err != nil) && (tt.expectedError == nil || err.Error() != tt.expectedError.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected result length %v, got %v", len(tt.expectedResult), len(result))
			}
			for i, v := range result {
				if v != tt.expectedResult[i] {
					t.Errorf("expected result %v, got %v", tt.expectedResult, result)
				}
			}
		})
	}
}
