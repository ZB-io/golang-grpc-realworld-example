package store

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





type mockDB struct {
	createFunc func(interface{}) *gorm.DB
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error 

*/
func (m *mockDB) Create(value interface{}) *gorm.DB {
	return m.createFunc(value)
}

func TestArticleStoreCreateComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		mockDB  *mockDB
		wantErr bool
	}{
		{
			name: "Successfully Create a New Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockDB: &mockDB{
				createFunc: func(value interface{}) *gorm.DB {
					return &gorm.DB{Error: nil}
				},
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			comment: &model.Comment{
				Body: "",
			},
			mockDB: &mockDB{
				createFunc: func(value interface{}) *gorm.DB {
					return &gorm.DB{Error: gorm.ErrRecordNotFound}
				},
			},
			wantErr: true,
		},
		{
			name: "Create a Comment with Maximum Length Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 65535)),
				UserID:    1,
				ArticleID: 1,
			},
			mockDB: &mockDB{
				createFunc: func(value interface{}) *gorm.DB {
					return &gorm.DB{Error: nil}
				},
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent Article",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 9999,
			},
			mockDB: &mockDB{
				createFunc: func(value interface{}) *gorm.DB {
					return &gorm.DB{Error: gorm.ErrRecordNotFound}
				},
			},
			wantErr: true,
		},
		{
			name: "Create a Comment with Special Characters in the Body",
			comment: &model.Comment{
				Body:      "Test comment with special characters: !@#$%^&*()_+ and emojis: 😀🎉",
				UserID:    1,
				ArticleID: 1,
			},
			mockDB: &mockDB{
				createFunc: func(value interface{}) *gorm.DB {
					return &gorm.DB{Error: nil}
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB,
			}

			err := s.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

