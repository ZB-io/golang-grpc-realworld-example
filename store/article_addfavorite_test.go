package store

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

// Mock DB to capture the transaction calls
type MockDB struct {
	mockModel        *model.Article
	mockAssociation  *gorm.Association
	mockUser         *model.User
	rollbackTriggered bool
	commitTriggered  bool
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	m.mockModel = value.(*model.Article)
	return &gorm.DB{}
}

func (m *MockDB) Update(column string, value ...interface{}) *gorm.DB {
	if column == "favorites_count" && m.mockModel != nil {
		m.mockModel.FavoritesCount++
	}
	return &gorm.DB{}
}

func (m *MockDB) Begin() *gorm.DB {
	return &gorm.DB{}
}

func (m *MockDB) Rollback() *gorm.DB {
	m.rollbackTriggered = true
	return &gorm.DB{}
}

func (m *MockDB) Commit() *gorm.DB {
	m.commitTriggered = true
	return &gorm.DB{}
}

func (m *MockDB) Association(column string) *gorm.Association {
	return m.mockAssociation
}

// Mock Association to capture append calls
type MockAssociation struct {
	appendError error
}

func (m *MockAssociation) Append(values ...interface{}) *gorm.Association {
	return &gorm.Association{Error: m.appendError}
}

type ArticleStore struct {
	db *MockDB
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


// Here is the implementation of TestArticleStoreAddFavorite test function
func TestArticleStoreAddFavorite(t *testing.T) {

	tests := []struct {
		name              string
		article           *model.Article
		user              *model.User
		appendError       error
		expectedErr       error
		expectedRollback  bool
		expectedCommit    bool
		expectedIncrement bool
	}{
		{
			name:              "Successful addition of a favorite article by a user",
			user:              &model.User{},
			article:           &model.Article{},
			appendError:       nil,
			expectedErr:       nil,
			expectedRollback:  false,
			expectedCommit:    true,
			expectedIncrement: true,
		},
		{
			name:              "Rollback transaction on error while appending to the Association FavoritedUsers",
			user:              &model.User{},
			article:           &model.Article{},
			appendError:       fmt.Errorf("some error"),
			expectedErr:       fmt.Errorf("some error"),
			expectedRollback:  true,
			expectedCommit:    false,
			expectedIncrement: false,
		},
		{
			name:              "Unsuccessful operation due to nil Article",
			user:              &model.User{},
			article:           nil,
			appendError:       nil,
			expectedErr:       errors.New("Article is nil"),
			expectedRollback:  false,
			expectedCommit:    false,
			expectedIncrement: false,
		},
		{
			name:              "Unsuccessful operation due to nil User",
			user:              nil,
			article:           &model.Article{},
			appendError:       nil,
			expectedErr:       errors.New("User is nil"),
			expectedRollback:  false,
			expectedCommit:    false,
			expectedIncrement: false,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			db := &MockDB{
				mockAssociation: &MockAssociation{appendError: test.appendError},
				mockModel:       test.article,
			}

			articleStore := ArticleStore{db: db}

			err := articleStore.AddFavorite(test.article, test.user)
			rollbackTriggered := db.rollbackTriggered
			commitTriggered := db.commitTriggered
			increment := false
			if test.article != nil {
				increment = test.article.FavoritesCount == 1
			}

			if test.expectedErr != nil {
				if err != nil && err.Error() != test.expectedErr.Error() {
					t.Errorf("Expected Error: %v, got: %v", test.expectedErr, err)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
					return
				}
			}

			if rollbackTriggered != test.expectedRollback {
				t.Errorf("Expected rollback flag: %v, got: %v", test.expectedRollback, rollbackTriggered)
				return
			}

			if commitTriggered != test.expectedCommit {
				t.Errorf("Expected commit flag: %v, got: %v", test.expectedCommit, commitTriggered)
				return
			}

			if increment != test.expectedIncrement {
				t.Errorf("Expected increment flag: %v, got: %v", test.expectedIncrement, increment)
				return
			}
		})
	}
}
