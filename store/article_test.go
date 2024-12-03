package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"time"
	"reflect"
	"github.com/stretchr/testify/require"
	"sync"
	_ "github.com/mattn/go-sqlite3"
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
			dbError: errors.New("title cannot be null"),
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
			dbError: errors.New("database connection failed"),
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
			name: "Maximum Field Length Article",
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
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Article created successfully")
			}

			mockDB.AssertExpectations(t)
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
		setup   func(*testing.T, *ArticleStore)
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
			setup: func(t *testing.T, store *ArticleStore) {

			},
			wantErr: false,
		},
		{
			name: "Empty Body Comment",
			comment: &model.Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			setup: func(t *testing.T, store *ArticleStore) {

			},
			wantErr: true,
			errMsg:  "not null constraint",
		},
		{
			name: "Non-Existent UserID",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    99999,
				ArticleID: 1,
			},
			setup: func(t *testing.T, store *ArticleStore) {

			},
			wantErr: true,
			errMsg:  "foreign key constraint",
		},
		{
			name: "Non-Existent ArticleID",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 99999,
			},
			setup: func(t *testing.T, store *ArticleStore) {

			},
			wantErr: true,
			errMsg:  "foreign key constraint",
		},
		{
			name: "Maximum Length Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 65535)),
				UserID:    1,
				ArticleID: 1,
			},
			setup: func(t *testing.T, store *ArticleStore) {

			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := setupTestDB(t)
			if err != nil {
				t.Fatalf("failed to setup test database: %v", err)
			}
			defer db.Close()

			store := &ArticleStore{db: db}

			if tt.setup != nil {
				tt.setup(t, store)
			}

			err = store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err != nil && tt.errMsg != "" && !containsError(err.Error(), tt.errMsg) {
					t.Errorf("CreateComment() error message = %v, want %v", err, tt.errMsg)
				}
				return
			}

			var savedComment model.Comment
			if err := db.First(&savedComment, tt.comment.ID).Error; err != nil {
				t.Errorf("Failed to retrieve created comment: %v", err)
				return
			}

			if savedComment.Body != tt.comment.Body {
				t.Errorf("Comment body = %v, want %v", savedComment.Body, tt.comment.Body)
			}
			if savedComment.UserID != tt.comment.UserID {
				t.Errorf("Comment UserID = %v, want %v", savedComment.UserID, tt.comment.UserID)
			}
			if savedComment.ArticleID != tt.comment.ArticleID {
				t.Errorf("Comment ArticleID = %v, want %v", savedComment.ArticleID, tt.comment.ArticleID)
			}
		})
	}
}

func TestCreateCommentConcurrent(t *testing.T) {

}

func containsError(err, msg string) bool {
	return strings.Contains(strings.ToLower(err), strings.ToLower(msg))
}

func setupTestDB(t *testing.T) (*gorm.DB, error) {
	return nil, nil
}

/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func (m *ArticleDeleteMockDB) Delete(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbError error
		wantErr bool
	}{
		{
			name: "Successfully delete existing article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name:    "Attempt to delete nil article",
			article: nil,
			dbError: errors.New("invalid article: nil pointer"),
			wantErr: true,
		},
		{
			name: "Database error during deletion",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
			},
			dbError: errors.New("database connection error"),
			wantErr: true,
		},
		{
			name: "Delete non-existent article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 9999,
				},
			},
			dbError: gorm.ErrRecordNotFound,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(ArticleDeleteMockDB)

			store := &ArticleStore{
				db: &gorm.DB{},
			}

			if tt.article != nil {
				mockDB.On("Delete", tt.article).Return(&gorm.DB{Error: tt.dbError})
			}

			err := store.Delete(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.dbError != nil {
					assert.Equal(t, tt.dbError, err)
				}
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
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
		},
		{
			name: "Delete Non-existent Comment",
			setup: func(db *gorm.DB) *model.Comment {
				return &model.Comment{
					Model: gorm.Model{ID: 99999},
				}
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := setupTestDBForDelete(t)
			assert.NoError(t, err)
			defer db.Close()

			store := &ArticleStore{db: db}
			comment := tt.setup(db)

			err = store.DeleteComment(comment)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				var deletedComment model.Comment
				result := db.Unscoped().First(&deletedComment, comment.ID)
				assert.NoError(t, result.Error)
				assert.NotNil(t, deletedComment.DeletedAt)
			}
		})
	}
}

func TestDeleteCommentConcurrent(t *testing.T) {
	db, err := setupTestDBForDelete(t)
	assert.NoError(t, err)
	defer db.Close()

	store := &ArticleStore{db: db}

	comments := make([]*model.Comment, 3)
	for i := range comments {
		comments[i] = &model.Comment{
			Body:      "Concurrent Test Comment",
			UserID:    uint(i + 1),
			ArticleID: 1,
		}
		db.Create(comments[i])
	}

	errChan := make(chan error, len(comments))
	for _, comment := range comments {
		go func(c *model.Comment) {
			errChan <- store.DeleteComment(c)
		}(comment)
	}

	for i := 0; i < len(comments); i++ {
		err := <-errChan
		assert.NoError(t, err)
	}

	for _, comment := range comments {
		var deletedComment model.Comment
		result := db.Unscoped().First(&deletedComment, comment.ID)
		assert.NoError(t, result.Error)
		assert.NotNil(t, deletedComment.DeletedAt)
	}
}

func setupTestDBForDelete(t *testing.T) (*gorm.DB, error) {

	return gorm.Open("sqlite3", ":memory:")
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func (m *MockArticleDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockArticleDB) Preload(column string) *MockArticleDB {
	args := m.Called(column)
	return args.Get(0).(*MockArticleDB)
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockSetup     func(*MockArticleDB)
		expectedError error
		expectedData  *model.Article
	}{
		{
			name: "Successfully retrieve article",
			id:   1,
			mockSetup: func(m *MockArticleDB) {
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
				m.On("Find", mock.Anything, uint(1)).Return(&gorm.DB{Error: nil}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.Article)
					*arg = *expectedArticle
				})
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
		},
		{
			name: "Article not found",
			id:   99999,
			mockSetup: func(m *MockArticleDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(99999)).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
		{
			name: "Database error",
			id:   2,
			mockSetup: func(m *MockArticleDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(2)).Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			expectedError: errors.New("database connection error"),
			expectedData:  nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			mockSetup: func(m *MockArticleDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.Anything, uint(0)).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockArticleDB)
			tt.mockSetup(mockDB)

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
func (m *MockCommentDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetCommentByID(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		setupMock       func(*MockCommentDB) *gorm.DB
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name: "Successfully retrieve existing comment",
			id:   1,
			setupMock: func(m *MockCommentDB) *gorm.DB {
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
				return &gorm.DB{Error: nil, Value: comment}
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
			setupMock: func(m *MockCommentDB) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Database connection error",
			id:   1,
			setupMock: func(m *MockCommentDB) *gorm.DB {
				return &gorm.DB{Error: errors.New("database connection error")}
			},
			expectedError:   errors.New("database connection error"),
			expectedComment: nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			setupMock: func(m *MockCommentDB) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockCommentDB)
			db := tt.setupMock(mockDB)
			mockDB.On("Find", mock.Anything, mock.Anything).Return(db)

			store := &ArticleStore{
				db: db,
			}

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
			t.Logf("Test completed: %s", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestGetComments(t *testing.T) {
	type testCase struct {
		name     string
		article  *Article
		setup    func(*gorm.DB)
		validate func(*testing.T, []Comment, error)
	}

	tests := []testCase{
		{
			name: "Successfully retrieve comments for an article",
			article: &Article{
				Model: Model{ID: 1},
			},
			setup: func(db *gorm.DB) {
				author := User{Model: Model{ID: 1}, Username: "testuser"}
				db.Create(&author)

				article := Article{
					Model:  Model{ID: 1},
					Title:  "Test Article",
					UserID: author.ID,
				}
				db.Create(&article)

				comments := []Comment{
					{
						Body:      "Comment 1",
						UserID:    author.ID,
						ArticleID: article.ID,
						Author:    author,
					},
					{
						Body:      "Comment 2",
						UserID:    author.ID,
						ArticleID: article.ID,
						Author:    author,
					},
				}
				for _, comment := range comments {
					db.Create(&comment)
				}
			},
			validate: func(t *testing.T, comments []Comment, err error) {
				assert.NoError(t, err)
				assert.Len(t, comments, 2)
				for _, comment := range comments {
					assert.Equal(t, uint(1), comment.ArticleID)
					assert.NotEmpty(t, comment.Author)
				}
			},
		},
		{
			name: "Article with no comments",
			article: &Article{
				Model: Model{ID: 2},
			},
			setup: func(db *gorm.DB) {
				article := Article{
					Model: Model{ID: 2},
					Title: "Article without comments",
				}
				db.Create(&article)
			},
			validate: func(t *testing.T, comments []Comment, err error) {
				assert.NoError(t, err)
				assert.Empty(t, comments)
			},
		},
		{
			name: "Non-existent article",
			article: &Article{
				Model: Model{ID: 999},
			},
			setup: func(db *gorm.DB) {},
			validate: func(t *testing.T, comments []Comment, err error) {
				assert.NoError(t, err)
				assert.Empty(t, comments)
			},
		},
		{
			name: "Invalid article ID",
			article: &Article{
				Model: Model{ID: 0},
			},
			setup: func(db *gorm.DB) {},
			validate: func(t *testing.T, comments []Comment, err error) {
				assert.NoError(t, err)
				assert.Empty(t, comments)
			},
		},
		{
			name: "Verify comment order",
			article: &Article{
				Model: Model{ID: 3},
			},
			setup: func(db *gorm.DB) {
				author := User{Model: Model{ID: 2}, Username: "orderuser"}
				db.Create(&author)

				article := Article{
					Model:  Model{ID: 3},
					Title:  "Article for order test",
					UserID: author.ID,
				}
				db.Create(&article)

				now := time.Now()
				comments := []Comment{
					{
						Model:     Model{CreatedAt: now.Add(-1 * time.Hour)},
						Body:      "Old comment",
						UserID:    author.ID,
						ArticleID: article.ID,
						Author:    author,
					},
					{
						Model:     Model{CreatedAt: now},
						Body:      "New comment",
						UserID:    author.ID,
						ArticleID: article.ID,
						Author:    author,
					},
				}
				for _, comment := range comments {
					db.Create(&comment)
				}
			},
			validate: func(t *testing.T, comments []Comment, err error) {
				assert.NoError(t, err)
				assert.Len(t, comments, 2)
				if len(comments) >= 2 {
					assert.True(t, comments[0].CreatedAt.Before(comments[1].CreatedAt))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, err := setupTestDB()
			if err != nil {
				t.Fatalf("Failed to setup test database: %v", err)
			}
			defer cleanupTestDB(db)

			tc.setup(db)

			store := &ArticleStore{db: db}
			comments, err := store.GetComments(tc.article)

			tc.validate(t, comments, err)
		})
	}
}

func cleanupTestDB(db *gorm.DB) {
}

func setupTestDB() (*gorm.DB, error) {
	return nil, nil
}

/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func (m *ArticleDBMock) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *ArticleDBMock) Limit(limit interface{}) *gorm.DB {
	args := m.Called(limit)
	return args.Get(0).(*gorm.DB)
}

func (m *ArticleDBMock) Offset(offset interface{}) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}

func (m *ArticleDBMock) Preload(column string) *gorm.DB {
	args := m.Called(column)
	return args.Get(0).(*gorm.DB)
}

func TestGetFeedArticles(t *testing.T) {
	tests := []struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(*ArticleDBMock)
		want      []model.Article
		wantErr   error
	}{
		{
			name:    "Successful retrieval - Single user",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(m *ArticleDBMock) {
				db := &gorm.DB{Error: nil}
				m.On("Preload", "Author").Return(db)
				m.On("Where", "user_id in (?)", mock.Anything).Return(db)
				m.On("Offset", int64(0)).Return(db)
				m.On("Limit", int64(10)).Return(db)
				m.On("Find", mock.Anything, mock.Anything).Return(db).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Article)
					*arg = []model.Article{{
						Model:  gorm.Model{ID: 1},
						Title:  "Test Article",
						UserID: 1,
						Author: model.User{Model: gorm.Model{ID: 1}},
					}}
				})
			},
			want: []model.Article{{
				Model:  gorm.Model{ID: 1},
				Title:  "Test Article",
				UserID: 1,
				Author: model.User{Model: gorm.Model{ID: 1}},
			}},
			wantErr: nil,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{999},
			limit:   10,
			offset:  0,
			mockSetup: func(m *ArticleDBMock) {
				db := &gorm.DB{Error: nil}
				m.On("Preload", "Author").Return(db)
				m.On("Where", "user_id in (?)", mock.Anything).Return(db)
				m.On("Offset", int64(0)).Return(db)
				m.On("Limit", int64(10)).Return(db)
				m.On("Find", mock.Anything, mock.Anything).Return(db)
			},
			want:    []model.Article{},
			wantErr: nil,
		},
		{
			name:    "Database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(m *ArticleDBMock) {
				db := &gorm.DB{Error: errors.New("database error")}
				m.On("Preload", "Author").Return(db)
				m.On("Where", "user_id in (?)", mock.Anything).Return(db)
				m.On("Offset", int64(0)).Return(db)
				m.On("Limit", int64(10)).Return(db)
				m.On("Find", mock.Anything, mock.Anything).Return(db)
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(ArticleDBMock)
			tt.mockSetup(mockDB)

			s := &ArticleStore{
				db: mockDB,
			}

			got, err := s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *ArticleDBMock) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func (m *TagDBMock) Find(dest interface{}) *gorm.DB {
	args := m.Called(dest)
	return args.Get(0).(*gorm.DB)
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*TagDBMock)
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Successfully retrieve tags",
			setupMock: func(mockDB *TagDBMock) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Tag)
					*arg = []model.Tag{
						{
							Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
							Name:  "golang",
						},
						{
							Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
							Name:  "testing",
						},
					}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedTags: []model.Tag{
				{
					Model: gorm.Model{ID: 1},
					Name:  "golang",
				},
				{
					Model: gorm.Model{ID: 2},
					Name:  "testing",
				},
			},
			expectedError: nil,
		},
		{
			name: "Empty database returns empty tag list",
			setupMock: func(mockDB *TagDBMock) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Tag)
					*arg = []model.Tag{}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database connection error",
			setupMock: func(mockDB *TagDBMock) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).Return(
					&gorm.DB{Error: errors.New("database connection failed")},
				)
			},
			expectedTags:  nil,
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "Database query timeout",
			setupMock: func(mockDB *TagDBMock) {
				mockDB.On("Find", mock.AnythingOfType("*[]model.Tag")).Return(
					&gorm.DB{Error: errors.New("context deadline exceeded")},
				)
			},
			expectedTags:  nil,
			expectedError: errors.New("context deadline exceeded"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(TagDBMock)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB,
			}

			tags, err := store.GetTags()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				if tt.expectedTags != nil {

					for i := range tt.expectedTags {
						assert.Equal(t, tt.expectedTags[i].ID, tags[i].ID)
						assert.Equal(t, tt.expectedTags[i].Name, tags[i].Name)
					}
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func (m *MockArticleDB) Count(value interface{}) *gorm.DB {
	return m.Called(value).Get(0).(*gorm.DB)
}

func (m *MockArticleDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestIsFavorited(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(*MockArticleDB)
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
			setupMock: func(m *MockArticleDB) {
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
			setupMock: func(m *MockArticleDB) {
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
			setupMock:   func(m *MockArticleDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User Parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			setupMock:   func(m *MockArticleDB) {},
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
			setupMock: func(m *MockArticleDB) {
				db := &gorm.DB{Error: errors.New("database error")}
				m.On("Table", "favorite_articles").Return(db)
				m.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Return(db)
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name: "Zero-Value IDs",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			user: &model.User{
				Model: gorm.Model{ID: 0},
			},
			setupMock: func(m *MockArticleDB) {
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
			mockDB := new(MockArticleDB)
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

func (m *MockArticleDB) Where(query interface{}, args ...interface{}) *gorm.DB {
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
		{
			name: "Scenario 3: Verify DB Reference Integrity",
			db: &gorm.DB{
				Error: nil,
			},
			wantNil:  false,
			scenario: "Ensure DB reference matches input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test scenario: %s", tt.scenario)

			got := NewArticleStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewArticleStore() nil check = %v, want %v", got == nil, tt.wantNil)
				return
			}

			if !tt.wantNil && got != nil {
				if !reflect.DeepEqual(got.db, tt.db) {
					t.Errorf("NewArticleStore() db reference mismatch = %v, want %v", got.db, tt.db)
				}
			}

			switch tt.name {
			case "Scenario 3: Verify DB Reference Integrity":
				if got.db != tt.db {
					t.Error("DB reference integrity check failed: references don't match")
				}
			}

			t.Logf("Test scenario completed successfully: %s", tt.scenario)
		})
	}

	t.Run("Scenario 4: Multiple ArticleStore Instances Independence", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}

		store1 := NewArticleStore(db1)
		store2 := NewArticleStore(db2)

		if store1.db == store2.db {
			t.Error("Multiple instances should maintain independent DB references")
		}

		t.Log("Successfully verified instance independence")
	})

	t.Run("Scenario 5: Type Safety and Interface Compliance", func(t *testing.T) {
		store := NewArticleStore(&gorm.DB{})

		if _, ok := interface{}(store).(*ArticleStore); !ok {
			t.Error("ArticleStore does not implement expected type")
		}

		t.Log("Successfully verified type safety")
	})
}

/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe


 */
func (m *ArticleDBMock) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbError error
		wantErr bool
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
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Empty Required Fields",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "",
			},
			dbError: errors.New("not null constraint violation"),
			wantErr: true,
		},
		{
			name: "Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
			},
			dbError: gorm.ErrRecordNotFound,
			wantErr: true,
		},
		{
			name: "Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Title",
			},
			dbError: errors.New("database connection failed"),
			wantErr: true,
		},
		{
			name: "Maximum Field Values",
			article: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 5000)),
				UserID:      1,
			},
			dbError: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(ArticleDBMock)
			store := &ArticleStore{
				db: mockDB,
			}

			mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
			mockDB.On("Update", mock.Anything).Return(&gorm.DB{Error: tt.dbError})

			err := store.Update(tt.article)

			if tt.wantErr {
				require.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				require.NoError(t, err)
				t.Log("Update successful")
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *ArticleDBMock) Update(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90


 */
func (m *MockArticleFavoriteDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.Association)
}

func (m *MockArticleFavoriteDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockArticleFavoriteDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockArticleFavoriteDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockArticleFavoriteDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestAddFavorite(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(*MockArticleFavoriteDB)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Addition of Favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mockDB *MockArticleFavoriteDB) {
				mockDB.On("Begin").Return(&gorm.DB{})
				mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{})
				mockDB.On("Commit").Return(&gorm.DB{})
			},
			expectError: false,
		},
		{
			name:    "Null Article Parameter",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectError: true,
			errorMsg:    "article cannot be nil",
		},
		{
			name: "Null User Parameter",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user:        nil,
			expectError: true,
			errorMsg:    "user cannot be nil",
		},
		{
			name: "Association Error",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mockDB *MockArticleFavoriteDB) {
				mockDB.On("Begin").Return(&gorm.DB{})
				mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
				mockDB.On("Association", "FavoritedUsers").Return(nil)
				mockDB.On("Rollback").Return(&gorm.DB{})
			},
			expectError: true,
			errorMsg:    "association error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockArticleFavoriteDB)
			if tt.setupMock != nil {
				tt.setupMock(mockDB)
			}

			store := &ArticleStore{
				db: mockDB,
			}

			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.article != nil {
					assert.Equal(t, int32(1), tt.article.FavoritesCount)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestAddFavoriteConcurrent(t *testing.T) {
	article := &model.Article{
		Model:          gorm.Model{ID: 1},
		FavoritesCount: 0,
	}

	users := make([]*model.User, 5)
	for i := range users {
		users[i] = &model.User{
			Model: gorm.Model{ID: uint(i + 1)},
		}
	}

	mockDB := new(MockArticleFavoriteDB)
	store := &ArticleStore{
		db: mockDB,
	}

	mockDB.On("Begin").Return(&gorm.DB{})
	mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
	mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
	mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{})
	mockDB.On("Commit").Return(&gorm.DB{})

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u *model.User) {
			defer wg.Done()
			err := store.AddFavorite(article, u)
			assert.NoError(t, err)
		}(user)
	}
	wg.Wait()

	assert.Equal(t, int32(5), article.FavoritesCount)
	mockDB.AssertExpectations(t)
}

func (m *MockArticleFavoriteDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func (m *MockDeleteFavoriteDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.Association)
}

func (m *MockDeleteFavoriteDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDeleteFavoriteDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDeleteFavoriteDB) Delete(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDeleteFavoriteDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDeleteFavoriteDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestDeleteFavorite(t *testing.T) {
	tests := []struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(*MockDeleteFavoriteDB)
		expectedError error
		expectedCount int32
	}{
		{
			name: "Successful deletion of favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mockDB *MockDeleteFavoriteDB) {
				tx := &gorm.DB{}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(tx)
				mockDB.On("Association", "FavoritedUsers").Return(nil)
				mockDB.On("Delete", mock.Anything).Return(nil)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(tx)
				mockDB.On("Commit").Return(tx)
			},
			expectedError: nil,
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
			setupMock: func(mockDB *MockDeleteFavoriteDB) {
				tx := &gorm.DB{}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(tx)
				mockDB.On("Association", "FavoritedUsers").Return(errors.New("association error"))
				mockDB.On("Rollback").Return(tx)
			},
			expectedError: errors.New("association error"),
			expectedCount: 2,
		},
		{
			name: "Failed favorites count update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mockDB *MockDeleteFavoriteDB) {
				tx := &gorm.DB{}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(tx)
				mockDB.On("Association", "FavoritedUsers").Return(nil)
				mockDB.On("Delete", mock.Anything).Return(nil)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(errors.New("update error"))
				mockDB.On("Rollback").Return(tx)
			},
			expectedError: errors.New("update error"),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDeleteFavoriteDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB,
			}

			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			mockDB.AssertExpectations(t)
		})
	}
}

func TestDeleteFavoriteConcurrent(t *testing.T) {
	article := &model.Article{
		Model:          gorm.Model{ID: 1},
		FavoritesCount: 10,
	}

	mockDB := new(MockDeleteFavoriteDB)
	store := &ArticleStore{
		db: mockDB,
	}

	tx := &gorm.DB{}
	mockDB.On("Begin").Return(tx)
	mockDB.On("Model", mock.Anything).Return(tx)
	mockDB.On("Association", "FavoritedUsers").Return(nil)
	mockDB.On("Delete", mock.Anything).Return(nil)
	mockDB.On("Update", "favorites_count", mock.Anything).Return(tx)
	mockDB.On("Commit").Return(tx)

	numOperations := 5
	var wg sync.WaitGroup
	wg.Add(numOperations)

	for i := 0; i < numOperations; i++ {
		go func(userID uint) {
			defer wg.Done()
			user := &model.User{Model: gorm.Model{ID: userID}}
			err := store.DeleteFavorite(article, user)
			assert.NoError(t, err)
		}(uint(i + 1))
	}

	wg.Wait()

	assert.Equal(t, int32(5), article.FavoritesCount)
	mockDB.AssertExpectations(t)
}

func (m *MockDeleteFavoriteDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestGetArticles(t *testing.T) {
	type testCase struct {
		name        string
		tagName     string
		username    string
		favoritedBy *User
		limit       int64
		offset      int64
		setupFunc   func(*gorm.DB)
		wantCount   int
		wantErr     bool
	}

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	store := &ArticleStore{db: db}

	db.AutoMigrate(&User{}, &Article{}, &Tag{})

	createTestData := func(db *gorm.DB) (*User, []Article) {
		user := &User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
		}
		db.Create(user)

		articles := []Article{
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "Test Article 1",
				Description: "Test Description 1",
				Body:        "Test Body 1",
				UserID:      user.ID,
				Author:      *user,
			},
			{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "Test Article 2",
				Description: "Test Description 2",
				Body:        "Test Body 2",
				UserID:      user.ID,
				Author:      *user,
			},
		}

		for _, article := range articles {
			db.Create(&article)
		}

		return user, articles
	}

	tests := []testCase{
		{
			name:   "Get All Articles Without Filters",
			limit:  10,
			offset: 0,
			setupFunc: func(db *gorm.DB) {
				_, _ = createTestData(db)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:     "Get Articles By Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			setupFunc: func(db *gorm.DB) {
				_, _ = createTestData(db)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:     "Get Articles With Invalid Username",
			username: "nonexistent",
			limit:    10,
			offset:   0,
			setupFunc: func(db *gorm.DB) {
				_, _ = createTestData(db)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "Test Pagination",
			limit:  1,
			offset: 1,
			setupFunc: func(db *gorm.DB) {
				_, _ = createTestData(db)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:   "Get Favorited Articles",
			limit:  10,
			offset: 0,
			setupFunc: func(db *gorm.DB) {
				user, articles := createTestData(db)
				db.Exec("INSERT INTO favorite_articles (user_id, article_id) VALUES (?, ?)", user.ID, articles[0].ID)
			},
			favoritedBy: &User{Model: gorm.Model{ID: 1}},
			wantCount:   1,
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM articles")
			db.Exec("DELETE FROM users")
			db.Exec("DELETE FROM favorite_articles")

			if tc.setupFunc != nil {
				tc.setupFunc(db)
			}

			articles, err := store.GetArticles(tc.tagName, tc.username, tc.favoritedBy, tc.limit, tc.offset)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantCount, len(articles))
			}

			t.Logf("Test: %s - Got %d articles", tc.name, len(articles))
		})
	}
}

