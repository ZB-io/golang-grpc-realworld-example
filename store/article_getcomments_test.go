package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)









func TestArticleStoreGetComments(t *testing.T) {

	testCases := []struct {
		name                  string
		prepareMock           func(mock sqlmock.Sqlmock, article *model.Article)
		expectedError         error
		expectedCommentsCount int
	}{
		{
			name: "Successful retrieval of comments for a given article",
			prepareMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((article_id = \\$1))").
					WithArgs(article.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "article_id", "user_id"}).
						AddRow(1, "comment1", article.ID, 1).
						AddRow(2, "comment2", article.ID, 1))
				mock.ExpectCommit()
			},
			expectedError:         nil,
			expectedCommentsCount: 2,
		},
		{
			name: "Retrieval of comments for an article with no comments",
			prepareMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((article_id = \\$1))").
					WithArgs(article.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "article_id", "user_id"}))
				mock.ExpectCommit()
			},
			expectedError:         nil,
			expectedCommentsCount: 0,
		},
		{
			name: "Retrieval of comments for a non-existent article",
			prepareMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((article_id = \\$1))").
					WithArgs(article.ID).
					WillReturnError(errors.New("record not found"))
				mock.ExpectRollback()
			},
			expectedError:         errors.New("record not found"),
			expectedCommentsCount: 0,
		},
		{
			name: "Retrieval of comments when the database is unreachable",
			prepareMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((article_id = \\$1))").
					WithArgs(article.ID).
					WillReturnError(errors.New("database is down"))
				mock.ExpectRollback()
			},
			expectedError:         errors.New("database is down"),
			expectedCommentsCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormDB, _ := gorm.Open("postgres", db)
			defer db.Close()

			article := &model.Article{Model: gorm.Model{ID: 1}}

			tc.prepareMock(mock, article)

			store := ArticleStore{db: gormDB}

			comments, err := store.GetComments(article)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCommentsCount, len(comments))
			}
		})
	}
}
