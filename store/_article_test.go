package store

import (
	fmt "fmt"
	debug "runtime/debug"
	sync "sync"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	sql "database/sql"
	errors "errors"
	time "time"
	math "math"
	assert "github.com/stretchr/testify/assert"
)





type mockDB struct {
	db   *gorm.DB
	mock sqlmock.Sqlmock
}


/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


*/
func TestArticleStoreAddFavorite(t *testing.T) {
	type testCase struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		concurrent  bool
	}

	baseArticle := &model.Article{
		Model: gorm.Model{ID: 1},
		Title: "Test Article",
	}
	baseUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
	}

	tests := []testCase{
		{
			name:    "Successful favorite addition",
			article: baseArticle,
			user:    baseUser,
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
			name:    "Nil article parameter",
			article: nil,
			user:    baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
		},
		{
			name:    "Nil user parameter",
			article: baseArticle,
			user:    nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
		},
		{
			name:    "Database association error",
			article: baseArticle,
			user:    baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(fmt.Errorf("association error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name:    "Database update error",
			article: baseArticle,
			user:    baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(fmt.Errorf("update error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name:       "Concurrent favorite operations",
			article:    baseArticle,
			user:       baseUser,
			concurrent: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				for i := 0; i < 3; i++ {
					mock.ExpectBegin()
					mock.ExpectExec("INSERT INTO `favorite_articles`").
						WithArgs(1, 1).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec("UPDATE `articles`").
						WithArgs(1, 1).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered: %v\n%s", r, string(debug.Stack()))
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

			if tc.concurrent {
				var wg sync.WaitGroup
				errChan := make(chan error, 3)

				for i := 0; i < 3; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						if err := store.AddFavorite(tc.article, tc.user); err != nil {
							errChan <- err
						}
					}()
				}

				wg.Wait()
				close(errChan)

				for err := range errChan {
					if !tc.expectError {
						t.Errorf("Unexpected error in concurrent operation: %v", err)
					}
				}
			} else {
				err = store.AddFavorite(tc.article, tc.user)

				if tc.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tc.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
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
			name: "Success - Create Valid Article",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
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
			name: "Failure - Missing Required Fields",
			article: &model.Article{

				Body:   "Test Body",
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
			name: "Success - Article with Tags",
			article: &model.Article{
				Title:       "Test Article with Tags",
				Description: "Test Description",
				Body:        "Test Body Content",
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
			name: "Failure - Database Connection Error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
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
			name: "Failure - Duplicate Entry",
			article: &model.Article{
				Title:       "Duplicate Title",
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Success - Maximum Field Lengths",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 65535)),
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

			tt.mock()

			err := store.Create(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case failed as expected: %v", err)
			} else {
				t.Logf("Test case succeeded")
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
		t.Fatalf("Failed to open GORM connection: %v", err)
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
				Author: model.User{
					Model: gorm.Model{ID: 1},
				},
				Article: model.Article{
					Model: gorm.Model{ID: 1},
				},
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failure - Missing Required Fields",
			comment: &model.Comment{

				UserID:    1,
				ArticleID: 1,
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
			name: "Failure - Non-existent Article Reference",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
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
			name: "Failure - Non-existent User Reference",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 1,
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
			name: "Failure - Database Connection Issue",
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

			err := store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case completed with expected error: %v", err)
			} else {
				t.Log("Test case completed successfully")
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
			name: "Successfully delete existing article",
			article: &model.Article{
				Model: gorm.Model{
					ID:        1,
					DeletedAt: nil,
				},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
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
			name: "Attempt to delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
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
			name:    "Delete with nil article parameter",
			article: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid article: nil pointer"),
		},
		{
			name: "Delete with database connection error",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
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

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Delete(tc.article)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if tc.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			} else if tc.expectedError != nil && err == nil {
				t.Errorf("Expected error: %v, got nil", tc.expectedError)
			} else if tc.expectedError != nil && err != nil && tc.expectedError.Error() != err.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
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
		name    string
		comment *model.Comment
		mockSQL func()
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Fail to delete non-existent comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(sqlmock.AnyArg()).
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

			t.Log("Starting test case:", tt.name)

			tt.mockSQL()

			err := store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Log("Test case completed successfully")
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
		setupMock     func(sqlmock.Sqlmock)
		expectedErr   error
		expectedCount int32
	}

	tests := []testCase{
		{
			name: "Successfully unfavorite article",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
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
			expectedErr:   nil,
			expectedCount: 1,
		},
		{
			name: "Failed association deletion",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(errors.New("association deletion failed"))
				mock.ExpectRollback()
			},
			expectedErr:   errors.New("association deletion failed"),
			expectedCount: 2,
		},
		{
			name: "Failed favorites count update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
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
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			expectedErr:   errors.New("update failed"),
			expectedCount: 2,
		},
		{
			name:    "Null article reference",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedErr:   errors.New("invalid article reference"),
			expectedCount: 0,
		},
		{
			name: "Null user reference",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user:          nil,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedErr:   errors.New("invalid user reference"),
			expectedCount: 2,
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

			store := &ArticleStore{
				db: gormDB,
			}

			initialCount := int32(0)
			if tc.article != nil {
				initialCount = tc.article.FavoritesCount
			}

			err = store.DeleteFavorite(tc.article, tc.user)

			if (err != nil && tc.expectedErr == nil) || (err == nil && tc.expectedErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

			if tc.article != nil && tc.article.FavoritesCount != tc.expectedCount {
				t.Errorf("Expected favorites count: %d, got: %d", tc.expectedCount, tc.article.FavoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
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

	db, mock, err := gosqlmock.New()
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
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test.jpg",
	}

	testArticle := model.Article{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:          "Test Article",
		Description:    "Test Description",
		Body:           "Test Body",
		UserID:         testUser.ID,
		FavoritesCount: 0,
		Author:         *testUser,
	}

	testCases := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(gosqlmock.Sqlmock)
		wantErr     bool
		wantLen     int
	}{
		{
			name:   "Get All Articles",
			limit:  10,
			offset: 0,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				rows := gosqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(testArticle.ID, testArticle.CreatedAt, testArticle.UpdatedAt, nil, testArticle.Title, testArticle.Description, testArticle.Body, testArticle.UserID, testArticle.FavoritesCount)

				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := gosqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).
					AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Bio, testUser.Image)
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
		{
			name:     "Get Articles By Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				rows := gosqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(testArticle.ID, testArticle.CreatedAt, testArticle.UpdatedAt, nil, testArticle.Title, testArticle.Description, testArticle.Body, testArticle.UserID, testArticle.FavoritesCount)

				mock.ExpectQuery("^SELECT (.+) FROM `articles` JOIN users").
					WithArgs("testuser").
					WillReturnRows(rows)

				authorRows := gosqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).
					AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Bio, testUser.Image)
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
		{
			name:    "Get Articles By Tag",
			tagName: "programming",
			limit:   10,
			offset:  0,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				rows := gosqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(testArticle.ID, testArticle.CreatedAt, testArticle.UpdatedAt, nil, testArticle.Title, testArticle.Description, testArticle.Body, testArticle.UserID, testArticle.FavoritesCount)

				mock.ExpectQuery("^SELECT (.+) FROM `articles` JOIN article_tags").
					WithArgs("programming").
					WillReturnRows(rows)

				authorRows := gosqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).
					AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Bio, testUser.Image)
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
		{
			name:        "Get Favorited Articles",
			favoritedBy: testUser,
			limit:       10,
			offset:      0,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				favoriteRows := gosqlmock.NewRows([]string{"article_id"}).
					AddRow(testArticle.ID)
				mock.ExpectQuery("^SELECT article_id FROM `favorite_articles`").
					WithArgs(testUser.ID).
					WillReturnRows(favoriteRows)

				rows := gosqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(testArticle.ID, testArticle.CreatedAt, testArticle.UpdatedAt, nil, testArticle.Title, testArticle.Description, testArticle.Body, testArticle.UserID, testArticle.FavoritesCount)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := gosqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).
					AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Bio, testUser.Image)
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tc.mockSetup(mock)

			articles, err := store.GetArticles(tc.tagName, tc.username, tc.favoritedBy, tc.limit, tc.offset)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetArticles() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if len(articles) != tc.wantLen {
				t.Errorf("GetArticles() got %d articles, want %d", len(articles), tc.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetByID_6fe18728fc
ROOST_METHOD_SIG_HASH=ArticleStore_GetByID_bb488e542f

FUNCTION_DEF=func (s *ArticleStore) GetByID(id uint) (*model.Article, error) // GetByID finds an article from id


*/
func TestArticleStoreGetById(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mock    func(mock sqlmock.Sqlmock)
		want    *model.Article
		wantErr error
	}{
		{
			name: "Success - Article Found",
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Title", "Test Description", "Test Body", 1, 0)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)

				tagRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "tag1")
				mock.ExpectQuery("SELECT").WillReturnRows(tagRows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").WillReturnRows(authorRows)
			},
			want: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Title",
				Tags:  []model.Tag{{Model: gorm.Model{ID: 1}, Name: "tag1"}},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
				},
			},
			wantErr: nil,
		},
		{
			name: "Failure - Article Not Found",
			id:   999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Failure - Database Error",
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: sql.ErrConnDone,
		},
		{
			name: "Failure - Zero ID",
			id:   0,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Edge Case - Max ID",
			id:   math.MaxUint32,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
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

			store, mock := mockArticleStore(t)
			defer store.db.Close()

			tt.mock(mock)

			got, err := store.GetByID(tt.id)

			if err != tt.wantErr {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ArticleStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			if tt.want == nil && got != nil {
				t.Errorf("ArticleStore.GetByID() = %v, want nil", got)
			}

			if tt.want != nil {
				if got == nil {
					t.Error("ArticleStore.GetByID() = nil, want non-nil Article")
				} else {

					if got.ID != tt.want.ID {
						t.Errorf("ArticleStore.GetByID() ID = %v, want %v", got.ID, tt.want.ID)
					}
					if got.Title != tt.want.Title {
						t.Errorf("ArticleStore.GetByID() Title = %v, want %v", got.Title, tt.want.Title)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}

func mockArticleStore(t *testing.T) (*ArticleStore, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return &ArticleStore{db: gormDB}, mock
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
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test comment 1", 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Test comment 2", 2, 1)

				mock.ExpectQuery(`SELECT (.+) FROM "comments"`).
					WithArgs(1).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).
					AddRow(1, "user1", "user1@test.com", "bio1", "image1.jpg").
					AddRow(2, "user2", "user2@test.com", "bio2", "image2.jpg")

				mock.ExpectQuery(`SELECT (.+) FROM "users"`).
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
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"})
				mock.ExpectQuery(`SELECT (.+) FROM "comments"`).
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
				mock.ExpectQuery(`SELECT (.+) FROM "comments"`).
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

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			comments, err := store.GetComments(tc.article)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(comments) != tc.expectedCount {
				t.Errorf("Expected %d comments, got %d", tc.expectedCount, len(comments))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
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
			name:    "Successfully retrieve feed articles",
			userIDs: []uint{1, 2},
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

				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{999, 1000},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
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
				mock.ExpectQuery(`SELECT`).WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "Empty userIDs array",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Zero limit",
			userIDs: []uint{1},
			limit:   0,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			wantLen: 0,
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
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("ArticleStore.GetFeedArticles() returned %d articles, want %d", len(got), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
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
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
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
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnError(sql.ErrConnDone)
			},
			expectedTags:  []model.Tag{},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Database timeout",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnError(errors.New("context deadline exceeded"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("context deadline exceeded"),
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

			store := &ArticleStore{
				db: gormDB,
			}

			tags, err := store.GetTags()

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, tags, len(tc.expectedTags))

				for i := range tags {
					assert.Equal(t, tc.expectedTags[i].ID, tags[i].ID)
					assert.Equal(t, tc.expectedTags[i].Name, tags[i].Name)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
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
		wantErr     bool
		description string
	}{
		{
			name: "Valid Article and User with Existing Favorite",
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
			wantErr:     false,
			description: "Article is favorited by user",
		},
		{
			name: "Valid Article and User with No Favorite",
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
			wantErr:     false,
			description: "Article is not favorited by user",
		},
		{
			name:        "Nil Article Parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			wantFav:     false,
			wantErr:     false,
			description: "Nil article should return false without error",
		},
		{
			name:        "Nil User Parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			wantFav:     false,
			wantErr:     false,
			description: "Nil user should return false without error",
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
			wantErr:     true,
			description: "Database error should be handled properly",
		},
		{
			name:        "Both Parameters Nil",
			article:     nil,
			user:        nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			wantFav:     false,
			wantErr:     false,
			description: "Both nil parameters should return false without error",
		},
		{
			name: "Invalid Article ID",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(0, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			wantFav:     false,
			wantErr:     false,
			description: "Invalid article ID should return false",
		},
		{
			name: "Invalid User ID",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 0},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `favorite_articles`").
					WithArgs(1, 0).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			wantFav:     false,
			wantErr:     false,
			description: "Invalid user ID should return false",
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

			t.Log(tt.description)

			tt.mockSetup(mock)

			gotFav, err := store.IsFavorited(tt.article, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFav != tt.wantFav {
				t.Errorf("IsFavorited() = %v, want %v", gotFav, tt.wantFav)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_Update_3cddacb803
ROOST_METHOD_SIG_HASH=ArticleStore_Update_e245edd177

FUNCTION_DEF=func (s *ArticleStore) Update(m *model.Article) error // Update updates an article


*/
func TestArticleStoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful Update",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Updated Title",
				Body:  "Updated Body",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Updated Title", "Updated Body", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "Nil Article",
			article: nil,
			mockFn: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Title",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty Required Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("validation error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Large Content Update",
			article: &model.Article{
				Model:     gorm.Model{ID: 1},
				Title:     "Title",
				Body:      string(make([]byte, 10000)),
				UpdatedAt: time.Now(),
				CreatedAt: time.Now(),
				DeletedAt: nil,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
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

			mockDB := newMockDB(t)
			defer mockDB.db.Close()

			tt.mockFn(mockDB.mock)

			store := &ArticleStore{
				db: mockDB.db,
			}

			err := store.Update(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mockDB.mock.ExpectationsWereMet(); err != nil {
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

func newMockDB(t *testing.T) *mockDB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	return &mockDB{
		db:   gormDB,
		mock: mock,
	}
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
			name: "Handle nil DB connection",
			db:   nil,
			validate: func(t *testing.T, store *ArticleStore) {
				if store == nil {
					t.Error("Expected non-nil ArticleStore even with nil DB")
					return
				}
				if store.db != nil {
					t.Error("Expected nil DB in store")
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

			t.Log("Starting test case:", tc.name)

			store := NewArticleStore(tc.db)

			tc.validate(t, store)

			t.Log("Successfully completed test case:", tc.name)
		})
	}
}

func TestNewArticleStoreIndependence(t *testing.T) {

	mockDB1, _, err := gosqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create first mock DB: %v", err)
	}
	defer mockDB1.Close()

	mockDB2, _, err := gosqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create second mock DB: %v", err)
	}
	defer mockDB2.Close()

	gormDB1, err := gorm.Open("mysql", mockDB1)
	if err != nil {
		t.Fatalf("Failed to create first GORM DB: %v", err)
	}
	defer gormDB1.Close()

	gormDB2, err := gorm.Open("mysql", mockDB2)
	if err != nil {
		t.Fatalf("Failed to create second GORM DB: %v", err)
	}
	defer gormDB2.Close()

	t.Run("Verify instance independence", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		store1 := NewArticleStore(gormDB1)
		store2 := NewArticleStore(gormDB2)

		if store1 == store2 {
			t.Error("Store instances should be different")
		}

		if store1.db == store2.db {
			t.Error("Store instances should have different DB connections")
		}
	})
}

