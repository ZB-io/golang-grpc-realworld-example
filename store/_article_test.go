package store

import (
	errors "errors"
	debug "runtime/debug"
	sync "sync"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	sql "database/sql"
	time "time"
	fmt "fmt"
	assert "github.com/stretchr/testify/assert"
)





type MockDB struct {
	*gorm.DB
	mock sqlmock.Sqlmock
}


/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


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
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Handle Duplicate Favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				FavoritesCount: 1,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("duplicate entry"),
		},
		{
			name: "Handle DB Transaction Failure During Association",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("db error"),
		},
		{
			name:    "Handle Null Input Parameters",
			article: nil,
			user:    nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid parameters"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

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

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.AddFavorite(tc.article, tc.user)

			if (err != nil && tc.expectedError == nil) || (err == nil && tc.expectedError != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case '%s' failed with error: %v", tc.name, err)
			} else {
				t.Logf("Test case '%s' passed successfully", tc.name)
			}
		})
	}

	t.Run("Concurrent Access", func(t *testing.T) {
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
		article := &model.Article{Model: gorm.Model{ID: 1}}

		numGoroutines := 5
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `favorite_articles`").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE `articles`").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			go func(userID uint) {
				defer wg.Done()
				user := &model.User{Model: gorm.Model{ID: userID}}
				err := store.AddFavorite(article, user)
				if err != nil {
					t.Errorf("Concurrent execution failed: %v", err)
				}
			}(uint(i + 1))
		}

		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled mock expectations in concurrent test: %v", err)
		}
	})
}


/*
ROOST_METHOD_HASH=ArticleStore_Create_1273475ade
ROOST_METHOD_SIG_HASH=ArticleStore_Create_a27282cad5

FUNCTION_DEF=func (s *ArticleStore) Create(m *model.Article) error // Create creates an article


*/
func TestArticleStoreCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		article *model.Article
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success - Create article with all required fields",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test"}},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failure - Missing required fields",
			article: &model.Article{

				Body:   "Test Body",
				UserID: 1,
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
			name: "Success - Create article with tags",
			article: &model.Article{
				Title:       "Test Article with Tags",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `article_tags`").
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failure - Database connection error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mockFn(mock)

			err := store.Create(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
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
ROOST_METHOD_HASH=ArticleStore_CreateComment_b16d4a71d4
ROOST_METHOD_SIG_HASH=ArticleStore_CreateComment_7475736b06

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error // CreateComment creates a comment of the article


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
		mock    func()
		wantErr bool
	}{
		{
			name: "Success - Create Valid Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test comment", uint(1), uint(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failure - Missing Required Fields",
			comment: &model.Comment{
				Body: "Test comment",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failure - Non-existent Foreign Keys",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 999,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failure - Database Connection Error",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Success - Maximum Length Content",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			t.Log("Starting test case:", tt.name)

			tt.mock()

			now := time.Now()
			tt.comment.CreatedAt = now
			tt.comment.UpdatedAt = now

			err := store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err == nil {
				t.Log("Successfully created comment")
			} else {
				t.Log("Failed to create comment:", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_Delete_8daad9ff19
ROOST_METHOD_SIG_HASH=ArticleStore_Delete_0e09651031

FUNCTION_DEF=func (s *ArticleStore) Delete(m *model.Article) error // Delete deletes an article


*/
func TestArticleStoreDelete(t *testing.T) {

	type testCase struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}

	tests := []testCase{
		{
			name: "Successfully Delete Existing Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Delete Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: connection is already closed",
		},
		{
			name: "Delete Article with Associated Records",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}},
				},
				Comments: []model.Comment{
					{Model: gorm.Model{ID: 1}},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Delete Article with Null Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
				Body:  "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

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

			tc.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Delete(tc.article)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
				t.Logf("Successfully caught expected error: %v", err)
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				t.Log("Successfully deleted article")
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_DeleteComment_effbcb38aa
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteComment_d3c99623e4

FUNCTION_DEF=func (s *ArticleStore) DeleteComment(m *model.Comment) error // DeleteComment deletes an comment


*/
func TestArticleStoreDeleteComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		comment *model.Comment
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successfully delete existing comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete non-existent comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database connection error",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Delete comment with null fields",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 2,
				},
				Body:      "",
				UserID:    0,
				ArticleID: 0,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mockFn(mock)

			err := store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
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
ROOST_METHOD_HASH=ArticleStore_DeleteFavorite_29c18a04a8
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteFavorite_53deb5e792

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error // DeleteFavorite unfavorite an article


*/
func TestArticleStoreDeleteFavorite(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successful unfavorite",
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
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Failed association removal",
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
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedError: fmt.Errorf("database error"),
		},
		{
			name: "Failed favorites count update",
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
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("update error"))
				mock.ExpectRollback()
			},
			expectedError: fmt.Errorf("update error"),
		},
		{
			name:    "Null parameter handling",
			article: nil,
			user:    nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: fmt.Errorf("invalid parameters"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

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

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.DeleteFavorite(tc.article, tc.user)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if (tc.expectedError == nil && err != nil) ||
				(tc.expectedError != nil && err == nil) ||
				(tc.expectedError != nil && err != nil && tc.expectedError.Error() != err.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			if err == nil && tc.article != nil {
				if tc.article.FavoritesCount != 0 {
					t.Errorf("Expected FavoritesCount to be 0, got: %d", tc.article.FavoritesCount)
				}
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=ArticleStore_GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error) 

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

	testUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Email:    "test@example.com",
	}

	testArticles := []model.Article{
		{
			Model:  gorm.Model{ID: 1},
			Title:  "Test Article 1",
			Body:   "Test Body 1",
			UserID: 1,
			Author: *testUser,
		},
		{
			Model:  gorm.Model{ID: 2},
			Title:  "Test Article 2",
			Body:   "Test Body 2",
			UserID: 1,
			Author: *testUser,
		},
	}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(sqlmock.Sqlmock)
		want        []model.Article
		wantErr     bool
	}{
		{
			name: "Get All Articles Without Filters",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "user_id", "created_at", "updated_at"})
				for _, article := range testArticles {
					rows.AddRow(article.ID, article.Title, article.Body, article.UserID, time.Now(), time.Now())
				}
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").WillReturnRows(rows)
			},
			want:    testArticles,
			wantErr: false,
		},
		{
			name:     "Filter By Username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "user_id", "created_at", "updated_at"})
				rows.AddRow(testArticles[0].ID, testArticles[0].Title, testArticles[0].Body, testArticles[0].UserID, time.Now(), time.Now())
				mock.ExpectQuery("^SELECT (.+) FROM `articles` JOIN users").WillReturnRows(rows)
			},
			want:    testArticles[:1],
			wantErr: false,
		},
		{
			name:    "Filter By Tag",
			tagName: "programming",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "user_id", "created_at", "updated_at"})
				mock.ExpectQuery("^SELECT (.+) FROM `articles` JOIN article_tags").WillReturnRows(rows)
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:   "Test Pagination",
			limit:  1,
			offset: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "user_id", "created_at", "updated_at"})
				rows.AddRow(testArticles[1].ID, testArticles[1].Title, testArticles[1].Body, testArticles[1].UserID, time.Now(), time.Now())
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").WillReturnRows(rows)
			},
			want:    testArticles[1:],
			wantErr: false,
		},
		{
			name: "Database Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrConnDone)
			},
			want:    []model.Article{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mockSetup(mock)

			got, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ArticleStore.GetArticles() got %d articles, want %d", len(got), len(tt.want))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetByID_6fe18728fc
ROOST_METHOD_SIG_HASH=ArticleStore_GetByID_bb488e542f

FUNCTION_DEF=func (s *ArticleStore) GetByID(id uint) (*model.Article, error) // GetByID finds an article from id


*/
func TestArticleStoreGetById(t *testing.T) {

	type testCase struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
		expectedData  *model.Article
	}

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

	sampleTime := time.Now()
	sampleArticle := &model.Article{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: sampleTime,
			UpdatedAt: sampleTime,
		},
		Title:       "Test Article",
		Description: "Test Description",
		Body:        "Test Body",
		Tags: []model.Tag{
			{Model: gorm.Model{ID: 1}, Name: "tag1"},
		},
		Author: model.User{
			Model:    gorm.Model{ID: 1},
			Username: "testuser",
			Email:    "test@example.com",
		},
		UserID:         1,
		FavoritesCount: 0,
	}

	tests := []testCase{
		{
			name: "Successful Retrieval",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, sampleTime, sampleTime, nil, "Test Article", "Test Description", "Test Body", 1, 0)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(1).
					WillReturnRows(rows)

				tagRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "tag1")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnRows(tagRows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedData:  sampleArticle,
			expectedError: nil,
		},
		{
			name: "Article Not Found",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedData:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedData:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Zero ID",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedData:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tc.mockSetup(mock)

			article, err := store.GetByID(tc.id)

			if err != tc.expectedError {
				if !errors.Is(err, tc.expectedError) {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
			}

			if tc.expectedData != nil {
				if article == nil {
					t.Error("Expected article data, got nil")
				} else {

					if article.ID != tc.expectedData.ID {
						t.Errorf("Expected article ID %d, got %d", tc.expectedData.ID, article.ID)
					}
					if article.Title != tc.expectedData.Title {
						t.Errorf("Expected article title %s, got %s", tc.expectedData.Title, article.Title)
					}

				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetComments_d7c78dda64
ROOST_METHOD_SIG_HASH=ArticleStore_GetComments_af08ddd59e

FUNCTION_DEF=func (s *ArticleStore) GetComments(m *model.Article) ([ // GetComments gets coments of the article
]model.Comment, error) 

*/
func TestArticleStoreGetComments(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
	}

	tests := []testCase{
		{
			name: "Successfully retrieve comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test comment 1", 1, 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Test comment 2", 2, 1, 2)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "user1", "user1@test.com").
					AddRow(2, "user2", "user2@test.com")

				mock.ExpectQuery(`SELECT \* FROM "comments".*WHERE.*article_id.*=.*`).
					WithArgs(1).
					WillReturnRows(rows)

				mock.ExpectQuery(`SELECT \* FROM "users".*WHERE.*id.*IN.*`).
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
				mock.ExpectQuery(`SELECT \* FROM "comments".*WHERE.*article_id.*=.*`).
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
				mock.ExpectQuery(`SELECT \* FROM "comments".*WHERE.*article_id.*=.*`).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

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

			tc.mockSetup(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			comments, err := store.GetComments(tc.article)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, comments)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tc.expectedCount)
				t.Logf("Successfully retrieved %d comments", len(comments))

				if len(comments) > 0 {
					for _, comment := range comments {
						assert.NotEmpty(t, comment.Author)
						assert.Equal(t, comment.UserID, comment.Author.ID)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetFeedArticles_a37e1934b6
ROOST_METHOD_SIG_HASH=ArticleStore_GetFeedArticles_f5f09c020e

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs [ // GetFeedArticles returns following users' articles
]uint, limit, offset int64) ([]model.Article, error) 

*/
func TestArticleStoreGetFeedArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		wantLen   int
		wantErr   bool
	}{
		{
			name:    "Success - Multiple Users Articles",
			userIDs: []uint{1, 2, 3},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
					"author_id", "author_username", "author_email", "author_bio", "author_image",
				}).
					AddRow(1, time.Now(), time.Now(), nil,
						"Title 1", "Desc 1", "Body 1", 1, 0,
						1, "user1", "user1@test.com", "bio1", "image1").
					AddRow(2, time.Now(), time.Now(), nil,
						"Title 2", "Desc 2", "Body 2", 2, 0,
						2, "user2", "user2@test.com", "bio2", "image2")

				mock.ExpectQuery(`SELECT`).
					WithArgs(1, 2, 3).
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Empty Result",
			userIDs: []uint{999},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery(`SELECT`).
					WithArgs(999).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Database Error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "Empty UserIDs Array",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery(`SELECT`).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Pagination Test",
			userIDs: []uint{1},
			limit:   5,
			offset:  10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				}).
					AddRow(11, time.Now(), time.Now(), nil,
						"Title 11", "Desc 11", "Body 11", 1, 0)
				mock.ExpectQuery(`SELECT`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mockSetup(mock)

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("GetFeedArticles() got %d articles, want %d", len(got), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetTags_45f5cdc4bb
ROOST_METHOD_SIG_HASH=ArticleStore_GetTags_fb0aefcdd2

FUNCTION_DEF=func (s *ArticleStore) GetTags() ([ // GetTags creates a article tag
]model.Tag, error) 

*/
func TestArticleStoreGetTags(t *testing.T) {

	type testCase struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedTags  []model.Tag
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully retrieve tags",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, time.Now(), time.Now(), nil, "golang").
					AddRow(2, time.Now(), time.Now(), nil, "testing")
				mock.ExpectQuery("SELECT \\* FROM `tags`").
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
				mock.ExpectQuery("SELECT \\* FROM `tags`").
					WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `tags`").
					WillReturnError(sql.ErrConnDone)
			},
			expectedTags:  []model.Tag{},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

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

			tags, err := store.GetTags()

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if (tc.expectedError == nil && err != nil) || (tc.expectedError != nil && err == nil) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			if len(tags) != len(tc.expectedTags) {
				t.Errorf("Expected %d tags, got %d", len(tc.expectedTags), len(tags))
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}

	t.Run("Concurrent access", func(t *testing.T) {

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

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, time.Now(), time.Now(), nil, "golang")
		mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
		mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
		mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)

		var wg sync.WaitGroup
		concurrentCalls := 3

		for i := 0; i < concurrentCalls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := store.GetTags()
				if err != nil {
					t.Errorf("Concurrent call failed: %v", err)
				}
			}()
		}

		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations in concurrent test: %v", err)
		}
	})
}


/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user


*/
func TestArticleStoreIsFavorited(t *testing.T) {

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
		name     string
		article  *model.Article
		user     *model.User
		mockFunc func()
		want     bool
		wantErr  bool
	}{
		{
			name: "Valid article and user with favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid article and user without favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
		},
		{
			name:     "Nil article",
			article:  nil,
			user:     &model.User{Model: gorm.Model{ID: 1}},
			mockFunc: func() {},
			want:     false,
			wantErr:  false,
		},
		{
			name:     "Nil user",
			article:  &model.Article{Model: gorm.Model{ID: 1}},
			user:     nil,
			mockFunc: func() {},
			want:     false,
			wantErr:  false,
		},
		{
			name: "Database error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    false,
			wantErr: true,
		},
		{
			name:     "Both nil parameters",
			article:  nil,
			user:     nil,
			mockFunc: func() {},
			want:     false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			t.Logf("Running test case: %s", tt.name)

			tt.mockFunc()

			got, err := store.IsFavorited(tt.article, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case completed successfully: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_Update_3cddacb803
ROOST_METHOD_SIG_HASH=ArticleStore_Update_e245edd177

FUNCTION_DEF=func (s *ArticleStore) Update(m *model.Article) error // Update updates an article


*/
func TestArticleStoreUpdate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		article *model.Article
		mock    func()
		wantErr bool
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
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Invalid Data",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title: "",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Update with Related Entities",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title:       "Title with Relations",
				Description: "Description",
				Body:        "Body",
				UserID:      1,
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "tag1"},
				},
				Comments: []model.Comment{
					{Model: gorm.Model{ID: 1}, Body: "comment1", UserID: 1},
				},
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			t.Log("Starting test case:", tt.name)

			tt.mock()

			err := store.Update(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Log("Test case completed successfully:", tt.name)
		})
	}

	t.Run("Concurrent Updates", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		article := &model.Article{
			Model: gorm.Model{ID: 1},
			Title: "Concurrent Test",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `articles`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		done := make(chan bool)

		go func() {
			err := store.Update(article)
			if err != nil {
				t.Errorf("Concurrent update failed: %v", err)
			}
			done <- true
		}()

		<-done

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations in concurrent test: %v", err)
		}
	})
}


/*
ROOST_METHOD_HASH=NewArticleStore_85784abca5
ROOST_METHOD_SIG_HASH=NewArticleStore_436ae9c986

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore // NewArticleStore returns a new ArticleStore


*/
func TestNewArticleStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		validate func(*testing.T, *ArticleStore)
	}

	mockDB, _, err := gosqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open("mysql", mockDB)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	defer gormDB.Close()

	tests := []testCase{
		{
			name: "Successful creation with valid DB",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				if store == nil {
					t.Error("Expected non-nil ArticleStore")
					return
				}
				if store.db != gormDB {
					t.Error("DB reference not properly set")
				}
			},
		},
		{
			name: "Creation with nil DB",
			db:   nil,
			validate: func(t *testing.T, store *ArticleStore) {
				if store == nil {
					t.Error("Expected non-nil ArticleStore even with nil DB")
					return
				}
				if store.db != nil {
					t.Error("Expected nil DB in ArticleStore")
				}
			},
		},
		{
			name: "Verify instance independence",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {

				store2 := NewArticleStore(gormDB)
				if store == store2 {
					t.Error("Expected different instances of ArticleStore")
				}
				if store.db != store2.db {
					t.Error("Expected same DB reference in both instances")
				}
			},
		},
		{
			name: "Verify DB connection persistence",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				if store.db != gormDB {
					t.Error("DB connection not persisted correctly")
				}
			},
		},
		{
			name: "Memory address uniqueness",
			db:   gormDB,
			validate: func(t *testing.T, store *ArticleStore) {
				store2 := NewArticleStore(gormDB)
				store3 := NewArticleStore(gormDB)

				if &store == &store2 || &store == &store3 || &store2 == &store3 {
					t.Error("Expected unique memory addresses for each instance")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			t.Log("Starting test:", tc.name)

			store := NewArticleStore(tc.db)

			tc.validate(t, store)

			t.Log("Successfully completed test:", tc.name)
		})
	}
}

