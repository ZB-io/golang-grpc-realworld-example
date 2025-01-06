package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestArticleStoreGetComments is a test function for the GetComments method of the ArticleStore type
func TestArticleStoreGetComments(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		prepare   func(mock sqlmock.Sqlmock, article *model.Article)
		article   *model.Article
		wantError bool
	}{
		{
			name: "Successful retrieval of comments for a given article",
			prepare: func(mock sqlmock.Sqlmock, article *model.Article) {
				rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
					AddRow(1, "Test comment 1", 1, article.ID).
					AddRow(2, "Test comment 2", 1, article.ID)
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE (.+)").WillReturnRows(rows)
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}},
			wantError: false,
		},
		{
			name: "No comments for a given article",
			prepare: func(mock sqlmock.Sqlmock, article *model.Article) {
				rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"})
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE (.+)").WillReturnRows(rows)
			},
			article:   &model.Article{Model: gorm.Model{ID: 2}},
			wantError: false,
		},
		{
			name: "Error in database query",
			prepare: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE (.+)").WillReturnError(gorm.ErrRecordNotFound)
			},
			article:   &model.Article{Model: gorm.Model{ID: 3}},
			wantError: true,
		},
		{
			name: "Article not found",
			prepare: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE (.+)").WillReturnError(gorm.ErrRecordNotFound)
			},
			article:   &model.Article{Model: gorm.Model{ID: 4}},
			wantError: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormDB, _ := gorm.Open("postgres", db)
			tt.prepare(mock, tt.article)
			store := &ArticleStore{db: gormDB}

			comments, err := store.GetComments(tt.article)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, len(tt.article.Comments))
			}
		})
	}
}
