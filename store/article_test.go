package store

import (
	"errors"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"database/sql"
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
func (m *MockDB) Delete(value interface{}) *gorm.DB {
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

			mockDB := new(MockDB)

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

			db, err := setupTestDB()
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
	db, err := setupTestDB()
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

func setupTestDB() (*gorm.DB, error) {

	return gorm.Open("sqlite3", ":memory:")
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(column string) *gorm.DB {
	args := m.Called(column)
	return args.Get(0).(*gorm.DB)
}

func TestGetByID(t *testing.T) {

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(*MockDB)
		expectedError error
		expectedData  *model.Article
	}{
		{
			name: "Successfully retrieve article",
			id:   1,
			mockSetup: func(m *MockDB) {
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
			mockSetup: func(m *MockDB) {
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
			mockSetup: func(m *MockDB) {
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
			mockSetup: func(m *MockDB) {
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

			mockDB := new(MockDB)
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
			setupMock: func(m *MockDB) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Database connection error",
			id:   1,
			setupMock: func(m *MockDB) *gorm.DB {
				return &gorm.DB{Error: errors.New("database connection error")}
			},
			expectedError:   errors.New("database connection error"),
			expectedComment: nil,
		},
		{
			name: "Zero ID input",
			id:   0,
			setupMock: func(m *MockDB) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)

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

