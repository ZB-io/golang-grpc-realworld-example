package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

const testArticleID = uint(1)
const testArticleID = uint(1)testUserID = uint(1)
type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestIsFavorited(t *testing.T) {
	tests := []struct {
		name          string
		setupMockDB   func(*gorm.DB)
		article       *model.Article
		user          *model.User
		expectedFav   bool
		expectedError error
	}{
		{
			name: "Article is favorited by the user",
			setupMockDB: func(db *gorm.DB) {
				db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (?, ?)", testArticleID, testUserID)
			},
			article:       &model.Article{Model: gorm.Model{ID: testArticleID}},
			user:          &model.User{Model: gorm.Model{ID: testUserID}},
			expectedFav:   true,
			expectedError: nil,
		},
		{
			name:          "Article is not favorited by the user",
			setupMockDB:   func(db *gorm.DB) {},
			article:       &model.Article{Model: gorm.Model{ID: testArticleID}},
			user:          &model.User{Model: gorm.Model{ID: testUserID}},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Nil article parameter",
			setupMockDB:   func(db *gorm.DB) {},
			article:       nil,
			user:          &model.User{Model: gorm.Model{ID: testUserID}},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Nil user parameter",
			setupMockDB:   func(db *gorm.DB) {},
			article:       &model.Article{Model: gorm.Model{ID: testArticleID}},
			user:          nil,
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name: "Database error occurs",
			setupMockDB: func(db *gorm.DB) {
				db.AddError(errors.New("database error"))
			},
			article:       &model.Article{Model: gorm.Model{ID: testArticleID}},
			user:          &model.User{Model: gorm.Model{ID: testUserID}},
			expectedFav:   false,
			expectedError: errors.New("database error"),
		},
		{
			name: "Multiple favorites for the same article-user pair",
			setupMockDB: func(db *gorm.DB) {
				db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (?, ?)", testArticleID, testUserID)
				db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (?, ?)", testArticleID, testUserID)
			},
			article:       &model.Article{Model: gorm.Model{ID: testArticleID}},
			user:          &model.User{Model: gorm.Model{ID: testUserID}},
			expectedFav:   true,
			expectedError: nil,
		},
		{
			name: "Large number of favorites in the database",
			setupMockDB: func(db *gorm.DB) {
				for i := 0; i < 1000; i++ {
					db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (?, ?)", i, i)
				}
				db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (1000, 1000)")
			},
			article:       &model.Article{Model: gorm.Model{ID: 1000}},
			user:          &model.User{Model: gorm.Model{ID: 1000}},
			expectedFav:   true,
			expectedError: nil,
		},
		{
			name: "Article and user exist but have no relationship",
			setupMockDB: func(db *gorm.DB) {
				db.Exec("INSERT INTO articles (id) VALUES (?)", testArticleID)
				db.Exec("INSERT INTO users (id) VALUES (?)", testUserID)
			},
			article:       &model.Article{Model: gorm.Model{ID: testArticleID}},
			user:          &model.User{Model: gorm.Model{ID: testUserID}},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Zero IDs for article and user",
			setupMockDB:   func(db *gorm.DB) {},
			article:       &model.Article{Model: gorm.Model{ID: 0}},
			user:          &model.User{Model: gorm.Model{ID: 0}},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Very large IDs for article and user",
			setupMockDB:   func(db *gorm.DB) {},
			article:       &model.Article{Model: gorm.Model{ID: 1<<63 - 1}},
			user:          &model.User{Model: gorm.Model{ID: 1<<63 - 1}},
			expectedFav:   false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, _ := gorm.Open("sqlite3", ":memory:")
			defer db.Close()
			db.AutoMigrate(&model.Article{}, &model.User{})
			db.Exec("CREATE TABLE favorite_articles (article_id int, user_id int)")
			tt.setupMockDB(db)

			store := &ArticleStore{db: db}

			isFavorited, err := store.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedFav, isFavorited)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
