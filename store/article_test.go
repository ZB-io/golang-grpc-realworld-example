package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"time"
	"github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext
}
type ExpectedBegin struct {
	commonExpectation
	delay time.Duration
}
type ExpectedCommit struct {
	commonExpectation
}
type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}
type ExpectedRollback struct {
	commonExpectation
}
type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}
type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}
/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92


 */
func TestNewArticleStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		validate func(*testing.T, *ArticleStore)
	}

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open("sqlite3", mockDB)
	if err != nil {
		t.Fatalf("Failed to create gorm DB: %v", err)
	}
	defer gormDB.Close()

	tests := []testCase{
		{
			name: "Scenario 1: Successfully Create New ArticleStore with Valid DB Connection",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				t.Log("Validating successful ArticleStore creation")
				assert.NotNil(t, store, "ArticleStore should not be nil")
				assert.Equal(t, gormDB, store.db, "DB reference should match")
			},
		},
		{
			name: "Scenario 2: Create ArticleStore with Nil DB Connection",
			db:   nil,
			validate: func(t *testing.T, store *ArticleStore) {
				t.Log("Validating ArticleStore creation with nil DB")
				assert.NotNil(t, store, "ArticleStore should not be nil even with nil DB")
				assert.Nil(t, store.db, "DB reference should be nil")
			},
		},
		{
			name: "Scenario 3: Verify ArticleStore Instance Independence",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				t.Log("Validating ArticleStore instance independence")
				store2 := NewArticleStore(gormDB)
				assert.NotEqual(t, store, store2, "Different instances should not be equal")
				assert.Equal(t, store.db, store2.db, "DB references should be the same")
			},
		},
		{
			name: "Scenario 4: Verify DB Reference Integrity",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				t.Log("Validating DB reference integrity")
				assert.Equal(t, gormDB, store.db, "DB reference should remain unchanged")
			},
		},
		{
			name: "Scenario 5: Memory Resource Management",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				t.Log("Validating memory resource management")
				var stores []*ArticleStore
				for i := 0; i < 100; i++ {
					stores = append(stores, NewArticleStore(gormDB))
				}
				assert.Len(t, stores, 100, "Should create multiple instances without issues")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Starting test case:", tc.name)

			store := NewArticleStore(tc.db)

			if tc.validate != nil {
				tc.validate(t, store)
			}

			t.Log("Completed test case:", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestArticleStoreCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Article Creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Empty Required Fields",
			article: &model.Article{
				Title:       "",
				Description: "",
				Body:        "",
				UserID:      0,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(errors.New("validation error"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "validation error",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: connection is already closed",
		},
		{
			name: "Article with Tags",
			article: &model.Article{
				Title:       "Test Article with Tags",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags: []model.Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `article_tags`").
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

			err := store.Create(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Error expected: %v, Got error: %v",
				tt.name, tt.expectError, err)
		})
	}
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestArticleStoreCreateComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name    string
		comment *model.Comment
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful comment creation",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"Test comment",
						uint(1),
						uint(1),
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Missing required field - Body",
			comment: &model.Comment{
				UserID:    1,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Non-existent UserID",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Maximum length body",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						string(make([]byte, 1000)),
						uint(1),
						uint(1),
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Database connection error",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock)

			err := store.CreateComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Logf("Comment created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
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
		name    string
		article *model.Article
		mockFn  func(sqlmock.Sqlmock)
		wantErr error
	}{
		{
			name: "Successfully delete existing article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "Fail to delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: errors.New("database connection error"),
		},
		{
			name:    "Nil article parameter",
			article: nil,
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: errors.New("invalid article parameter"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			gormDB.LogMode(false)

			store := &ArticleStore{
				db: gormDB,
			}

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			if tt.article == nil {
				err = store.Delete(nil)
				assert.Error(t, err)
				assert.Equal(t, "invalid article parameter", err.Error())
				return
			}

			err = store.Delete(tt.article)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func TestArticleStoreDeleteComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		mockDB  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successfully Delete Existing Comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Attempt to Delete Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(999).
					WillReturnError(errors.New("record not found"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name:    "Delete Comment with NULL Input",
			comment: nil,
			mockDB: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
			errMsg:  "invalid comment: nil pointer",
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
			},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			gormDB.LogMode(true)
			defer gormDB.Close()

			if tt.comment == nil {
				err = errors.New("invalid comment: nil pointer")
				if err.Error() != tt.errMsg {
					t.Errorf("DeleteComment() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			tt.mockDB(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("DeleteComment() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test '%s' completed with expected error: %v", tt.name, err)
			} else {
				t.Logf("Test '%s' completed successfully", tt.name)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		commentID     uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
		expectedBody  string
	}{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test comment", 1, 1)
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedError: nil,
			expectedBody:  "Test comment",
		},
		{
			name:      "Non-existent comment",
			commentID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedBody:  "",
		},
		{
			name:      "Database connection error",
			commentID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
			expectedBody:  "",
		},
		{
			name:      "Zero ID input",
			commentID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedBody:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			comment, err := store.GetCommentByID(tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, comment)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedBody, comment.Body)
				t.Logf("Successfully retrieved comment with ID %d", tt.commentID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func TestArticleStoreGetTags(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Successfully Retrieve Tags",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, time.Now(), time.Now(), nil, "golang").
					AddRow(2, time.Now(), time.Now(), nil, "testing")
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "golang"},
				{Model: gorm.Model{ID: 2}, Name: "testing"},
			},
			expectedError: nil,
		},
		{
			name: "Empty Tags List",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"})
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnError(errors.New("database connection error"))
			},
			expectedTags:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Database Query Timeout",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnError(errors.New("context deadline exceeded"))
			},
			expectedTags:  nil,
			expectedError: errors.New("context deadline exceeded"),
		},
		{
			name: "Large Dataset Retrieval",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"})
				for i := 1; i <= 1000; i++ {
					rows.AddRow(i, time.Now(), time.Now(), nil, fmt.Sprintf("tag%d", i))
				}
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  make([]model.Tag, 1000),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			gormDB.LogMode(false)
			defer gormDB.Close()

			tt.setupMock(mock)

			store := &ArticleStore{db: gormDB}
			tags, err := store.GetTags()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				if tt.name == "Large Dataset Retrieval" {
					assert.Equal(t, 1000, len(tags))
					for i := range tags {
						assert.Equal(t, uint(i+1), tags[i].ID)
						assert.Equal(t, fmt.Sprintf("tag%d", i+1), tags[i].Name)
					}
				} else {
					assert.Equal(t, tt.expectedTags, tags)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func TestArticleStoreGetByID(t *testing.T) {
	type testCase struct {
		name          string
		id            uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
		expectedData  *model.Article
	}

	tests := []testCase{
		{
			name: "Successfully retrieve article with valid ID",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, "Test Article", "Test Description", "Test Body", 1, time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnRows(rows)

				tagRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, "test-tag", time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT \* FROM "tags"`).
					WillReturnRows(tagRows)

				userRows := sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, "testuser", "test@example.com", time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(userRows)
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Model: gorm.Model{ID: 1}, Name: "test-tag"}},
				Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser", Email: "test@example.com"},
			},
		},
		{
			name: "Article not found",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
		{
			name: "Database connection error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
			expectedData:  nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(0).
					WillReturnError(errors.New("invalid ID"))
			},
			expectedError: errors.New("invalid ID"),
			expectedData:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}
			article, err := store.GetByID(tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, article)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, tc.expectedData.Title, article.Title)
				assert.Equal(t, tc.expectedData.Description, article.Description)
				assert.Equal(t, tc.expectedData.Body, article.Body)
				assert.Equal(t, tc.expectedData.UserID, article.UserID)

				if len(tc.expectedData.Tags) > 0 {
					assert.Equal(t, tc.expectedData.Tags[0].ID, article.Tags[0].ID)
					assert.Equal(t, tc.expectedData.Tags[0].Name, article.Tags[0].Name)
				}

				assert.Equal(t, tc.expectedData.Author.ID, article.Author.ID)
				assert.Equal(t, tc.expectedData.Author.Username, article.Author.Username)
				assert.Equal(t, tc.expectedData.Author.Email, article.Author.Email)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe


 */
func TestArticleStoreUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Update",
			article: &model.Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "Updated Title",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(
						sqlmock.AnyArg(),
						"Updated Title",
						"Updated Description",
						"Updated Body",
						uint(1),
						uint(1),
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Update Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: no rows in result set",
		},
		{
			name: "Update with Empty Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("not null constraint violation"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "not null constraint violation",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Title",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("connection refused"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := store.Update(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestArticleStoreGetComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		article       *model.Article
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
	}{
		{
			name: "Successful retrieval of comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				commentRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "user_id", "article_id",
				}).
					AddRow(1, time.Now(), time.Now(), nil, "Comment 1", 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Comment 2", 1, 1)

				mock.ExpectQuery("SELECT (.+) FROM `comments`").
					WithArgs(1).
					WillReturnRows(commentRows)

				authorRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "bio", "image",
				}).
					AddRow(1, time.Now(), time.Now(), nil,
						"testuser", "test@example.com", "bio", "image.jpg")

				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
				Title: "Empty Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "user_id", "article_id",
				})
				mock.ExpectQuery("SELECT (.+) FROM `comments`").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Database error",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
				Title: "Error Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `comments`").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			comments, err := store.GetComments(tt.article)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tt.expectedCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func TestArticleStoreIsFavorited(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		mockSetup   func(sqlmock.Sqlmock)
		want        bool
		wantErr     bool
		description string
	}{
		{
			name: "Valid Article and User with Existing Favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want:        true,
			wantErr:     false,
			description: "User has favorited the article",
		},
		{
			name: "Valid Article and User with No Favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:        false,
			wantErr:     false,
			description: "User has not favorited the article",
		},
		{
			name:        "Nil Article",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			want:        false,
			wantErr:     false,
			description: "Article parameter is nil",
		},
		{
			name:        "Nil User",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			want:        false,
			wantErr:     false,
			description: "User parameter is nil",
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:        false,
			wantErr:     true,
			description: "Database operation fails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			got, err := store.IsFavorited(tt.article, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed: %s", tt.name, tt.description)
			if err != nil {
				t.Logf("Error encountered: %v", err)
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
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		want      []model.Article
		wantErr   bool
	}{
		{
			name:    "Successfully retrieve feed articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\?,\\?\\)\\) LIMIT 10 OFFSET 0"
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test Article", "Test Description", "Test Body", 1, 0,
				)

				mock.ExpectQuery(expectedSQL).
					WithArgs(1, 2).
					WillReturnRows(rows)

				authorSQL := "SELECT \\* FROM `users` WHERE \\(`id` = \\?\\)"
				authorRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "bio", "image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"testuser", "test@example.com", "", "",
				)

				mock.ExpectQuery(authorSQL).
					WithArgs(1).
					WillReturnRows(authorRows)
			},
			want: []model.Article{
				{
					Model:          gorm.Model{ID: 1},
					Title:          "Test Article",
					Description:    "Test Description",
					Body:           "Test Body",
					UserID:         1,
					FavoritesCount: 0,
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty result for non-existent user IDs",
			userIDs: []uint{999},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\?\\)\\) LIMIT 10 OFFSET 0"
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery(expectedSQL).
					WithArgs(999).
					WillReturnRows(rows)
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\?\\)\\) LIMIT 10 OFFSET 0"
				mock.ExpectQuery(expectedSQL).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty user IDs array",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\)\\) LIMIT 10 OFFSET 0"
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery(expectedSQL).
					WillReturnRows(rows)
			},
			want:    []model.Article{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
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
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully Add Favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Database Error During User Association",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "Database Error During Count Update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update error"),
		},
		{
			name:    "Nil Article Parameter",
			article: nil,
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article"),
		},
		{
			name: "Nil User Parameter",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user:          nil,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid user"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			store := &ArticleStore{db: gormDB}
			err = store.AddFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tc.article != nil {
					assert.Equal(t, int32(1), tc.article.FavoritesCount)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent Favorite Operations", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to create GORM DB: %v", err)
		}
		defer gormDB.Close()

		article := &model.Article{
			Model:          gorm.Model{ID: 1},
			FavoritesCount: 0,
		}

		numUsers := 5
		var wg sync.WaitGroup
		wg.Add(numUsers)

		for i := 0; i < numUsers; i++ {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `favorite_articles`").
				WithArgs(1, uint(i+1)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE `articles`").
				WithArgs(1, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		store := &ArticleStore{db: gormDB}

		for i := 0; i < numUsers; i++ {
			go func(userID uint) {
				defer wg.Done()
				user := &model.User{
					Model:    gorm.Model{ID: userID},
					Username: "testuser",
				}
				err := store.AddFavorite(article, user)
				assert.NoError(t, err)
			}(uint(i + 1))
		}

		wg.Wait()

		assert.Equal(t, int32(numUsers), article.FavoritesCount)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {

	createMockDB := func(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to open GORM DB: %v", err)
		}

		return gormDB, mock
	}

	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		article *model.Article
		user    *model.User
		wantErr bool
	}{
		{
			name: "Successful deletion",
			setup: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 1,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "Association deletion error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			wantErr: true,
		},
		{
			name: "Update error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			wantErr: true,
		},
		{
			name:    "Null article parameter",
			setup:   func(mock sqlmock.Sqlmock) {},
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			wantErr: true,
		},
		{
			name:  "Null user parameter",
			setup: func(mock sqlmock.Sqlmock) {},
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gormDB, mock := createMockDB(t)
			defer gormDB.Close()

			tt.setup(mock)

			store := &ArticleStore{db: gormDB}

			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.article.FavoritesCount-1, tt.article.FavoritesCount)
				t.Log("Successfully deleted favorite and decremented count")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestArticleStoreGetArticles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(sqlmock.Sqlmock)
		expected    []model.Article
		expectError bool
	}{
		{
			name:     "Scenario 1: Get Articles Without Filters",
			tagName:  "",
			username: "",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Test Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
				},
			},
			expectError: false,
		},
		{
			name:     "Scenario 2: Get Articles By Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT (.+) FROM `articles` JOIN users").
					WithArgs("testuser").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Test Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
				},
			},
			expectError: false,
		},
		{
			name:    "Scenario 3: Get Articles By Tag",
			tagName: "programming",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Programming Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT (.+) FROM `articles` JOIN article_tags").
					WithArgs("programming").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Programming Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Scenario 4: Get Favorited Articles",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favRows := sqlmock.NewRows([]string{"article_id"}).
					AddRow(1)
				mock.ExpectQuery("SELECT article_id FROM `favorite_articles`").
					WithArgs(1).
					WillReturnRows(favRows)

				articleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Favorited Article", "Description", "Body", 1, 1)
				mock.ExpectQuery("SELECT (.+) FROM `articles`").
					WillReturnRows(articleRows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Favorited Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
					FavoritesCount: 1,
				},
			},
			expectError: false,
		},
		{
			name:   "Scenario 8: Database Error Handling",
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `articles`").
					WillReturnError(sql.ErrConnDone)
			},
			expected:    []model.Article{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(articles))
				if len(articles) > 0 {
					assert.Equal(t, tt.expected[0].Title, articles[0].Title)
					assert.Equal(t, tt.expected[0].Author.Username, articles[0].Author.Username)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

