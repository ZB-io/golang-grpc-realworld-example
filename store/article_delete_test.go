package store_test

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

/*
type ArticleStore struct {
	db *gorm.DB
}
*/

func TestDelete(t *testing.T) {
	const table = "articles"
	tests := []struct {
		name   string
		setup  func(mock sqlmock.Sqlmock)
		model  *model.Article
		wantErr bool
	}{
		{
			name: "Correct deletion of existing article",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(fmt.Sprintf("DELETE FROM %s", table)).
					WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			model:  &model.Article{Model: gorm.Model{ID: 1}},
			wantErr: false,
		},
		{
			name: "Attempt deletion of non-existent article",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(fmt.Sprintf("DELETE FROM %s", table)).
					WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			model:  &model.Article{Model: gorm.Model{ID: 2}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setup(mock)

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}
			defer gdb.Close()

			s := store.ArticleStore{db: gdb}

			if err := s.Delete(tt.model); (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
func (s *ArticleStore) Delete(m *model.Article) error {
	return s.db.Delete(m).Error
}
*/
