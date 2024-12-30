package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestGetByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening stub database connection: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("unexpected error when opening gorm DB: %s", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	t.Run("Successful Retrieval of an Existing Article", func(t *testing.T) {

		articleID := uint(1)
		expectedArticle := model.Article{
			Title:       "Test Article",
			Description: "Description",
			Body:        "Body",
			UserID:      1,
		}

		rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
			AddRow(articleID, expectedArticle.Title, expectedArticle.Description, expectedArticle.Body, expectedArticle.UserID)

		mock.ExpectQuery("SELECT (.+) FROM \"articles\" WHERE (.+)").
			WithArgs(articleID).
			WillReturnRows(rows)

		article, err := store.GetByID(articleID)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if article.Title != expectedArticle.Title || article.Description != expectedArticle.Description ||
			article.Body != expectedArticle.Body || article.UserID != expectedArticle.UserID {
			t.Errorf("article details do not match expected values")
		}
	})

	t.Run("Article Not Found", func(t *testing.T) {

		articleID := uint(999)

		mock.ExpectQuery("SELECT (.+) FROM \"articles\" WHERE (.+)").
			WithArgs(articleID).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := store.GetByID(articleID)

		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("expected record not found error, got: %s", err)
		}
	})

	t.Run("Database Connection Error", func(t *testing.T) {

		articleID := uint(1)

		mock.ExpectQuery("SELECT (.+) FROM \"articles\" WHERE (.+)").
			WithArgs(articleID).
			WillReturnError(errors.New("connection refused"))

		_, err := store.GetByID(articleID)

		if err == nil || err.Error() != "connection refused" {
			t.Errorf("expected connection refused error, got: %s", err)
		}
	})

	t.Run("Retrieval with Preloaded Associations", func(t *testing.T) {

		articleID := uint(2)
		expectedTitle := "Article with Relations"

		rows := sqlmock.NewRows([]string{"id", "title"}).
			AddRow(articleID, expectedTitle)
		mock.ExpectQuery("SELECT (.+) FROM \"articles\" WHERE (.+)").
			WithArgs(articleID).
			WillReturnRows(rows)

		mock.ExpectQuery("SELECT (.+) FROM \"tags\"").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Technology"))

		mock.ExpectQuery("SELECT (.+) FROM \"users\"").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Author Name"))

		article, err := store.GetByID(articleID)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if article.Title != expectedTitle {
			t.Errorf("expected title: %s, got: %s", expectedTitle, article.Title)
		}

		if len(article.Tags) != 1 || article.Tags[0].Name != "Technology" {
			t.Errorf("expected tags to include 'Technology'")
		}

		if article.Author.UserName != "Author Name" {
			t.Errorf("expected author name: %s, got: %s", "Author Name", article.Author.UserName)
		}
	})

	t.Run("Invalid ID Type", func(t *testing.T) {

		invalidID := uint(0)

		_, err := store.GetByID(invalidID)

		if err == nil {
			t.Errorf("expected error for invalid ID, got nil")
		}
	})

	t.Run("Article with No Tags or Comments", func(t *testing.T) {

		articleID := uint(3)
		expectedArticle := model.Article{
			Title:       "Article without Tags or Comments",
			Description: "Description without extras",
			Body:        "Body content",
			UserID:      3,
		}

		rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
			AddRow(articleID, expectedArticle.Title, expectedArticle.Description, expectedArticle.Body, expectedArticle.UserID)

		mock.ExpectQuery("SELECT (.+) FROM \"articles\" WHERE (.+)").
			WithArgs(articleID).
			WillReturnRows(rows)

		article, err := store.GetByID(articleID)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if article.Title != expectedArticle.Title || article.Description != expectedArticle.Description ||
			article.Body != expectedArticle.Body || article.UserID != expectedArticle.UserID {
			t.Errorf("article details do not match expected values")
		}

		if len(article.Tags) != 0 {
			t.Errorf("expected zero tags, found: %d", len(article.Tags))
		}
		if len(article.Comments) != 0 {
			t.Errorf("expected zero comments, found: %d", len(article.Comments))
		}
	})
}
