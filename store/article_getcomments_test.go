package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// Test getting comments with valid article
func TestArticleStoreGetComments(t *testing.T) {
	// Initialize DB & Mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %s", err)
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm db: %s", err)
	}

	// Initialize ArticleStore with mock DB so its methods can be tested
	articleStore := NewArticleStore(gdb)

	// Define dummy Article and Comments
	article := &model.Article{Model: gorm.Model{ID: 1}}
	comments := []model.Comment{{Model: gorm.Model{ID: 1}, Body: "Test Comment", UserID: 123, ArticleID: 1}}

	// Create Rows for mocking DB response
	rows := sqlmock.
		NewRows([]string{"id", "body", "user_id", "article_id"}).
		AddRow(comments[0].ID, comments[0].Body, comments[0].UserID, comments[0].ArticleID)

	mock.ExpectQuery(`SELECT(.+)FROM "comments"`).WithArgs(article.ID).WillReturnRows(rows)

	// Perform test and check result
	result, err := articleStore.GetComments(article)
	if err != nil {
		t.Errorf("Error while getting comments: %s", err)
	}
	assert.Equal(t, comments, result)
}

// Test getting comments with article that has no comments
func TestArticleStoreGetCommentsNoComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %s", err)
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm db: %s", err)
	}

	articleStore := NewArticleStore(gdb)

	// It returns nothing as there are no comments
	rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"})

	article := &model.Article{Model: gorm.Model{ID: 1}}

	mock.ExpectQuery(`SELECT(.+)FROM "comments"`).WithArgs(article.ID).WillReturnRows(rows)

	result, err := articleStore.GetComments(article)
	if err != nil {
		t.Errorf("Error while getting comments: %s", err)
	}
	assert.Empty(t, result)
}

// Test getting comments with invalid article
func TestArticleStoreGetCommentsInvalidArticle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %s", err)
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm db: %s", err)
	}

	articleStore := NewArticleStore(gdb)

	invalidArticle := &model.Article{Model: gorm.Model{ID: 0}}
	mock.ExpectQuery(`SELECT(.+)FROM "comments"`).WithArgs(invalidArticle.ID).WillReturnError(gorm.ErrRecordNotFound)

	_, err = articleStore.GetComments(invalidArticle)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// Test GetComments when database is not accessible
func TestArticleStoreGetCommentsDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %s", err)
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm db: %s", err)
	}

	articleStore := NewArticleStore(gdb)

	// Create some random database error
	dbError := errors.New("DB connection failed")

	article := &model.Article{Model: gorm.Model{ID: 1}}
	mock.ExpectQuery(`SELECT(.+)FROM "comments"`).WithArgs(article.ID).WillReturnError(dbError)

	_, err = articleStore.GetComments(article)
	assert.Equal(t, dbError, err)
}
