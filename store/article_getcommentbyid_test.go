package store

import (
	"testing"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)






func TestArticleStoreGetCommentByID(t *testing.T) {

	testCases := []struct {
		name            string
		mockDbSetup     func(mock sqlmock.Sqlmock, id uint)
		id              uint
		expectedComment *model.Comment
		expectedError   error
	}{
		{
			name: "Get Comment by Valid ID",
			mockDbSetup: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
					AddRow(id, "Test body", 1, 1)
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE (.+)").
					WithArgs(id).WillReturnRows(rows)
			},
			id: 1,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test body",
				UserID:    1,
				ArticleID: 1,
			},
			expectedError: nil,
		},
		{
			name: "Get Comment by Invalid ID",
			mockDbSetup: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE (.+)").
					WithArgs(id).WillReturnError(gorm.ErrRecordNotFound)
			},
			id:              999,
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
		{
			name: "Get Comment by ID with Empty Database",
			mockDbSetup: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE (.+)").
					WithArgs(id).WillReturnError(gorm.ErrRecordNotFound)
			},
			id:              1,
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm db: %v", err)
			}

			tc.mockDbSetup(mock, tc.id)

			store := &ArticleStore{db: gormDB}

			comment, err := store.GetCommentByID(tc.id)

			if !equalComments(comment, tc.expectedComment) {
				t.Errorf("Expected comment %v, got %v", tc.expectedComment, comment)
			}
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
func equalComments(c1, c2 *model.Comment) bool {
	if c1 == nil || c2 == nil {
		return c1 == c2
	}
	return c1.ID == c2.ID && c1.Body == c2.Body && c1.UserID == c2.UserID && c1.ArticleID == c2.ArticleID
}
