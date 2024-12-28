package store

import (
	"reflect"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestArticleStoreGetComments(t *testing.T) {

	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		article model.Article
		want    []model.Comment
		wantErr bool
	}{
		{
			name: "Successful Retrieval of Comments",
			setup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("^SELECT \\* FROM comments WHERE (article_id = ?)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "author_id", "article_id"}).
						AddRow(1, "Test comment", 1, 1, 1).
						AddRow(2, "Another comment", 2, 2, 1))
			},
			article: model.Article{Model: gorm.Model{ID: 1}},
			want: []model.Comment{
				{Body: "Test comment", UserID: 1, ArticleID: 1},
				{Body: "Another comment", UserID: 2, ArticleID: 1},
			},
			wantErr: false,
		},
		{
			name: "Article Without Comments",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM comments WHERE (article_id = ?)").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "author_id", "article_id"}))
			},
			article: model.Article{Model: gorm.Model{ID: 2}},
			want:    []model.Comment{},
			wantErr: false,
		},
		{
			name: "Non-existent Article",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM comments WHERE (article_id = ?)").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "author_id", "article_id"}))
			},
			article: model.Article{Model: gorm.Model{ID: 3}},
			want:    []model.Comment{},
			wantErr: false,
		},
		{
			name: "Database Error Occurrence",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM comments WHERE (article_id = ?)").
					WithArgs(4).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			article: model.Article{Model: gorm.Model{ID: 4}},
			want:    []model.Comment{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlmock", db)
			if err != nil {
				t.Fatalf("failed to initialize GORM: %v", err)
			}

			articleStore := ArticleStore{db: gormDB}

			tt.setup(mock)

			comments, err := articleStore.GetComments(&tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(comments, tt.want) {
				t.Errorf("GetComments() = %v, want %v", comments, tt.want)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' successful. Retrieve: %+v with errors: %v", tt.name, comments, err)
		})
	}
}
