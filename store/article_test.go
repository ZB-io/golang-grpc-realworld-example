package github

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)








/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


*/
func (m *mockDB) Append(values ...interface{}) error {
	m.appendCalled = true
	if m.appendError != nil {
		return m.appendError
	}
	m.favoritedUsers = append(m.favoritedUsers, values[0].(*model.User))
	return nil
}

func (m *mockDB) Association(column string) *gorm.Association {
	m.associationCalled = true
	return &gorm.Association{DB: &gorm.DB{Value: m}}
}

func (m *mockDB) Begin() *gorm.DB {
	m.beginCalled = true
	return &gorm.DB{Value: m}
}

func (m *mockDB) Commit() *gorm.DB {
	m.commitCalled = true
	return &gorm.DB{Value: m}
}

func (m *mockDB) Model(value interface{}) *gorm.DB {
	return &gorm.DB{Value: m}
}

func (m *mockDB) Rollback() *gorm.DB {
	m.rollbackCalled = true
	return &gorm.DB{Value: m}
}

func (m *mockDB) Update(column string, value interface{}) *gorm.DB {
	m.updateCalled = true
	if m.updateError != nil {
		return &gorm.DB{Error: m.updateError}
	}
	m.favoritesCount++
	return &gorm.DB{Value: m}
}

