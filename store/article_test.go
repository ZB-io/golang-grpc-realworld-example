package store

import (
	"errors"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"context"
	"sync"
	"github.com/stretchr/testify/require"
	"reflect"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestCreate(t *testing.T) {

	tests := []struct {
		name    string
		article *model.Article
		dbError error
		wantErr bool
	}{
		{
			name: "Successful article creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Missing required fields",
			article: &model.Article{

				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbError: errors.New("title cannot be null"),
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
			dbError: errors.New("database connection failed"),
			wantErr: true,
		},
		{
			name: "Article with tags and author",
			article: &model.Article{
				Title:       "Test Article with Relations",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test-tag"}},
				Author:      model.User{},
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Article with maximum field lengths",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 65535)),
				UserID:      1,
			},
			dbError: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &MockDB{}

			db := &gorm.DB{
				Error: tt.dbError,
			}

			mockDB.On("Create", mock.Anything).Return(db)

			store := &ArticleStore{
				db: db,
			}

			err := store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.dbError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)

			t.Logf("Test case '%s' completed. Error: %v", tt.name, err)
		})
	}
}

/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestCreateComment(t *testing.T) {

	tests := []struct {
		name    string
		comment *model.Comment
		setup   func(*gorm.DB)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successfully create valid comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setup: func(db *gorm.DB) {

				db.Exec("INSERT INTO users (id) VALUES (1)")
				db.Exec("INSERT INTO articles (id, user_id) VALUES (1, 1)")
			},
			wantErr: false,
		},
		{
			name: "Fail with empty body",
			comment: &model.Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			setup:   func(db *gorm.DB) {},
			wantErr: true,
			errMsg:  "not null constraint",
		},
		{
			name: "Fail with non-existent user",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 1,
			},
			setup:   func(db *gorm.DB) {},
			wantErr: true,
			errMsg:  "foreign key constraint",
		},
		{
			name: "Fail with non-existent article",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 999,
			},
			setup:   func(db *gorm.DB) {},
			wantErr: true,
			errMsg:  "foreign key constraint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := setupTestDB()
			if err != nil {
				t.Fatalf("Failed to setup test database: %v", err)
			}
			defer db.Close()

			tt.setup(db)

			store := &ArticleStore{db: db}

			err = store.CreateComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				var savedComment model.Comment
				err = db.First(&savedComment, "body = ?", tt.comment.Body).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.comment.Body, savedComment.Body)
				assert.Equal(t, tt.comment.UserID, savedComment.UserID)
				assert.Equal(t, tt.comment.ArticleID, savedComment.ArticleID)
			}
		})
	}
}

func TestCreateCommentConcurrent(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	db.Exec("INSERT INTO users (id) VALUES (1)")
	db.Exec("INSERT INTO articles (id, user_id) VALUES (1, 1)")

	store := &ArticleStore{db: db}
	numGoroutines := 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			comment := &model.Comment{
				Body:      fmt.Sprintf("Concurrent comment %d", i),
				UserID:    1,
				ArticleID: 1,
			}
			err := store.CreateComment(comment)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	var count int64
	db.Model(&model.Comment{}).Count(&count)
	assert.Equal(t, int64(numGoroutines), count)
}

func setupTestDB() (*gorm.DB, error) {

	return nil, nil
}

/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func TestDelete(t *testing.T) {

	tests := []struct {
		name    string
		article *model.Article
		dbSetup func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
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
			name: "Attempt to delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "Delete with nil article",
			article: nil,
			dbSetup: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
			errMsg:  "article cannot be nil",
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("database connection lost"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection lost",
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

			if tt.dbSetup != nil {
				tt.dbSetup(mock)
			}

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.Delete(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
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
func TestDeleteComment(t *testing.T) {

	type testCase struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		input         *model.Comment
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully Delete Existing Comment",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test Comment",
			},
			expectedError: nil,
		},
		{
			name: "Delete Non-existent Comment",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			input: &model.Comment{
				Model: gorm.Model{ID: 999},
				Body:  "Non-existent Comment",
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			input: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test Comment",
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
			defer gormDB.Close()

			tc.setupMock(mock)

			store := &ArticleStore{db: gormDB}

			err = store.DeleteComment(tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteCommentConcurrent(t *testing.T) {

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

	for i := 1; i <= 3; i++ {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `comments` SET").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(int64(i), 1))
		mock.ExpectCommit()
	}

	store := &ArticleStore{db: gormDB}

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			comment := &model.Comment{
				Model: gorm.Model{ID: id},
				Body:  "Concurrent Test Comment",
			}

			select {
			case <-ctx.Done():
				t.Error("Operation timed out")
				return
			default:
				if err := store.DeleteComment(comment); err != nil {
					t.Errorf("Concurrent delete failed: %v", err)
				}
			}
		}(uint(i))
	}

	wg.Wait()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *MockDB {
	args := m.Called(out, where)
	return args.Get(0).(*MockDB)
}

func (m *MockDB) Preload(column string) *MockDB {
	args := m.Called(column)
	return args.Get(0).(*MockDB)
}

func TestGetByID(t *testing.T) {

	tests := []struct {
		name          string
		id            uint
		setupMock     func(*MockDB)
		expectedError error
		expectedData  *model.Article
	}{
		{
			name: "Successfully retrieve article",
			id:   1,
			setupMock: func(m *MockDB) {
				expectedArticle := &model.Article{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
					Tags:        []model.Tag{{Name: "test"}},
					Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				}

				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(1)).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*model.Article)
						*arg = *expectedArticle
					}).
					Return(m)
				m.Error = nil
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title: "Test Article",
			},
		},
		{
			name: "Article not found",
			id:   99999,
			setupMock: func(m *MockDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(99999)).Return(m)
				m.Error = gorm.ErrRecordNotFound
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
		{
			name: "Database error",
			id:   1,
			setupMock: func(m *MockDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(1)).Return(m)
				m.Error = errors.New("database connection error")
			},
			expectedError: errors.New("database connection error"),
			expectedData:  nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			setupMock: func(m *MockDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(0)).Return(m)
				m.Error = gorm.ErrRecordNotFound
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &MockDB{}
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB,
			}

			article, err := store.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, article)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)

				assert.Equal(t, tt.expectedData.ID, article.ID)
				assert.Equal(t, tt.expectedData.Title, article.Title)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetCommentByID(t *testing.T) {

	tests := []struct {
		name            string
		id              uint
		setupMock       func(*MockDB) *gorm.DB
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name: "Successfully retrieve existing comment",
			id:   1,
			setupMock: func(m *MockDB) *gorm.DB {
				expectedComment := &model.Comment{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				}

				m.On("Find", mock.AnythingOfType("*model.Comment"), []interface{}{uint(1)}).
					Return(&gorm.DB{Error: nil}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*model.Comment)
						*arg = *expectedComment
					})

				return &gorm.DB{Error: nil}
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
			name: "Non-existent comment",
			id:   999,
			setupMock: func(m *MockDB) *gorm.DB {
				m.On("Find", mock.AnythingOfType("*model.Comment"), []interface{}{uint(999)}).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Database connection error",
			id:   1,
			setupMock: func(m *MockDB) *gorm.DB {
				dbError := errors.New("database connection failed")
				m.On("Find", mock.AnythingOfType("*model.Comment"), []interface{}{uint(1)}).
					Return(&gorm.DB{Error: dbError})
				return &gorm.DB{Error: dbError}
			},
			expectedError:   errors.New("database connection failed"),
			expectedComment: nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			setupMock: func(m *MockDB) *gorm.DB {
				m.On("Find", mock.AnythingOfType("*model.Comment"), []interface{}{uint(0)}).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			store := &ArticleStore{
				db: mockDB,
			}

			tt.setupMock(mockDB)

			comment, err := store.GetCommentByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedComment.ID, comment.ID)
				assert.Equal(t, tt.expectedComment.Body, comment.Body)
				assert.Equal(t, tt.expectedComment.UserID, comment.UserID)
				assert.Equal(t, tt.expectedComment.ArticleID, comment.ArticleID)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestGetComments(t *testing.T) {

	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	store := &ArticleStore{db: db}

	tests := []struct {
		name         string
		setupFunc    func(*testing.T, *gorm.DB) *model.Article
		wantErr      bool
		validateFunc func(*testing.T, []model.Comment, error)
	}{
		{
			name: "Successfully retrieve comments for an article",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				article := &model.Article{
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
				}
				if err := db.Create(article).Error; err != nil {
					t.Fatalf("Failed to create test article: %v", err)
				}

				comments := []model.Comment{
					{
						Body:      "Comment 1",
						UserID:    1,
						ArticleID: article.ID,
					},
					{
						Body:      "Comment 2",
						UserID:    2,
						ArticleID: article.ID,
					},
				}
				for _, comment := range comments {
					if err := db.Create(&comment).Error; err != nil {
						t.Fatalf("Failed to create test comment: %v", err)
					}
				}
				return article
			},
			wantErr: false,
			validateFunc: func(t *testing.T, comments []model.Comment, err error) {
				assert.NoError(t, err)
				assert.Len(t, comments, 2)
				assert.NotNil(t, comments[0].Author)
				assert.NotNil(t, comments[1].Author)
			},
		},
		{
			name: "Article with no comments",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				article := &model.Article{
					Title:       "Empty Article",
					Description: "No Comments",
					Body:        "Test Body",
					UserID:      1,
				}
				if err := db.Create(article).Error; err != nil {
					t.Fatalf("Failed to create test article: %v", err)
				}
				return article
			},
			wantErr: false,
			validateFunc: func(t *testing.T, comments []model.Comment, err error) {
				assert.NoError(t, err)
				assert.Empty(t, comments)
			},
		},
		{
			name: "Non-existent article",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				return &model.Article{Model: gorm.Model{ID: 99999}}
			},
			wantErr: false,
			validateFunc: func(t *testing.T, comments []model.Comment, err error) {
				assert.NoError(t, err)
				assert.Empty(t, comments)
			},
		},
		{
			name: "Invalid article ID",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				return &model.Article{Model: gorm.Model{ID: 0}}
			},
			wantErr: false,
			validateFunc: func(t *testing.T, comments []model.Comment, err error) {
				assert.NoError(t, err)
				assert.Empty(t, comments)
			},
		},
		{
			name: "Verify comment order",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				article := &model.Article{
					Title:       "Ordered Comments",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
				}
				if err := db.Create(article).Error; err != nil {
					t.Fatalf("Failed to create test article: %v", err)
				}

				comments := []model.Comment{
					{
						Body:      "First Comment",
						UserID:    1,
						ArticleID: article.ID,
						Model:     gorm.Model{CreatedAt: time.Now().Add(-1 * time.Hour)},
					},
					{
						Body:      "Second Comment",
						UserID:    1,
						ArticleID: article.ID,
						Model:     gorm.Model{CreatedAt: time.Now()},
					},
				}
				for _, comment := range comments {
					if err := db.Create(&comment).Error; err != nil {
						t.Fatalf("Failed to create test comment: %v", err)
					}
				}
				return article
			},
			wantErr: false,
			validateFunc: func(t *testing.T, comments []model.Comment, err error) {
				assert.NoError(t, err)
				assert.Len(t, comments, 2)
				assert.True(t, comments[0].CreatedAt.Before(comments[1].CreatedAt))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cleanupTestDB(t, db)

			article := tt.setupFunc(t, db)

			comments, err := store.GetComments(article)

			tt.validateFunc(t, comments, err)
		})
	}
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {

	t.Helper()
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM articles")
	db.Exec("DELETE FROM users")
}

func setupTestDB() (*gorm.DB, error) {

	return nil, nil
}

/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func TestGetFeedArticles(t *testing.T) {

	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	store := &ArticleStore{db: db}

	tests := []struct {
		name     string
		setup    func(*testing.T, *gorm.DB) []uint
		userIDs  []uint
		limit    int64
		offset   int64
		expected struct {
			count int
			err   bool
		}
	}{
		{
			name: "Scenario 1: Successfully Retrieve Feed Articles for Single User",
			setup: func(t *testing.T, db *gorm.DB) []uint {
				user := createTestUser(t, db, "user1")
				createTestArticles(t, db, user.ID, 5)
				return []uint{user.ID}
			},
			limit:  10,
			offset: 0,
			expected: struct {
				count int
				err   bool
			}{
				count: 5,
				err:   false,
			},
		},
		{
			name: "Scenario 2: Successfully Retrieve Feed Articles for Multiple Users",
			setup: func(t *testing.T, db *gorm.DB) []uint {
				user1 := createTestUser(t, db, "user2")
				user2 := createTestUser(t, db, "user3")
				createTestArticles(t, db, user1.ID, 3)
				createTestArticles(t, db, user2.ID, 3)
				return []uint{user1.ID, user2.ID}
			},
			limit:  20,
			offset: 0,
			expected: struct {
				count int
				err   bool
			}{
				count: 6,
				err:   false,
			},
		},
		{
			name: "Scenario 3: Pagination Testing",
			setup: func(t *testing.T, db *gorm.DB) []uint {
				user := createTestUser(t, db, "user4")
				createTestArticles(t, db, user.ID, 15)
				return []uint{user.ID}
			},
			limit:  5,
			offset: 10,
			expected: struct {
				count int
				err   bool
			}{
				count: 5,
				err:   false,
			},
		},
		{
			name: "Scenario 4: Empty Result Set",
			setup: func(t *testing.T, db *gorm.DB) []uint {
				return []uint{999}
			},
			limit:  10,
			offset: 0,
			expected: struct {
				count int
				err   bool
			}{
				count: 0,
				err:   false,
			},
		},
		{
			name: "Scenario 6: Invalid Input - Empty UserIDs",
			setup: func(t *testing.T, db *gorm.DB) []uint {
				return []uint{}
			},
			limit:  10,
			offset: 0,
			expected: struct {
				count int
				err   bool
			}{
				count: 0,
				err:   false,
			},
		},
		{
			name: "Scenario 8: Deleted Article Handling",
			setup: func(t *testing.T, db *gorm.DB) []uint {
				user := createTestUser(t, db, "user5")
				articles := createTestArticles(t, db, user.ID, 3)

				require.NoError(t, db.Delete(&articles[0]).Error)
				return []uint{user.ID}
			},
			limit:  10,
			offset: 0,
			expected: struct {
				count int
				err   bool
			}{
				count: 2,
				err:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cleanupDB(t, db)

			userIDs := tt.setup(t, db)
			if tt.userIDs == nil {
				tt.userIDs = userIDs
			}

			articles, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.expected.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected.count, len(articles))

			for _, article := range articles {
				assert.NotEmpty(t, article.Author)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}

func cleanupDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	db.Unscoped().Delete(&model.Article{})
	db.Unscoped().Delete(&model.User{})
}

func createTestArticles(t *testing.T, db *gorm.DB, userID uint, count int) []model.Article {
	t.Helper()
	var articles []model.Article
	for i := 0; i < count; i++ {
		article := model.Article{
			Title:       "Test Article",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      userID,
		}
		require.NoError(t, db.Create(&article).Error)
		articles = append(articles, article)
	}
	return articles
}

func createTestUser(t *testing.T, db *gorm.DB, username string) model.User {
	t.Helper()
	user := model.User{
		Username: username,
		Email:    username + "@test.com",
	}
	require.NoError(t, db.Create(&user).Error)
	return user
}

func setupTestDB() (*gorm.DB, error) {

	return nil, nil
}

/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func (m *MockDB) Find(dest interface{}) *gorm.DB {
	args := m.Called(dest)
	return args.Get(0).(*gorm.DB)
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockDB)
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Successfully retrieve tags",
			setupMock: func(mockDB *MockDB) {
				tags := []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "golang"},
					{Model: gorm.Model{ID: 2}, Name: "testing"},
				}
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: nil}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*[]model.Tag)
						*arg = tags
					})
			},
			expectedTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "golang"},
				{Model: gorm.Model{ID: 2}, Name: "testing"},
			},
			expectedError: nil,
		},
		{
			name: "Empty database returns empty tag list",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: nil}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*[]model.Tag)
						*arg = []model.Tag{}
					})
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database connection error",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: errors.New("connection error")})
			},
			expectedTags:  nil,
			expectedError: errors.New("connection error"),
		},
		{
			name: "Database query timeout",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: errors.New("context deadline exceeded")})
			},
			expectedTags:  nil,
			expectedError: errors.New("context deadline exceeded"),
		},
		{
			name: "Large dataset handling",
			setupMock: func(mockDB *MockDB) {
				largeTags := make([]model.Tag, 1000)
				for i := 0; i < 1000; i++ {
					largeTags[i] = model.Tag{
						Model: gorm.Model{ID: uint(i + 1)},
						Name:  fmt.Sprintf("tag-%d", i),
					}
				}
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: nil}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*[]model.Tag)
						*arg = largeTags
					})
			},
			expectedTags:  make([]model.Tag, 1000),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: &gorm.DB{},
			}

			start := time.Now()
			tags, err := store.GetTags()
			duration := time.Since(start)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, tags)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTags, tags)

				if tt.name == "Large dataset handling" {
					assert.Less(t, duration, 1*time.Second, "Query took too long")
				}
			}

			mockDB.AssertExpectations(t)

			t.Logf("Test '%s' completed. Duration: %v", tt.name, duration)
		})
	}
}

/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func (m *MockDB) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestIsFavorited(t *testing.T) {

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(*MockDB)
		expected    bool
		expectedErr error
	}{
		{
			name: "Valid Article and User with Existing Favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(m *MockDB) {
				db := &gorm.DB{Error: nil}
				m.On("Table", "favorite_articles").Return(db)
				m.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 1
				}).Return(&gorm.DB{Error: nil})
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name: "Valid Article and User with No Favorite",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(m *MockDB) {
				db := &gorm.DB{Error: nil}
				m.On("Table", "favorite_articles").Return(db)
				m.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 0
				}).Return(&gorm.DB{Error: nil})
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:    "Nil Article Parameter",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "Nil User Parameter",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user:        nil,
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(m *MockDB) {
				db := &gorm.DB{Error: errors.New("database error")}
				m.On("Table", "favorite_articles").Return(db)
				m.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Return(db)
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:        "Both Parameters Nil",
			article:     nil,
			user:        nil,
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "Zero-Value IDs",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			user: &model.User{
				Model: gorm.Model{ID: 0},
			},
			setupMock: func(m *MockDB) {
				db := &gorm.DB{Error: nil}
				m.On("Table", "favorite_articles").Return(db)
				m.On("Where", "article_id = ? AND user_id = ?", uint(0), uint(0)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 0
				}).Return(&gorm.DB{Error: nil})
			},
			expected:    false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: &gorm.DB{},
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)

			mockDB.AssertExpectations(t)

			t.Logf("Test '%s' completed. Expected: %v, Got: %v", tt.name, tt.expected, result)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := append([]interface{}{query}, args...)
	return m.Called(callArgs...).Get(0).(*gorm.DB)
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
			name:     "Scenario 1: Successfully Create ArticleStore with Valid DB Connection",
			db:       &gorm.DB{},
			wantNil:  false,
			scenario: "Verify successful initialization with valid DB",
		},
		{
			name:     "Scenario 2: Create ArticleStore with Nil DB Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Verify behavior with nil DB connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Executing:", tt.scenario)

			got := NewArticleStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewArticleStore() = %v, want nil: %v", got, tt.wantNil)
				return
			}

			if got != nil && !reflect.DeepEqual(got.db, tt.db) {
				t.Errorf("NewArticleStore() db reference mismatch = %v, want %v", got.db, tt.db)
			}

			if tt.db != nil {

				if got.db != tt.db {
					t.Error("DB reference not maintained correctly")
				}
			}

			t.Log("Test case completed successfully")
		})
	}

	t.Run("Multiple Instance Independence", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}

		store1 := NewArticleStore(db1)
		store2 := NewArticleStore(db2)

		if store1.db == store2.db {
			t.Error("Different ArticleStore instances should maintain independent DB references")
		}
	})

	t.Run("DB Configuration Preservation", func(t *testing.T) {
		mockDB := &gorm.DB{

			Error: nil,
		}
		store := NewArticleStore(mockDB)

		if !reflect.DeepEqual(store.db, mockDB) {
			t.Error("DB configuration not preserved in ArticleStore")
		}
	})

	t.Run("Type Safety Verification", func(t *testing.T) {
		store := NewArticleStore(&gorm.DB{})

		if _, ok := interface{}(store).(*ArticleStore); !ok {
			t.Error("ArticleStore does not implement expected type")
		}
	})
}

/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe


 */
func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(*MockDB)
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
			setupMock: func(mockDB *MockDB) {
				db := &gorm.DB{Error: nil}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			expectError: false,
		},
		{
			name: "Update Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			setupMock: func(mockDB *MockDB) {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Update with Empty Required Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			setupMock: func(mockDB *MockDB) {
				db := &gorm.DB{Error: errors.New("not null constraint violation")}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			expectError: true,
			errorMsg:    "not null constraint violation",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Valid Title",
			},
			setupMock: func(mockDB *MockDB) {
				db := &gorm.DB{Error: errors.New("database connection failed")}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			expectError: true,
			errorMsg:    "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB,
			}

			err := store.Update(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Update(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90


 */
func TestAddFavorite(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successful favorite addition",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("UPDATE").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:    "Nil article parameter",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid article"),
		},
		{
			name: "Association error",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").
					WithArgs(1, 1).
					WillReturnError(errors.New("association error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("association error"),
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
				assert.Equal(t, tc.expectedError.Error(), err.Error())
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
}

func TestAddFavoriteConcurrent(t *testing.T) {

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

	article := &model.Article{
		Model:          gorm.Model{ID: 1},
		FavoritesCount: 0,
	}

	numUsers := 5
	var wg sync.WaitGroup
	store := &ArticleStore{db: gormDB}

	mock.ExpectBegin()
	for i := 0; i < numUsers; i++ {
		mock.ExpectExec("INSERT INTO").
			WithArgs(1, uint(i+1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
	}

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(userID uint) {
			defer wg.Done()
			user := &model.User{Model: gorm.Model{ID: userID}}
			err := store.AddFavorite(article, user)
			assert.NoError(t, err)
		}(uint(i + 1))
	}

	wg.Wait()

	assert.Equal(t, int32(numUsers), article.FavoritesCount)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestDeleteFavorite(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockDB)
		article       *model.Article
		user          *model.User
		expectedErr   error
		expectedCount int32
	}{
		{
			name: "Successful deletion",
			setupMock: func(m *MockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model").Return(tx)
				m.On("Association").Return(tx)
				m.On("Delete").Return(tx)
				m.On("Update").Return(tx)
				m.On("Commit").Return(tx)
			},
			article: &model.Article{
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectedErr:   nil,
			expectedCount: 0,
		},
		{
			name: "Failed association deletion",
			setupMock: func(m *MockDB) {
				tx := &gorm.DB{Error: errors.New("association deletion failed")}
				m.On("Begin").Return(tx)
				m.On("Model").Return(tx)
				m.On("Association").Return(tx)
				m.On("Delete").Return(tx)
				m.On("Rollback").Return(tx)
			},
			article: &model.Article{
				FavoritesCount: 1,
			},
			user:          &model.User{},
			expectedErr:   errors.New("association deletion failed"),
			expectedCount: 1,
		},
		{
			name: "Failed favorites count update",
			setupMock: func(m *MockDB) {
				tx := &gorm.DB{}
				txError := &gorm.DB{Error: errors.New("update failed")}
				m.On("Begin").Return(tx)
				m.On("Model").Return(tx)
				m.On("Association").Return(tx)
				m.On("Delete").Return(tx)
				m.On("Update").Return(txError)
				m.On("Rollback").Return(tx)
			},
			article: &model.Article{
				FavoritesCount: 1,
			},
			user:          &model.User{},
			expectedErr:   errors.New("update failed"),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Running test case:", tt.name)

			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB.DB,
			}

			initialCount := tt.article.FavoritesCount
			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Equal(t, initialCount, tt.article.FavoritesCount, "FavoritesCount should not change on error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount, "FavoritesCount should be decremented")
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestDeleteFavoriteConcurrent(t *testing.T) {
	article := &model.Article{
		FavoritesCount: 5,
	}

	mockDB := new(MockDB)
	store := &ArticleStore{
		db: mockDB.DB,
	}

	tx := &gorm.DB{}
	mockDB.On("Begin").Return(tx)
	mockDB.On("Model").Return(tx)
	mockDB.On("Association").Return(tx)
	mockDB.On("Delete").Return(tx)
	mockDB.On("Update").Return(tx)
	mockDB.On("Commit").Return(tx)

	numGoroutines := 5
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(userID uint) {
			defer wg.Done()
			user := &model.User{
				Model: gorm.Model{ID: userID},
			}
			err := store.DeleteFavorite(article, user)
			assert.NoError(t, err)
		}(uint(i + 1))
	}

	wg.Wait()
	assert.Equal(t, int32(0), article.FavoritesCount)
	mockDB.AssertExpectations(t)
}

/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestGetArticles(t *testing.T) {

	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	store := &ArticleStore{db: db}

	testUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	testArticles := []model.Article{
		{
			Model:       gorm.Model{ID: 1},
			Title:       "Test Article 1",
			Description: "Test Description 1",
			Body:        "Test Body 1",
			UserID:      testUser.ID,
			Author:      *testUser,
			Tags:        []model.Tag{{Name: "testtag"}},
		},
		{
			Model:       gorm.Model{ID: 2},
			Title:       "Test Article 2",
			Description: "Test Description 2",
			Body:        "Test Body 2",
			UserID:      testUser.ID,
			Author:      *testUser,
		},
	}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		wantLen     int
		wantErr     bool
	}{
		{
			name:        "Get Articles Without Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantLen:     2,
			wantErr:     false,
		},
		{
			name:        "Get Articles By Username",
			tagName:     "",
			username:    "testuser",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantLen:     2,
			wantErr:     false,
		},
		{
			name:        "Get Articles By Tag",
			tagName:     "testtag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantLen:     1,
			wantErr:     false,
		},
		{
			name:        "Get Favorited Articles",
			tagName:     "",
			username:    "",
			favoritedBy: testUser,
			limit:       10,
			offset:      0,
			wantLen:     0,
			wantErr:     false,
		},
		{
			name:        "Test Pagination",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       1,
			offset:      1,
			wantLen:     1,
			wantErr:     false,
		},
		{
			name:        "Empty Results Test",
			tagName:     "nonexistenttag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantLen:     0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tt.name)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantLen, len(articles))

			if tt.username != "" {
				for _, article := range articles {
					assert.Equal(t, tt.username, article.Author.Username)
				}
			}

			if tt.tagName != "" {
				for _, article := range articles {
					hasTag := false
					for _, tag := range article.Tags {
						if tag.Name == tt.tagName {
							hasTag = true
							break
						}
					}
					assert.True(t, hasTag)
				}
			}

			t.Logf("Test case completed successfully: %s", tt.name)
		})
	}
}

func setupTestDB() (*gorm.DB, error) {

	return nil, nil
}

