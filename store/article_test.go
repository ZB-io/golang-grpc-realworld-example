package store

import (
		"testing"
		"github.com/DATA-DOG/go-sqlmock"
		"github.com/jinzhu/gorm"
		"github.com/stretchr/testify/assert"
		"database/sql"
		"time"
		sqlmock "github.com/DATA-DOG/go-sqlmock"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"errors"
		"sync"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext
}
/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92


 */
func TestNewArticleStore(t *testing.T) {

	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name:     "Successful creation with valid DB",
			db:       setupTestDB(t),
			wantNil:  false,
			scenario: "Scenario 1: Successfully Create ArticleStore with Valid DB Connection",
		},
		{
			name:     "Creation with nil DB",
			db:       nil,
			wantNil:  false,
			scenario: "Scenario 2: Create ArticleStore with Nil DB Connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("Executing: %s", tt.scenario)

			store := NewArticleStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, store, "Expected nil ArticleStore")
			} else {
				assert.NotNil(t, store, "Expected non-nil ArticleStore")
				assert.Equal(t, tt.db, store.db, "DB reference mismatch")
			}

			if tt.db != nil {
				store1 := NewArticleStore(tt.db)
				store2 := NewArticleStore(tt.db)
				assert.NotSame(t, store1, store2, "Instances should be independent")
				assert.Same(t, store1.db, store2.db, "DB reference should be the same")
				t.Log("Verified instance independence")
			}

			if tt.db != nil {
				assert.Same(t, tt.db, store.db, "DB reference should remain unchanged")
				t.Log("Verified DB reference integrity")
			}

			func() {
				var stores []*ArticleStore
				for i := 0; i < 100; i++ {
					stores = append(stores, NewArticleStore(tt.db))
				}

				stores = nil
				t.Log("Completed memory management test")
			}()

			t.Log("Test case completed successfully")
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled mock expectations: %v", err)
		}
	})

	return gormDB
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
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name    string
		article *model.Article
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful article creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Model: gorm.Model{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"Test Article",
						"Test Description",
						"Test Body",
						uint(1),
						int32(0),
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Missing required fields",
			article: &model.Article{
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Article with tags",
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
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"Test Article with Tags",
						"Test Description",
						"Test Body",
						uint(1),
						int32(0),
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `article_tags`").
					WithArgs(
						uint(1), "tag1",
						uint(1), "tag2",
					).WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock)

			err := store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
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
			name: "Missing required fields",
			comment: &model.Comment{
				Body: "",
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
	type testCase struct {
		name    string
		article *model.Article
		mockDB  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}

	tests := []testCase{
		{
			name: "Successfully delete existing article",
			article: &model.Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles` WHERE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Attempt to delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
			},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles` WHERE").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
			},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles` WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:    "Nil article input",
			article: nil,
			mockDB:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "invalid article: nil pointer",
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
			gormDB.LogMode(true)
			defer gormDB.Close()

			if tc.mockDB != nil {
				tc.mockDB(mock)
			}

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Delete(tc.article)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
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
		comment     *model.Comment
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
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
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Delete Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 999},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: no rows in result set",
		},
		{
			name:    "Delete with NULL Input",
			comment: nil,
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "comment cannot be nil",
		},
		{
			name: "Delete with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnError(errors.New("database connection lost"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "database connection lost",
		},
		{
			name: "Delete with Foreign Key Constraint",
			comment: &model.Comment{
				Model: gorm.Model{ID: 2},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(2).
					WillReturnError(errors.New("foreign key constraint fails"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "foreign key constraint fails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			err := store.DeleteComment(tt.comment)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
		})
	}

	t.Run("Concurrent Comment Deletion", func(t *testing.T) {
		comment := &model.Comment{
			Model: gorm.Model{ID: 3},
		}

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `comments`").
			WithArgs(3).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		var wg sync.WaitGroup
		wg.Add(2)

		for i := 0; i < 2; i++ {
			go func() {
				defer wg.Done()
				err := store.DeleteComment(comment)
				if err != nil {
					t.Errorf("Concurrent deletion error: %v", err)
				}
			}()
		}

		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled mock expectations in concurrent test: %v", err)
		}
	})
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {

	type testCase struct {
		name            string
		commentID       uint
		mockSetup       func(sqlmock.Sqlmock)
		expectedError   error
		expectedComment *model.Comment
	}

	tests := []testCase{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE `comments`.`id` = \\? AND `comments`.`deleted_at` IS NULL").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "Test comment", 1, 1))
			},
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			expectedError: nil,
		},
		{
			name:      "Non-existent comment",
			commentID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
		{
			name:      "Database connection error",
			commentID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedComment: nil,
			expectedError:   sql.ErrConnDone,
		},
		{
			name:      "Zero ID input",
			commentID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments`").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
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

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			comment, err := store.GetCommentByID(tc.commentID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tc.expectedComment.ID, comment.ID)
				assert.Equal(t, tc.expectedComment.Body, comment.Body)
				assert.Equal(t, tc.expectedComment.UserID, comment.UserID)
				assert.Equal(t, tc.expectedComment.ArticleID, comment.ArticleID)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
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
			name: "Successfully retrieve tags",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, time.Now(), time.Now(), nil, "golang").
					AddRow(2, time.Now(), time.Now(), nil, "testing")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "golang"},
				{Model: gorm.Model{ID: 2}, Name: "testing"},
			},
			expectedError: nil,
		},
		{
			name: "Empty tags list",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"})
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database connection error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnError(errors.New("database connection error"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Database query timeout",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnError(errors.New("context deadline exceeded"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("context deadline exceeded"),
		},
		{
			name: "Large dataset handling",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"})

				for i := 1; i <= 1000; i++ {
					rows.AddRow(i, time.Now(), time.Now(), nil, "tag"+string(rune(i)))
				}
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnRows(rows)
			},
			expectedTags: func() []model.Tag {
				var tags []model.Tag
				for i := 1; i <= 1000; i++ {
					tags = append(tags, model.Tag{
						Model: gorm.Model{ID: uint(i)},
						Name:  "tag" + string(rune(i)),
					})
				}
				return tags
			}(),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB instance: %v", err)
			}
			defer gormDB.Close()

			tt.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			tags, err := store.GetTags()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedTags), len(tags))

				if len(tags) > 0 {
					for i := range tags {
						assert.Equal(t, tt.expectedTags[i].ID, tags[i].ID)
						assert.Equal(t, tt.expectedTags[i].Name, tags[i].Name)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func TestArticleStoreGetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
		expectArticle bool
		validateFunc  func(*testing.T, *model.Article)
	}{
		{
			name: "Successfully retrieve article with valid ID",
			id:   1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE "articles"."id" = \? AND "articles"."deleted_at" IS NULL`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
						AddRow(1, "Test Article", "Test Description", "Test Body", 1))

				mock.ExpectQuery(`SELECT \* FROM "tags" INNER JOIN "article_tags" ON "article_tags"."tag_id" = "tags"."id" WHERE "article_tags"."article_id" = \?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "golang").
						AddRow(2, "testing"))

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \? AND "users"."deleted_at" IS NULL`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
						AddRow(1, "testuser", "test@example.com"))
			},
			expectedError: nil,
			expectArticle: true,
			validateFunc: func(t *testing.T, article *model.Article) {
				assert.NotNil(t, article)
				assert.Equal(t, uint(1), article.ID)
				assert.Equal(t, "Test Article", article.Title)
				assert.Equal(t, "Test Description", article.Description)
				assert.Equal(t, "Test Body", article.Body)
				assert.Len(t, article.Tags, 2)
				assert.NotNil(t, article.Author)
			},
		},
		{
			name: "Article not found",
			id:   999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectArticle: false,
			validateFunc:  nil,
		},
		{
			name: "Database connection error",
			id:   1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
			expectArticle: false,
			validateFunc:  nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectArticle: false,
			validateFunc:  nil,
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

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			store := &ArticleStore{db: gormDB}

			article, err := store.GetByID(tt.id)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectArticle {
				assert.NotNil(t, article)
				if tt.validateFunc != nil {
					tt.validateFunc(t, article)
				}
			} else {
				assert.Nil(t, article)
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
			name: "Update with Empty Required Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("NOT NULL constraint failed"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "NOT NULL constraint failed",
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
					WillReturnError(errors.New("database connection lost"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "database connection lost",
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
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test comment 1", 1, 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Test comment 2", 1, 1, 1)

				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE .*article_id = \\?.*").
					WithArgs(1).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "testuser")
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE .*").
					WillReturnRows(authorRows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id", "author_id"})
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE .*article_id = \\?.*").
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
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE .*article_id = \\?.*").
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

			t.Logf("Test '%s' completed. Comments count: %d, Error: %v",
				tt.name, len(comments), err)
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

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		mockSetup   func(sqlmock.Sqlmock)
		want        bool
		wantErr     bool
		expectedErr error
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
			want:    true,
			wantErr: false,
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
			want:    false,
			wantErr: false,
		},
		{
			name:    "Nil Article Parameter",
			article: nil,
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      false,
			wantErr:   false,
		},
		{
			name: "Nil User Parameter",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user:      nil,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      false,
			wantErr:   false,
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
					WillReturnError(errors.New("database error"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name: "Invalid Article ID",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(0, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Multiple Favorite Relationships",
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
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			want:    true,
			wantErr: false,
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

			if tt.wantErr && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("ArticleStore.IsFavorited() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
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
			name:    "Successful retrieval of articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test Article", "Test Description", "Test Body", 1, 0,
				)

				authorRows := sqlmock.NewRows([]string{
					"id", "username", "email",
				}).AddRow(1, "testuser", "test@example.com")

				mock.ExpectQuery("SELECT").WillReturnRows(rows)
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
			},
			want: []model.Article{
				{
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
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
			name:    "Empty result set",
			userIDs: []uint{99},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
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
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Pagination boundary - zero limit",
			userIDs: []uint{1},
			limit:   0,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
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
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
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
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
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
			name: "Database Error During Association",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
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
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
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
				Email:    "test@example.com",
				Password: "password",
			},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article"),
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

			tc.setupMock(mock)

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

		store := &ArticleStore{db: gormDB}
		article := &model.Article{
			Model:          gorm.Model{ID: 1},
			Title:          "Test Article",
			Description:    "Test Description",
			Body:           "Test Body",
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

		for i := 0; i < numUsers; i++ {
			go func(userID uint) {
				defer wg.Done()
				user := &model.User{
					Model:    gorm.Model{ID: userID},
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
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
	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successful deletion of favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(0, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Association deletion error",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(errors.New("association deletion error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("association deletion error"),
		},
		{
			name: "Update favorites count error",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(0, 1).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update error"),
		},
		{
			name:          "Null article parameter",
			article:       nil,
			user:          &model.User{Model: gorm.Model{ID: 1}},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article"),
		},
		{
			name:          "Null user parameter",
			article:       &model.Article{Model: gorm.Model{ID: 1}},
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

			store := &ArticleStore{db: gormDB}

			if tc.article == nil {
				err = store.DeleteFavorite(tc.article, tc.user)
				assert.Error(t, err)
				assert.Equal(t, "invalid article", err.Error())
				return
			}

			if tc.user == nil {
				err = store.DeleteFavorite(tc.article, tc.user)
				assert.Error(t, err)
				assert.Equal(t, "invalid user", err.Error())
				return
			}

			tc.setupMock(mock)

			err = store.DeleteFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(0), tc.article.FavoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
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
		t.Fatalf("Failed to open GORM DB: %v", err)
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
		wantErr     bool
		expected    []model.Article
	}{
		{
			name:     "Scenario 1: Successfully Retrieve Articles Without Filters",
			tagName:  "",
			username: "",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				articleRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test Article", "Description", "Body", 1, 0,
				)

				mock.ExpectQuery("SELECT").WillReturnRows(articleRows)

				authorRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "password", "bio", "image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"testuser", "test@example.com", "password", "bio", "image",
				)
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
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
		},
		{
			name:     "Scenario 2: Filter Articles by Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expected: []model.Article{},
		},
		{
			name:    "Scenario 3: Filter Articles by Tag",
			tagName: "programming",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expected: []model.Article{},
		},
		{
			name: "Scenario 4: Filter Articles by Favorite User",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favRows := sqlmock.NewRows([]string{"article_id"}).AddRow(1)
				mock.ExpectQuery("SELECT article_id FROM").WillReturnRows(favRows)

				articleRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(articleRows)
			},
			expected: []model.Article{},
		},
		{
			name: "Scenario 8: Database Error Handling",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, articles)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

