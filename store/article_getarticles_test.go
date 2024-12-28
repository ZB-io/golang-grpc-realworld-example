package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"database/sql/driver"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetArticles(t *testing.T) {
	// Set up common variables and mock environment
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	store := &ArticleStore{db: gormDB}

	// Define table-driven tests for various scenarios
	tests := []struct {
		name          string
		tagName       string
		username      string
		favoritedBy   *model.User
		limit         int64
		offset        int64
		mockQueries   func()
		expectedCount int
		expectError   bool
	}{
		{
			name:     "Retrieve Articles by Author Username",
			username: "author1",
			mockQueries: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2)
				mock.ExpectQuery("join users on articles.user_id = users.id").
					WithArgs("author1").
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:    "Retrieve Articles by Tag Name",
			tagName: "go",
			mockQueries: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(3)
				mock.ExpectQuery("join tags on tags.id = article_tags.tag_id").
					WithArgs("go").
					WillReturnRows(rows)
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:        "Retrieve Articles Favorited by a User",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			mockQueries: func() {
				rows := sqlmock.NewRows([]string{"article_id"}).
					AddRow(4).
					AddRow(5)
				mock.ExpectQuery("select article_id from favorite_articles").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:   "Retrieve Articles with Pagination",
			limit:  10,
			offset: 5,
			mockQueries: func() {
				rows := sqlmock.NewRows([]string{"id", "title"})
				for i := 1; i <= 10; i++ {
					rows.AddRow(i, "Article Title")
				}
				mock.ExpectQuery("limit 10 offset 5").
					WillReturnRows(rows)
			},
			expectedCount: 10,
			expectError:   false,
		},
		{
			name:          "No Articles Match Criteria",
			username:      "nonexistent_user",
			mockQueries: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("join users on articles.user_id = users.id").
					WithArgs("nonexistent_user").
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError   false,
		},
		{
			name: "Error Handling with Database Failures",
			mockQueries: func() {
				mock.ExpectQuery("select article_id from favorite_articles").
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQueries()
			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if !tt.expectError {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(articles))
			} else {
				assert.Error(t, err)
			}
		})
	}
}
