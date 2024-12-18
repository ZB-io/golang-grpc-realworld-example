package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

func TestArticleStoreUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("failed to open the stub database connection: %s", err)
	}
	// creating a new store
	store := ArticleStore{db: gdb}

	// defining test cases
	testCases := [] struct {
		desc string
		id int
		article model.Article
		expectErr bool
	}{
		{
			desc: "Successful Update of Existing Article",
			id: 1,
			article: model.Article{Title: "Updated Title"},
			expectErr: false,
		},
		{
			desc: "Update Attempt on Non-Existent Article",
			id: 99999, // non-existent ID
			article: model.Article{Title: "New Title"},
			expectErr: true,
		},
		{
			desc: "Attempt to Update Article without Required Fields",
			id: 1,
			article: model.Article{Title: ""},
			expectErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Log(tC.desc)

			// Mocking an expected call to DB
			mock.ExpectExec("UPDATE").WithArgs(tC.article.Title).WillReturnResult(sqlmock.NewResult(1, 1))

			err := store.Update(&tC.article)

			if tC.expectErr {
				if err == nil {
					t.Errorf("expected an error but did not get one")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %s", err)
				}
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
