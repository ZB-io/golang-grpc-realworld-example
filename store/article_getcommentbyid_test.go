package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestArticleStoreGetCommentByID(t *testing.T) {

	tests := []struct {
		name           string
		commentID      uint
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult *model.Comment
		expectedError  error
	}{
		{
			name:      "Retrieve a Comment Successfully",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedRows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
					AddRow(1, "Test Comment", 1, 1)
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE `comments`.`id` = ?").
					WithArgs(1).
					WillReturnRows(expectedRows)
			},
			expectedResult: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			expectedError: nil,
		},
		{
			name:      "Comment Not Found Case",
			commentID: 2,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE `comments`.`id` = ?").
					WithArgs(2).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedResult: nil,
			expectedError:  gorm.ErrRecordNotFound,
		},
		{
			name:      "Database Connection Error Handling",
			commentID: 3,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE `comments`.`id` = ?").
					WithArgs(3).
					WillReturnError(errors.New("connection error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("connection error"),
		},
		{
			name:      "Retrieving a Comment with Relationships Loaded",
			commentID: 4,
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedRows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
					AddRow(4, "Detailed Comment", 1, 1)
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE `comments`.`id` = ?").
					WithArgs(4).
					WillReturnRows(expectedRows)
			},
			expectedResult: &model.Comment{
				Model:     gorm.Model{ID: 4},
				Body:      "Detailed Comment",
				UserID:    1,
				ArticleID: 1,

				Author:  model.User{},
				Article: model.Article{},
			},
			expectedError: nil,
		},
		{
			name:      "Handling of Non-existent Foreign Keys",
			commentID: 5,
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedRows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
					AddRow(5, "Foreign Key Issue Comment", 999, 999)
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE `comments`.`id` = ?").
					WithArgs(5).
					WillReturnRows(expectedRows)
			},
			expectedResult: &model.Comment{
				Model:     gorm.Model{ID: 5},
				Body:      "Foreign Key Issue Comment",
				UserID:    999,
				ArticleID: 999,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)

			gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			store := &ArticleStore{db: gormDB}
			tt.setupMock(mock)

			result, err := store.GetCommentByID(tt.commentID)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
