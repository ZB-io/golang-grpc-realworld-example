package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)






func TestIsFavorited(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(sqlmock.Sqlmock)
		article      *model.Article
		user         *model.User
		expectedBool bool
		expectError  bool
	}{
		{
			name: "Scenario 1: User has favorited the article",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count(.+) FROM favorite_articles WHERE").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			article:      &model.Article{ID: 1},
			user:         &model.User{ID: 1},
			expectedBool: true,
			expectError:  false,
		},
		{
			name: "Scenario 2: User has not favorited the article",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count(.+) FROM favorite_articles WHERE").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			article:      &model.Article{ID: 1},
			user:         &model.User{ID: 2},
			expectedBool: false,
			expectError:  false,
		},
		{
			name: "Scenario 3: Article is nil",
			setupMocks: func(mock sqlmock.Sqlmock) {

			},
			article:      nil,
			user:         &model.User{ID: 1},
			expectedBool: false,
			expectError:  false,
		},
		{
			name: "Scenario 4: User is nil",
			setupMocks: func(mock sqlmock.Sqlmock) {

			},
			article:      &model.Article{ID: 1},
			user:         nil,
			expectedBool: false,
			expectError:  false,
		},
		{
			name: "Scenario 5: Database error occurs",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count(.+) FROM favorite_articles WHERE").
					WithArgs(1, 1).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			article:      &model.Article{ID: 1},
			user:         &model.User{ID: 1},
			expectedBool: false,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setupMocks(mock)

			store := &ArticleStore{db: &gorm.DB{DB: db}}

			result, err := store.IsFavorited(tt.article, tt.user)

			if tt.expectError {
				assert.Error(t, err, "Expected an error but did not get one")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
			}

			assert.Equal(t, tt.expectedBool, result, "Result did not match expected value")

			assert.NoError(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations")
		})
	}
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
