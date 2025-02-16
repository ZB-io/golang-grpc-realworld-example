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
	*gorm.DB
}
type MockArticleStore struct {
	mock.Mock
}
type MockDB struct {
	mock.Mock
}
type MockGormDB struct {
	mock.Mock
}

/*
ROOST_METHOD_HASH=GetFeedArticles_a37e1934b6
ROOST_METHOD_SIG_HASH=GetFeedArticles_f5f09c020e

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs [ // GetFeedArticles returns following users' articles
]uint, limit, offset int64) ([]model.Article, error)
*/
func (m *mockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m.DB
}

func TestArticleStoreGetFeedArticles(t *testing.T) {
	tests := []struct {
		name          string
		userIDs       []uint
		limit         int64
		offset        int64
		mockSetup     func(*mockDB)
		expectedError error
		expectedLen   int
	}{
		{
			name:    "Successful retrieval of feed articles",
			userIDs: []uint{1, 2, 3},
			limit:   10,
			offset:  0,
			mockSetup: func(m *mockDB) {
				articles := []model.Article{
					{Model: gorm.Model{ID: 1}, UserID: 1, Title: "Article 1"},
					{Model: gorm.Model{ID: 2}, UserID: 2, Title: "Article 2"},
					{Model: gorm.Model{ID: 3}, UserID: 3, Title: "Article 3"},
				}
				m.DB.Error = nil
				m.DB.Value = &articles
			},
			expectedError: nil,
			expectedLen:   3,
		},
		{
			name:    "Empty result when no followed users have articles",
			userIDs: []uint{4, 5, 6},
			limit:   10,
			offset:  0,
			mockSetup: func(m *mockDB) {
				var articles []model.Article
				m.DB.Error = nil
				m.DB.Value = &articles
			},
			expectedError: nil,
			expectedLen:   0,
		},
		{
			name:    "Error handling for database issues",
			userIDs: []uint{1, 2, 3},
			limit:   10,
			offset:  0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = errors.New("database error")
			},
			expectedError: errors.New("database error"),
			expectedLen:   0,
		},
		{
			name:    "Limit exceeding available articles",
			userIDs: []uint{1, 2},
			limit:   5,
			offset:  0,
			mockSetup: func(m *mockDB) {
				articles := []model.Article{
					{Model: gorm.Model{ID: 1}, UserID: 1, Title: "Article 1"},
					{Model: gorm.Model{ID: 2}, UserID: 2, Title: "Article 2"},
				}
				m.DB.Error = nil
				m.DB.Value = &articles
			},
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name:    "Offset beyond available articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  5,
			mockSetup: func(m *mockDB) {
				var articles []model.Article
				m.DB.Error = nil
				m.DB.Value = &articles
			},
			expectedError: nil,
			expectedLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{DB: &gorm.DB{}}
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB.DB}

			articles, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedLen, len(articles))

			if tt.expectedLen > 0 {
				for _, article := range articles {
					assert.NotNil(t, article.Author)
					assert.Contains(t, tt.userIDs, article.UserID)
				}
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

/*
ROOST_METHOD_HASH=IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
func (m *MockDB) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	args := m.Called(a, u)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		user           *model.User
		setupMock      func(*MockArticleStore)
		expectedResult bool
		expectedError  error
	}{
		{
			name:    "Article is favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mas *MockArticleStore) {
				mas.On("IsFavorited", mock.AnythingOfType("*model.Article"), mock.AnythingOfType("*model.User")).Return(true, nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:    "Article is not favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mas *MockArticleStore) {
				mas.On("IsFavorited", mock.AnythingOfType("*model.Article"), mock.AnythingOfType("*model.User")).Return(false, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:    "Database error occurs",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mas *MockArticleStore) {
				mas.On("IsFavorited", mock.AnythingOfType("*model.Article"), mock.AnythingOfType("*model.User")).Return(false, errors.New("database error"))
			},
			expectedResult: false,
			expectedError:  errors.New("database error"),
		},
		{
			name:    "Nil article parameter",
			article: nil,
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mas *MockArticleStore) {
				mas.On("IsFavorited", nil, mock.AnythingOfType("*model.User")).Return(false, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:    "Nil user parameter",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    nil,
			setupMock: func(mas *MockArticleStore) {
				mas.On("IsFavorited", mock.AnythingOfType("*model.Article"), nil).Return(false, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:    "Both article and user parameters are nil",
			article: nil,
			user:    nil,
			setupMock: func(mas *MockArticleStore) {
				mas.On("IsFavorited", nil, nil).Return(false, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleStore := new(MockArticleStore)
			tt.setupMock(mockArticleStore)

			result, err := mockArticleStore.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockArticleStore.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}
