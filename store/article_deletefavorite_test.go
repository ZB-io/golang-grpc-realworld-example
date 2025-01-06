package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model" // Corrected the import path
	"github.com/stretchr/testify/assert"
)

type ArticleStore struct {
	db *gorm.DB
}

func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error {
	tx := s.db.Begin()

	err := tx.Model(a).Association("FavoritedUsers").
		Delete(u).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(a).
		Update("favorites_count", gorm.Expr("favorites_count - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	a.FavoritesCount--

	return nil
}

func TestArticleStoreDeleteFavorite(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)

	defer func() {
		_ = db.Close()
	}()

	store := &ArticleStore{db: gormDB}

	testCases := []struct {
		name           string
		setupMock      func()
		expectedError  error
		expectedResult int32
	}{
		{
			name: "Successful Deletion of Favorite",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM favorite_articles").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError:  nil,
			expectedResult: 1,
		},
		{
			name: "Unsuccessful Deletion due to Database Error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM favorite_articles").
					WithArgs(1, 1).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError:  errors.New("db error"),
			expectedResult: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			a := &model.Article{
				Model:           gorm.Model{ID: 1},
				Title:           "Test",
				Description:     "Test",
				Body:            "Test",
				FavoritesCount:  2,
				FavoritedUsers:  []model.User{{Model: gorm.Model{ID: 1}}},
			}

			u := &model.User{
				Model: gorm.Model{ID: 1},
			}

			err := store.DeleteFavorite(a, u)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, a.FavoritesCount)
		})
	}
}
