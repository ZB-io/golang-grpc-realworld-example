package store

import (
	fmt "fmt"
	debug "runtime/debug"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	sql "database/sql"
	time "time"
	errors "errors"
	math "math"
	assert "github.com/stretchr/testify/assert"
)








/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


*/
func TestArticleStoreAddFavorite(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("_", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
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
			expectError: false,
		},
		{
			name: "Association failure",
			article: &model.Article{
				Model:          gorm.Model{ID: 2},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 2},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(fmt.Errorf("association error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Update count failure",
			article: &model.Article{
				Model:          gorm.Model{ID: 3},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 3},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(3, 3).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(fmt.Errorf("update error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name:    "Nil article",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 4},
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
		},
		{
			name: "Nil user",
			article: &model.Article{
				Model:          gorm.Model{ID: 5},
				FavoritesCount: 0,
			},
			user:        nil,
			setupMock:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
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

			tt.setupMock(mock)

			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			t.Logf("Test case completed: %s", tt.name)
		})
	}
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
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name    string
		article *model.Article
		mock    func()
		wantErr bool
	}{
		{
			name: "Success - Create article with all required fields",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mock: func() {
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
				Title: "Test Article",

				UserID: 1,
			},
			mock: func() {
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
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "tag1"},
					{Model: gorm.Model{ID: 2}, Name: "tag2"},
				},
			},
			mock: func() {
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
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Success - Create article with author relationship",
			article: &model.Article{
				Title:       "Test Article with Author",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
				},
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
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

			err := store.Create(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err == nil {
				t.Log("Successfully created article")
			} else {
				t.Log("Failed to create article:", err)
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
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success - Create Valid Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
				Model: gorm.Model{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
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
				Body: "",
				Model: gorm.Model{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
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
			name: "Failure - Non-existent User Reference",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 1,
				Model: gorm.Model{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
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
			name: "Failure - Database Connection Issue",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
				Model: gorm.Model{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
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

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mockFn(mock)

			err := store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case '%s' failed with error: %v", tt.name, err)
			} else {
				t.Logf("Test case '%s' passed successfully", tt.name)
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
		name          string
		article       *model.Article
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
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
			expectedError: nil,
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
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Delete with Database Error",
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
			expectedError: sql.ErrConnDone,
		},
		{
			name:          "Delete with Nil Article",
			article:       nil,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article: nil pointer"),
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

			if tc.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tc.expectedError != nil && err == nil {
				t.Errorf("Expected error %v, got nil", tc.expectedError)
			} else if tc.expectedError != nil && err != nil && tc.expectedError.Error() != err.Error() {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_DeleteComment_effbcb38aa
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteComment_d3c99623e4

FUNCTION_DEF=func (s *ArticleStore) DeleteComment(m *model.Comment) error // DeleteComment deletes an comment


*/
func TestArticleStoreDeleteComment(t *testing.T) {

	type testCase struct {
		name          string
		comment       *model.Comment
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError bool
	}

	tests := []testCase{
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
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Attempt to delete non-existent comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name: "Database connection error",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: true,
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

			err = store.DeleteComment(tc.comment)

			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
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
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully unfavorite article",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
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
			name: "Association deletion error",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("association deletion error"))
				mock.ExpectRollback()
			},
			expectedError: fmt.Errorf("association deletion error"),
		},
		{
			name: "Update count error",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("update count error"))
				mock.ExpectRollback()
			},
			expectedError: fmt.Errorf("update count error"),
		},
		{
			name:    "Nil article parameter",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedError: fmt.Errorf("invalid article"),
		},
		{
			name: "Nil user parameter",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user:          nil,
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedError: fmt.Errorf("invalid user"),
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

			tc.mockSetup(mock)

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
				expectedCount := tc.article.FavoritesCount - 1
				if tc.article.FavoritesCount != expectedCount {
					t.Errorf("Expected favorites count: %d, got: %d", expectedCount, tc.article.FavoritesCount)
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

	mockUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
	}

	mockArticles := []model.Article{
		{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Title:       "Test Article 1",
			Description: "Test Description 1",
			Body:        "Test Body 1",
			UserID:      1,
			Author:      *mockUser,
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
		wantErr     bool
		wantLen     int
	}{
		{
			name:   "Success - No Filters",
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article 1", "Test Description 1", "Test Body 1", 1)

				mock.ExpectQuery("SELECT").WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:     "Success - Filter by Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article 1", "Test Description 1", "Test Body 1", 1)

				mock.ExpectQuery("SELECT").WithArgs("testuser").WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "Success - Filter by Tag",
			tagName: "programming",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article 1", "Test Description 1", "Test Body 1", 1)

				mock.ExpectQuery("SELECT").WithArgs("programming").WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:        "Success - Filter by Favorited",
			favoritedBy: mockUser,
			limit:       10,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favRows := sqlmock.NewRows([]string{"article_id"}).AddRow(1)
				mock.ExpectQuery("SELECT article_id FROM favorite_articles").WithArgs(mockUser.ID).WillReturnRows(favRows)

				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article 1", "Test Description 1", "Test Body 1", 1)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:   "Error - Database Error",
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
			wantLen: 0,
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

			if len(got) != tt.wantLen {
				t.Errorf("ArticleStore.GetArticles() got %v articles, want %v", len(got), tt.wantLen)
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
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
		expectedData  *model.Article
	}

	validArticle := &model.Article{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:          "Test Article",
		Description:    "Test Description",
		Body:           "Test Body",
		UserID:         1,
		FavoritesCount: 0,
		Tags:           []model.Tag{{Model: gorm.Model{ID: 1}, Name: "test-tag"}},
		Author:         model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
	}

	tests := []testCase{
		{
			name: "Successful retrieval",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).
						AddRow(1, "Test Article", "Test Description", "Test Body"))

				mock.ExpectQuery(`SELECT \* FROM "tags"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "test-tag"))

				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(1, "testuser"))
			},
			expectedData:  validArticle,
			expectedError: nil,
		},
		{
			name: "Non-existent article",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedData:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expectedData:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Zero ID input",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedData:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Maximum ID value",
			id:   math.MaxUint32,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(math.MaxUint32).
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

			t.Logf("Executing test case: %s", tc.name)
			result, err := store.GetByID(tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedData.ID, result.ID)
				assert.Equal(t, tc.expectedData.Title, result.Title)
				t.Logf("Successfully retrieved article with ID: %d", result.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
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
		name         string
		article      *model.Article
		mockSetup    func(sqlmock.Sqlmock)
		wantComments []model.Comment
		wantErr      bool
	}{
		{
			name: "Successfully retrieve comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "user_id", "article_id",
					"author.id", "author.username", "author.email",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test comment", 1, 1,
					1, "testuser", "test@example.com",
				)

				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			wantComments: []model.Comment{
				{
					Model:     gorm.Model{ID: 1},
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
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
			name: "Article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "user_id", "article_id",
					"author.id", "author.username", "author.email",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			wantComments: []model.Comment{},
			wantErr:      false,
		},
		{
			name: "Database error",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			wantComments: nil,
			wantErr:      true,
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

			tt.mockSetup(mock)

			gotComments, err := store.GetComments(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(gotComments) != len(tt.wantComments) {
					t.Errorf("ArticleStore.GetComments() returned %d comments, want %d",
						len(gotComments), len(tt.wantComments))
					return
				}

				if len(tt.wantComments) > 0 {
					for i, want := range tt.wantComments {
						got := gotComments[i]
						if got.ID != want.ID || got.Body != want.Body ||
							got.UserID != want.UserID || got.ArticleID != want.ArticleID {
							t.Errorf("Comment mismatch at index %d", i)
						}
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Log("Test case completed successfully:", tt.name)
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
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

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
			name:    "Successful retrieval of articles",
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
						"Test Article 1", "Description 1", "Body 1", 1, 0,
						1, "user1", "user1@test.com", "bio1", "image1").
					AddRow(2, time.Now(), time.Now(), nil,
						"Test Article 2", "Description 2", "Body 2", 2, 0,
						2, "user2", "user2@test.com", "bio2", "image2")

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1, 2, 3).
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{4, 5},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(4, 5).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
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

			t.Logf("Running test case: %s", tt.name)

			tt.mockSetup(mock)

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("ArticleStore.GetFeedArticles() returned %d articles, want %d", len(got), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}

			t.Logf("Test case completed successfully: %s", tt.name)
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
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "name"}
				mock.ExpectQuery("SELECT (.+) FROM `tags`").
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "golang").
						AddRow(2, time.Now(), time.Now(), nil, "testing"))
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
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "name"}
				mock.ExpectQuery("SELECT (.+) FROM `tags`").
					WillReturnRows(sqlmock.NewRows(columns))
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `tags`").
					WillReturnError(sql.ErrConnDone)
			},
			expectedTags:  []model.Tag{},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Large dataset handling",
			setupMock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "name"}
				rows := sqlmock.NewRows(columns)
				for i := 1; i <= 1000; i++ {
					rows.AddRow(i, time.Now(), time.Now(), nil, fmt.Sprintf("tag-%d", i))
				}
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: func() []model.Tag {
				tags := make([]model.Tag, 1000)
				for i := range tags {
					tags[i] = model.Tag{
						Model: gorm.Model{ID: uint(i + 1)},
						Name:  fmt.Sprintf("tag-%d", i+1),
					}
				}
				return tags
			}(),
			expectedError: nil,
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
}


/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user


*/
func TestArticleStoreIsFavorited(t *testing.T) {

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
		name        string
		article     *model.Article
		user        *model.User
		mockSetup   func(sqlmock.Sqlmock)
		wantFav     bool
		wantErr     error
		description string
	}{
		{
			name:    "Nil Article and User",
			article: nil,
			user:    nil,
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			wantFav:     false,
			wantErr:     nil,
			description: "Should handle nil inputs gracefully",
		},
		{
			name: "Article Not Favorited",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			wantFav:     false,
			wantErr:     nil,
			description: "Should return false when article is not favorited",
		},
		{
			name: "Article Favorited",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			wantFav:     true,
			wantErr:     nil,
			description: "Should return true when article is favorited",
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			wantFav:     false,
			wantErr:     sql.ErrConnDone,
			description: "Should handle database errors properly",
		},
		{
			name: "Multiple Favorites",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			wantFav:     true,
			wantErr:     nil,
			description: "Should handle multiple favorite entries",
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

			t.Log("Test Description:", tt.description)

			tt.mockSetup(mock)

			gotFav, gotErr := store.IsFavorited(tt.article, tt.user)

			if (gotErr != nil && tt.wantErr == nil) || (gotErr == nil && tt.wantErr != nil) {
				t.Errorf("IsFavorited() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if gotErr != nil && tt.wantErr != nil && !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("IsFavorited() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if gotFav != tt.wantFav {
				t.Errorf("IsFavorited() = %v, want %v", gotFav, tt.wantFav)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_Update_3cddacb803
ROOST_METHOD_SIG_HASH=ArticleStore_Update_e245edd177

FUNCTION_DEF=func (s *ArticleStore) Update(m *model.Article) error // Update updates an article


*/
func TestArticleStoreUpdate(t *testing.T) {

	type testCase struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}

	tests := []testCase{
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
						int32(0),
						uint(1),
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Update Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
				Title:       "Non-existent",
				Description: "Description",
				Body:        "Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Update with Invalid Data",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title:       "",
				Description: "Description",
				Body:        "Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
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

			err = store.Update(tc.article)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=NewArticleStore_85784abca5
ROOST_METHOD_SIG_HASH=NewArticleStore_436ae9c986

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore // NewArticleStore returns a new ArticleStore


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
			db:       &gorm.DB{},
			wantNil:  false,
			scenario: "Create ArticleStore with valid DB connection",
		},
		{
			name:     "Handle nil DB connection",
			db:       nil,
			wantNil:  false,
			scenario: "Create ArticleStore with nil DB connection",
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

			t.Logf("Testing scenario: %s", tt.scenario)

			var mockDB *gorm.DB
			if tt.db != nil {

				db, _, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock DB: %v", err)
				}
				mockDB = &gorm.DB{
					Value: db,
				}
			}

			got := NewArticleStore(mockDB)

			if got == nil && !tt.wantNil {
				t.Errorf("NewArticleStore() returned nil, want non-nil")
			}

			if got != nil {

				if got.db != mockDB {
					t.Errorf("NewArticleStore().db = %v, want %v", got.db, mockDB)
				}
				t.Logf("Successfully created ArticleStore with expected DB reference")
			}

			if tt.db != nil {
				store1 := NewArticleStore(mockDB)
				store2 := NewArticleStore(mockDB)

				if store1 == store2 {
					t.Error("Different ArticleStore instances should not be equal")
				}
				t.Log("Verified instance independence")
			}
		})
	}
}

