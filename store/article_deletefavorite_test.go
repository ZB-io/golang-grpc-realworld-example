package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/require"
)

func TestArticleStoreDeleteFavorite(t *testing.T) {
	cases := []struct {
		name string
		shouldFail                  bool
		associationDeleteWillFail   bool
		favoritesCountUpdateWillFail bool
	}{
		{"Successful Deletion of Favorite", false, false, false},
		{"Failed Deletion Due to Database Error", true, true, false},
		{"Rollback on Failed Favorite Count Update", true, false, true},
		{"Check Correct Commit when there is No Error", false, false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			require.NoError(t, err)

			as := &ArticleStore{gormDB}

			// Begin and Commit operations are performed within an actual transaction.
			// We simulate this by ExpectBegin and ExpectCommit 
			mock.ExpectBegin()  

			if c.associationDeleteWillFail {
				mock.ExpectExec("DELETE FROM favorited_users").WillReturnError(errors.New("Association deletion failed error"))
			} else {
				mock.ExpectExec("DELETE FROM favorited_users").WillReturnResult(sqlmock.NewResult(1, 1))
			}

			if c.favoritesCountUpdateWillFail {
				mock.ExpectExec("UPDATE favorites_count").WillReturnError(errors.New("Failed Favorites Count Update"))
			} else {
				mock.ExpectExec("UPDATE favorites_count").WillReturnResult(sqlmock.NewResult(1, 1))
			}

			mock.ExpectCommit() 

			// create mock Article and User
			article := &model.Article{
				FavoritesCount: 5,
			}

			user := &model.User{}

			err = as.DeleteFavorite(article, user)

			if c.shouldFail {
				require.Error(t, err, "Expected an error but got none")
			} else {
				require.NoError(t, err, "Unexpected error occurred")
			}
		})
	}
}

// This function is redeclared in this block, therefore it is commented out to avoid a compilation error.
// Similarly, the ArticleStore and User structs are also redeclared and commented out.
/*
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
type ArticleStore struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}

type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
}
*/
