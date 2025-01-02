package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)





func TestArticleStoreGetFeedArticles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("failed to open gorm DB: %s", err)
	}

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name    string
		userIDs []uint
		limit   int64
		offset  int64
		setup   func()
		check   func(articles []model.Article, err error)
	}{
		{
			name:    "Retrieve valid feed articles",
			userIDs: []uint{1, 2, 3},
			limit:   5,
			offset:  0,
			setup: func() {

			},
			check: func(articles []model.Article, err error) {

			},
		},
		{
			name:    "Retrieve feed articles with an empty list of userIDs",
			userIDs: []uint{},
			limit:   5,
			offset:  0,
			setup: func() {

			},
			check: func(articles []model.Article, err error) {
				assert.Nil(t, err)
				assert.Empty(t, articles)
			},
		},
		{
			name:    "Retrieve feed articles beyond the available limit",
			userIDs: []uint{1, 2, 3},
			limit:   100,
			offset:  0,
			setup: func() {

			},
			check: func(articles []model.Article, err error) {

			},
		},
		{
			name:    "Error retrieving feed articles",
			userIDs: []uint{1, 2, 3},
			limit:   5,
			offset:  0,
			setup: func() {

			},
			check: func(articles []model.Article, err error) {
				assert.Error(t, err)
				assert.Nil(t, articles)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			articles, err := store.GetFeedArticles(test.userIDs, test.limit, test.offset)
			test.check(articles, err)
			mock.ExpectationsWereMet()
		})
	}
}
