package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
)








/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore 

 */
func TestNewArticleStore(t *testing.T) {

	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name:     "Scenario 1: Successfully Create New ArticleStore with Valid DB Connection",
			db:       setupTestDB(t),
			wantNil:  false,
			scenario: "Valid DB connection should create a proper ArticleStore instance",
		},
		{
			name:     "Scenario 2: Create ArticleStore with Nil DB Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB connection should still create ArticleStore instance but with nil DB reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("Testing scenario: %s", tt.scenario)

			store := NewArticleStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, store, "ArticleStore should be nil")
			} else {
				assert.NotNil(t, store, "ArticleStore should not be nil")
				assert.Equal(t, tt.db, store.db, "DB reference should match input DB")
			}

			if tt.db != nil {

				assert.Equal(t, tt.db, store.db, "DB reference should maintain integrity")

				t.Log("Successfully verified DB reference integrity")
			}

		})
	}

	t.Run("Multiple ArticleStore Instances Independence", func(t *testing.T) {
		db1 := setupTestDB(t)
		db2 := setupTestDB(t)

		store1 := NewArticleStore(db1)
		store2 := NewArticleStore(db2)

		assert.NotEqual(t, store1.db, store2.db, "Different stores should have independent DB references")
		t.Log("Successfully verified multiple instance independence")
	})

	t.Run("Memory Management", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			db := setupTestDB(t)
			store := NewArticleStore(db)
			assert.NotNil(t, store, "Store creation should succeed")
		}
		t.Log("Successfully completed memory management test")
	})
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
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	return gormDB
}


/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377

FUNCTION_DEF=func (s *ArticleStore) Create(m *model.Article) error 

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
		name        string
		article     *model.Article
		mockSetup   func(sqlmock.Sqlmock)
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
				Tags:        []model.Tag{{Name: "test-tag"}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `article_tags`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Missing Required Fields",
			article: &model.Article{

				Body:   "Test Body",
				UserID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(errors.New("validation failed"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: connection is already closed",
		},
		{
			name: "Article with Tags and Author",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test-tag"}},
				Author:      model.User{Model: gorm.Model{ID: 1}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `article_tags`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			err := store.Create(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Article created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

