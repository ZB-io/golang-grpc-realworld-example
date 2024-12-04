package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {

	testCases := []struct {
		name          string
		mockArticle   *model.Article
		mockDBError   error
		expectedError error
	}{
		{
			name:          "Successfully Creating a New Article",
			mockArticle:   &model.Article{},
			mockDBError:   nil,
			expectedError: nil,
		},
		{
			name:          "Fail to Create an Article When Required Fields Are Missing",
			mockArticle:   nil,
			mockDBError:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "Failing to Create an Article Due to Database Error",
			mockArticle:   &model.Article{},
			mockDBError:   errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			defer mockDB.Close()

			gdb, err := gorm.Open("postgres", mockDB)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
			}

			store := &ArticleStore{
				db: gdb,
			}

			mock.ExpectBegin()
			mock.ExpectQuery("^INSERT*").WillReturnError(tc.mockDBError)
			mock.ExpectCommit()

			err = store.Create(tc.mockArticle)
			switch tc.expectedError {
			case nil:
				assert.NoError(t, err)
			default:
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}

