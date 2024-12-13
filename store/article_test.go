package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
	"fmt"
	"log"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"sync"
	"os"
	"bytes"
	"context"
	"database/sql"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestArticleStoreCreate(t *testing.T) {
	type testCase struct {
		name        string
		article     model.Article
		mockDBSetup func(sqlmock.Sqlmock)
		shouldError bool
	}

	tests := []testCase{
		{
			name: "Scenario 1: Successfully Create a New Article",
			article: model.Article{
				Title:       "Sample Title",
				Description: "Sample Description",
				Body:        "Sample Body",
				UserID:      1,
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"articles\"").
					WithArgs(sqlmock.AnyArg(), "Sample Title", "Sample Description", "Sample Body", 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			shouldError: false,
		},
		{
			name: "Scenario 2: Attempt to Create an Article with Missing Required Fields",
			article: model.Article{
				Title:       "",
				Description: "Sample Description",
				Body:        "Sample Body",
				UserID:      1,
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"articles\"").
					WithArgs(sqlmock.AnyArg(), "", "Sample Description", "Sample Body", 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("missing required fields"))
				mock.ExpectRollback()
			},
			shouldError: true,
		},
		{
			name: "Scenario 3: Database Connection Failure",
			article: model.Article{
				Title:       "Database Failure Title",
				Description: "Database Failure Description",
				Body:        "Database Failure Body",
				UserID:      1,
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("database connection failure"))
			},
			shouldError: true,
		},
		{
			name: "Scenario 4: Creating an Article with Pre-existing Title",
			article: model.Article{
				Title:       "Duplicate Title",
				Description: "Sample Description",
				Body:        "Sample Body",
				UserID:      1,
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"articles\"").
					WithArgs(sqlmock.AnyArg(), "Duplicate Title", "Sample Description", "Sample Body", 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("duplicate title"))
				mock.ExpectRollback()
			},
			shouldError: true,
		},
		{
			name: "Scenario 5: Validate Database Rollback on Failure",
			article: model.Article{
				Title:       "Rollback Title",
				Description: "Rollback Description",
				Body:        "Rollback Body",
				UserID:      1,
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"articles\"").
					WithArgs(sqlmock.AnyArg(), "Rollback Title", "Rollback Description", "Rollback Body", 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("internal error"))
				mock.ExpectRollback()
			},
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create a new mock database: %v", err)
			}
			defer db.Close()

			tc.mockDBSetup(mock)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
				Logger:         logger.Default.LogMode(logger.Silent),
				NamingStrategy: schema.NamingStrategy{},
			})
			if err != nil {
				t.Fatalf("Failed to open gorm DB: %v", err)
			}

			articleStore := ArticleStore{db: gormDB}

			err = articleStore.Create(&tc.article)

			if tc.shouldError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unmet expectations: %v", err)
			}

			t.Logf("%s - success, error state: %v", tc.name, err != nil)
		})
	}
}

/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestArticleStoreCreateComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name      string
		prepare   func()
		input     *model.Comment
		wantError bool
		errorMsg  string
	}{
		{
			name: "Successfully Create a Comment",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO`).WithArgs().
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: &model.Comment{
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantError: false,
		},
		{
			name: "Fail to Create a Comment Due to Database Error",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO`).WithArgs().
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			input: &model.Comment{
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantError: true,
			errorMsg:  "database error",
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO`).WithArgs().
					WillReturnError(fmt.Errorf("validation error: missing required fields"))
				mock.ExpectRollback()
			},
			input:     &model.Comment{},
			wantError: true,
			errorMsg:  "validation error: missing required fields",
		},
		{
			name: "Creating a Comment with a User Not Existing in Database",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO`).WithArgs().
					WillReturnError(fmt.Errorf("foreign key constraint fails for UserID"))
				mock.ExpectRollback()
			},
			input: &model.Comment{
				Body:      "Test Comment",
				UserID:    9999,
				ArticleID: 1,
			},
			wantError: true,
			errorMsg:  "foreign key constraint fails for UserID",
		},
		{
			name: "Creating a Comment for an Article Not Existing in Database",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO`).WithArgs().
					WillReturnError(fmt.Errorf("foreign key constraint fails for ArticleID"))
				mock.ExpectRollback()
			},
			input: &model.Comment{
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 9999,
			},
			wantError: true,
			errorMsg:  "foreign key constraint fails for ArticleID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			err := store.CreateComment(tt.input)

			if tt.wantError {
				assert.Error(t, err, tt.errorMsg)
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				assert.NoError(t, err)
				t.Log("Test passed for:", tt.name)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func TestArticleStoreDelete(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		input         *model.Article
		expectedError error
	}{
		{
			name: "Successfully Delete an Existing Article",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\" WHERE \"articles\".\"id\" = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			input:         &model.Article{Model: gorm.Model{ID: 1}},
			expectedError: nil,
		},
		{
			name: "Delete Non-Existent Article",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\" WHERE \"articles\".\"id\" = ?").
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			input:         &model.Article{Model: gorm.Model{ID: 2}},
			expectedError: nil,
		},
		{
			name:          "Delete with Nil Article Object",
			setupMock:     func(mock sqlmock.Sqlmock) {},
			input:         nil,
			expectedError: gorm.ErrInvalidValue,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM \"articles\" WHERE \"articles\".\"id\" = ?").
					WithArgs(1).
					WillReturnError(errors.New("connection error"))
			},
			input:         &model.Article{Model: gorm.Model{ID: 1}},
			expectedError: errors.New("connection error"),
		},
		{
			name: "Attempt to Delete Article with Constraints",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\" WHERE \"articles\".\"id\" = ?").
					WithArgs(3).
					WillReturnError(errors.New("constraint violation"))
				mock.ExpectRollback()
			},
			input:         &model.Article{Model: gorm.Model{ID: 3}},
			expectedError: errors.New("constraint violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			assert.NoError(t, err)

			articleStore := &ArticleStore{db: gormDB}

			tt.setupMock(mock)

			err = articleStore.Delete(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Log("Test scenario executed successfully for:", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func TestArticleStoreDeleteComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %s", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		scenario      string
		comment       *model.Comment
		mockFunc      func()
		expectedError bool
	}{
		{
			scenario: "Successful Deletion of an Existing Comment",
			comment:  &model.Comment{Model: gorm.Model{ID: 1}},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			scenario: "Attempt to Delete a Non-Existing Comment",
			comment:  &model.Comment{Model: gorm.Model{ID: 2}},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			scenario: "Handling Database Connectivity Issues",
			comment:  &model.Comment{Model: gorm.Model{ID: 3}},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").
					WithArgs(3).
					WillReturnError(errors.New("database connection failure"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			scenario: "Concurrent Deletion Requests for the Same Comment",
			comment:  &model.Comment{Model: gorm.Model{ID: 4}},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").
					WithArgs(4).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE \"comments\".\"id\" = $1").
					WithArgs(4).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			tt.mockFunc()

			err := store.DeleteComment(tt.comment)

			if tt.expectedError && err == nil {
				t.Errorf("%s: expected error but got none", tt.scenario)
			} else if !tt.expectedError && err != nil {
				t.Errorf("%s: did not expect error but got %s", tt.scenario, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed with expected error status: %v", tt.scenario, tt.expectedError)
		})
	}
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func TestArticleStoreGetByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}
	articleStore := &ArticleStore{db: gormDB}

	tests := []struct {
		name      string
		setupMock func()
		inputID   uint

		expectErr error
		verify    func(t *testing.T, article *model.Article, err error)
	}{
		{
			name: "Successfully Retrieve Article by ID",
			setupMock: func() {
				mock.ExpectQuery("^SELECT .+ FROM \"articles\" .+").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}).
						AddRow(1, "Test Article", 1))
				mock.ExpectQuery("^SELECT .+ FROM \"tags\" .+").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "Test Tag"))
				mock.ExpectQuery("^SELECT .+ FROM \"authors\" .+").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "Test Author"))
			},
			inputID:   1,
			expectErr: nil,
			verify: func(t *testing.T, article *model.Article, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, "Test Article", article.Title)
				assert.Equal(t, 1, article.AuthorID)
			},
		},
		{
			name: "Article Not Found in Database",
			setupMock: func() {
				mock.ExpectQuery("^SELECT .+ FROM \"articles\" .+ WHERE .+\"id\" = \\$1").
					WithArgs(99).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			inputID:   99,
			expectErr: gorm.ErrRecordNotFound,
			verify: func(t *testing.T, article *model.Article, err error) {
				assert.Nil(t, article)
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "Database Connection Error",
			setupMock: func() {
				mock.ExpectQuery("^SELECT .+ FROM \"articles\" .+").
					WillReturnError(errors.New("connection error"))
			},
			inputID:   1,
			expectErr: errors.New("connection error"),
			verify: func(t *testing.T, article *model.Article, err error) {
				assert.Equal(t, article, nil)
				assert.EqualError(t, err, "connection error")
			},
		},
		{
			name: "Article ID is Zero",
			setupMock: func() {
				mock.ExpectQuery("^SELECT .+ FROM \"articles\" .+ WHERE .+\"id\" = \\$1").
					WithArgs(0).
					WillReturnError(errors.New("invalid ID zero"))
			},
			inputID:   0,
			expectErr: errors.New("invalid ID zero"),
			verify: func(t *testing.T, article *model.Article, err error) {
				assert.Nil(t, article)
				assert.EqualError(t, err, "invalid ID zero")
			},
		},
		{
			name: "Article with Multiple Tags and A Complete Author Record",
			setupMock: func() {
				mock.ExpectQuery("^SELECT .+ FROM \"articles\" .+").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}).
						AddRow(2, "Another Article", 2))
				mock.ExpectQuery("^SELECT .+ FROM \"tags\" .+").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "Tag1").AddRow(2, "Tag2"))
				mock.ExpectQuery("^SELECT .+ FROM \"authors\" .+").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(2, "Another Author"))
			},
			inputID:   2,
			expectErr: nil,
			verify: func(t *testing.T, article *model.Article, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, "Another Article", article.Title)
				assert.Equal(t, []model.Tag{model.Tag{Name: "Tag1"}, model.Tag{Name: "Tag2"}}, article.Tags)
				assert.Equal(t, "Another Author", article.Author.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			article, err := articleStore.GetByID(tt.inputID)
			tt.verify(t, article, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {

	tests := []struct {
		name            string
		mockSetup       func(sqlmock.Sqlmock)
		inputID         uint
		expectedErr     error
		expectedComment *model.Comment
	}{
		{
			name: "Successfully retrieves a comment by ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND \\(\\d+\\)$").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
						AddRow(1, "Test comment", 1, 1))
			},
			inputID:         1,
			expectedErr:     nil,
			expectedComment: &model.Comment{Model: gorm.Model{ID: 1}, Body: "Test comment", UserID: 1, ArticleID: 1},
		},
		{
			name: "Attempt to retrieve a comment with non-existent ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND \\(\\d+\\)$").
					WithArgs(99).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			inputID:         99,
			expectedErr:     gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Handles a database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND \\(\\d+\\)$").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			inputID:         1,
			expectedErr:     errors.New("database error"),
			expectedComment: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Expected no error on DB open, got %v", err)
			}

			test.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			comment, err := store.GetCommentByID(test.inputID)

			if err != test.expectedErr {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if comment != nil && test.expectedComment != nil {
				if comment.ID != test.expectedComment.ID || comment.Body != test.expectedComment.Body {
					t.Errorf("Expected comment: %v, got: %v", test.expectedComment, comment)
				}
			}

			if comment == nil && test.expectedComment != nil || comment != nil && test.expectedComment == nil {
				t.Errorf("Expected comment: %v, got: %v", test.expectedComment, comment)
			}

			t.Logf("Test case '%s' completed successfully", test.name)
		})
	}

}

/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestArticleStoreGetComments(t *testing.T) {
	t.Run("Scenario 1: Retrieve Comments for an Article with Multiple Comments", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", db)
		require.NoError(t, err)

		mock.ExpectQuery(`SELECT * FROM "comments" WHERE article_id = \?`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
				AddRow(1, "Great article!", 1, 1).
				AddRow(2, "Very informative.", 2, 1))

		articleStore := &ArticleStore{db: gormDB}

		article := &model.Article{Model: gorm.Model{ID: 1}}
		comments, err := articleStore.GetComments(article)

		require.NoError(t, err)
		require.Len(t, comments, 2, "expected two comments")

		t.Log("Successfully retrieved comments for an article with multiple comments.")
	})

	t.Run("Scenario 2: Retrieve Comments for an Article with No Comments", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", db)
		require.NoError(t, err)

		mock.ExpectQuery(`SELECT * FROM "comments" WHERE article_id = \?`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}))

		articleStore := &ArticleStore{db: gormDB}
		article := &model.Article{Model: gorm.Model{ID: 2}}
		comments, err := articleStore.GetComments(article)

		require.NoError(t, err)
		require.Len(t, comments, 0, "expected no comments")

		t.Log("Successfully verified that no comments are returned for an article with no comments.")
	})

	t.Run("Scenario 3: Database Error During Comment Retrieval", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", db)
		require.NoError(t, err)

		mock.ExpectQuery(`SELECT * FROM "comments" WHERE article_id = \?`).
			WillReturnError(fmt.Errorf("database error"))

		articleStore := &ArticleStore{db: gormDB}
		article := &model.Article{Model: gorm.Model{ID: 3}}
		comments, err := articleStore.GetComments(article)

		require.Error(t, err, "expected a database error")
		require.Len(t, comments, 0, "expected no comments due to error")
		t.Log("Successfully handled database errors during comment retrieval.")
	})

	t.Run("Scenario 4: Preloaded Author Details Are Correct for Comments", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", db)
		require.NoError(t, err)

		mock.ExpectQuery(`SELECT * FROM "comments" WHERE article_id = \?`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
				AddRow(1, "Great article!", 1, 4))

		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE "id" = \?`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
				AddRow(1, "JohnDoe", "john@example.com"))

		articleStore := &ArticleStore{db: gormDB}
		article := &model.Article{Model: gorm.Model{ID: 4}}
		comments, err := articleStore.GetComments(article)

		require.NoError(t, err)
		require.Len(t, comments, 1, "expected one comment")
		require.Equal(t, "JohnDoe", comments[0].Author.Username, "expected author username to be 'JohnDoe'")
		t.Log("Successfully retrieved comments with correct preloaded author details.")
	})

	t.Run("Scenario 5: Large Number of Comments Retrieval Efficiency", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", db)
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"})
		for i := 1; i <= 1000; i++ {
			rows = rows.AddRow(i, fmt.Sprintf("Comment %d", i), 1, 5)
		}

		mock.ExpectQuery(`SELECT * FROM "comments" WHERE article_id = \?`).WillReturnRows(rows)

		articleStore := &ArticleStore{db: gormDB}
		article := &model.Article{Model: gorm.Model{ID: 5}}
		comments, err := articleStore.GetComments(article)

		require.NoError(t, err)
		require.Len(t, comments, 1000, "expected 1000 comments")

		t.Log("Successfully retrieved a large number of comments efficiently.")
	})

	t.Run("Scenario 6: Retrieve Comments when Article ID is Invalid", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", db)
		require.NoError(t, err)

		mock.ExpectQuery(`SELECT * FROM "comments" WHERE article_id = \?`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}))

		articleStore := &ArticleStore{db: gormDB}
		article := &model.Article{Model: gorm.Model{ID: 9999}}
		comments, err := articleStore.GetComments(article)

		require.NoError(t, err)
		require.Len(t, comments, 0, "expected no comments for invalid article ID")

		t.Log("Gracefully handled invalid article ID with no errors and no comments.")
	})
}

/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func TestArticleStoreGetTags(t *testing.T) {

	t.Run("Success: Retrieve Tags", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening mock database: %v", err)
		}
		defer db.Close()

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("error opening gorm database: %v", err)
		}

		mockTags := []model.Tag{
			{Model: gorm.Model{ID: 1}, TagName: "go"},
			{Model: gorm.Model{ID: 2}, TagName: "programming"},
		}

		rows := sqlmock.NewRows([]string{"id", "tag_name"}).
			AddRow(mockTags[0].ID, mockTags[0].TagName).
			AddRow(mockTags[1].ID, mockTags[1].TagName)

		mock.ExpectQuery("^SELECT \\* FROM \"tags\"").WillReturnRows(rows)

		store := &ArticleStore{db: gdb}
		tags, err := store.GetTags()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tags) != len(mockTags) {
			t.Errorf("expected %d tags, got %d", len(mockTags), len(tags))
		}
	})

	t.Run("Empty: No Tags Available", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening mock database: %v", err)
		}
		defer db.Close()

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("error opening gorm database: %v", err)
		}

		mock.ExpectQuery("^SELECT \\* FROM \"tags\"").WillReturnRows(sqlmock.NewRows([]string{"id", "tag_name"}))

		store := &ArticleStore{db: gdb}
		tags, err := store.GetTags()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("Error: Database Failure", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening mock database: %v", err)
		}
		defer db.Close()

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("error opening gorm database: %v", err)
		}

		mock.ExpectQuery("^SELECT \\* FROM \"tags\"").WillReturnError(gorm.ErrInvalidSQL)

		store := &ArticleStore{db: gdb}
		tags, err := store.GetTags()

		if err == nil {
			t.Error("expected an error, got nil")
		}
		if tags != nil {
			t.Errorf("expected nil or empty tags, got %v", tags)
		}
	})

	t.Run("Concurrent: Simultaneous Retrievals", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening mock database: %v", err)
		}
		defer db.Close()

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("error opening gorm database: %v", err)
		}

		mockTags := []model.Tag{
			{Model: gorm.Model{ID: 1}, TagName: "go"},
			{Model: gorm.Model{ID: 2}, TagName: "cloud"},
		}

		rows := sqlmock.NewRows([]string{"id", "tag_name"}).
			AddRow(mockTags[0].ID, mockTags[0].TagName).
			AddRow(mockTags[1].ID, mockTags[1].TagName)

		mock.ExpectQuery("^SELECT \\* FROM \"tags\"").WillReturnRows(rows)

		store := &ArticleStore{db: gdb}

		concurrency := 5
		results := make(chan []model.Tag, concurrency)
		errors := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				tags, err := store.GetTags()
				results <- tags
				errors <- err
			}()
		}

		for i := 0; i < concurrency; i++ {
			if err := <-errors; err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			tags := <-results
			if len(tags) != len(mockTags) {
				t.Errorf("expected %d tags, got %d", len(mockTags), len(tags))
			}
		}
	})

	t.Run("Performance: Large Tag Set", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening mock database: %v", err)
		}
		defer db.Close()

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("error opening gorm database: %v", err)
		}

		largeNumber := 1000
		mockRows := sqlmock.NewRows([]string{"id", "tag_name"})
		for i := 1; i <= largeNumber; i++ {
			mockRows.AddRow(i, "tag"+string(i))
		}
		mock.ExpectQuery("^SELECT \\* FROM \"tags\"").WillReturnRows(mockRows)

		store := &ArticleStore{db: gdb}
		tags, err := store.GetTags()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tags) != largeNumber {
			t.Errorf("expected %d tags, got %d", largeNumber, len(tags))
		}
	})
}

/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func TestArticleStoreIsFavorited(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, %s", err)
	}
	defer db.Close()

	mockDB, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("failed to open gorm DB, %s", err)
	}

	store := &ArticleStore{
		db: mockDB,
	}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		expectFavor bool
		expectErr   bool
		mockQuery   func()
	}{
		{
			name: "Identifying Favorited Article by a User",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectFavor: true,
			expectErr:   false,
			mockQuery: func() {
				mock.ExpectQuery("^SELECT count(.+) FROM \"favorite_articles\" WHERE").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
		},
		{
			name: "Article Not Favorited by the User",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectFavor: false,
			expectErr:   false,
			mockQuery: func() {
				mock.ExpectQuery("^SELECT count(.+) FROM \"favorite_articles\" WHERE").
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
		},
		{
			name:        "Handling Nil Article Input",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			expectFavor: false,
			expectErr:   false,
			mockQuery:   nil,
		},
		{
			name:        "Handling Nil User Input",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			expectFavor: false,
			expectErr:   false,
			mockQuery:   nil,
		},
		{
			name: "Handling Database Error Scenarios",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectFavor: false,
			expectErr:   true,
			mockQuery: func() {
				mock.ExpectQuery("^SELECT count(.+) FROM \"favorite_articles\" WHERE").
					WithArgs(3, 1).
					WillReturnError(errors.New("database error"))
			},
		},
		{
			name:        "Uninitialized Database in ArticleStore",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			expectFavor: false,
			expectErr:   true,
			mockQuery: func() {

				storeNilDB := &ArticleStore{
					db: nil,
				}
				_, _ = storeNilDB.IsFavorited(&model.Article{Model: gorm.Model{ID: 1}}, &model.User{Model: gorm.Model{ID: 1}})
			},
		},
		{
			name: "Zero ID Inputs",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			user: &model.User{
				Model: gorm.Model{ID: 0},
			},
			expectFavor: false,
			expectErr:   false,
			mockQuery: func() {
				mock.ExpectQuery("^SELECT count(.+) FROM \"favorite_articles\" WHERE").
					WithArgs(0, 0).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
		},
		{
			name: "Valid Inputs with Zero Favorited Count",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
			},
			user: &model.User{
				Model: gorm.Model{ID: 2},
			},
			expectFavor: false,
			expectErr:   false,
			mockQuery: func() {
				mock.ExpectQuery("^SELECT count(.+) FROM \"favorite_articles\" WHERE").
					WithArgs(4, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockQuery != nil {
				tt.mockQuery()
			}

			favorited, err := store.IsFavorited(tt.article, tt.user)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if favorited != tt.expectFavor {
				t.Errorf("Expected favorited to be %v but got %v", tt.expectFavor, favorited)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92


 */
func TestNewArticleStore(t *testing.T) {
	t.Run("Successful Initialization with a Non-nil DB", func(t *testing.T) {

		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("Error creating gorm DB: %v", err)
		}
		defer gormDB.Close()

		articleStore := NewArticleStore(gormDB)

		assert.NotNil(t, articleStore, "Expected non-nil ArticleStore")
		assert.Equal(t, gormDB, articleStore.db, "Expected db to be set correctly")
		t.Log("Test passed: ArticleStore successfully initialized with a non-nil DB")
	})

	t.Run("Initialization with Nil DB", func(t *testing.T) {

		var nilDB *gorm.DB = nil

		articleStore := NewArticleStore(nilDB)

		assert.NotNil(t, articleStore, "ArticleStore should be non-nil even with a nil DB")
		assert.Nil(t, articleStore.db, "Expected db field to be nil when initialized with a nil DB")
		t.Log("Test passed: ArticleStore initialized with a nil DB behaves as expected")
	})

	t.Run("Thread-Safety Check for Concurrent Initialization", func(t *testing.T) {

		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("Error creating gorm DB: %v", err)
		}
		defer gormDB.Close()

		var wg sync.WaitGroup
		instances := make([]*ArticleStore, 0)
		mu := &sync.Mutex{}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				store := NewArticleStore(gormDB)
				mu.Lock()
				instances = append(instances, store)
				mu.Unlock()
			}()
		}
		wg.Wait()

		for _, instance := range instances {
			assert.Equal(t, gormDB, instance.db)
		}
		t.Log("Test passed: All ArticleStore instances have correct DB set under concurrent initialization")
	})

	t.Run("Validation of ArticleStore's Functionality Post-Initialization", func(t *testing.T) {

		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("Error creating gorm DB: %v", err)
		}
		defer gormDB.Close()

		articleStore := NewArticleStore(gormDB)

		t.Log("Test passed: ArticleStore post-initialization functionality validated")
	})
}

/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe


 */
func TestArticleStoreUpdate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		log.Fatalf("An error '%s' was not expected when wrapping mock db with gorm", err)
	}
	store := ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		article       *model.Article
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Successfully Update an Article",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Valid Title",
				Description:    "Description",
				Body:           "Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			mockSetup: func() {

				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Update Article with Empty Title",
			article: &model.Article{
				Model:          gorm.Model{ID: 2},
				Title:          "",
				Description:    "Description",
				Body:           "Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			mockSetup: func() {

				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Model:          gorm.Model{ID: 3},
				Title:          "Title",
				Description:    "Description",
				Body:           "Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			mockSetup: func() {

				mock.ExpectBegin().WillReturnError(errors.New("connection failed"))
			},
			expectedError: errors.New("connection failed"),
		},
		{
			name: "Updating a Non-Existent Article",
			article: &model.Article{
				Model:          gorm.Model{ID: 999},
				Title:          "Non-existent",
				Description:    "Description",
				Body:           "Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Update Fails Due to Validation Error",
			article: &model.Article{
				Model:          gorm.Model{ID: 4},
				Title:          "Duplicate Title",
				Description:    "Description",
				Body:           "Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			mockSetup: func() {

				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").WillReturnError(errors.New("duplicate key"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("duplicate key"),
		},
		{
			name: "Ensure Update Only Alters Specified Fields",
			article: &model.Article{
				Model:          gorm.Model{ID: 5},
				Title:          "Update Title",
				Description:    "Persistent Description",
				Body:           "Persistent Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").WithArgs("Update Title", 5).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup()

			err := store.Update(tt.article)

			assert.Equal(t, tt.expectedError, err)

			if err == nil {
				t.Logf("%s passed: Article with ID %v was updated as expected.", tt.name, tt.article.ID)
			} else {
				t.Logf("%s failed: Expected error %v, got %v", tt.name, tt.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func TestArticleStoreGetFeedArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("unexpected error when initializing gorm with mock: %s", err)
	}
	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name       string
		userIDs    []uint
		limit      int64
		offset     int64
		mockExpect func()
		want       []model.Article
		wantErr    bool
	}{
		{
			name:    "Fetch Articles for Existing Users",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  0,
			mockExpect: func() {
				mock.ExpectQuery(`SELECT * FROM "articles" WHERE \(user_id in \(\?, \?\)\) LIMIT \? OFFSET \?`).
					WithArgs(1, 2, 2, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).
						AddRow(1, "Title1", "Description1", "Body1").
						AddRow(2, "Title2", "Description2", "Body2"))
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Title1", Description: "Description1", Body: "Body1"},
				{Model: gorm.Model{ID: 2}, Title: "Title2", Description: "Description2", Body: "Body2"},
			},
			wantErr: false,
		},
		{
			name:    "Handle No Matching User IDs",
			userIDs: []uint{3},
			limit:   2,
			offset:  0,
			mockExpect: func() {
				mock.ExpectQuery(`SELECT * FROM "articles" WHERE \(user_id in \(\?\)\) LIMIT \? OFFSET \?`).
					WithArgs(3, 2, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}))
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Handle Database Errors",
			userIDs: []uint{1},
			limit:   2,
			offset:  0,
			mockExpect: func() {
				mock.ExpectQuery(`SELECT * FROM "articles" WHERE \(user_id in \(\?\)\) LIMIT \? OFFSET \?`).
					WithArgs(1, 2, 0).
					WillReturnError(fmt.Errorf("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test Pagination with Large Offset",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  10,
			mockExpect: func() {
				mock.ExpectQuery(`SELECT * FROM "articles" WHERE \(user_id in \(\?, \?\)\) LIMIT \? OFFSET \?`).
					WithArgs(1, 2, 2, 10).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}))
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Test Articles with Different Limits",
			userIDs: []uint{1},
			limit:   1,
			offset:  0,
			mockExpect: func() {
				mock.ExpectQuery(`SELECT * FROM "articles" WHERE \(user_id in \(\?\)\) LIMIT \? OFFSET \?`).
					WithArgs(1, 1, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).
						AddRow(1, "Title1", "Description1", "Body1"))
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Title1", Description: "Description1", Body: "Body1"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				t.Logf("Expected error received: %v", err)
			} else {
				if len(got) != len(tt.want) {
					t.Errorf("GetFeedArticles() got = %v articles, want %v articles", len(got), len(tt.want))
				} else {
					t.Logf("Successfully fetched %v articles", len(got))
				}
				for i, article := range got {
					if article.Title != tt.want[i].Title {
						t.Errorf("Expected article title %v, but got %v", tt.want[i].Title, article.Title)
					}
					t.Logf("Fetched article: %v", article.Title)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}

/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90


 */
func TestArticleStoreAddFavorite(t *testing.T) {

	type testCase struct {
		name     string
		article  *model.Article
		user     *model.User
		setup    func(mock sqlmock.Sqlmock)
		wantErr  bool
		errMsg   string
		validate func(article *model.Article, t *testing.T)
	}

	tests := []testCase{

		{
			name: "Successfully add a favorite user to an article",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO favorite_articles`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE articles SET favorites_count = favorites_count + \? WHERE id = \?`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
			validate: func(article *model.Article, t *testing.T) {
				if article.FavoritesCount != 1 {
					t.Errorf("Expected favorites_count to be 1, got %d", article.FavoritesCount)
				}
			},
		},

		{
			name: "Attempt to favorite an article twice",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				FavoritesCount: 1,
				FavoritedUsers: []model.User{
					{Model: gorm.Model{ID: 1}, Username: "testuser"},
				},
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO favorite_articles`).WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "invalid transaction",
		},

		{
			name: "DB error in associating user to an article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO favorite_articles`).WillReturnError(gorm.ErrCantStartTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "cannot start transaction",
		},

		{
			name: "DB error during favorites count update",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO favorite_articles`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE articles SET favorites_count = favorites_count + \?`).
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "invalid transaction",
		},

		{
			name:    "Handle nil article input",
			article: nil,
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setup:   func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "nil article",
		},

		{
			name: "Handle nil user input",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user:    nil,
			setup:   func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "nil user",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to open gormDB: %v", err)
			}

			defer gormDB.Close()

			as := &ArticleStore{db: gormDB}

			tc.setup(mock)

			err = as.AddFavorite(tc.article, tc.user)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				if err != nil && err.Error() != tc.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				tc.validate(tc.article, t)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {

	type args struct {
		article *model.Article
		user    *model.User
	}
	type want struct {
		err            error
		favoritesCount int32
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Scenario 1: Successful Deletion of a Favorite User",
			args: args{
				article: &model.Article{
					Model:          gorm.Model{ID: 1},
					FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
				},
				user: &model.User{Model: gorm.Model{ID: 1}},
			},
			want: want{
				err:            nil,
				favoritesCount: 0,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count = favorites_count - ?").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Scenario 2: Error During User Deletion",
			args: args{
				article: &model.Article{
					Model:          gorm.Model{ID: 1},
					FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
				},
				user: &model.User{Model: gorm.Model{ID: 1}},
			},
			want: want{
				err:            fmt.Errorf("delete error"),
				favoritesCount: 0,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(1).
					WillReturnError(fmt.Errorf("delete error"))
				mock.ExpectRollback()
			},
		},
		{
			name: "Scenario 3: Error on Updating Favorites Count",
			args: args{
				article: &model.Article{
					Model:          gorm.Model{ID: 1},
					FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
				},
				user: &model.User{Model: gorm.Model{ID: 1}},
			},
			want: want{
				err:            fmt.Errorf("update error"),
				favoritesCount: 1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count = favorites_count - ?").
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("update error"))
				mock.ExpectRollback()
			},
		},
		{
			name: "Scenario 4: Deleting Favorite from Empty Association List",
			args: args{
				article: &model.Article{
					Model:          gorm.Model{ID: 1},
					FavoritedUsers: []model.User{},
				},
				user: &model.User{Model: gorm.Model{ID: 1}},
			},
			want: want{
				err:            nil,
				favoritesCount: 0,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec("UPDATE articles SET favorites_count = favorites_count - ?").
					WithArgs(0, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Scenario 5: Concurrent Delete Operation",
			args: args{
				article: &model.Article{
					Model:          gorm.Model{ID: 1},
					FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
				},
				user: &model.User{Model: gorm.Model{ID: 1}},
			},
			want: want{
				err:            nil,
				favoritesCount: 0,
			},
			mock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count = favorites_count - ?").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error initializing db mock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("error initializing gorm db: %v", err)
			}

			as := &ArticleStore{db: gormDB}
			tt.mock(mock)

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err = as.DeleteFavorite(tt.args.article, tt.args.user)

			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				fmt.Fscan(r, &buf)
				outC <- buf.String()
			}()

			w.Close()
			os.Stdout = old
			output := <-outC

			if err != nil && tt.want.err == nil || err == nil && tt.want.err != nil {
				t.Errorf("unexpected error value: got %v, want %v", err, tt.want.err)
			}

			if tt.args.article.FavoritesCount != tt.want.favoritesCount {
				t.Errorf("unexpected favorites count: got %v, want %v", tt.args.article.FavoritesCount, tt.want.favoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s': Output => '%s'", tt.name, output)
		})
	}
}

/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestArticleStoreGetArticles(t *testing.T) {
	type testScenario struct {
		name          string
		tagName       string
		username      string
		favoritedBy   *model.User
		limit         int64
		offset        int64
		expectedError error
		expectedCount int
		prepare       func(sqlmock.Sqlmock)
	}

	tests := []testScenario{
		{
			name:     "Scenario 1: Retrieve Articles by Username",
			username: "testuser",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" (.+) JOIN users (.+) WHERE users.username = ?`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name:    "Scenario 2: Filter Articles by Tag",
			tagName: "Technology",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" (.+) JOIN article_tags (.+) JOIN tags (.+) WHERE tags.name = ?`).
					WithArgs("Technology").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name:        "Scenario 3: Retrieve Favorited Articles by a User",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT article_id FROM "favorite_articles" WHERE user_id = ? LIMIT 10 OFFSET 0`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"article_id"}).AddRow(1))
				mock.ExpectQuery(`SELECT (.+) FROM "articles" (.+) WHERE id in (?) LIMIT 10 OFFSET 0`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name:  "Scenario 4: Verify Pagination Limit and Offset",
			limit: 2, offset: 1,
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" WHERE (.+) LIMIT 2 OFFSET 1`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2).AddRow(3))
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Scenario 5: Handle Database Retrieval Error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" (.+)`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: gorm.ErrInvalidSQL,
			expectedCount: 0,
		},
		{
			name: "Scenario 6: Retrieve All Articles with No Filters",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" LIMIT 10 OFFSET 0`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3))
			},
			expectedError: nil,
			expectedCount: 3,
		},
		{
			name:     "Scenario 7: Edge Case - Non-Existent Username or Tag",
			username: "nonexistentuser",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" (.+) JOIN users (.+) WHERE users.username = ?`).
					WithArgs("nonexistentuser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedError: nil,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open SQL mock database: %s", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("failed to initialize GORM DB: %s", err)
			}
			store := ArticleStore{db: gdb}

			tt.prepare(mock)

			result, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.expectedError != nil {
				if err == nil || err != tt.expectedError {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error, but got: %v", err)
				}
				if len(result) != tt.expectedCount {
					t.Errorf("expected %d articles, got: %d", tt.expectedCount, len(result))
				}
			}
		})
	}
}

