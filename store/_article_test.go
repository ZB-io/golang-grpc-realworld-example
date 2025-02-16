package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





type mockDB struct {
	createError error
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error 

*/
func (m *mockDB) Create(value interface{}) *gorm.DB {
	return &gorm.DB{Error: m.createError}
}

func TestArticleStoreCreateComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		dbError error
		wantErr bool
	}{
		{
			name: "Successfully Create a New Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			comment: &model.Comment{
				Body:      "",
				UserID:    0,
				ArticleID: 1,
			},
			dbError: errors.New("missing required fields"),
			wantErr: true,
		},
		{
			name: "Create a Comment with Maximum Length Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent Article",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 9999,
			},
			dbError: errors.New("foreign key constraint violation"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{createError: tt.dbError}
			store := &ArticleStore{db: mockDB}

			err := store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("Create Multiple Comments in Quick Succession", func(t *testing.T) {
		mockDB := &mockDB{createError: nil}
		store := &ArticleStore{db: mockDB}

		comments := []*model.Comment{
			{Body: "Comment 1", UserID: 1, ArticleID: 1},
			{Body: "Comment 2", UserID: 2, ArticleID: 1},
			{Body: "Comment 3", UserID: 3, ArticleID: 2},
		}

		for _, comment := range comments {
			err := store.CreateComment(comment)
			if err != nil {
				t.Errorf("CreateComment() failed for comment: %v, error: %v", comment, err)
			}
		}
	})
}

