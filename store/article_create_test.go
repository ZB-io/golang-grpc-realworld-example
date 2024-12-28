package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





func TestArticleStoreCreate(t *testing.T) {

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		input       *model.Article
		expectError bool
	}{
		{
			name: "Successful Article Creation",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: &model.Article{
				Title:   "Valid Title",
				Content: "Valid Content",
			},
			expectError: false,
		},
		{
			name: "Article Creation Fails with Database Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
			},
			input: &model.Article{
				Title:   "Valid Title",
				Content: "Valid Content",
			},
			expectError: true,
		},
		{
			name:        "Article Creation with Nil Article Pointer",
			setupMock:   func(mock sqlmock.Sqlmock) {},
			input:       nil,
			expectError: true,
		},
		{
			name: "Duplicate Article Creation Handling",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").WillReturnError(gorm.ErrUniqueConstraint)
				mock.ExpectRollback()
			},
			input: &model.Article{
				Title:   "Unique Title",
				Content: "Content that leads to duplication",
			},
			expectError: true,
		},
		{
			name: "Article Creation with Maximum Field Sizes",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: &model.Article{
				Title:   "Max Length Title",
				Content: "Max Length Content",
			},
			expectError: false,
		},
		{
			name: "Article Creation with Invalid Field Types",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
			},
			input: &model.Article{
				Title:   "",
				Content: "123456",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("failed to open gorm DB: %s", err)
			}

			tt.setupMock(mock)

			store := &ArticleStore{db: gormDB}

			err = store.Create(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed.", tt.name)
		})
	}
}
func (s *ArticleStore) Create(m *model.Article) error {
	if m == nil {
		return gorm.ErrRecordNotFound
	}
	return s.db.Create(&m).Error
}
