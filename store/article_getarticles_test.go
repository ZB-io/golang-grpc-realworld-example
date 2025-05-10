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
func TestGetArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}
	defer gormDB.Close()

	articleStore := &ArticleStore{db: gormDB}

	tests := []struct {
		name           string
		tagName        string
		username       string
		favoritedBy    *model.User
		limit          int64
		offset         int64
		arrangeMocks   func()
		expectedError  error
		expectedResult []model.Article
	}{
		{
			name:     "Retrieve Articles by Author Username",
			username: "author1",
			arrangeMocks: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles" JOIN users on articles.user_id = users.id WHERE users.username = \?`).
					WithArgs("author1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).AddRow(1, "Title1", "Description1", "Body1"))
			},
			expectedError: nil,
			expectedResult: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Title1", Description: "Description1", Body: "Body1"},
			},
		},
		{
			name:    "Retrieve Articles by Tag Name",
			tagName: "tag1",
			arrangeMocks: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles" JOIN article_tags on articles.id = article_tags.article_id JOIN tags on tags.id = article_tags.tag_id WHERE tags.name = \?`).
					WithArgs("tag1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).AddRow(2, "Title2", "Description2", "Body2"))
			},
			expectedError: nil,
			expectedResult: []model.Article{
				{Model: gorm.Model{ID: 2}, Title: "Title2", Description: "Description2", Body: "Body2"},
			},
		},
		{
			name: "Retrieve Articles Favorited by a User",
			favoritedBy: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
			},
			arrangeMocks: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id in \(\?\)`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).AddRow(3, "Title3", "Description3", "Body3"))
			},
			expectedError: nil,
			expectedResult: []model.Article{
				{Model: gorm.Model{ID: 3}, Title: "Title3", Description: "Description3", Body: "Body3"},
			},
		},
		{
			name:   "Retrieve Articles with Pagination",
			limit:  1,
			offset: 1,
			arrangeMocks: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles" LIMIT ? OFFSET ?`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).AddRow(4, "Title4", "Description4", "Body4"))
			},
			expectedError: nil,
			expectedResult: []model.Article{
				{Model: gorm.Model{ID: 4}, Title: "Title4", Description: "Description4", Body: "Body4"},
			},
		},
		{
			name: "No Articles Match Criteria",
			arrangeMocks: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedError:  nil,
			expectedResult: []model.Article{},
		},
		{
			name: "Error Handling with Database Failures",
			arrangeMocks: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WillReturnError(errors.New("simulated db error"))
			},
			expectedError:  errors.New("simulated db error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.arrangeMocks()

			result, err := articleStore.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected result length %v, got %v", len(tt.expectedResult), len(result))
			}
			for i, article := range result {
				if article.Title != tt.expectedResult[i].Title {
					t.Errorf("expected article title %v, got %v", tt.expectedResult[i].Title, article.Title)
				}
			}
			t.Log("Test executed successfully")
		})
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
