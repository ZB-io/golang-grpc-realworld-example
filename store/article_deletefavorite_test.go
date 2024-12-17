package store

import (
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreDeleteFavorite(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}

	tests := []struct {
		name    string
		article *model.Article
		user    *model.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Successful Deletion of Favorite Article",
			article: &model.Article{
				FavoritesCount: 1,
			},
			user: &model.User{},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Unable to delete an Article/User that do not exist",
			article: &model.Article{
				FavoritesCount: 1,
			},
			user:    &model.User{},
			mock:    func() {
				// TODO: Replace myDatabaseError
				myDBError := fmt.Errorf("DB error")
				mock.ExpectBegin()
      			mock.ExpectExec("^DELETE").WillReturnError(myDBError)
			},
			wantErr: true,
		},
		{
			name: "Failure of Update operation due to DB errors",
			article: &model.Article{
				FavoritesCount: 1,
			},
			user: &model.User{},
			mock: func() {
				// TODO: Replace myDatabaseError
				myDBError := fmt.Errorf("DB error")
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE").WillReturnError(myDBError)
			},
			wantErr: true,
		},
		{
			name: "Trying to delete a non-favorite article",
			article: &model.Article{
				FavoritesCount: 0,
			},
			user: &model.User{},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE").WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec("^UPDATE").WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			as := &ArticleStore{db: gormDB}
			tc.mock()
			
			err := as.DeleteFavorite(tc.article, tc.user)
			if tc.wantErr {
				assert.Error(t, err)
				t.Log(fmt.Sprintf("Case %s completed successfully: ", tc.name), err.Error())
				return
			}
			
			assert.NoError(t, err)
			t.Log(fmt.Sprintf("Case %s completed successfully ", tc.name))

		})
	}
}
