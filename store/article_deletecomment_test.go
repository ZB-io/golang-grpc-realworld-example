package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// Assuming ArticleStore is defined in the same package
type ArticleStore struct {
	db *gorm.DB
}

func (s *ArticleStore) DeleteComment(m *model.Comment) error {
	return s.db.Delete(m).Error
}

func TestArticleStoreDeleteComment(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		comment        *model.Comment
		mockDB         func() (*gorm.DB, sqlmock.Sqlmock)
		expectedErr    error
	}{
		{
			name: "Successful Deletion of a Comment",
			comment: &model.Comment{
				Body: "Test Comment",
			},
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
				gormDB, _ := gorm.Open("postgres", db)
				return gormDB, mock
			},
			expectedErr: nil,
		},
		{
			name: "Error When Deleting a Comment",
			comment: &model.Comment{
				Body: "Test Comment",
			},
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("DELETE").WillReturnError(errors.New("delete error"))
				gormDB, _ := gorm.Open("postgres", db)
				return gormDB, mock
			},
			expectedErr: errors.New("delete error"),
		},
		{
			name: "Deletion of a Non-Existing Comment",
			comment: &model.Comment{
				Body: "Test Comment",
			},
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 0))
				gormDB, _ := gorm.Open("postgres", db)
				return gormDB, mock
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := tc.mockDB()
			store := &ArticleStore{
				db: db,
			}

			err := store.DeleteComment(tc.comment)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err.Error() != tc.expectedErr.Error() {
					t.Fatalf("expected %v, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
