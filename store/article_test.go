Here's the corrected and complete test file content with the ROOST_METHOD_HASH and ROOST_METHOD_SIG_HASH comments:

```go
package store

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDB struct {
	mock.Mock
}

type MockAssociation struct {
	mock.Mock
	Error error
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Joins(query string, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Limit(limit interface{}) *gorm.DB {
	args := m.Called(limit)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Offset(offset interface{}) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rows() (*sql.Rows, error) {
	args := m.Called()
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *MockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func (m *MockAssociation) Append(values ...interface{}) error {
	args := m.Called(values...)
	return args.Error(0)
}

func (m *MockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92
*/
func TestNewArticleStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *ArticleStore
	}{
		{
			name: "Create ArticleStore with Valid DB Connection",
			db:   &gorm.DB{},
			want: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name: "Create ArticleStore with Nil DB Connection",
			db:   nil,
			want: &ArticleStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleStore(tt.db)

			if got == nil {
				t.Errorf("NewArticleStore() returned nil, want non-nil ArticleStore")
			}

			if !reflect.DeepEqual(got.db, tt.want.db) {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.want.db)
			}

			got2 := NewArticleStore(tt.db)
			if got == got2 {
				t.Errorf("NewArticleStore() returned same instance, want different instances")
			}

			if reflect.TypeOf(got) != reflect.TypeOf(&ArticleStore{}) {
				t.Errorf("NewArticleStore() returned %T, want *ArticleStore", got)
			}
		})
	}

	t.Run("Check ArticleStore with Different DB Connections", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}

		got1 := NewArticleStore(db1)
		got2 := NewArticleStore(db2)

		if got1.db == got2.db {
			t.Errorf("NewArticleStore() with different DB connections returned same DB instance")
		}
	})
}

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377
*/
func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbError error
		wantErr bool
	}{
		{
			name: "Successfully Create a New Article",
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
			name: "Attempt to Create an Article with Missing Required Fields",
			article: &model.Article{
				UserID: 1,
			},
			dbError: errors.New("missing required fields"),
			wantErr: true,
		},
		{
			name: "Handle Database Connection Error During Article Creation",
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
			name: "Create Article with Associated Tags",
			article: &model.Article{
				Title:       "Test Article with Tags",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Create Article with Very Long Content",
			article: &model.Article{
				Title:       string(make([]byte, 1000)),
				Description: string(make([]byte, 5000)),
				Body:        string(make([]byte, 10000)),
				UserID:      1,
			},
			dbError: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockDB.On("Create", mock.AnythingOfType("*model.Article")).Return(&gorm.DB{Error: tt.dbError})

			store := &ArticleStore{
				db: mockDB,
			}

			err := store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.dbError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertCalled(t, "Create", tt.article)
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
		setup   func(*gorm.DB)
		comment *model.Comment
		wantErr bool
	}{
		{
			name: "Successfully Create a Comment",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			comment: &model.Comment{
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
		},
		{
			name: "Create Comment with Maximum Length Body",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent Article",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
			},
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 9999,
			},
			wantErr: true,
		},
		{
			name: "Create Comment with Special Characters in Body",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			comment: &model.Comment{
				Body:      "Test comment with special characters: !@#$%^&*()_+ and emojis: 😀🎉",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent User",
			setup: func(db *gorm.DB) {
				db.Create(&model.Article{ID: 1})
			},
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    9999,
				ArticleID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			err = db.AutoMigrate(&model.Comment{}, &model.User{}, &model.Article{}).Error
			require.NoError(t, err)

			tt.setup(db)

			store := &ArticleStore{db: db}

			err = store.CreateComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var createdComment model.Comment
				err = db.First(&createdComment, "body = ?", tt.comment.Body).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.comment.Body, createdComment.Body)
				assert.Equal(t, tt.comment.UserID, createdComment.UserID)
				assert.Equal(t, tt.comment.ArticleID, createdComment.ArticleID)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1
*/
func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbSetup func(*MockDB)
		wantErr bool
	}{
		{
			name: "Successfully Delete an Existing Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			dbSetup: func(mockDB *MockDB) {
				mockDB.On("Delete", mock.Anything).Return(&gorm.DB{Error: nil})
			},
			wantErr: false,
		},
		{
			name: "Attempt to Delete a Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			dbSetup: func(mockDB *MockDB) {
				mockDB.On("Delete", mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			wantErr: true,
		},
		{
			name: "Database Connection Error During Deletion",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
				Title: "Error Article",
			},
			dbSetup: func(mockDB *MockDB) {
				mockDB.On("Delete", mock.Anything).Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.dbSetup(mockDB)

			s := &ArticleStore{
				db: mockDB,
			}

			err := s.Delete(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
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
		comment *model.Comment
		dbError error
		wantErr bool
	}{
		{
			name: "Successfully Delete an Existing Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Delete a Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 999},
				Body:  "Non-existent comment",
			},
			dbError: gorm.ErrRecordNotFound,
			wantErr: true,
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 2},
				Body:  "Another test comment",
			},
			dbError: errors.New("database connection error"),
			wantErr: true,
		},
		{
			name:    "Delete Comment with Null Input",
			comment: nil,
			dbError: errors.New("invalid input: comment is nil"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			store := &ArticleStore{db: mockDB}

			if tt.comment != nil {
				mockDB.On("Delete", tt.comment).Return(&gorm.DB{Error: tt.dbError})
			}

			err := store.DeleteComment(tt.comment)

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
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b
*/
func TestGetCommentByID(t *testing.T) {
	tests := []struct {
		name            string
		commentID       uint
		setupMock       func(*MockDB)
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name:      "Successfully retrieve an existing comment",
			commentID: 1,
			setupMock: func(m *MockDB) {
				m.On("Find", mock.AnythingOfType("*model.Comment"), uint(1)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.Comment)
					*arg = model.Comment{
						Model:     gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Body:      "Test comment",
						UserID:    1,
						ArticleID: 1,
					}
				}).Return(&gorm.DB{Error: nil})
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
			name:      "Attempt to retrieve a non-existent comment",
			commentID: 999,
			setupMock: func(m *MockDB) {
				m.On("Find", mock.AnythingOfType("*model.Comment"), uint(999)).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name:      "Handle database connection error",
			commentID: 2,
			setupMock: func(m *MockDB) {
				m.On("Find", mock.AnythingOfType("*model.Comment"), uint(2)).Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			expectedError:   errors.New("database connection error"),
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{db: mockDB}

			comment, err := store.GetCommentByID(tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedComment != nil {
				assert.Equal(t, tt.expectedComment.ID, comment.ID)
				assert.Equal(t, tt.expectedComment.Body, comment.Body)
				assert.Equal(t, tt.expectedComment.UserID, comment.UserID)
				assert.Equal(t, tt.expectedComment.ArticleID, comment.ArticleID)
			} else {
				assert.Nil(t, comment)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0
*/
func TestGetTags(t *testing.T) {
	tests := []struct {
		name    string
		dbSetup func(*gorm.DB)
		want    []model.Tag
		wantErr bool
	}{
		{
			name: "Successfully retrieve all tags",
			dbSetup: func(db *gorm.DB) {
				tags := []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "tag1"},
					{Model: gorm.Model{ID: 2}, Name: "tag2"},
					{Model: gorm.Model{ID: 3}, Name: "tag3"},
				}
				db.Create(&tags)
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
				{Model: gorm.Model{ID: 3}, Name: "tag3"},
			},
			wantErr: false,
		},
		{
			name:    "Empty tag list",
			dbSetup: func(db *gorm.DB) {},
			want:    []model.Tag{},
			wantErr: false,
		},
		{
			name: "Database connection error",
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Tag{})

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			got, err := s.GetTags()

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52
*/
func TestGetByID(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		mockSetup       func(*MockDB)
		expectedError   error
		expectedArticle *model.Article
	}{
		{
			name: "Successfully retrieve an existing article",
			id:   1,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.AnythingOfType("*model.Article"), uint(1)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.Article)
					*arg = model.Article{
						Model:       gorm.Model{ID: 1},
						Title:       "Test Article",
						Description: "Test Description",
						Body:        "Test Body",
						UserID:      1,
						Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
					}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedError: nil,
			expectedArticle: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
			},
		},
		{
			name: "Attempt to retrieve a non-existent article",
			id:   999,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Tags").Return(m)
				m.On("Preload", "Author").Return(m)
				m.On("Find", mock.AnythingOfType("*model.Article"), uint(999)).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedArticle: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB}
			article, err := store.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedArticle, article)
			mockDB.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe
*/
func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		setupDB     func(*gorm.DB)
		input       *model.Article
		expectedErr error
		validate    func(*testing.T, *gorm.DB, *model.Article)
	}{
		{
			name: "Successfully Update an Existing Article",
			setupDB: func(db *gorm.DB) {
				article := &model.Article{
					Model:       gorm.Model{ID: 1},
					Title:       "Original Title",
					Description: "Original Description",
					Body:        "Original Body",
					UserID:      1,
				}
				db.Create(article)
			},
			input: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Updated Title",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      1,
			},
			expectedErr: nil,
			validate: func(t *testing.T, db *gorm.DB, a *model.Article) {
				var updatedArticle model.Article
				err := db.First(&updatedArticle, a.ID).Error
				require.NoError(t, err)
				assert.Equal(t, "Updated Title", updatedArticle.Title)
				assert.Equal(t, "Updated Description", updatedArticle.Description)
				assert.Equal(t, "Updated Body", updatedArticle.Body)
			},
		},
		{
			name:    "Attempt to Update a Non-existent Article",
			setupDB: func(db *gorm.DB) {},
			input: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			expectedErr: gorm.ErrRecordNotFound,
			validate:    func(t *testing.T, db *gorm.DB, a *model.Article) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.Tag{}, &model.User{})

			tt.setupDB(db)

			store := &ArticleStore{db: db}
			err = store.Update(tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, db, tt.input)
		})
	}
}

/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e
*/
func TestGetComments(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		mockSetup      func(*MockDB)
		expectedResult []model.Comment
		expectedError  error
	}{
		{
			name: "Successfully retrieve comments for an article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(1)).Return(mockDB)
				mockDB.On("Find", mock.AnythingOfType("*[]model.Comment")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Comment)
					*arg = []model.Comment{
						{Model: gorm.Model{ID: 1}, Body: "Comment 1", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
						{Model: gorm.Model{ID: 2}, Body: "Comment 2", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
					}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 1}, Body: "Comment 1", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Body: "Comment 2", ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			expectedError: nil,
		},
		{
			name: "No comments for the article",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			mockSetup: func(mockDB *Mock