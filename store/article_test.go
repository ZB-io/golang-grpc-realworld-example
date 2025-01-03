package store

import (
		"testing"
		"github.com/DATA-DOG/go-sqlmock"
		"github.com/jinzhu/gorm"
		"github.com/stretchr/testify/assert"
		"errors"
		"regexp"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"database/sql"
		"time"
		sqlmock "github.com/DATA-DOG/go-sqlmock"
		"sync"
		_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB struct {
	sync.RWMutex
	Value        interface{}
	Error        error
	RowsAffected int64

	// single db
	db                SQLCommon
	blockGlobalUpdate bool
	logMode           logModeValue
	logger            logger
	search            *search
	values            sync.Map

	// global db
	parent        *DB
	callbacks     *Callback
	dialect       Dialect
	singularTable bool

	// function to be used to override the creating of a new timestamp
	nowFuncOverride func() time.Time
}
type ArticleStore struct {
	db *gorm.DB
}
type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
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
type Tag struct {
	gorm.Model
	Name string `gorm:"not null"`
}
type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
}
type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
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
	tests := []struct {
		name    string
		article *model.Article
		dbErr   error
		wantErr bool
	}{
		{
			name: "Successfully create article",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
				Tags: []model.Tag{
					{Name: "test-tag"},
				},
			},
			dbErr:   nil,
			wantErr: false,
		},
		{
			name: "Create article with missing required fields",
			article: &model.Article{
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
			},
			dbErr:   errors.New("title cannot be null"),
			wantErr: true,
		},
		{
			name: "Create article with maximum field lengths",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 255)),
				Body:        string(make([]byte, 1000)),
				UserID:      1,
			},
			dbErr:   nil,
			wantErr: false,
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
			},
			dbErr:   errors.New("database connection lost"),
			wantErr: true,
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
			gormDB.LogMode(true)
			defer gormDB.Close()

			store := &ArticleStore{
				db: gormDB,
			}

			if tt.dbErr != nil {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `articles`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						tt.article.Title, tt.article.Description, tt.article.Body, tt.article.UserID, 0).
					WillReturnError(tt.dbErr)
				mock.ExpectRollback()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `articles`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						tt.article.Title, tt.article.Description, tt.article.Body, tt.article.UserID, 0).
					WillReturnResult(sqlmock.NewResult(1, 1))

				if len(tt.article.Tags) > 0 {
					mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `article_tags`")).
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			}

			err = store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
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
		t.Fatalf("Failed to open GORM connection: %v", err)
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
			name: "Successfully Create Comment",
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
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"Test comment",
						uint(1),
						uint(1),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Missing Required Fields",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 2,
				},
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
						sqlmock.AnyArg(),
						"",
						uint(1),
						uint(1),
					).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Non-existent User ID",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 3,
				},
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"Test comment",
						uint(999),
						uint(1),
					).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 5,
				},
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
						sqlmock.AnyArg(),
						"Test comment",
						uint(1),
						uint(1),
					).
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
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


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
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Delete Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:    "Delete with Nil Article",
			article: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid article: nil pointer"),
		},
		{
			name: "Delete with Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: sql.ErrConnDone,
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

			tc.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			if tc.article == nil {
				err = store.Delete(nil)
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				}
				return
			}

			err = store.Delete(tc.article)

			if tc.expectedError != nil {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				}
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
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func TestArticleStoreDeleteComment(t *testing.T) {
	type testCase struct {
		name          string
		comment       *model.Comment
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
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
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
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
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:    "Delete with nil comment",
			comment: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid comment: nil pointer"),
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
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			assert.NoError(t, err)
			gormDB.LogMode(true)
			defer gormDB.Close()

			tc.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			if tc.comment == nil {
				err = store.DeleteComment(nil)
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			err = store.DeleteComment(tc.comment)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {
	tests := []struct {
		name          string
		commentID     uint
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
		expectedBody  string
	}{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE `comments`.`id` = \\? AND `comments`.`deleted_at` IS NULL").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, "2023-01-01 00:00:00", "2023-01-01 00:00:00", nil, "Test comment", 1, 1))
			},
			expectedError: nil,
			expectedBody:  "Test comment",
		},
		{
			name:      "Non-existent comment",
			commentID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE `comments`.`id` = \\? AND `comments`.`deleted_at` IS NULL").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedBody:  "",
		},
		{
			name:      "Database connection error",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE `comments`.`id` = \\? AND `comments`.`deleted_at` IS NULL").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expectedError: errors.New("database connection error"),
			expectedBody:  "",
		},
		{
			name:      "Zero ID value",
			commentID: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `comments` WHERE `comments`.`id` = \\? AND `comments`.`deleted_at` IS NULL").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedBody:  "",
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

			tt.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			comment, err := store.GetCommentByID(tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedBody, comment.Body)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed", tt.name)
			if err != nil {
				t.Logf("Error: %v", err)
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
			name: "Successfully retrieve tags",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "golang", time.Now(), time.Now()).
					AddRow(2, "testing", time.Now(), time.Now())
				mock.ExpectQuery("^SELECT (.+) FROM `tags`$").WillReturnRows(rows)
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
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				mock.ExpectQuery("^SELECT (.+) FROM `tags`$").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database connection error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`$").WillReturnError(errors.New("connection refused"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("connection refused"),
		},
		{
			name: "Query timeout",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`$").WillReturnError(errors.New("context deadline exceeded"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("context deadline exceeded"),
		},
		{
			name: "Large dataset retrieval",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				for i := 1; i <= 1000; i++ {
					rows.AddRow(uint(i), "tag"+string(rune(i)), time.Now(), time.Now())
				}
				mock.ExpectQuery("^SELECT (.+) FROM `tags`$").WillReturnRows(rows)
			},
			expectedTags: func() []model.Tag {
				tags := make([]model.Tag, 1000)
				for i := range tags {
					tags[i] = model.Tag{
						Model: gorm.Model{ID: uint(i + 1)},
						Name:  "tag" + string(rune(i+1)),
					}
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
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
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
				if tt.name == "Large dataset retrieval" {
					assert.Equal(t, len(tt.expectedTags), len(tags))
				} else {
					assert.Equal(t, tt.expectedTags, tags)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed", tt.name)
			if err != nil {
				t.Logf("Error: %v", err)
			}
			t.Logf("Retrieved %d tags", len(tags))
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
		expectArticle bool
	}

	tests := []testCase{
		{
			name: "Successfully retrieve article with valid ID",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE "articles"."id" = \? AND "articles"."deleted_at" IS NULL`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "created_at", "updated_at", "deleted_at",
						"title", "description", "body", "user_id", "favorites_count",
					}).AddRow(
						1, time.Now(), time.Now(), nil,
						"Test Article", "Test Description", "Test Body", 1, 0,
					))

				mock.ExpectQuery(`SELECT \* FROM "tags" INNER JOIN "article_tags" ON "article_tags"."tag_id" = "tags"."id" WHERE "article_tags"."article_id" = \?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "tag1"))

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "username", "email",
					}).AddRow(1, "testuser", "test@example.com"))
			},
			expectedError: nil,
			expectArticle: true,
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
			expectArticle: false,
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
			expectArticle: false,
		},
		{
			name: "Zero ID handling",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectArticle: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			assert.NoError(t, err)
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			article, err := store.GetByID(tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, article)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)

			t.Logf("Test case '%s' completed", tc.name)
			if err != nil {
				t.Logf("Error: %v", err)
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
				Title:          "Updated Title",
				Description:    "Updated Description",
				Body:           "Updated Body",
				UserID:         1,
				FavoritesCount: 0,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Empty Title Update",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("title cannot be null"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "title cannot be null",
		},
		{
			name: "Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Test Title",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
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
				Title: "Test Title",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "database connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

			err := store.Update(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Update successful")
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
	type testCase struct {
		name          string
		article       *model.Article
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
	}

	tests := []testCase{
		{
			name: "Successfully retrieve comments for an article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				commentRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "article_id", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Comment 1", 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Comment 2", 1, 2)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE (article_id = ?)`)).
					WithArgs(1).
					WillReturnRows(commentRows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "user1", "user1@example.com").
					AddRow(2, "user2", "user2@example.com")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" IN (?)`)).
					WithArgs(1, 2).
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
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "article_id", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE (article_id = ?)`)).
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE (article_id = ?)`)).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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

			gormDB.LogMode(true)

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			comments, err := store.GetComments(tc.article)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tc.expectedCount)
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

	type testCase struct {
		name        string
		article     *model.Article
		user        *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expectFav   bool
		expectError bool
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
			expectFav:   true,
			expectError: false,
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
			expectFav:   false,
			expectError: false,
		},
		{
			name:        "Nil Article Parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectFav:   false,
			expectError: false,
		},
		{
			name:        "Nil User Parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectFav:   false,
			expectError: false,
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
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnError(errors.New("database error"))
			},
			expectFav:   false,
			expectError: true,
		},
		{
			name:        "Both Parameters Nil",
			article:     nil,
			user:        nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectFav:   false,
			expectError: false,
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
			expectFav:   false,
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
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

			isFav, err := store.IsFavorited(tc.article, tc.user)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if isFav != tc.expectFav {
				t.Errorf("Expected favorited status %v but got %v", tc.expectFav, isFav)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
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
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		expected  struct {
			articles []model.Article
			err      error
		}
	}{
		{
			name:    "Successful retrieval with multiple users",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "description", "body"}).
					AddRow(1, time.Now(), time.Now(), nil, 1, "Article 1", "Desc 1", "Body 1").
					AddRow(2, time.Now(), time.Now(), nil, 2, "Article 2", "Desc 2", "Body 2")

				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\?,\\?\\)\\) LIMIT 10 OFFSET 0").
					WithArgs(1, 2).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "user1", "user1@example.com").
					AddRow(2, "user2", "user2@example.com")
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE `id` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(authorRows)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{
					{
						Model:       gorm.Model{ID: 1},
						UserID:      1,
						Title:       "Article 1",
						Description: "Desc 1",
						Body:        "Body 1",
					},
					{
						Model:       gorm.Model{ID: 2},
						UserID:      2,
						Title:       "Article 2",
						Description: "Desc 2",
						Body:        "Body 2",
					},
				},
				err: nil,
			},
		},
		{
			name:    "Empty result set",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "description", "body"})
				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\?\\)\\) LIMIT 10 OFFSET 0").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{},
				err:      nil,
			},
		},
		{
			name:    "Database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `articles`").
					WillReturnError(sql.ErrConnDone)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: nil,
				err:      sql.ErrConnDone,
			},
		},
		{
			name:    "Empty userIDs slice",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "description", "body"})
				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE \\(user_id in \\(\\)\\) LIMIT 10 OFFSET 0").
					WillReturnRows(rows)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{},
				err:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			articles, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected.articles), len(articles))

				if len(tt.expected.articles) > 0 {
					for i, expectedArticle := range tt.expected.articles {
						assert.Equal(t, expectedArticle.Model.ID, articles[i].Model.ID)
						assert.Equal(t, expectedArticle.UserID, articles[i].UserID)
						assert.Equal(t, expectedArticle.Title, articles[i].Title)
						assert.Equal(t, expectedArticle.Description, articles[i].Description)
						assert.Equal(t, expectedArticle.Body, articles[i].Body)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
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
		expected    int
	}{
		{
			name:     "Scenario 1: Get Articles Without Filters",
			tagName:  "",
			username: "",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 1)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected: 1,
		},
		{
			name:     "Scenario 2: Get Articles By Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "User Article", "Description", "Body", 1, 1)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected: 1,
		},
		{
			name:    "Scenario 3: Get Articles By Tag",
			tagName: "technology",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Tech Article", "Description", "Body", 1, 1)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected: 1,
		},
		{
			name: "Scenario 4: Get Favorited Articles",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favRows := sqlmock.NewRows([]string{"article_id"}).AddRow(1)
				mock.ExpectQuery("SELECT article_id FROM").WillReturnRows(favRows)

				articleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "author_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Favorited Article", "Description", "Body", 1, 1)
				mock.ExpectQuery("SELECT").WillReturnRows(articleRows)
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected: 1,
		},
		{
			name:    "Scenario 7: Database Error",
			wantErr: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			expected: 0,
		},
		{
			name:   "Scenario 8: Zero Limit and Offset",
			limit:  0,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "author_id"})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expected: 0,
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
				assert.Len(t, articles, tt.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

