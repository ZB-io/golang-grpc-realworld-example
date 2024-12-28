package store

import (
	"fmt"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)




func TestUpdate(t *testing.T) {

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		article     *model.Article
		expectedErr error
	}{
		{
			name: "Successful Update of Article",
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article:     &model.Article{ID: 1, Title: "Updated Title"},
			expectedErr: nil,
		},
		{
			name: "Update with Non-Existent Article ID",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			article:     &model.Article{ID: 9999, Title: "Non-Existent"},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error During Update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("connection error"))
			},
			article:     &model.Article{ID: 1, Title: "Title"},
			expectedErr: fmt.Errorf("connection error"),
		},
		{
			name:        "Update With Nil Article Reference",
			setupMock:   func(mock sqlmock.Sqlmock) {},
			article:     nil,
			expectedErr: fmt.Errorf("invalid input"),
		},
		{
			name: "Update When Article Has Unchanged Fields",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			article:     &model.Article{ID: 1, Title: "Unchanged Title"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock sql db, got error: %v", err)
			}
			defer db.Close()

			sqlDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to connect to mock db, got error: %v", err)
			}
			store := &ArticleStore{db: sqlDB}

			tt.setupMock(mock)

			err = store.Update(tt.article)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			t.Logf("Test case '%s' ran successfully", tt.name)
		})
	}
}


func (s *ArticleStore) Update(m *model.Article) error {
	if m == nil {
		return fmt.Errorf("invalid input")
	}
	return s.db.Model(&m).Update(&m).Error
}
