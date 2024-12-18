package store_test

import (
	"fmt"
	"testing"
	"reflect"
	
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/require"
)

// TestArticleStoreGetByID tests the GetByID method of the ArticleStore
func TestArticleStoreGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	require.NoError(t, err)

	s := store.ArticleStore{DB: gormDB}    // Modification: It's case sensitive and the field name is `DB` not `db`.

	testCases := []struct {
		desc     string
		id       uint
		setup    func(mock sqlmock.Sqlmock, id uint)
		expected *model.Article
		hasError bool
	}{
		{
			desc: "Normal operation",
			id:   1,
			setup: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"ID", "title"})
        			expected := &model.Article{ID: int32(1), Title: "Test Article"} // Modification: Assumed `ID` is of type `int32` not `uint`.
        			rows.AddRow(expected.ID, expected.Title)
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WillReturnRows(rows)
			},
			expected: &model.Article{ID: int32(1), Title: "Test Article"}, // Modification: Assumed `ID` is of type `int32` not `uint`.
			hasError: false,
		},
		{
			desc: "Non-existing ID",
			id:   2,
			setup: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WillReturnError(gorm.ErrRecordNotFound)
			},
			expected: nil,
			hasError: true,
		},
		{
			desc: "Database error",
			id:   3,
			setup: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WillReturnError(fmt.Errorf("database error"))
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.setup(mock, tC.id)
			actual, err := s.GetByID(tC.id)

			if (err != nil) != tC.hasError {
				t.Fatalf("expected error: %v, got: %v, error: %v", tC.hasError, (err != nil), err)
			}

			if !reflect.DeepEqual(actual, tC.expected) {
				t.Fatalf("expected: %#v, got: %#v", tC.expected, actual)
			}
		})
	}
}
