package store

import (
	"testing"
	"database/sql"
	"github.com/jinzhu/gorm"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreGetComments(t *testing.T) {
	// setting up the sql mock database
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open("postgres", sqlDB)
	if err != nil {
		t.Fatalf("can't open gorm connection: %s", err)
	}
	defer gormDB.Close()

	// creating an instance of our ArticleStore with mock db
	store := ArticleStore{db: gormDB}

	rows := sqlmock.NewRows([]string{"id", "body", "article_id", "author"}).
		AddRow(1, "comment body", 1, "author")

	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	tests := []struct {
		name    string
		article model.Article
		wantErr bool
	}{
		{
			"Scenario 1: Retrieving Comments of a Valid Article",
			model.Article{ID: 1},
			false,
		},
		{
			"Scenario 2: Retrieving Comments of an Article with No Comments",
			model.Article{ID: 2},
			false,
		},
		{
			"Scenario 3: Retrieving Comments of an Invalid Article",
			model.Article{ID: 3},
			true,
		}, 
		{
			"Scenario 4: Database Error During Retrieval of Comments",
			model.Article{ID: 4},  
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.GetComments(&tt.article)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
