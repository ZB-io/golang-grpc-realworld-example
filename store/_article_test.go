package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)





type mockDB struct {
	mock.Mock
}


/*
ROOST_METHOD_HASH=ArticleStore_CreateComment_b16d4a71d4
ROOST_METHOD_SIG_HASH=ArticleStore_CreateComment_7475736b06

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error // CreateComment creates a comment of the article


*/
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
			name: "Fail to Create Comment Due to Database Error",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: errors.New("database error"),
			wantErr: true,
		},
		{
			name:    "Attempt to Create Comment with Nil Model",
			comment: nil,
			dbError: nil,
			wantErr: true,
		},
		{
			name: "Create Comment with Maximum Length Content",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create Duplicate Comment",
			comment: &model.Comment{
				Body:      "Duplicate comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: errors.New("unique constraint violation"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			store := &ArticleStore{
				db: mockDB,
			}

			if tt.comment != nil {
				mockDB.On("Create", tt.comment).Return(&gorm.DB{Error: tt.dbError})
			}

			err := store.CreateComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.dbError != nil {
					assert.Equal(t, tt.dbError, err)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.comment != nil {
				mockDB.AssertCalled(t, "Create", tt.comment)
			} else {
				mockDB.AssertNotCalled(t, "Create")
			}
		})
	}
}

func (m *mockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

