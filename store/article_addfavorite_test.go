
// ********RoostGPT********
/*

roost_feedback [12/24/2024, 10:13:56 AM]:Fix these errors \n```\n./article_addfavorite_test.go:59:6: ArticleStore redeclared in this block\n\t./article.go:9:6: other declaration of ArticleStore\n./article_addfavorite_test.go:63:24: method ArticleStore.AddFavorite already declared at ./article.go:124:24\n./article_addfavorite_test.go:146:22: cannot use &MockAssociation{…} (value of type *MockAssociation) as *gorm.Association value in struct literal\n./article_addfavorite_test.go:150:37: cannot use db (variable of type *MockDB) as *gorm.DB value in struct literal\n```

roost_feedback [12/24/2024, 10:35:29 AM]:ERROR:\n```\n./article_addfavorite_test.go:157:22: cannot use &TestAssociation{…} (value of type *TestAssociation) as *gorm.Association value in struct literal\n```\nFix the above error
*/

// ********RoostGPT********

package store

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

type MockDB struct {
	mockModel        *model.Article
	mockAssociation  *MockAssociation
	mockUser         *model.User
	rollbackTriggered bool
	commitTriggered  bool
}

func (m *MockDB) Model(value interface{}) *MockDB {
	m.mockModel = value.(*model.Article)
	return m
}

func (m *MockDB) Update(column string, value ...interface{}) *MockDB {
	if column == "favorites_count" && m.mockModel != nil {
		m.mockModel.FavoritesCount++
	}
	return m
}

func (m *MockDB) Begin() *MockDB {
	return m
}

func (m *MockDB) Rollback() *MockDB {
	m.rollbackTriggered = true
	return m
}

func (m *MockDB) Commit() *MockDB {
	m.commitTriggered = true
	return m
}

func (m *MockDB) Association(column string) *MockAssociation {
	return m.mockAssociation
}

type MockAssociation struct {
	appendError error
}

func (m *MockAssociation) Append(values ...interface{}) *MockAssociation {
	return &MockAssociation{appendError: m.appendError}
}

type TestArticleStore struct {
	db *MockDB
}

func (s *TestArticleStore) AddFavorite(a *model.Article, u *model.User) error {
	if a == nil || u == nil {
		return errors.New("Article or User is nil")
	}
	
	tx := s.db.Begin()

	err := tx.Model(a).Association("FavoritedUsers").
		Append(u).appendError
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(a).
		Update("favorites_count", gorm.Expr("favorites_count + ?", 1)).appendError
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	a.FavoritesCount++

	return nil
}

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
			appendError:       errors.New("some error"),
			expectedErr:       errors.New("some error"),
			expectedRollback:  true,
			expectedCommit:    false,
			expectedIncrement: false,
		},
		{
			name:              "Unsuccessful operation due to nil Article",
			user:              &model.User{},
			article:           nil,
			appendError:       nil,
			expectedErr:       errors.New("Article or User is nil"),
			expectedRollback:  false,
			expectedCommit:    false,
			expectedIncrement: false,
		},
		{
			name:              "Unsuccessful operation due to nil User",
			user:              nil,
			article:           &model.Article{},
			appendError:       nil,
			expectedErr:       errors.New("Article or User is nil"),
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

			articleStore := TestArticleStore{db: db}

			err := articleStore.AddFavorite(test.article, test.user)
			rollbackTriggered := db.rollbackTriggered
			commitTriggered := db.commitTriggered
			increment := false
			if test.article != nil {
				increment = test.article.FavoritesCount == 1
			}

			if test.expectedErr != nil {
				if err == nil || err.Error() != test.expectedErr.Error() {
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
