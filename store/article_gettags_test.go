package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)



func TestArticleStoreGetTags(t *testing.T) {

	testCases := []struct {
		name       string
		mockDBFunc func(mock sqlmock.Sqlmock)
		wantTags   []model.Tag
		wantErr    bool
	}{
		{
			name: "Successful Retrieval of Tags",
			mockDBFunc: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows([]string{"ID", "Name"}).
					AddRow(1, "tag1").
					AddRow(2, "tag2")
				mock.ExpectQuery("^SELECT (.+) FROM \"tags\"").WillReturnRows(rows)
			},
			wantTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
			},
			wantErr: false,
		},
		{
			name: "Database Error during Tag Retrieval",
			mockDBFunc: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("^SELECT (.+) FROM \"tags\"").WillReturnError(errors.New("db error"))
			},
			wantTags: nil,
			wantErr:  true,
		},
		{
			name: "Empty Database",
			mockDBFunc: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows([]string{"ID", "Name"})
				mock.ExpectQuery("^SELECT (.+) FROM \"tags\"").WillReturnRows(rows)
			},
			wantTags: []model.Tag{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			gormDB, err := gorm.Open("postgres", db)
			assert.NoError(t, err)
			tc.mockDBFunc(mock)

			store := &ArticleStore{db: gormDB}
			gotTags, err := store.GetTags()

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantTags, gotTags)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
