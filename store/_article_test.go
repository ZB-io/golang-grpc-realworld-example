package store

import (
	sql "database/sql"
	debug "runtime/debug"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	time "time"
	assert "github.com/stretchr/testify/assert"
	errors "errors"
	sync "sync"
)





type mockDB struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}


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
		name    string
		article *model.Article
		user    *model.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Success - Add favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error - Association fails",
			article: &model.Article{
				Model:          gorm.Model{ID: 2},
				Title:          "Test Article 2",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "testuser2",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Error - Update fails",
			article: &model.Article{
				Model:          gorm.Model{ID: 3},
				Title:          "Test Article 3",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "testuser3",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "Error - Nil article",
			article: nil,
			user: &model.User{
				Model:    gorm.Model{ID: 4},
				Username: "testuser4",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Error - Nil user",
			article: &model.Article{
				Model:          gorm.Model{ID: 5},
				Title:          "Test Article 5",
				FavoritesCount: 0,
			},
			user:    nil,
			mock:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered: %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mock()

			err := store.AddFavorite(tt.article, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddFavorite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if !tt.wantErr && tt.article != nil {
				if tt.article.FavoritesCount != 1 {
					t.Errorf("FavoritesCount not incremented, got: %d, want: %d",
						tt.article.FavoritesCount, 1)
				}
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

	type testCase struct {
		name    string
		article *model.Article
		dbSetup func(mock sqlmock.Sqlmock)
		wantErr bool
	}

	tests := []testCase{
		{
			name: "Successfully create article",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test-tag"}},
			},
			dbSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Create article with missing required fields",
			article: &model.Article{

				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Create article with database connection error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Create article with maximum field lengths",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 5000)),
				UserID:      1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
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

			if tt.dbSetup != nil {
				tt.dbSetup(mock)
			}

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Create(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if !tt.wantErr && err == nil {
				t.Logf("Successfully created article: %v", tt.article.Title)
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

	type testCase struct {
		name          string
		comment       *model.Comment
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError bool
	}

	tests := []testCase{
		{
			name: "Successfully create comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test comment", uint(1), uint(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Fail to create comment - missing required fields",
			comment: &model.Comment{
				Body: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name: "Fail to create comment - non-existent article reference",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 999,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name: "Fail to create comment - database connection error",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
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
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.CreateComment(tc.comment)

			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case '%s' failed with error: %v", tc.name, err)
			} else {
				t.Logf("Test case '%s' passed successfully", tc.name)
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
		name    string
		article *model.Article
		dbSetup func(mock sqlmock.Sqlmock)
		wantErr bool
	}

	tests := []testCase{
		{
			name: "Successfully delete existing article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "Delete with nil article",
			article: nil,
			dbSetup: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
		},
		{
			name: "Delete with DB connection error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Delete article with associated records",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Comments: []model.Comment{
					{Model: gorm.Model{ID: 1}},
				},
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}},
				},
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
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

			if tt.dbSetup != nil {
				tt.dbSetup(mock)
			}

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Delete(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case '%s' completed with expected error: %v", tt.name, err)
			} else {
				t.Logf("Test case '%s' completed successfully", tt.name)
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
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Attempt to delete non-existent comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(999).
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
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Delete comment with null values in non-required fields",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 2,
				},
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(2).
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

			tt.mockFn(mock)

			err := store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "Successful unfavorite",
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
			name: "Association deletion fails",
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
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedErr:   sql.ErrConnDone,
			expectedCount: 2,
		},
		{
			name: "Update favorites count fails",
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
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedErr:   sql.ErrConnDone,
			expectedCount: 2,
		},
		{
			name:    "Nil article",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedErr:   gorm.ErrRecordNotFound,
			expectedCount: 0,
		},
		{
			name: "Nil user",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user:          nil,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedErr:   gorm.ErrRecordNotFound,
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

			initialCount := tc.article.FavoritesCount
			err = store.DeleteFavorite(tc.article, tc.user)

			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			if tc.article != nil && tc.article.FavoritesCount != tc.expectedCount {
				t.Errorf("Expected favorites count %d, got %d", tc.expectedCount, tc.article.FavoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test case '%s' failed with error: %v", tc.name, err)
			} else {
				t.Logf("Test case '%s' passed. FavoritesCount changed from %d to %d",
					tc.name, initialCount, tc.article.FavoritesCount)
			}
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

	testArticle := model.Article{
		Model:       gorm.Model{ID: 1},
		Title:       "Test Article",
		Description: "Test Description",
		Body:        "Test Body",
		UserID:      testUser.ID,
		Author:      *testUser,
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
			name: "Get All Articles",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Test Description", "Test Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)
			},
			wantLen: 1,
		},
		{
			name:     "Get Articles By Username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Test Description", "Test Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` join users").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			wantLen: 1,
		},
		{
			name:    "Get Articles By Tag",
			tagName: "testtag",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Test Description", "Test Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` join article_tags").
					WithArgs("testtag").
					WillReturnRows(rows)
			},
			wantLen: 1,
		},
		{
			name:        "Get Favorited Articles",
			favoritedBy: testUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favoriteRows := sqlmock.NewRows([]string{"article_id"}).AddRow(1)
				mock.ExpectQuery("^SELECT article_id FROM `favorite_articles`").
					WithArgs(testUser.ID).
					WillReturnRows(favoriteRows)

				articleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Test Description", "Test Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(articleRows)
			},
			wantLen: 1,
		},
		{
			name: "Database Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrConnDone)
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
				t.Errorf("ArticleStore.GetArticles() returned %d articles, want %d", len(got), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
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
		Tags: []model.Tag{
			{
				Model: gorm.Model{ID: 1},
				Name:  "test-tag",
			},
		},
		Author: model.User{
			Model:    gorm.Model{ID: 1},
			Username: "testuser",
			Email:    "test@example.com",
		},
	}

	tests := []testCase{
		{
			name: "Successful Article Retrieval",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "created_at", "updated_at", "deleted_at",
						"title", "description", "body", "user_id", "favorites_count",
					}).AddRow(
						validArticle.ID, validArticle.CreatedAt, validArticle.UpdatedAt, nil,
						validArticle.Title, validArticle.Description, validArticle.Body,
						validArticle.UserID, validArticle.FavoritesCount,
					))

				mock.ExpectQuery(`SELECT \* FROM "tags"`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name",
					}).AddRow(1, "test-tag"))

				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "username", "email",
					}).AddRow(1, "testuser", "test@example.com"))
			},
			expectedData:  validArticle,
			expectedError: nil,
		},
		{
			name: "Non-existent Article",
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
			name: "Database Error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedData:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Zero ID Input",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
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

			article, err := store.GetByID(tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, article)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, tc.expectedData.ID, article.ID)
				assert.Equal(t, tc.expectedData.Title, article.Title)
				assert.Equal(t, tc.expectedData.Description, article.Description)
				assert.Equal(t, tc.expectedData.Body, article.Body)
				t.Log("Article retrieved successfully")
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
		name          string
		article       *model.Article
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
	}{
		{
			name: "Successfully retrieve comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test comment 1", 1, 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Test comment 2", 2, 1, 2)

				mock.ExpectQuery("SELECT").WillReturnRows(rows)
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
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
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
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectError:   true,
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

			comments, err := store.GetComments(tt.article)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(comments) != tt.expectedCount {
				t.Errorf("Expected %d comments, got %d", tt.expectedCount, len(comments))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
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
			name:    "Empty result for non-existent users",
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
			name:    "Empty user IDs array",
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
			name:    "Pagination limit enforcement",
			userIDs: []uint{1},
			limit:   5,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
					"author_id", "author_username", "author_email", "author_bio", "author_image",
				})
				for i := 1; i <= 5; i++ {
					rows.AddRow(i, time.Now(), time.Now(), nil,
						"Test Article", "Description", "Body", 1, 0,
						1, "user1", "user1@test.com", "bio1", "image1")
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			wantLen: 5,
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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("ArticleStore.GetFeedArticles() got %v articles, want %v", len(got), tt.wantLen)
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
			name: "Database connection error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnError(sql.ErrConnDone)
			},
			expectedTags:  []model.Tag{},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Database timeout error",
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

			t.Log("Executing GetTags for scenario:", tc.name)
			tags, err := store.GetTags()

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedTags), len(tags))
				if len(tc.expectedTags) > 0 {
					for i, expectedTag := range tc.expectedTags {
						assert.Equal(t, expectedTag.ID, tags[i].ID)
						assert.Equal(t, expectedTag.Name, tags[i].Name)
					}
				}
				t.Log("Tags retrieved successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user


*/
func TestArticleStoreIsFavorited(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedBool  bool
		expectedError error
	}

	tests := []testCase{
		{
			name: "Valid Article and User with Existing Favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedBool:  true,
			expectedError: nil,
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
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedBool:  false,
			expectedError: nil,
		},
		{
			name:          "Nil Article Parameter",
			article:       nil,
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedBool:  false,
			expectedError: nil,
		},
		{
			name:          "Nil User Parameter",
			article:       &model.Article{Model: gorm.Model{ID: 1}},
			user:          nil,
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedBool:  false,
			expectedError: nil,
		},
		{
			name: "Database Error Condition",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedBool:  false,
			expectedError: sql.ErrConnDone,
		},
		{
			name:          "Both Nil Parameters",
			article:       nil,
			user:          nil,
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedBool:  false,
			expectedError: nil,
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
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(0, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedBool:  false,
			expectedError: nil,
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
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 0).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedBool:  false,
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

			tc.mockSetup(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			result, err := store.IsFavorited(tc.article, tc.user)

			if result != tc.expectedBool {
				t.Errorf("Expected bool %v, got %v", tc.expectedBool, result)
			}

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
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
		name      string
		article   *model.Article
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
		errMsg    string
	}

	tests := []testCase{
		{
			name: "Successful Update",
			article: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Updated Title",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:      "Invalid Article (nil)",
			article:   nil,
			setupMock: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
			errMsg:    "invalid article model",
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Title",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database error",
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

			tc.setupMock(mock)

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Update(tc.article)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tc.errMsg != "" && err.Error() != tc.errMsg {
					t.Errorf("Expected error message %q but got %q", tc.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
		})
	}

	t.Run("Concurrent Updates", func(t *testing.T) {

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

		const numGoroutines = 5
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE `articles`").
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				article := &model.Article{
					Model:       gorm.Model{ID: 1},
					Title:       "Concurrent Update",
					Description: "Test concurrent update",
					Body:        "Test body",
					UserID:      uint(i),
				}
				if err := store.Update(article); err != nil {
					t.Errorf("Concurrent update %d failed: %v", i, err)
				}
			}(i)
		}

		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled mock expectations in concurrent test: %v", err)
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

	mockDB, _, err := sqlmock.New()
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
				assert.NotNil(t, store, "Store should not be nil")
				assert.Equal(t, gormDB, store.db, "DB reference should match")
			},
		},
		{
			name: "Creation with nil DB",
			db:   nil,
			validate: func(t *testing.T, store *ArticleStore) {
				assert.NotNil(t, store, "Store should not be nil even with nil DB")
				assert.Nil(t, store.db, "DB reference should be nil")
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

			if tc.validate != nil {
				tc.validate(t, store)
			}

			t.Log("Successfully completed test case:", tc.name)
		})
	}
}

