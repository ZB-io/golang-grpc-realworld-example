package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

func TestUserStoreIsFollowing(t *testing.T) {
	testCases := []struct {
		name            string
		userA           *model.User
		userB           *model.User
		mock            func(sqlmock.Sqlmock, *model.User, *model.User)
		error           error
		expectedResult  bool
	}{
		{
			name: "valid users and following exists",
			userA: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
			},
			userB: &model.User{
				Model: gorm.Model{
					ID: 2,
				},
			},
			mock: func(mock sqlmock.Sqlmock, userA *model.User, userB *model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT count(*) FROM follows WHERE from_user_id = ? AND to_user_id = ?").
					WithArgs(userA.ID, userB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedResult: true,
		},

		//other test cases can be include in this manner
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gdb, _ := gorm.Open("postgres", db)
			sqlmock.NewRows([]string{"column"}).AddRow(1)
			testCase.mock(mock, testCase.userA, testCase.userB)

			userStore := &UserStore{db: gdb}
			result, err := userStore.IsFollowing(testCase.userA, testCase.userB)

			if testCase.expectedResult != result {
				t.Errorf("Expected result does not match with actual result. Expected: %t, got: %t", testCase.expectedResult, result)
			}

			if testCase.error != err {
				t.Errorf("Expected error does not match with actual error. Expected: %v, got: %v", testCase.error, err)
			}
		})
	}
}
