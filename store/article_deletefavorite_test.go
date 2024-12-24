package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
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
	var testCases = []struct {
		name    string
		article *model.Article
		user    *model.User
		db      *gorm.DB
		setupDB func(mock sqlmock.Sqlmock, article *model.Article, user *model.User)
		hasErr  bool
	}{
		{
			name:    "Successful Deletion of Favorite",
			article: &model.Article{Model: gorm.Model{ID: 999}, FavoritesCount: 5},
			user:    &model.User{Model: gorm.Model{ID: 123}},
			setupDB: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("^DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			hasErr: false,
		},
		{
			name:    "Deletion of Nonexistent Favorite",
			article: &model.Article{Model: gorm.Model{ID: 999}, FavoritesCount: 5},
			user:    &model.User{Model: gorm.Model{ID: 123}},
			setupDB: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectExec("^DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			hasErr: false,
		},
		{
			name:    "Deletion when FavoritesCount is Zero",
			article: &model.Article{Model: gorm.Model{ID: 999}, FavoritesCount: 0},
			user:    &model.User{Model: gorm.Model{ID: 123}},
			setupDB: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectExec("^DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			hasErr: false,
		},
		{
			name:    "Database Rollback in Case of Error",
			article: &model.Article{Model: gorm.Model{ID: 999}, FavoritesCount: 5},
			user:    &model.User{Model: gorm.Model{ID: 123}},
			setupDB: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("error executing query"))
				mock.ExpectRollback()
			},
			hasErr: true,
		},
	}

	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	gdb, _ := gorm.Open("postgres", mockDB)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			// Preparing the mock DB
			tc.setupDB(mock, tc.article, tc.user)
	
			// Creating the test scenario
			store := ArticleStore{db: gdb}
			err := store.DeleteFavorite(tc.article, tc.user)

			// Assertion
			if (err != nil) != tc.hasErr {
				t.Errorf("DeleteFavorite() error = %v, hasErr %v", err, tc.hasErr)
				return
			}
		})
	}	
}
