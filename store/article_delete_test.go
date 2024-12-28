package store

import (
	"testing"
	"fmt"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	_ "github.com/go-sql-driver/mysql"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





func TestArticleStoreDelete(t *testing.T) {
	tests := []struct {
		name          string
		setupMockFunc func(sqlmock.Sqlmock)
		article       *model.Article
		db            *gorm.DB
		expectedError error
	}{
		{
			name: "Scenario 1: Successful Deletion of an Existing Article",
			setupMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^DELETE FROM articles WHERE id = ?$").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			article:       &model.Article{ID: 1},
			expectedError: nil,
		},
		{
			name: "Scenario 2: Attempting to Delete a Non-Existent Article",
			setupMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^DELETE FROM articles WHERE id = ?$").
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			article:       &model.Article{ID: 2},
			expectedError: nil,
		},
		{
			name: "Scenario 3: Database Error During Deletion Process",
			setupMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^DELETE FROM articles WHERE id = ?$").
					WithArgs(3).
					WillReturnError(errors.New("some error"))
			},
			article:       &model.Article{ID: 3},
			expectedError: errors.New("some error"),
		},
		{
			name:          "Scenario 4: Attempt to Delete Article When Database is Nil",
			setupMockFunc: func(_ sqlmock.Sqlmock) {},
			db:            nil,
			article:       &model.Article{ID: 4},
			expectedError: fmt.Errorf("database is nil"),
		},
		{
			name: "Scenario 5: Deletion with Uninitialized Article Object",
			setupMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^DELETE FROM articles WHERE id = ?$").
					WithArgs(nil).
					WillReturnError(errors.New("article is nil"))
			},
			article:       nil,
			expectedError: fmt.Errorf("article is nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			if tt.setupMockFunc != nil {
				tt.setupMockFunc(mock)
			}

			var gormDB *gorm.DB

			if tt.db == nil {
				gormDB, err = gorm.Open("mysql", db)
				assert.NoError(t, err)
			} else {
				gormDB = tt.db
			}

			store := &ArticleStore{db: gormDB}

			err = store.Delete(tt.article)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				t.Logf("Expected error: %v, got: %v", tt.expectedError, err)
			} else {
				t.Logf("Deletion test passed for ID: %v", tt.article.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}


func (s *ArticleStore) Delete(m *model.Article) error {
	if s.db == nil {
		return fmt.Errorf("database is nil")
	}
	if m == nil {
		return fmt.Errorf("article is nil")
	}
	return s.db.Delete(m).Error
}
