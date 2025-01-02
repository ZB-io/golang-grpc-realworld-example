package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)






func TestArticleStoreGetByID(t *testing.T) {

	var testCases = []struct {
		name            string
		mockDB          func() (*gorm.DB, sqlmock.Sqlmock)
		id              uint
		expectedArticle *model.Article
		expectedError   error
	}{
		{
			name: "retrieve existing article by id",
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("postgres", db)

				article := &model.Article{Model: gorm.Model{ID: 1}}

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = ?$").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				return gormDB, mock
			},
			id:              1,
			expectedArticle: &model.Article{Model: gorm.Model{ID: 1}},
			expectedError:   nil,
		},
		{
			name: "retrieve non-existing article by id",
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("postgres", db)

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = ?$").
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)

				return gormDB, mock
			},
			id:              1,
			expectedArticle: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
		{
			name: "database connection error",
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("postgres", db)

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = ?$").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))

				return gormDB, mock
			},
			id:              1,
			expectedArticle: nil,
			expectedError:   errors.New("database connection error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, _ := tc.mockDB()
			store := &ArticleStore{db: db}
			article, err := store.GetByID(tc.id)

			assert.Equal(t, tc.expectedArticle, article)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
