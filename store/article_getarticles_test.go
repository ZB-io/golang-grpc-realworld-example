package store

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetArticles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, _ := gorm.Open("postgres", db)

	articleStore := &ArticleStore{db: gormDB}
	defer db.Close()

	testCases := []struct {
		name          string
		tagName       string
		username      string
		favoritedBy   *model.User
		limit         int64
		offset        int64
		mockFunc      func()
		expected      []model.Article
		expectedError error
	}{
		{
			name: "Retrieve all articles",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"title", "description", "body"}).
					AddRow("title1", "description1", "body1").
					AddRow("title2", "description2", "body2")
				mock.ExpectQuery("^SELECT (.+) FROM articles$").WillReturnRows(rows)
			},
			expected: []model.Article{
				{Title: "title1", Description: "description1", Body: "body1"},
				{Title: "title2", Description: "description2", Body: "body2"},
			},
			expectedError: nil,
		},
		// TODO: Add more test cases here for other scenarios
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			articles, err := articleStore.GetArticles(tc.tagName, tc.username, tc.favoritedBy, tc.limit, tc.offset)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expected, articles)
		})
	}
}
