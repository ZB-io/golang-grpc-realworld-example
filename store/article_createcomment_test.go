package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
)




func TestArticleStoreCreateComment(t *testing.T) {

	testCases := []struct {
		name      string
		comment   *model.Comment
		wantError bool
	}{
		{
			name: "Successful Comment Creation",
			comment: &model.Comment{
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantError: false,
		},
		{
			name: "Comment Creation with Missing Required Fields",
			comment: &model.Comment{
				Body: "Test Comment",
			},
			wantError: true,
		},
		{
			name: "Comment Creation with Invalid User ID",
			comment: &model.Comment{
				Body:      "Test Comment",
				UserID:    0,
				ArticleID: 1,
			},
			wantError: true,
		},
		{
			name: "Comment Creation with Invalid Article ID",
			comment: &model.Comment{
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 0,
			},
			wantError: true,
		},
		{
			name: "Comment Creation with Empty Body",
			comment: &model.Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			wantError: true,
		},
	}

	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			store := &ArticleStore{db: gormDB}

			if tc.wantError {
				mock.ExpectBegin()
				mock.ExpectRollback()
			} else {
				mock.ExpectBegin()
				mock.ExpectCommit()
			}

			err := store.CreateComment(tc.comment)

			if tc.wantError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}
			}
		})
	}
}
