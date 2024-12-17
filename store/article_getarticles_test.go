package store

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// TestArticleStoreGetArticles function implements the unit test for GetArticles function in store package
func TestArticleStoreGetArticles(t *testing.T) {

	tests := []struct {
		name          string
		tagName       string
		username      string
		favoritedBy   *model.User
		limit         int64
		offset        int64
		mockDB        func(mock sqlmock.Sqlmock, rows *sqlmock.Rows)
		expectedError string
	}{
		{name: "Successful Retrieval of Articles By Username",
			tagName:  "",
			username: "testUser",
			limit:    1,
			offset:   0,
			mockDB: func(mock sqlmock.Sqlmock, rows *sqlmock.Rows) {
				mock.ExpectQuery("^select (.+) from `articles` join users on articles.user_id = users.id where users.username = .*").
					WithArgs("testUser").
					WillReturnRows(rows)
			},
			expectedError: "",
		},
		// TODO: Add more test scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gdb, _ := gorm.Open("postgres", db)

			articles := []model.Article{
				{
					// TODO: Initialize your model.Article with the desired test data
				},
				// TODO: Add more articles if needed
			}

			rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "tag", "author", "user_id", "favorites_count", "favoritedUsers", "comments"})
			rowIDs := make([]uint, len(articles))
			for i, article := range articles {
				row := []driver.Value{
					article.ID,
					article.Title,
					article.Description,
					article.Body,
					article.Tags,
					article.Author,
					article.UserID,
					article.FavoritesCount,
					article.FavoritedUsers,
					article.Comments,
				}
				rows = rows.AddRow(row...)
				rowIDs[i] = article.ID
			}

			tt.mockDB(mock, rows)
			articleStore := &ArticleStore{db: gdb}

			result, err := articleStore.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			// Handle error test case
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				// Assert results for success test case
				assert.Nil(t, err)
				// TODO: For accuracy you will need to compare the result with your desired result
			}
		})
	}
}
