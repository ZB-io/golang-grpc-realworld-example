package store

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"sync"
	"github.com/stretchr/testify/suite"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	articleStore := store.NewArticleStore(db)

	type testData struct {
		description string
		article     *model.Article
		mockSetup   func()
		wantError   bool
	}

	tests := []testData{

		{
			description: "Create a new article",
			article: &model.Article{
				Title:       "Test Article Title",
				Description: "Test Article Description",
				Body:        "Test Article Body",
				UserID:      1,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Test Article Title", "Test Article Description", "Test Article Body", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantError: false,
		},

		{
			description: "Create an article with empty fields",
			article: &model.Article{
				Title:       "",
				Description: "",
				Body:        "",
				UserID:      0,
			},
			mockSetup: func() {},
			wantError: true,
		},

		{
			description: "Create an article with invalid field values",
			article: &model.Article{
				Title:       "Test Article Title",
				Description: "Test Article Description",
				Body:        "Test Article Body",
				UserID:      -1,
			},
			mockSetup: func() {},
			wantError: true,
		},

		{
			description: "Create an article with duplicate title",
			article: &model.Article{
				Title:       "Duplicate Title",
				Description: "Test Article Description",
				Body:        "Test Article Body",
				UserID:      1,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Duplicate Title", "Test Article Description", "Test Article Body", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
		},

		{
			description: "Create an article with non-existent user ID",
			article: &model.Article{
				Title:       "Test Article Title",
				Description: "Test Article Description",
				Body:        "Test Article Body",
				UserID:      9999,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Test Article Title", "Test Article Description", "Test Article Body", 9999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
		},

		{
			description: "Create an article with nil pointer",
			article:     nil,
			mockSetup:   func() {},
			wantError:   true,
		},

		{
			description: "Create an article with a database error",
			article: &model.Article{
				Title:       "Test Article Title",
				Description: "Test Article Description",
				Body:        "Test Article Body",
				UserID:      1,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Test Article Title", "Test Article Description", "Test Article Body", 1).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			mock.ExpectedCalls = nil

			tt.mockSetup()

			err := articleStore.Create(tt.article)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())

			os.Stdout = old

			out, err := io.ReadAll(w)
			require.NoError(t, err)

			if tt.wantError {
				assert.Contains(t, string(out), "error creating article")
			} else {
				assert.NotContains(t, string(out), "error creating article")
			}
		})
	}
}

/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestCreateComment(t *testing.T) {

	t.Log("TEST SCENARIO 1: Successful Comment Creation")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := store.NewArticleStore(db)
	mockComment := model.Comment{
		Model: gorm.Model{
			ID: 1,
		},
		Body: "This is a test comment",
		Author: &model.User{
			Model: gorm.Model{
				ID: 1,
			},
			Username: "testuser",
		},
		ArticleID: 1,
	}

	mock.ExpectExec("INSERT INTO `comments` (`body`,`author_id`,`article_id`) VALUES (?,?,?)").
		WithArgs(mockComment.Body, mockComment.Author.ID, mockComment.ArticleID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = s.CreateComment(&mockComment)

	if err != nil {
		t.Errorf("an error '%s' occurred when creating a comment, expected successful insertion", err)
	}

	t.Log("TEST SCENARIO 2: Empty Comment Creation")

	mock.ExpectExec("INSERT INTO `comments` (`body`,`author_id`,`article_id`) VALUES (?,?,?)").
		WithArgs(nil, nil, nil).
		WillReturnError(fmt.Errorf("empty comment"))
	emptyComment := model.Comment{}

	err = s.CreateComment(&emptyComment)

	if err == nil {
		t.Error("expected an error when creating an empty comment, but none occurred")
	} else if err.Error() != "empty comment" {
		t.Errorf("expected error 'empty comment', got '%s'", err.Error())
	}

	t.Log("TEST SCENARIO 3: Unique Comment Creation")

	mock.ExpectExec("INSERT INTO `comments` (`body`,`author_id`,`article_id`) VALUES (?,?,?)").
		WithArgs(mockComment.Body, mockComment.Author.ID, mockComment.ArticleID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	uniqueComment := model.Comment{
		Model: gorm.Model{
			ID: 2,
		},
		Body: "This is a unique comment",
		Author: &model.User{
			Model: gorm.Model{
				ID: 1,
			},
			Username: "testuser",
		},
		ArticleID: 1,
	}

	err = s.CreateComment(&uniqueComment)

	if err != nil {
		t.Errorf("an error '%s' occurred when creating a unique comment, expected successful insertion", err)
	}

	t.Log("TEST SCENARIO 4: Comment Creation with Invalid Data")

	mock.ExpectExec("INSERT INTO `comments` (`body`,`author_id`,`article_id`) VALUES (?,?,?)").
		WithArgs(mockComment.Body, mockComment.Author.ID, mockComment.ArticleID).
		WillReturnError(fmt.Errorf("invalid comment data"))
	invalidComment := model.Comment{
		Model: gorm.Model{
			ID: 1,
		},
		Body: "This comment has invalid data",
		Author: &model.User{
			Model: gorm.Model{
				ID: 1,
			},
			Username: "testuser",
		},
		ArticleID: 1,
	}

	err = s.CreateComment(&invalidComment)

	if err == nil {
		t.Error("expected an error when creating a comment with invalid data, but none occurred")
	} else if err.Error() != "invalid comment data" {
		t.Errorf("expected error 'invalid comment data', got '%s'", err.Error())
	}

	t.Log("TEST SCENARIO 5: Comment Creation with Database Error")

	mock.ExpectExec("INSERT INTO `comments` (`body`,`author_id`,`article_id`) VALUES (?,?,?)").
		WithArgs(mockComment.Body, mockComment.Author.ID, mockComment.ArticleID).
		WillReturnError(fmt.Errorf("database error"))

	err = s.CreateComment(&mockComment)

	if err == nil {
		t.Error("expected an error when creating a comment with a database error, but none occurred")
	} else if err.Error() != "database error" {
		t.Errorf("expected error 'database error', got '%s'", err.Error())
	}
}

/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func (suite *ArticleStoreTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ArticleStoreTestSuite) SetupSuite() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	require.NoError(suite.T(), err)
}

func TestDelete(t *testing.T) {
	suite.Run(t, new(ArticleStoreTestSuite))
}

func (suite *ArticleStoreTestSuite) TestDelete_ArticleWithComments() {
	suite.T().Log("Test Case: Deleting Article with Comments")

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM articles WHERE id = \\?").WithArgs(suite.articles[0].ID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec("DELETE FROM comments WHERE article_id = \\?").WithArgs(suite.articles[0].ID).WillReturnResult(sqlmock.NewResult(2, 2))
	suite.mock.ExpectCommit()

	err := (&ArticleStore{suite.db}).Delete(&suite.articles[0])
	require.NoError(suite.T(), err)
}

func (suite *ArticleStoreTestSuite) TestDelete_BoundaryConditions() {
	suite.T().Log("Test Case: Boundary Conditions")

}

func (suite *ArticleStoreTestSuite) TestDelete_ConcurrentAccess() {
	suite.T().Log("Test Case: Simultaneous Article Deletion")

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			err := (&ArticleStore{suite.db}).Delete(&suite.articles[0])
			require.NoError(suite.T(), err)
		}(i)
	}

	wg.Wait()
}

func (suite *ArticleStoreTestSuite) TestDelete_DatabaseError() {
	suite.T().Log("Test Case: Error Handling for Database Issues")

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM articles WHERE id = \\?").WithArgs(suite.articles[0].ID).WillReturnError(gorm.ErrRecordNotFound)
	suite.mock.ExpectRollback()

	err := (&ArticleStore{suite.db}).Delete(&suite.articles[0])
	require.Error(suite.T(), err)
}

func (suite *ArticleStoreTestSuite) TestDelete_EmptyArticleModel() {
	suite.T().Log("Test Case: Empty Article Model")

}

func (suite *ArticleStoreTestSuite) TestDelete_NonExistentArticle() {
	suite.T().Log("Test Case: Deleting Non-Existent Article")

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM articles WHERE id = \\?").WithArgs(suite.articles[len(suite.articles)-1].ID + 1).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := (&ArticleStore{suite.db}).Delete(&model.Article{Model: gorm.Model{ID: suite.articles[len(suite.articles)-1].ID + 1}})
	require.NoError(suite.T(), err)
}

func (suite *ArticleStoreTestSuite) TestDelete_Success() {
	suite.T().Log("Test Case: Successful Article Deletion")

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM articles WHERE id = \\?").WithArgs(suite.articles[0].ID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := (&ArticleStore{suite.db}).Delete(&suite.articles[0])
	require.NoError(suite.T(), err)
}

func (suite *ArticleStoreTestSuite) TestDelete_UnexpectedInput() {
	suite.T().Log("Test Case: Unexpected Input")

}

