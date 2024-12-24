package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// setupMock returns a mock SQL driver and a gorm connection to the mock
func setupMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		panic("an error occurred while opening a mock database connection: " + err.Error())
	}

	gormDB, err := gorm.Open("postgres", db)

	if err != nil {
		panic("an error occurred while opening gorm connection: " + err.Error())
	}

	return gormDB, mock
}

type ArticleStore struct {
	db *gorm.DB
}

func (s *ArticleStore) Update(m *model.Article) error {
	return s.db.Model(m).Update(m).Error
}

// TestArticleStoreUpdate tests update functionality of article in the store
func TestArticleStoreUpdate(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		expectError bool
	}{
		{
			name: "Successful Article Update",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			expectError: false,
		},
		{
			name: "Article Update With Invalid Article Parameters",
			article: &model.Article{
				Title:       "",
				Description: "",
				Body:        "",
				UserID:      0,
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db, mock := setupMock()
			defer db.Close()

			articleStore := &ArticleStore{db: db}

			if test.expectError {
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("mock error"))
			} else {
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := articleStore.Update(test.article)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
