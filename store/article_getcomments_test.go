package store

import (
	"errors"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)




type Call struct {
	Parent *Mock

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// Holds the caller info for the On() call
	callerInfo []string

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Amount of times this call has been called
	totalCalls int

	// Call to this method can be optional
	optional bool

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	RunFn func(Arguments)
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}



type MockDB struct {
	mock.Mock
}








type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}
func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}
func TestGetComments(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		mockSetup      func(*MockDB)
		expectedResult []model.Comment
		expectedError  error
	}{
		{
			name: "Successfully retrieve comments for an article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(1)).Return(mockDB)
				mockDB.On("Find", mock.AnythingOfType("*[]model.Comment")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Comment)
					*arg = []model.Comment{
						{Model: gorm.Model{ID: 1}, Body: "Comment 1", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
						{Model: gorm.Model{ID: 2}, Body: "Comment 2", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
					}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 1}, Body: "Comment 1", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Body: "Comment 2", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			dbWrapper := struct {
				*MockDB
			}{mockDB}

			store := &ArticleStore{db: &dbWrapper}

			result, err := store.GetComments(tt.article)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)

			mockDB.AssertExpectations(t)
		})
	}
}
func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*gorm.DB)
}
