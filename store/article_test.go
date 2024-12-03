package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/DATA-DOG/go-sqlmock"
	"sync"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
	"fmt"
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
			name: "Successful Article Creation",
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
			name: "Missing Required Fields",
			article: &model.Article{
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbError: errors.New("validation error: Title is required"),
			wantErr: true,
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbError: errors.New("database connection error"),
			wantErr: true,
		},
		{
			name: "Article with Tags and Author",
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
			name: "Maximum Field Lengths",
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
			mockDB := new(MockDB)

			store := &ArticleStore{
				db: &gorm.DB{},
			}

			mockDB.On("Create", mock.AnythingOfType("*model.Article")).Return(&gorm.DB{
				Error: tt.dbError,
			})

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

	type testCase struct {
		name          string
		comment       *model.Comment
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully Create Valid Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test comment", uint(1), uint(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Empty Comment Body",
			comment: &model.Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(errors.New("not null constraint violation"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("not null constraint violation"),
		},
		{
			name: "Non-Existent UserID",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    999,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(errors.New("foreign key constraint violation"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("foreign key constraint violation"),
		},
		{
			name: "Database Connection Error",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Maximum Length Comment Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
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

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
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
func TestDelete(t *testing.T) {

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
			name: "Fail to delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnError(errors.New("record not found"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("record not found"),
		},
		{
			name: "Database connection error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("database connection lost"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database connection lost"),
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

			err = store.Delete(tc.article)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
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
	tests := []struct {
		name    string
		setup   func(*gorm.DB) *model.Comment
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successfully Delete Existing Comment",
			setup: func(db *gorm.DB) *model.Comment {
				comment := &model.Comment{
					Body:      "Test Comment",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(comment)
				return comment
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Delete Non-existent Comment",
			setup: func(db *gorm.DB) *model.Comment {
				return &model.Comment{Model: gorm.Model{ID: 99999}}
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Delete Comment with NULL Fields",
			setup: func(db *gorm.DB) *model.Comment {
				comment := &model.Comment{
					Body:      "",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(comment)
				return comment
			},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := setupTestDB()
			if err != nil {
				t.Fatalf("failed to setup test database: %v", err)
			}
			defer db.Close()

			store := &ArticleStore{db: db}
			comment := tt.setup(db)

			if tt.name == "Successfully Delete Existing Comment" {
				var wg sync.WaitGroup
				wg.Add(2)

				errChan := make(chan error, 2)

				for i := 0; i < 2; i++ {
					go func() {
						defer wg.Done()
						errChan <- store.DeleteComment(comment)
					}()
				}

				wg.Wait()
				close(errChan)

				successCount := 0
				for err := range errChan {
					if err == nil {
						successCount++
					}
				}

				assert.Equal(t, 1, successCount, "Only one deletion should succeed")
			} else {

				err := store.DeleteComment(comment)

				if tt.wantErr {
					assert.Error(t, err)
					if tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)

					var found model.Comment
					result := db.Unscoped().First(&found, comment.ID)
					assert.NoError(t, result.Error)
					assert.NotNil(t, found.DeletedAt, "DeletedAt should be populated")
				}
			}
		})
	}
}

func TestDeleteCommentDBError(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test database: %v", err)
	}

	db.Close()

	store := &ArticleStore{db: db}
	comment := &model.Comment{Model: gorm.Model{ID: 1}}

	err = store.DeleteComment(comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sql: database is closed")
}

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Comment{}, &model.User{}, &model.Article{})

	return db, nil
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func (m *ExtendedMockDB) Find(out interface{}, where ...interface{}) *ExtendedMockDB {
	args := m.MockDB.Called(out, where)
	return args.Get(0).(*ExtendedMockDB)
}

func (m *ExtendedMockDB) Preload(column string) *ExtendedMockDB {
	m.MockDB.Called(column)
	return m
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		setupMock     func(*ExtendedMockDB)
		expectedError error
		expectedData  *model.Article
	}{
		{
			name: "Successfully retrieve article",
			id:   1,
			setupMock: func(mockDB *ExtendedMockDB) {
				expectedArticle := &model.Article{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					Tags:        []model.Tag{{Name: "test"}},
					Author:      model.User{Model: gorm.Model{ID: 1}},
					UserID:      1,
				}

				mockDB.On("Preload", "Tags").Return(mockDB)
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Find", mock.Anything, uint(1)).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*model.Article)
						*arg = *expectedArticle
					}).
					Return(mockDB)
				mockDB.Error = nil
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				Tags:        []model.Tag{{Name: "test"}},
				Author:      model.User{Model: gorm.Model{ID: 1}},
				UserID:      1,
			},
		},
		{
			name: "Article not found",
			id:   999,
			setupMock: func(mockDB *ExtendedMockDB) {
				mockDB.On("Preload", "Tags").Return(mockDB)
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Find", mock.Anything, uint(999)).Return(mockDB)
				mockDB.Error = gorm.ErrRecordNotFound
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
		{
			name: "Database error",
			id:   1,
			setupMock: func(mockDB *ExtendedMockDB) {
				mockDB.On("Preload", "Tags").Return(mockDB)
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Find", mock.Anything, uint(1)).Return(mockDB)
				mockDB.Error = errors.New("database connection error")
			},
			expectedError: errors.New("database connection error"),
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &ExtendedMockDB{
				MockDB: &MockDB{},
			}
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
				assert.Equal(t, tt.expectedData.Description, article.Description)
				assert.Equal(t, tt.expectedData.Body, article.Body)
				assert.Equal(t, tt.expectedData.UserID, article.UserID)
				assert.Len(t, article.Tags, len(tt.expectedData.Tags))
			}

			mockDB.MockDB.AssertExpectations(t)
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
		commentID       uint
		setupMock       func(*MockDB) *gorm.DB
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			setupMock: func(mockDB *MockDB) *gorm.DB {
				comment := &model.Comment{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				}

				db := &gorm.DB{Error: nil}
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"),
					mock.AnythingOfType("uint")).Return(db).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.Comment)
					*arg = *comment
				})

				return db
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name:      "Non-existent comment",
			commentID: 99999,
			setupMock: func(mockDB *MockDB) *gorm.DB {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"),
					mock.AnythingOfType("uint")).Return(db)
				return db
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name:      "Database connection error",
			commentID: 1,
			setupMock: func(mockDB *MockDB) *gorm.DB {
				dbError := errors.New("database connection error")
				db := &gorm.DB{Error: dbError}
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"),
					mock.AnythingOfType("uint")).Return(db)
				return db
			},
			expectedError:   errors.New("database connection error"),
			expectedComment: nil,
		},
		{
			name:      "Zero ID parameter",
			commentID: 0,
			setupMock: func(mockDB *MockDB) *gorm.DB {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"),
					mock.AnythingOfType("uint")).Return(db)
				return db
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			store := &ArticleStore{
				db: tt.setupMock(mockDB),
			}

			comment, err := store.GetCommentByID(tt.commentID)

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
		name          string
		setupFunc     func(*testing.T, *gorm.DB) *model.Article
		expectedCount int
		expectError   bool
	}{
		{
			name: "Successfully retrieve comments for article",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				article := &model.Article{
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
				}

				err := db.Create(article).Error
				assert.NoError(t, err)

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
					err := db.Create(&comment).Error
					assert.NoError(t, err)
				}

				return article
			},
			expectedCount: 2,
			expectError:   false,
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

				err := db.Create(article).Error
				assert.NoError(t, err)
				return article
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Non-existent article",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				return &model.Article{Model: gorm.Model{ID: 99999}}
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Database connection error",
			setupFunc: func(t *testing.T, db *gorm.DB) *model.Article {
				db.Close()
				return &model.Article{Model: gorm.Model{ID: 1}}
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := tt.setupFunc(t, db)

			comments, err := store.GetComments(article)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tt.expectedCount)

				if tt.expectedCount > 0 {
					for _, comment := range comments {
						assert.NotZero(t, comment.Author.ID)
					}
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func TestGetFeedArticles(t *testing.T) {
	type testCase struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		expected  struct {
			articles []model.Article
			err      error
		}
	}

	now := time.Now()
	testArticles := []model.Article{
		{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Title:       "Test Article 1",
			Description: "Test Description 1",
			UserID:      1,
			Author: model.User{
				Model: gorm.Model{ID: 1},
			},
		},
		{
			Model: gorm.Model{
				ID:        2,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Title:       "Test Article 2",
			Description: "Test Description 2",
			UserID:      2,
			Author: model.User{
				Model: gorm.Model{ID: 2},
			},
		},
	}

	tests := []testCase{
		{
			name:    "Successfully retrieve articles for single user",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE \(user_id in \(\?\)\) LIMIT 10 OFFSET 0`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "user_id"}).
						AddRow(1, "Test Article 1", "Test Description 1", 1))

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{testArticles[0]},
				err:      nil,
			},
		},
		{
			name:    "Successfully retrieve articles for multiple users",
			userIDs: []uint{1, 2},
			limit:   20,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE \(user_id in \(\?,\?\)\) LIMIT 20 OFFSET 0`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "user_id"}).
						AddRow(1, "Test Article 1", "Test Description 1", 1).
						AddRow(2, "Test Article 2", "Test Description 2", 2))

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1).
						AddRow(2))
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: testArticles,
				err:      nil,
			},
		},
		{
			name:    "Handle empty result set",
			userIDs: []uint{999},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE`).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "user_id"}))
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
			name:    "Handle database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE`).
					WillReturnError(errors.New("database error"))
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: nil,
				err:      errors.New("database error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm connection: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			articles, err := store.GetFeedArticles(tc.userIDs, tc.limit, tc.offset)

			if tc.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expected.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expected.articles), len(articles))

				if len(articles) > 0 {
					for i, article := range articles {
						assert.Equal(t, tc.expected.articles[i].Title, article.Title)
						assert.Equal(t, tc.expected.articles[i].UserID, article.UserID)
						assert.Equal(t, tc.expected.articles[i].Author.ID, article.Author.ID)
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
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func TestGetTags(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockDB)
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Successfully retrieve tags",
			setupMock: func(m *MockDB) {
				tags := []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "golang"},
					{Model: gorm.Model{ID: 2}, Name: "testing"},
				}
				m.On("Find", mock.AnythingOfType("*[]model.Tag")).
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
			setupMock: func(m *MockDB) {
				m.On("Find", mock.AnythingOfType("*[]model.Tag")).
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
			setupMock: func(m *MockDB) {
				m.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: errors.New("connection error")})
			},
			expectedTags:  nil,
			expectedError: errors.New("connection error"),
		},
		{
			name: "Database query timeout",
			setupMock: func(m *MockDB) {
				m.On("Find", mock.AnythingOfType("*[]model.Tag")).
					Return(&gorm.DB{Error: errors.New("context deadline exceeded")})
			},
			expectedTags:  nil,
			expectedError: errors.New("context deadline exceeded"),
		},
		{
			name: "Large dataset handling",
			setupMock: func(m *MockDB) {
				largeTags := make([]model.Tag, 1000)
				for i := 0; i < 1000; i++ {
					largeTags[i] = model.Tag{
						Model: gorm.Model{ID: uint(i + 1)},
						Name:  fmt.Sprintf("tag-%d", i),
					}
				}
				m.On("Find", mock.AnythingOfType("*[]model.Tag")).
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
				db: &gorm.DB{
					Value: mockDB,
				},
			}

			startTime := time.Now()
			tags, err := store.GetTags()
			duration := time.Since(startTime)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, tags)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedTags), len(tags))
				assert.Equal(t, tt.expectedTags, tags)

				if len(tags) > 100 {
					assert.Less(t, duration, 1*time.Second, "Query took too long for large dataset")
				}
			}

			mockDB.AssertExpectations(t)
			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
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
			name:        "Nil Article Parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User Parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
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
				m.On("Count", mock.AnythingOfType("*int")).Return(&gorm.DB{Error: errors.New("database error")})
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB,
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
		})
	}
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
			scenario: "Valid DB connection should create valid ArticleStore",
		},
		{
			name:     "Scenario 2: Create ArticleStore with Nil DB Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB connection should still create ArticleStore instance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Scenario:", tt.scenario)

			got := NewArticleStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewArticleStore() = %v, want nil: %v", got, tt.wantNil)
				return
			}

			if got != nil && !reflect.DeepEqual(got.db, tt.db) {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.db)
			}

			t.Log("Test completed successfully")
		})
	}

	t.Run("Scenario 3: Verify ArticleStore Instance Independence", func(t *testing.T) {
		mockDB := &gorm.DB{}
		store1 := NewArticleStore(mockDB)
		store2 := NewArticleStore(mockDB)

		if store1 == store2 {
			t.Error("Expected different instances, got same instance")
		}
		if store1.db != store2.db {
			t.Error("Expected same DB reference, got different references")
		}
		t.Log("Instance independence verified successfully")
	})

	t.Run("Scenario 4: Verify DB Reference Integrity", func(t *testing.T) {
		mockDB := &gorm.DB{}
		store := NewArticleStore(mockDB)

		if !reflect.DeepEqual(store.db, mockDB) {
			t.Error("DB reference integrity not maintained")
		}
		t.Log("DB reference integrity verified successfully")
	})

	t.Run("Scenario 5: Memory Resource Management", func(t *testing.T) {
		mockDB := &gorm.DB{}
		var stores []*ArticleStore

		for i := 0; i < 1000; i++ {
			stores = append(stores, NewArticleStore(mockDB))
		}

		for i, store := range stores {
			if store == nil {
				t.Errorf("Instance %d is nil", i)
			}
			if store.db != mockDB {
				t.Errorf("Instance %d has incorrect DB reference", i)
			}
		}
		t.Log("Memory resource management verified successfully")
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

func NewArticleStore(db DBInterface) *ArticleStore {
	return &ArticleStore{db: db}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		setupFn func(*MockDB)
		wantErr bool
		errMsg  string
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
			setupFn: func(mockDB *MockDB) {
				db := &gorm.DB{Error: nil}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			wantErr: false,
		},
		{
			name: "Update Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			setupFn: func(mockDB *MockDB) {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Update with Empty Required Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			setupFn: func(mockDB *MockDB) {
				db := &gorm.DB{Error: errors.New("not null constraint violation")}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			wantErr: true,
			errMsg:  "not null constraint violation",
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Valid Title",
			},
			setupFn: func(mockDB *MockDB) {
				db := &gorm.DB{Error: errors.New("database connection failed")}
				mockDB.On("Model", mock.Anything).Return(db)
				mockDB.On("Update", mock.Anything).Return(db)
			},
			wantErr: true,
			errMsg:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupFn(mockDB)

			store := NewArticleStore(mockDB)

			err := store.Update(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Update successful")
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
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully add favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
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
			name: "Error during user association",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(errors.New("association error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("association error"),
		},
		{
			name: "Error during favorites count update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update error"),
		},
		{
			name:    "Nil article",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article"),
		},
		{
			name: "Nil user",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user:          nil,
			mockSetup:     func(mock sqlmock.Sqlmock) {},
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

			tc.mockSetup(mock)

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

	t.Run("Concurrent operations", func(t *testing.T) {
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

		for i := 0; i < numGoroutines; i++ {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `favorite_articles`").
				WithArgs(1, uint(i+1)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE `articles`").
				WithArgs(1, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(userID uint) {
				defer wg.Done()
				user := &model.User{Model: gorm.Model{ID: userID}}
				err := store.AddFavorite(article, user)
				assert.NoError(t, err)
			}(uint(i + 1))
		}

		wg.Wait()
		assert.Equal(t, int32(numGoroutines), article.FavoritesCount)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})
}

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func (m *mockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return &gorm.Association{Error: args.Get(0).(*gorm.DB).Error}
}

func (m *mockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestDeleteFavorite(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*mockDB)
		article     *model.Article
		user        *model.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful deletion",
			setupMock: func(mockDB *mockDB) {
				tx := &gorm.DB{Error: nil}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(tx)
				mockDB.On("Association", "FavoritedUsers").Return(tx)
				mockDB.On("Commit").Return(tx)
			},
			article: &model.Article{
				FavoritesCount: 1,
				FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectError: false,
		},
		{
			name: "Association deletion error",
			setupMock: func(mockDB *mockDB) {
				tx := &gorm.DB{Error: errors.New("association error")}
				mockDB.On("Begin").Return(&gorm.DB{Error: nil})
				mockDB.On("Model", mock.Anything).Return(tx)
				mockDB.On("Association", "FavoritedUsers").Return(tx)
				mockDB.On("Rollback").Return(&gorm.DB{Error: nil})
			},
			article: &model.Article{
				FavoritesCount: 1,
			},
			user:        &model.User{},
			expectError: true,
			errorMsg:    "association error",
		},
		{
			name:        "Nil article",
			article:     nil,
			user:        &model.User{},
			expectError: true,
			errorMsg:    "invalid article",
		},
		{
			name:        "Nil user",
			article:     &model.Article{},
			user:        nil,
			expectError: true,
			errorMsg:    "invalid user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(mockDB)
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			store := &ArticleStore{
				db: mock.db,
			}

			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.article.FavoritesCount-1, tt.article.FavoritesCount)
			}
		})
	}
}

func TestDeleteFavoriteConcurrent(t *testing.T) {
	article := &model.Article{
		FavoritesCount: 5,
		FavoritedUsers: make([]model.User, 5),
	}

	mock := new(mockDB)
	tx := &gorm.DB{Error: nil}
	mock.On("Begin").Return(tx)
	mock.On("Model", mock.Anything).Return(tx)
	mock.On("Association", "FavoritedUsers").Return(tx)
	mock.On("Commit").Return(tx)

	store := &ArticleStore{
		db: mock.db,
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			user := &model.User{Model: gorm.Model{ID: uint(idx)}}
			err := store.DeleteFavorite(article, user)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, int32(0), article.FavoritesCount)
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

	if err := db.Create(testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	articles := []model.Article{
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

	for _, article := range articles {
		if err := db.Create(&article).Error; err != nil {
			t.Fatalf("Failed to create test article: %v", err)
		}
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
			got, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantLen, len(got))

			if tt.username != "" {
				for _, article := range got {
					assert.Equal(t, tt.username, article.Author.Username)
				}
			}

			if tt.tagName != "" {
				for _, article := range got {
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
		})
	}
}

