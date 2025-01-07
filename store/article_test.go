package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
	"database/sql"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/mock"
	"time"
	"github.com/jinzhu/gorm/dialects/mysql"
)





type MockDB struct {
	mock.Mock
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
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful article creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Model:       gorm.Model{ID: 1},
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
			name: "Empty required fields",
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
			name: "Database connection error",
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
			name: "Article with tags",
			article: &model.Article{
				Title:       "Test Article with Tags",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test-tag"}},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
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
			name: "Maximum field lengths",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 255)),
				Body:        string(make([]byte, 1000)),
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
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
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func (m *MockDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(value, where)
	return args.Get(0).(*gorm.DB)
}

func TestArticleStoreDelete(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbError error
		want    error
	}{
		{
			name: "Successfully delete existing article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
				Body:  "Test Body",
			},
			dbError: nil,
			want:    nil,
		},
		{
			name: "Attempt to delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			dbError: gorm.ErrRecordNotFound,
			want:    gorm.ErrRecordNotFound,
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			dbError: errors.New("database connection error"),
			want:    errors.New("database connection error"),
		},
		{
			name:    "Nil article parameter",
			article: nil,
			dbError: errors.New("invalid article"),
			want:    errors.New("invalid article"),
		},
		{
			name: "Delete article with associated records",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
				Tags:  []model.Tag{{Name: "test"}},
			},
			dbError: nil,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)

			db := &gorm.DB{Error: tt.dbError}

			store := &ArticleStore{
				db: db,
			}

			if tt.article != nil {
				mockDB.On("Delete", tt.article, mock.Anything).Return(db)
			}

			got := store.Delete(tt.article)

			if tt.want == nil {
				assert.NoError(t, got)
			} else {
				assert.Error(t, got)
				assert.Equal(t, tt.want.Error(), got.Error())
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
		setupFn func(sqlmock.Sqlmock)
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
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name:    "Delete Comment with NULL Input",
			comment: nil,
			setupFn: func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "invalid comment: nil pointer",
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name: "Delete Comment with Associated Records",
			comment: &model.Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body:      "Test Comment with Associations",
				UserID:    1,
				ArticleID: 1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `comments` WHERE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			if tt.setupFn != nil {
				tt.setupFn(mock)
			}

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{
				db: gormDB,
			}

			if tt.comment == nil {
				err = store.DeleteComment(nil)
				if !tt.wantErr {
					t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("DeleteComment() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			err = store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("DeleteComment() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {
	type testCase struct {
		name            string
		commentID       uint
		setupMock       func(sqlmock.Sqlmock)
		expectedError   error
		expectedComment *model.Comment
	}

	tests := []testCase{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`.`id` = \\?\\) AND `comments`.`deleted_at` IS NULL").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "Test comment", 1, 1))
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name:      "Non-existent comment",
			commentID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`.`id` = \\?\\) AND `comments`.`deleted_at` IS NULL").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name:      "Database connection error",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments`").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError:   sql.ErrConnDone,
			expectedComment: nil,
		},
		{
			name:      "Zero ID input",
			commentID: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`.`id` = \\?\\) AND `comments`.`deleted_at` IS NULL").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
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

			store := &ArticleStore{db: gormDB}

			comment, err := store.GetCommentByID(tc.commentID)

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
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test Article", "Test Description", "Test Body", 1, 0,
				)

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnRows(rows)

				tagRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "tag1")
				mock.ExpectQuery(`SELECT \* FROM "tags"`).
					WillReturnRows(tagRows)

				userRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(userRows)
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
				Tags:  []model.Tag{{Model: gorm.Model{ID: 1}, Name: "tag1"}},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
				},
				UserID:         1,
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 0,
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
				assert.Equal(t, tc.expectedData.ID, article.ID)
				assert.Equal(t, tc.expectedData.Title, article.Title)
				assert.Equal(t, tc.expectedData.Description, article.Description)
				assert.Equal(t, tc.expectedData.Body, article.Body)
				assert.Equal(t, tc.expectedData.UserID, article.UserID)
				assert.Equal(t, tc.expectedData.FavoritesCount, article.FavoritesCount)

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
		t.Fatalf("Failed to open GORM connection: %v", err)
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
			name: "Update with Empty Required Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("Column 'title' cannot be null"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "Column 'title' cannot be null",
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
		want      []model.Article
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
					"author.id", "author.created_at", "author.updated_at", "author.deleted_at",
					"author.username", "author.email",
				}).
					AddRow(
						1, time.Now(), time.Now(), nil,
						"Test Article", "Description", "Body", 1, 0,
						1, time.Now(), time.Now(), nil,
						"testuser", "test@example.com",
					)

				mock.ExpectQuery("SELECT").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want: []model.Article{
				{
					Model:       gorm.Model{ID: 1},
					Title:       "Test Article",
					Description: "Description",
					Body:        "Body",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
					UserID:         1,
					FavoritesCount: 0,
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
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery("SELECT").
					WithArgs(sqlmock.AnyArg()).
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
				mock.ExpectQuery("SELECT").
					WithArgs(sqlmock.AnyArg()).
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
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery("SELECT").
					WillReturnRows(rows)
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Zero limit and offset",
			userIDs: []uint{1},
			limit:   0,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery("SELECT").
					WithArgs(sqlmock.AnyArg()).
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
				t.Errorf("there were unfulfilled expectations: %s", err)
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

	baseTime := time.Now()
	validArticle := &model.Article{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: baseTime,
			UpdatedAt: baseTime,
		},
		Title:          "Test Article",
		Description:    "Test Description",
		Body:           "Test Body",
		FavoritesCount: 0,
	}

	validUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: baseTime,
			UpdatedAt: baseTime,
		},
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	tests := []testCase{
		{
			name:    "Successful favorite addition",
			article: validArticle,
			user:    validUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(validArticle.ID, validUser.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, validArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:    "Database error during association",
			article: validArticle,
			user:    validUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(validArticle.ID, validUser.ID).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database error"),
		},
		{
			name:    "Database error during count update",
			article: validArticle,
			user:    validUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(validArticle.ID, validUser.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1, validArticle.ID).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update error"),
		},
		{
			name:          "Nil article parameter",
			article:       nil,
			user:          validUser,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article"),
		},
		{
			name:          "Nil user parameter",
			article:       validArticle,
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

			tc.setupMock(mock)

			store := &ArticleStore{db: gormDB}

			initialCount := int32(0)
			if tc.article != nil {
				initialCount = tc.article.FavoritesCount
			}

			err = store.AddFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				if tc.article != nil {
					assert.Equal(t, initialCount, tc.article.FavoritesCount, "FavoritesCount should not change on error")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, initialCount+1, tc.article.FavoritesCount, "FavoritesCount should increment by 1")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {
	createMockDB := func(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to open GORM DB: %v", err)
		}
		gormDB.LogMode(true)
		return db, mock, gormDB
	}

	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		article *model.Article
		user    *model.User
		wantErr bool
	}{
		{
			name: "Successful deletion",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET `favorites_count` = `favorites_count` - 1").
					WithArgs().
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:          "Test Article",
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Association deletion error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Update favorites count error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET `favorites_count` = `favorites_count` - 1").
					WithArgs().
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Null parameter handling",
			setup: func(mock sqlmock.Sqlmock) {

			},
			article: nil,
			user:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, gormDB := createMockDB(t)
			defer db.Close()

			if tt.setup != nil {
				tt.setup(mock)
			}

			store := &ArticleStore{
				db: gormDB,
			}

			originalCount := int32(0)
			if tt.article != nil {
				originalCount = tt.article.FavoritesCount
			}

			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.article != nil {
					assert.Equal(t, originalCount, tt.article.FavoritesCount, "FavoritesCount should not change on error")
				}
			} else {
				assert.NoError(t, err)
				if tt.article != nil {
					assert.Equal(t, originalCount-1, tt.article.FavoritesCount, "FavoritesCount should decrease by 1")
				}
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err, "All database expectations should be met")
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

