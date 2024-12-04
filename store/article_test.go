package store

import (
	"database/sql"
	"testing"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestDeleteFavorite(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
		expectedCount int32
	}

	baseTime := time.Now()
	testArticle := &model.Article{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: baseTime,
			UpdatedAt: baseTime,
		},
		Title:          "Test Article",
		Description:    "Test Description",
		Body:           "Test Body",
		FavoritesCount: 1,
	}

	testUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: baseTime,
			UpdatedAt: baseTime,
		},
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		Bio:      "Test Bio",
		Image:    "test.jpg",
	}

	tests := []testCase{
		{
			name:    "Successful Deletion",
			article: testArticle,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(testUser.ID, testArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, testArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name:    "Failed Association Deletion",
			article: testArticle,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(testUser.ID, testArticle.ID).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: sql.ErrConnDone,
			expectedCount: 1,
		},
		{
			name:    "Failed FavoritesCount Update",
			article: testArticle,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(testUser.ID, testArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, testArticle.ID).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: sql.ErrConnDone,
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gdb, err := gorm.Open("mysql", db)
			assert.NoError(t, err)
			defer gdb.Close()

			tc.setupMock(mock)

			store := &ArticleStore{db: gdb}

			startCount := tc.article.FavoritesCount
			err = store.DeleteFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, startCount, tc.article.FavoritesCount, "FavoritesCount should not change on error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, tc.article.FavoritesCount, "FavoritesCount should be decremented")
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

