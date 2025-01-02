package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)



func MockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}
func TestNewArticleStore(t *testing.T) {

	testCases := []struct {
		name        string
		db          *gorm.DB
		expectedErr error
	}{
		{
			name: "Successful creation of ArticleStore",
			db:   &gorm.DB{},
		},
		{
			name: "Passing nil as db parameter",
			db:   nil,
		},
		{
			name: "Passing a db instance with existing data",
			db:   &gorm.DB{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result := NewArticleStore(tc.db)

			assert.IsType(t, &ArticleStore{}, result, "Expected ArticleStore type")

			if tc.db == nil {
				assert.Nil(t, result.db, "Expected db to be nil")
			} else {
				assert.Equal(t, tc.db, result.db, "Expected db to be equal to the provided db instance")
			}
		})
	}
}
