package store

import (
	"bytes"
	"errors"
	"os"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name        string
		comment     *model.Comment
		setupDB     func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "Successful Deletion of an Existing Comment",
			comment: &model.Comment{
				ID: 1,
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM comments WHERE id = ?").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "Attempt to Delete a Non-Existent Comment",
			comment: &model.Comment{
				ID: 2,
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM comments WHERE id = ?").WithArgs(2).
					WillReturnResult(sqlmock.NewResult(2, 0))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "Database Error During Deletion",
			comment: &model.Comment{
				ID: 3,
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM comments WHERE id = ?").WithArgs(3).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("database error"),
		},
		{
			name:        "Malformed Comment Input",
			comment:     &model.Comment{},
			setupDB:     func(mock sqlmock.Sqlmock) {},
			expectedErr: errors.New("record not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupDB(mock)

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
			}

			store := &ArticleStore{db: gormDB}

			var out bytes.Buffer
			old := os.Stdout
			defer func() { os.Stdout = old }()
			os.Stdout = &out

			err = store.DeleteComment(tt.comment)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
			}

			t.Log(out.String())
		})
	}

	t.Run("Concurrent Deletion Scenarios", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		comment := &model.Comment{ID: 4}
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM comments WHERE id = ?").WithArgs(4).
			WillReturnResult(sqlmock.NewResult(4, 1))
		mock.ExpectCommit()

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
		}

		store := &ArticleStore{db: gormDB}
		var wg sync.WaitGroup

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = store.DeleteComment(comment)
				if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
					t.Errorf("unexpected error occurred during concurrent deletion: %v", err)
				}
			}()
		}

		wg.Wait()
	})
}

func (s *ArticleStore) DeleteComment(m *model.Comment) error {
	return s.db.Delete(m).Error
}
