package store

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		// {
		// 	name:          "Fail to Create an Article When Required Fields Are Missing",
		// 	mockArticle:   nil,
		// 	mockDBError:   nil,
		// 	expectedError: gorm.ErrRecordNotFound,
		// },
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

/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


*/
func TestCreateComment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc      string
		comment   *model.Comment
		dbErr     error
		expectErr bool
	}{
		{
			desc:      "Successfully Creating a Comment",
			comment:   &model.Comment{Body: "This is a test comment"},
			dbErr:     nil,
			expectErr: false,
		},
		{
			desc:      "Failing to Create Comment due to a Null Model",
			comment:   nil,
			dbErr:     nil,
			expectErr: true,
		},
		{
			desc:      "Failing to Create Comment due to Database Write Failures",
			comment:   &model.Comment{Body: "This is a test comment"},
			dbErr:     errors.New("database error"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			mockDB, mock, err := sqlmock.New()

			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db, err := gorm.Open("postgres", mockDB)
			if err != nil {
				t.Fatalf("failed to create gorm.DB from sqlmock: %v", err)
			}

			if tc.comment != nil && tc.dbErr == nil {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			as := NewArticleStore(db)
			err = as.CreateComment(tc.comment)

			if (err != nil) != tc.expectErr {
				t.Errorf("ArticleStore.CreateComment(%v) returned error %v, expected error? %v",
					tc.comment, err, tc.expectErr)
			}

		})
	}
}

/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


*/
func TestDelete(t *testing.T) {

	scenarios := []struct {
		name                 string
		expectedError        error
		mockArticleExists    bool
		mockDatabaseDown     bool
		mockDatabaseOpFailed bool
	}{
		{
			name:              "Normal Operation - Article Successfully Deleted from the Database",
			expectedError:     nil,
			mockArticleExists: true,
		},
		{
			name:              "Edge Case - Attempting to Delete a Non-Existent Article",
			expectedError:     errors.New("record not found"),
			mockArticleExists: false,
		},
		{
			name:             "Error Case - Database is Down or Unreachable",
			expectedError:    errors.New("database is down"),
			mockDatabaseDown: true,
		},
		{
			name:                 "Error Case - Database Operation Results in Error",
			expectedError:        errors.New("database operation failed"),
			mockDatabaseOpFailed: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {

			db, mock, _ := sqlmock.New()
			gormDB, _ := gorm.Open("postgres", db)

			articleStore := &ArticleStore{
				db: gormDB,
			}

			article := new(model.Article)

			if s.mockDatabaseDown {
				mock.ExpectExec("DELETE").WillReturnError(s.expectedError)
			} else if !s.mockDatabaseOpFailed {
				if s.mockArticleExists {
					mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
				} else {
					mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(0, 0))
				}
			} else {
				mock.ExpectExec("DELETE").WillReturnError(s.expectedError)
			}

			err := articleStore.Delete(article)

			if s.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, s.expectedError, err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


*/
func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name      string
		comment   *model.Comment
		mock      func(mock sqlmock.Sqlmock, comment *model.Comment)
		wantError bool
	}{
		{
			name:    "DeleteComment Successfully Deletes a Comment",
			comment: &model.Comment{},
			mock: func(mock sqlmock.Sqlmock, comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(comment.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantError: false,
		},
		// {
		// 	name:    "DeleteComment Returns an Error when the Comment Object is Nil",
		// 	comment: nil,
		// 	mock: func(mock sqlmock.Sqlmock, comment *model.Comment) {
		// 		mock.ExpectBegin()
		// 		mock.ExpectExec("DELETE").
		// 			WithArgs(nil).
		// 			WillReturnError(gorm.ErrRecordNotFound)
		// 		mock.ExpectRollback()
		// 	},
		// 	wantError: true,
		// },
		{
			name:    "DeleteComment Returns an Error when the DB Operation Fails",
			comment: &model.Comment{},
			mock: func(mock sqlmock.Sqlmock, comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(comment.ID).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			gormDB, err := gorm.Open("postgres", db)
			assert.NoError(t, err)

			articleStore := &ArticleStore{db: gormDB}

			test.mock(mock, test.comment)

			err = articleStore.DeleteComment(test.comment)

			if test.wantError {
				assert.Error(t, err)
				t.Log("Error expected and returned")
			} else {
				assert.NoError(t, err)
				t.Log("No error expected and none returned")
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


*/
func TestGetByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		setupMock   func(sqlmock.Sqlmock, uint)
		wantArticle *model.Article
		wantErr     error
	}{
		{
			name: "Correct ID",
			id:   1,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"id", "Title", "Description"}).
					AddRow(1, "Test article", "Test description")

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WithArgs(id).WillReturnRows(rows)
			},
			wantArticle: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Test article",
				Description: "Test description",
			},
			wantErr: nil,
		},
		{
			name: "Incorrect ID",
			id:   20,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"id", "Title", "Description"})

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WithArgs(id).WillReturnRows(rows)
			},
			wantArticle: nil,
			wantErr:     gorm.ErrRecordNotFound,
		},
		{
			name: "Preloading relations",
			id:   1,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				rows := sqlmock.NewRows([]string{"id", "Title", "Description", "Tags", "Author"}).
					AddRow(1, "Test article", "Test description", "Test tags", "Test author")

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WithArgs(id).WillReturnRows(rows)
			},
			wantArticle: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Test article",
				Description: "Test description",
				Tags:        []model.Tag{model.Tag{Name: "Test tags"}},
				Author:      model.User{Username: "Test author"},
			},
			wantErr: nil,
		},
		{
			name: "Database is unavailable",
			id:   1,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").WithArgs(id).WillReturnError(gorm.ErrInvalidSQL)
			},
			wantArticle: nil,
			wantErr:     gorm.ErrInvalidSQL,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			gDb, err := gorm.Open("postgres", db)
			require.NoError(t, err)

			mock.ExpectBegin()
			test.setupMock(mock, test.id)
			mock.ExpectCommit()

			store := ArticleStore{db: gDb}
			article, err := store.GetByID(test.id)

			if test.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, test.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.wantArticle, article)
			}
		})
	}
}
