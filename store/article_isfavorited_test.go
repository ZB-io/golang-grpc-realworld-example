package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"regexp" // Added missing "regexp" package
)

// TestArticleStoreIsFavorited is a test function for the IsFavorited method of the ArticleStore
func TestArticleStoreIsFavorited(t *testing.T) {
	// Initialize the mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("failed to open gorm database: %s", err)
	}
	store := &ArticleStore{db : gdb}

	testCases := []struct {
		description      string
		user             *model.User
		article          *model.Article
		prepareMock      func()
		expectedFavorited bool
		expectedError    bool
	}{
		{
			description: "favorited article returns true",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
			},
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
			},
			prepareMock: func() {
				// prepare db mock
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "favorite_articles" WHERE (article_id = $1 AND user_id = $2)`)).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedFavorited: true,
			expectedError:    false,
		},  
		// TODO: Add other test cases for each scenario following the format above.
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			tt.prepareMock()
			result, err := store.IsFavorited(tt.article, tt.user)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFavorited, result)
			}
		})
	}
}
