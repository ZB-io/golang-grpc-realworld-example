package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"errors"
	"time"
)








/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore 

 */
func TestNewArticleStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open("mysql", mockDB)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	tests := []testCase{
		{
			name:     "Scenario 1: Successfully Create ArticleStore with Valid DB Connection",
			db:       gormDB,
			wantNil:  false,
			scenario: "Verify successful initialization with valid DB",
		},
		{
			name:     "Scenario 2: Create ArticleStore with Nil DB Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Verify handling of nil DB connection",
		},
		{
			name:     "Scenario 3: Verify DB Reference Integrity",
			db:       gormDB,
			wantNil:  false,
			scenario: "Ensure DB reference matches input",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Logf("Running scenario: %s", tc.scenario)

			store := NewArticleStore(tc.db)

			if tc.wantNil {
				assert.Nil(t, store, "Expected nil ArticleStore")
			} else {
				assert.NotNil(t, store, "Expected non-nil ArticleStore")
				if tc.db != nil {
					assert.Equal(t, tc.db, store.db, "DB reference mismatch")
				}
			}

			if tc.db != nil {
				assert.Same(t, tc.db, store.db, "DB reference should be the same instance")
			}

			t.Logf("Test completed successfully for: %s", tc.name)
		})
	}

	t.Run("Scenario 4: Multiple ArticleStore Instances Independence", func(t *testing.T) {
		store1 := NewArticleStore(gormDB)
		store2 := NewArticleStore(gormDB)

		assert.NotSame(t, store1, store2, "Different instances should not be the same")
		assert.Equal(t, store1.db, store2.db, "DB references should be equal")

		t.Log("Successfully verified multiple instance independence")
	})

	t.Run("Scenario 5: ArticleStore with Configured DB Connection", func(t *testing.T) {
		configuredDB := gormDB.LogMode(true)
		store := NewArticleStore(configuredDB)

		assert.Equal(t, configuredDB, store.db, "Should maintain DB configuration")
		t.Log("Successfully verified configured DB connection")
	})

	t.Run("Scenario 6: Memory Usage and Resource Management", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			store := NewArticleStore(gormDB)
			assert.NotNil(t, store, "Store creation should succeed")
		}
		t.Log("Successfully completed memory usage test")
	})

	t.Run("Scenario 7: Type Safety and Interface Compliance", func(t *testing.T) {
		store := NewArticleStore(gormDB)

		assert.IsType(t, &ArticleStore{}, store, "Should be of type *ArticleStore")
		t.Log("Successfully verified type safety")
	})
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error 

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
		errMsg  string
	}{
		{
			name: "Successfully Create Valid Comment",
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
			name: "Create Comment with Missing Body",
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
			errMsg:  "sql: no rows in result set",
		},
		{
			name: "Create Comment with Non-Existent UserID",
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
			errMsg:  "record not found",
		},
		{
			name: "Create Comment with Maximum Length Body",
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
		{
			name: "Database Connection Error",
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
			errMsg:  "sql: connection is already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mock()

			err := store.CreateComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Comment created successfully")
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

FUNCTION_DEF=func (s *ArticleStore) GetTags() ([]model.Tag, error) 

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
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
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
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnError(errors.New("database connection error"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Database Query Timeout",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnError(errors.New("context deadline exceeded"))
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("context deadline exceeded"),
		},
		{
			name: "Large Dataset Retrieval",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"})

				for i := 1; i <= 1000; i++ {
					rows.AddRow(i, time.Now(), time.Now(), nil, "tag"+string(rune(i)))
				}
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
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
		{
			name: "Malformed Tag Data",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, "invalid_time", time.Now(), nil, "golang")
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: errors.New("sql: Scan error"),
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

			store := &ArticleStore{db: gormDB}

			tags, err := store.GetTags()

			t.Logf("Executing test case: %s", tt.name)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedTags), len(tags))
				if len(tt.expectedTags) > 0 {

					assert.Equal(t, tt.expectedTags[0].ID, tags[0].ID)
					assert.Equal(t, tt.expectedTags[0].Name, tags[0].Name)
				}
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

FUNCTION_DEF=func (s *ArticleStore) GetComments(m *model.Article) ([]model.Comment, error) 

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
		errorMessage  string
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

				mock.ExpectQuery("^SELECT (.+) FROM `comments`").
					WithArgs(1).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
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
				mock.ExpectQuery("^SELECT (.+) FROM `comments`").
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
				mock.ExpectQuery("^SELECT (.+) FROM `comments`").
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectedCount: 0,
			expectError:   true,
			errorMessage:  "database error",
		},
		{
			name: "Invalid article ID",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id", "author_id"})
				mock.ExpectQuery("^SELECT (.+) FROM `comments`").
					WithArgs(0).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			comments, err := store.GetComments(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tt.expectedCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Expected count: %d, Got count: %d",
				tt.name, tt.expectedCount, len(comments))
		})
	}
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) 

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

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) 

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
		wantLen     int
	}{
		{
			name:   "Scenario 1: Get Articles Without Filters",
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
		{
			name:     "Scenario 2: Get Articles By Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
		{
			name:    "Scenario 3: Get Articles By Tag",
			tagName: "programming",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
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
				mock.ExpectQuery("SELECT article_id FROM").
					WillReturnRows(favRows)

				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("SELECT").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").
					WillReturnRows(authorRows)
			},
			wantLen: 1,
		},
		{
			name:    "Scenario 7: Empty Result Set",
			tagName: "nonexistenttag",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id", "favorites_count"})
				mock.ExpectQuery("SELECT").
					WillReturnRows(rows)
			},
			wantLen: 0,
		},
		{
			name:  "Scenario 8: Database Error",
			limit: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
			wantLen: 0,
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
				assert.Len(t, articles, tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

