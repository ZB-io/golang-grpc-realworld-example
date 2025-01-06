package store

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Association(column string) *gorm.Association {
	args := m.Called()
	return args.Get(0).(*gorm.Association)
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

type ArticleStore struct {
	db gorm.SQLCommon
}

func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error {
	tx := s.db.Begin()

	err := tx.Model(a).Association("FavoritedUsers").
		Append(u).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(a).
		Update("favorites_count", gorm.Expr("favorites_count + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	a.FavoritesCount++

	return nil
}

func TestArticleStoreAddFavorite(t *testing.T) {
	// Initialize test scenarios
	testCases := []struct {
		name     string
		article  model.Article
		user     model.User
		mockFunc func(mockDB *MockDB)
		wantErr  bool
	}{
		{
			name:    "Adding a new favorite article successfully",
			article: model.Article{FavoritesCount: 0},
			user:    model.User{},
			mockFunc: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
				mockDB.On("Update", "favorites_count", gorm.Expr("favorites_count + ?", 1)).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
			},
			wantErr: false,
		},
		{
			name:    "Error in adding the user to the article's FavoritedUsers",
			article: model.Article{FavoritesCount: 0},
			user:    model.User{},
			mockFunc: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{Error: errors.New("error")})
			},
			wantErr: true,
		},
		{
			name:    "Error in updating the article's favorites count",
			article: model.Article{FavoritesCount: 0},
			user:    model.User{},
			mockFunc: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
				mockDB.On("Update", "favorites_count", gorm.Expr("favorites_count + ?", 1)).Return(&gorm.DB{Error: errors.New("error")})
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockDB := new(MockDB)
			tc.mockFunc(mockDB)
			store := &ArticleStore{db: mockDB}

			// Act
			err := store.AddFavorite(&tc.article, &tc.user)

			// Assert
			if tc.wantErr {
				assert.Error(t, err)
				t.Log("Error is expected when", tc.name)
			} else {
				assert.NoError(t, err)
				t.Log("No error is expected when", tc.name)
				assert.Equal(t, tc.article.FavoritesCount, int32(1))
				t.Log("FavoritesCount is expected to be incremented by 1 when", tc.name)
			}
		})
	}
}
