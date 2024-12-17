package store

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

func TestArticleStoreAddFavorite(t *testing.T) {
	db, mock, err := sqlmock.New() // Initialize mock db
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gormDB, _ := gorm.Open("postgres", db)

	// Define test data
	testData := []struct {
		key         string
		favCount    int
		errMsg      string
		expectedErr error
	}{
		{"normal", 1, "", nil},
		{"FavoritedUsers association error", 0, "error at FavoritedUsers association", fmt.Errorf("error at FavoritedUsers association")},
		{"favorites_count update error", 0, "error updating favorites_count", fmt.Errorf("error updating favorites_count")},
	}

	// Implement table-driven tests
	for _, data := range testData {
		t.Logf("Exercising: %s", data.key)

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO \"favorite_articles\"").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"a", "u"}).AddRow(data.favCount, data.favCount))
		mock.ExpectCommit()

		article := &model.Article{
			Title:       "Test Article",
			Description: "Test Description",
			Body:        "Test Body",
		}

		user := &model.User{
			Username: "testuser",
			Email:    "testuser@test.com",
			Password: "testpassword",
		}

		as := ArticleStore{gormDB}
		err := as.AddFavorite(article, user)

		// If there is no error expected and we get no error or if the err message matches
		// the expected error message, log the success
		if (data.expectedErr == nil && err == nil) || (err != nil && err.Error() == data.errMsg) {
			t.Logf("Test %s passed", data.key)
		} else {
			t.Fatalf("Test %s failed: %s", data.key, err) // Log detailed failure reason for diagnostic clarity
		}
	}
}
