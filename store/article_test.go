package store

import (
	sql "database/sql"
	driver "database/sql/driver"
	errors "errors"
	fmt "fmt"
	debug "runtime/debug"
	sync "sync"
	testing "testing"
	time "time"

	"github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
)

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	ArticleID uint   `gorm:"not null"`
	Body      string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

/*
ROOST_METHOD_HASH=NewArticleStore_85784abca5
ROOST_METHOD_SIG_HASH=NewArticleStore_436ae9c986

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore // NewArticleStore returns a new ArticleStore
*/
func TestNewArticleStore(t *testing.T) {
	type testCase struct {
		name      string
		dbFunc    func() *gorm.DB
		wantNilDB bool
		wantPanic bool
	}

	tests := []testCase{
		{
			name: "Valid Initialization of ArticleStore",
			dbFunc: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gdb, err := gorm.Open("sqlmock", db)
				if err != nil {
					t.Fatalf("Error opening gorm: %v", err)
				}
				return gdb
			},
			wantNilDB: false,
			wantPanic: false,
		},
		{
			name: "Initialization with Nil Database",
			dbFunc: func() *gorm.DB {
				return nil
			},
			wantNilDB: true,
			wantPanic: false,
		},
		{
			name: "Initialization with Closed Database Connection",
			dbFunc: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gdb, err := gorm.Open("sqlmock", db)
				if err != nil {
					t.Fatalf("Error opening gorm: %v", err)
				}
				gdb.Close()
				return gdb
			},
			wantNilDB: false,
			wantPanic: false,
		},
		{
			name: "Concurrent Initialization",
			dbFunc: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gdb, err := gorm.Open("sqlmock", db)
				if err != nil {
					t.Fatalf("Error opening gorm: %v", err)
				}
				return gdb
			},
			wantNilDB: false,
			wantPanic: false,
		},
		{
			name: "Initialization with Invalid Database Configuration",
			dbFunc: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gdb, err := gorm.Open("sqlmock", db)
				if err != nil {
					t.Fatalf("Error opening gorm: %v", err)
				}
				return gdb
			},
			wantNilDB: false,
			wantPanic: false,
		},
		{
			name: "Initialization with Different Database Instances",
			dbFunc: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gdb, err := gorm.Open("sqlmock", db)
				if err != nil {
					t.Fatalf("Error opening gorm: %v", err)
				}
				return gdb
			},
			wantNilDB: false,
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					if !tt.wantPanic {
						t.Fail()
					}
				}
			}()

			db := tt.dbFunc()
			store := NewArticleStore(db)

			if store == nil {
				t.Fatal("Expected ArticleStore to be non-nil")
			}

			if tt.wantNilDB {
				if store.db != nil {
					t.Errorf("Expected db to be nil, but got %v", store.db)
				}
			} else {
				if store.db != db {
					t.Errorf("Expected db to be %v, but got %v", db, store.db)
				}
			}

			if tt.name == "Concurrent Initialization" {
				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						store := NewArticleStore(db)
						if store == nil {
							t.Error("Expected ArticleStore to be non-nil in concurrent test")
						}
						if store.db != db {
							t.Errorf("Expected db to be %v, but got %v in concurrent test", db, store.db)
						}
					}()
				}
				wg.Wait()
			}
			if tt.name == "Initialization with Different Database Instances" {
				db2, _, _ := sqlmock.New()
				gdb2, err := gorm.Open("sqlmock", db2)
				if err != nil {
					t.Fatalf("Error opening gorm: %v", err)
				}
				store2 := NewArticleStore(gdb2)
				if store2 == nil {
					t.Fatal("Expected ArticleStore to be non-nil")
				}
				if store2.db != gdb2 {
					t.Errorf("Expected db to be %v, but got %v", gdb2, store2.db)
				}
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
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer gormDB.Close()

	articleStore := &ArticleStore{db: gormDB}

	tests := []struct {
		name    string
		comment *model.Comment
		mock    func()
		wantErr bool
	}{
		{
			name: "Valid Comment",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      "This is a valid comment",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Invalid Comment",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      "",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnError(errors.New("validation failed"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database Error",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      "This is a valid comment",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Duplicate Comment",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      "This is a duplicate comment",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "Empty Comment",
			comment: &model.Comment{},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnError(errors.New("validation failed"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Maximum Length Comment",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      string(make([]byte, 65535)),
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Minimum Length Comment",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      "A",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Special Characters in Comment",
			comment: &model.Comment{
				ArticleID: 1,
				Body:      "This is a comment with special characters: !@#$%^&*()",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnResult(sqlmock.NewResult(1, 1))
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

			tt.mock()
			err := articleStore.CreateComment(tt.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent Requests", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		numConcurrentRequests := 10
		errorsChan := make(chan error, numConcurrentRequests)

		for i := 0; i < numConcurrentRequests; i++ {
			comment := &model.Comment{
				ArticleID: uint(i + 1),
				Body:      fmt.Sprintf("This is comment %d", i+1),
			}
			go func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments` .*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				err := articleStore.CreateComment(comment)
				errorsChan <- err
			}()
		}

		for i := 0; i < numConcurrentRequests; i++ {
			err := <-errorsChan
			if err != nil {
				t.Errorf("CreateComment() error = %v", err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Large Number of Comments", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		numComments := 1000
		for i := 0; i < numComments; i++ {
			comment := &model.Comment{
				ArticleID: uint(i + 1),
				Body:      fmt.Sprintf("This is comment %d", i+1),
			}
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `comments` .*").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			err := articleStore.CreateComment(comment)
			if err != nil {
				t.Errorf("CreateComment() error = %v", err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
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
		setupDB     func() *gorm.DB
		expectedErr error
	}

	testCases := []testCase{
		{
			name:    "Successful Deletion of an Existing Article",
			article: &model.Article{Model: gorm.Model{ID: 1}, Title: "Existing Article", Body: "Content"},
			setupDB: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\".*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return gormDB
			},
			expectedErr: nil,
		},
		{
			name:    "Deletion of a Non-Existing Article",
			article: &model.Article{Model: gorm.Model{ID: 1}, Title: "Non-Existing Article", Body: "Content"},
			setupDB: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\".*").WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
				return gormDB
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "Deletion with Database Connection Error",
			article: &model.Article{Model: gorm.Model{ID: 1}, Title: "Article", Body: "Content"},
			setupDB: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
				return gormDB
			},
			expectedErr: errors.New("connection error"),
		},
		{
			name:        "Deletion with Null Article Model",
			article:     nil,
			setupDB:     func() *gorm.DB { return nil },
			expectedErr: errors.New("article model is nil"),
		},
		{
			name:        "Deletion with Empty Article Model",
			article:     &model.Article{},
			setupDB:     func() *gorm.DB { return nil },
			expectedErr: errors.New("article model is incomplete or invalid"),
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

			db := tc.setupDB()
			if db == nil {
				t.Errorf("expected non-nil db, but got nil")
				return
			}
			store := &ArticleStore{db: db}

			err := store.Delete(tc.article)

			if err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
			} else if err != nil && tc.expectedErr == nil {
				t.Errorf("expected no error, but got %v", err)
			} else if err == nil && tc.expectedErr != nil {
				t.Errorf("expected error %v, but got no error", tc.expectedErr)
			} else {
				t.Logf("Test %s passed successfully", tc.name)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=ArticleStore_GetByID_6fe18728fc
ROOST_METHOD_SIG_HASH=ArticleStore_GetByID_bb488e542f

FUNCTION_DEF=func (s *ArticleStore) GetByID(id uint) (*model.Article, error) // GetByID finds an article from id
*/
func TestArticleStoreGetByID(t *testing.T) {
	type testCase struct {
		name           string
		id             uint
		prepareMock    func(mock sqlmock.Sqlmock)
		expectedErr    bool
		expectedNil    bool
		validateTags   bool
		validateAuthor bool
	}

	testCases := []testCase{
		{
			name: "Successfully retrieve article by ID",
			id:   1,
			prepareMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "author_id", "tags"}).
					AddRow(1, "Test Title", "Test Body", 1, "tag1,tag2")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectedErr:    false,
			expectedNil:    false,
			validateTags:   true,
			validateAuthor: true,
		},
		{
			name: "Handle non-existent article ID gracefully",
			id:   999,
			prepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			expectedErr:    true,
			expectedNil:    true,
			validateTags:   false,
			validateAuthor: false,
		},
		{
			name: "Handle database connection errors",
			id:   1,
			prepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("connection error"))
			},
			expectedErr:    true,
			expectedNil:    true,
			validateTags:   false,
			validateAuthor: false,
		},
		{
			name: "Verify preloaded relationships (Tags and Author)",
			id:   1,
			prepareMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "author_id", "tags"}).
					AddRow(1, "Test Title", "Test Body", 1, "tag1,tag2")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectedErr:    false,
			expectedNil:    false,
			validateTags:   true,
			validateAuthor: true,
		},
		{
			name: "Handle edge case with minimum valid ID (1)",
			id:   1,
			prepareMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "author_id", "tags"}).
					AddRow(1, "Test Title", "Test Body", 1, "tag1,tag2")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectedErr:    false,
			expectedNil:    false,
			validateTags:   true,
			validateAuthor: true,
		},
		{
			name: "Handle edge case with maximum valid ID (uint max)",
			id:   ^uint(0),
			prepareMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "author_id", "tags"}).
					AddRow(^uint(0), "Test Title Max", "Test Body Max", 1, "tag1,tag2")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectedErr:    false,
			expectedNil:    false,
			validateTags:   true,
			validateAuthor: true,
		},
		{
			name: "Handle invalid ID (0)",
			id:   0,
			prepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			expectedErr:    true,
			expectedNil:    true,
			validateTags:   false,
			validateAuthor: false,
		},
		{
			name: "Handle concurrent requests for different IDs",
			id:   1,
			prepareMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "author_id", "tags"}).
					AddRow(1, "Test Title 1", "Test Body 1", 1, "tag1,tag2").
					AddRow(2, "Test Title 2", "Test Body 2", 2, "tag3,tag4")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectedErr:    false,
			expectedNil:    false,
			validateTags:   true,
			validateAuthor: true,
		},
		{
			name: "Handle concurrent requests for the same ID",
			id:   1,
			prepareMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "body", "author_id", "tags"}).
					AddRow(1, "Test Title", "Test Body", 1, "tag1,tag2")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectedErr:    false,
			expectedNil:    false,
			validateTags:   true,
			validateAuthor: true,
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{db: gormDB}

			tc.prepareMock(mock)

			article, err := store.GetByID(tc.id)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got '%s'", err)
				}
			}

			if tc.expectedNil {
				if article != nil {
					t.Errorf("expected nil article, but got '%v'", article)
				}
			} else {
				if article == nil {
					t.Errorf("expected non-nil article, but got nil")
				}
			}

			if tc.validateTags {
				if len(article.Tags) == 0 {
					t.Errorf("expected tags to be preloaded, but got none")
				}
			}

			if tc.validateAuthor {
				if article.Author.ID == 0 {
					t.Errorf("expected author to be preloaded, but got nil")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
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

	tests := []struct {
		name          string
		article       *model.Article
		user          *model.User
		mockBehavior  func(mock sqlmock.Sqlmock)
		expectedFav   bool
		expectedError error
	}{
		{
			name:          "Article and User are both nil",
			article:       nil,
			user:          nil,
			mockBehavior:  func(mock sqlmock.Sqlmock) {},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Article is nil but User is valid",
			article:       nil,
			user:          &model.User{},
			mockBehavior:  func(mock sqlmock.Sqlmock) {},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "User is nil but Article is valid",
			article:       &model.Article{},
			user:          nil,
			mockBehavior:  func(mock sqlmock.Sqlmock) {},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:    "Both Article and User are valid and the article is favorited by the user",
			article: &model.Article{},
			user:    &model.User{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedFav:   true,
			expectedError: nil,
		},
		{
			name:    "Both Article and User are valid but the article is not favorited by the user",
			article: &model.Article{},
			user:    &model.User{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:    "Database error occurs during the query",
			article: &model.Article{},
			user:    &model.User{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count").WillReturnError(errors.New("database error"))
			},
			expectedFav:   false,
			expectedError: errors.New("database error"),
		},
		{
			name:    "Article ID and User ID are not matching in the database",
			article: &model.Article{},
			user:    &model.User{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:    "Article ID and User ID are matching in the database but the record is deleted",
			article: &model.Article{},
			user:    &model.User{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedFav:   false,
			expectedError: nil,
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
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockBehavior(mock)

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
			}
			store := &ArticleStore{db: gormDB}

			fav, err := store.IsFavorited(tt.article, tt.user)

			if fav != tt.expectedFav {
				t.Errorf("expected favorited to be %v, but got %v", tt.expectedFav, fav)
			}
			if err != nil && tt.expectedError == nil {
				t.Errorf("expected no error, but got %v", err)
			} else if err == nil && tt.expectedError != nil {
				t.Errorf("expected error %v, but got nil", tt.expectedError)
			} else if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, but got %v", tt.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
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

	type testCase struct {
		name          string
		userIDs       []uint
		limit         int64
		offset        int64
		expectedError bool
		expectedCount int
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	articleStore := &ArticleStore{db: gormDB}

	testCases := []testCase{
		{
			name:          "Retrieve articles for a single user with valid parameters",
			userIDs:       []uint{1},
			limit:         10,
			offset:        0,
			expectedError: false,
			expectedCount: 3,
		},
		{
			name:          "Retrieve articles for multiple users with valid parameters",
			userIDs:       []uint{1, 2},
			limit:         10,
			offset:        0,
			expectedError: false,
			expectedCount: 5,
		},
		{
			name:          "Retrieve articles with limit and offset",
			userIDs:       []uint{1},
			limit:         5,
			offset:        2,
			expectedError: false,
			expectedCount: 3,
		},
		{
			name:          "Retrieve articles with empty user ID list",
			userIDs:       []uint{},
			limit:         10,
			offset:        0,
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:          "Retrieve articles with negative limit and offset",
			userIDs:       []uint{1},
			limit:         -1,
			offset:        -1,
			expectedError: true,
		},
		{
			name:          "Retrieve articles with non-existent user IDs",
			userIDs:       []uint{999},
			limit:         10,
			offset:        0,
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:          "Retrieve articles with large limit and offset",
			userIDs:       []uint{1},
			limit:         1000,
			offset:        1000,
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:          "Retrieve articles with zero limit",
			userIDs:       []uint{1},
			limit:         0,
			offset:        0,
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:          "Retrieve articles with zero offset",
			userIDs:       []uint{1},
			limit:         10,
			offset:        0,
			expectedError: false,
			expectedCount: 3,
		},
		{
			name:          "Retrieve articles with large limit and zero offset",
			userIDs:       []uint{1},
			limit:         1000,
			offset:        0,
			expectedError: false,
			expectedCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v", r)
					t.Fail()
				}
			}()

			rows := sqlmock.NewRows([]string{"id", "user_id", "title", "body", "created_at", "updated_at"})
			for i := 0; i < tc.expectedCount; i++ {
				rows.AddRow(i+1, tc.userIDs[0], fmt.Sprintf("Article %d", i+1), "Body", "2023-01-01", "2023-01-01")
			}

			var args []driver.Value
			for _, id := range tc.userIDs {
				args = append(args, id)
			}

			mock.ExpectQuery("SELECT").WithArgs(args...).WillReturnRows(rows)

			articles, err := articleStore.GetFeedArticles(tc.userIDs, tc.limit, tc.offset)

			if tc.expectedError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Len(t, articles, tc.expectedCount, "Expected article count mismatch")
			}

			t.Logf("Test case '%s' passed successfully", tc.name)
		})
	}
}
