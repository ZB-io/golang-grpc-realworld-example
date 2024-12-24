package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreIsFavorited(t *testing.T) {
	// Initialize test cases
	var testCases = []struct {
		name        string
		article     *model.Article
		user        *model.User
		count       int
		expectError bool
	}{
		{"Favorited Article with Valid Input", &model.Article{Model: gorm.Model{ID: 1}}, &model.User{Model: gorm.Model{ID: 1}}, 1, false},
		{"Unfavorited Article with Valid Input", &model.Article{Model: gorm.Model{ID: 1}}, &model.User{Model: gorm.Model{ID: 1}}, 0, false},
		{"Error due to DB communication issues", &model.Article{Model: gorm.Model{ID: 1}}, &model.User{Model: gorm.Model{ID: 1}}, 0, true},
		{"Error due to Nil Article", nil, &model.User{Model: gorm.Model{ID: 1}}, 0, false},
		{"Error due to Nil User", &model.Article{Model: gorm.Model{ID: 1}}, nil, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize mock DB
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			gdb, _ := gorm.Open("postgres", db)

			store := &ArticleStore{
				db: gdb,
			}

			// Prepare for success or error scenario
			if !tc.expectError {
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"favorite_articles\" WHERE \\(article_id = \\$1 AND user_id = \\$2\\)$").
					WithArgs(tc.article.ID, tc.user.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			} else {
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"favorite_articles\" WHERE \\(article_id = \\$1 AND user_id = \\$2\\)$").
					WithArgs(tc.article.ID, tc.user.ID).
					WillReturnError(errors.New("database error"))
			}

			// Run the function
			favorite, err := store.IsFavorited(tc.article, tc.user)

			// Assert success scenario
			if !tc.expectError {
				assert.Nil(t, err) // no error
				assert.NotNil(t, favorite) // user is favorited
			} else {
				assert.NotNil(t, err) // expect error
			}
		})
	}
}
