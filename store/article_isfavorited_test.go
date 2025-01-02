package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)


func TestArticleStoreIsFavorited(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
	}

	testCases := []struct {
		desc    string
		article *model.Article
		user    *model.User
		mock    func()
		want    bool
		wantErr bool
	}{
		{
			desc: "Test when both parameters are nil",
			mock: func() {},
			want: false,
		},
		{
			desc: "Test when the article is nil and a valid user is passed",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mock: func() {},
			want: false,
		},
		{
			desc:    "Test when a valid article is passed and the user is nil",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			mock:    func() {},
			want:    false,
		},
		{
			desc:    "Test when both the article and user are valid but the user has not favorited the article",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mock: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles" WHERE \(article_id = \$1 AND user_id = \$2\)`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want: false,
		},
		{
			desc:    "Test when both the article and user are valid and the user has favorited the article",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mock: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles" WHERE \(article_id = \$1 AND user_id = \$2\)`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want: true,
		},
		{
			desc:    "Test when a database error occurs",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mock: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles" WHERE \(article_id = \$1 AND user_id = \$2\)`).
					WithArgs(1, 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.mock()

			s := &ArticleStore{db: gormDB}
			got, err := s.IsFavorited(tC.article, tC.user)
			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.want, got)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
