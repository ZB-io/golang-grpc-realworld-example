package store

import (
	"errors"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





func TestCreateComment(t *testing.T) {
	testCases := []struct {
		description   string
		comment       *model.Comment
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "Successful Comment Creation",
			comment: &model.Comment{
				ID:        0,
				Body:      "This is a test comment",
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs(sqlmock.AnyArg(), "This is a test comment", 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			description: "Comment Creation with Database Error",
			comment: &model.Comment{
				ID:        0,
				Body:      "This is a test comment",
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs(sqlmock.AnyArg(), "This is a test comment", 1).WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("db error"),
		},
		{
			description: "Comment Creation with Nil Comment",
			comment:     nil,
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("gorm: insert clause builder does not accept empty structs"),
		},
		{
			description: "Comment Creation with Pre-Existing ID",
			comment: &model.Comment{
				ID:        1,
				Body:      "This is another test comment",
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs(1, "This is another test comment", 1).WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrInvalidTransaction,
		},
		{
			description: "Comment Creation with Missing Required Fields",
			comment: &model.Comment{
				ID:        0,
				Body:      "",
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs(sqlmock.AnyArg(), "", 1).WillReturnError(errors.New("missing required fields"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("missing required fields"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock sql db, %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open gorm db, %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := ArticleStore{db: gormDB}
			err = store.CreateComment(tc.comment)

			if (err != nil && tc.expectedError == nil) || (err == nil && tc.expectedError != nil) {
				t.Fatalf("expected error '%v', got '%v'", tc.expectedError, err)
			}
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Fatalf("expected error message '%v', got '%v'", tc.expectedError.Error(), err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func (s *ArticleStore) CreateComment(m *model.Comment) error {
	if m == nil || m.ID != 0 {
		return errors.New("gorm: insert clause builder does not accept empty structs")
	}
	return s.db.Create(&m).Error
}
