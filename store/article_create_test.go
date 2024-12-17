package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var createTestCases = []struct {
	name		string
	article		*model.Article
	expected	error
}{
	{
		name:		"Valid Article Input",
		article:	&model.Article{Title: "First Article", Description: "This is a test article", Body: "Test Body", UserID: 1},
		expected:	nil,
	},
	{
		name:		"Null values for not null fields",
		article:	&model.Article{Description: "This is a test article", Body: "Test Body", UserID: 1},
		expected:	gorm.ErrRecordNotFound,
	},
	{
		name:		"Invalid UserID Input",
		article:	&model.Article{Title: "First Article", Description: "This is a test article", Body: "Test Body", UserID: -1},
		expected:	gorm.ErrRecordNotFound,
	},
	{
		name:		"Unique constraints violation",
		article:	&model.Article{Title: "First Article", Description: "This is a test article", Body: "Test Body", UserID: 1},
		expected:	gorm.ErrRecordNotFound,
	},
}func TestDataStore_CreateArticle(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("mysql", db)
	defer func() {
		_ = db.Close()
	}()
	store := ArticleStore{db: gormDB}
	mock.ExpectBegin()
	for _, test := range createTestCases {
		t.Run(test.name, func(t *testing.T) {
			t.Log(test.name)
			err := store.Create(test.article)
			if assert.Error(t, err) {
				assert.Equal(t, test.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	mock.ExpectCommit()
}
