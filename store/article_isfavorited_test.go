package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArticleStoreIsFavorited(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open("postgres", db)
	defer db.Close()

	articleStore := &ArticleStore{db: gdb}

	var testID uint = 1
	var testArticle = &model.Article{Model: gorm.Model{ID: testID}}
	var testUser = &model.User{Model: gorm.Model{ID: testID}}

	tests := []struct {
		name       string
		article    *model.Article
		user       *model.User
		mock       func()
		wantResult bool
		wantErr    error
	}{
		{
			name:       "Article is Favorited by User",
			article:    testArticle,
			user:       testUser,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"favorite_articles\"*").
					WillReturnRows(rows)
			},
			wantResult: true,
			wantErr:    nil,
		},
		{
			name:       "Article is not Favorited by User",
			article:    testArticle,
			user:       testUser,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"favorite_articles\"*").
					WillReturnRows(rows)
			},
			wantResult: false,
			wantErr:    nil,
		},
		{
			name:       "Article or User Parameters are Nil",
			article:    nil,
			user:       nil,
			mock:       func() {},
			wantResult: false,
			wantErr:    nil,
		},
		{
			name:       "DataBase Error",
			article:    testArticle,
			user:       testUser,
			mock: func() {
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"favorite_articles\"*").
					WillReturnError(errors.New("database error"))
			},
			wantResult: false,
			wantErr:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			result, err := articleStore.IsFavorited(tt.article, tt.user)

			if (err != nil) && (err.Error() != tt.wantErr.Error()) {
				t.Errorf("IsFavorited() returned err: %v, want: %v", err, tt.wantErr)
				return
			}
			if result != tt.wantResult {
				t.Errorf("IsFavorited() returned result: %v, want: %v", result, tt.wantResult)
			}

			// Assert that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
