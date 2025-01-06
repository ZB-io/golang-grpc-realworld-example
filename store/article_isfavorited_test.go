package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type ArticleStore struct {
	db *gorm.DB
}

func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	if a == nil || u == nil {
		return false, nil
	}

	var count int
	err := s.db.Table("favorite_articles").
		Where("article_id = ? AND user_id = ?", a.ID, u.ID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func TestArticleStoreIsFavorited(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		article       *model.Article
		user          *model.User
		mockDBFunc    func(mock sqlmock.Sqlmock)
		expected      bool
		expectedError bool
	}{
		{
			name: "Valid Article and User",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockDBFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(.*) FROM `favorite_articles` WHERE `article_id` = ? AND `user_id` = ?").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: false,
		},
		{
			name: "Valid Article and User but User has not favorited the Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockDBFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(.*) FROM `favorite_articles` WHERE `article_id` = ? AND `user_id` = ?").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: false,
		},
		{
			name:     "Nil Article",
			article:  nil,
			user:     &model.User{},
			expected: false,
		},
		{
			name:     "Nil User",
			article:  &model.Article{},
			user:     nil,
			expected: false,
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockDBFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(.*) FROM `favorite_articles` WHERE `article_id` = ? AND `user_id` = ?").
					WithArgs(1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormDB, _ := gorm.Open("postgres", db)

			if tc.mockDBFunc != nil {
				tc.mockDBFunc(mock)
			}

			store := &ArticleStore{db: gormDB}
			result, err := store.IsFavorited(tc.article, tc.user)

			assert.Equal(t, tc.expected, result)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
