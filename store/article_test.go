package store

import (
		"reflect"
		"testing"
		"github.com/jinzhu/gorm"
		"errors"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/mock"
		"github.com/DATA-DOG/go-sqlmock"
		"github.com/stretchr/testify/require"
		"time"
		"fmt"
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

type T struct {
	common
	isEnvSet bool
	context  *testContext
}T struct {
	common
	isEnvSet bool
	context  *testContext
}
type Call struct {
	Parent *Mock

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// Holds the caller info for the On() call
	callerInfo []string

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Amount of times this call has been called
	totalCalls int

	// Call to this method can be optional
	optional bool

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	RunFn func(Arguments)

	// PanicMsg holds msg to be used to mock panic on the function call
	//  if the PanicMsg is set to a non nil string the function call will panic
	// irrespective of other settings
	PanicMsg *string

	// Calls which must be satisfied before this call can be
	requires []*Call
}

type Call struct {
	Parent *Mock

	Method string

	Arguments Arguments

	ReturnArguments Arguments

	callerInfo []string

	Repeatability int

	totalCalls int

	optional bool

	WaitFor <-chan time.Time

	waitTime time.Duration

	RunFn func(Arguments)
}
type MockDB struct {
	mock.Mock
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

type Comment struct {
	gorm.Model
	Body      string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Author    User   `gorm:"foreignkey:UserID"`
	ArticleID uint   `gorm:"not null"`
	Article   Article
}

type DBInterface interface {
	Preload(column string, conditions ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
}
type Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}

type Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}
type MockDB struct {
	mock.Mock
	Error error
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

			if _, ok := interface{}(got).(*ArticleStore); !ok {
				t.Errorf("NewArticleStore() returned type %T, want *ArticleStore", got)
			}
		})
	}

	t.Run("Create Multiple ArticleStores with Different DB Connections", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}

		store1 := NewArticleStore(db1)
		store2 := NewArticleStore(db2)

		if store1.db != db1 {
			t.Errorf("First ArticleStore has incorrect DB connection")
		}

		if store2.db != db2 {
			t.Errorf("Second ArticleStore has incorrect DB connection")
		}

		if store1.db == store2.db {
			t.Errorf("ArticleStores should have different DB connections")
		}
	})
}


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
			name: "Create an Article with Maximum Length Content",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 65535)),
				UserID:      1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Create an Article with Associated Tags",
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
			name: "Attempt to Create an Article with Database Connection Error",
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
			name: "Create an Article with Duplicate Title",
			article: &model.Article{
				Title:       "Duplicate Title",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbError: errors.New("duplicate title"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockDB.On("Create", mock.AnythingOfType("*model.Article")).Return(&gorm.DB{Error: tt.dbError})

			store := &ArticleStore{
				db: &MockGormDB{MockDB: mockDB},
			}

			err := store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.dbError, err)
			} else {
				assert.NoError(t, err)
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
		dbSetup func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successfully Create a Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs("Test comment", 1, 1, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			comment: &model.Comment{
				UserID:    1,
				ArticleID: 1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs("", 1, 1, sqlmock.AnyArg()).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Create Comment with Maximum Length Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
				UserID:    1,
				ArticleID: 1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs(string(make([]byte, 1000)), 1, 1, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent Article",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 9999,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs("Test comment", 1, 9999, sqlmock.AnyArg()).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Create Multiple Comments in Quick Succession",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WithArgs("Test comment", 1, 1, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Database Connection Issues",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			require.NoError(t, err)

			tt.dbSetup(mock)

			s := &ArticleStore{db: gormDB}

			err = s.CreateComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.NotZero(t, tt.comment.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
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
		setup   func(*gorm.DB) *model.Article
		wantErr bool
	}{
		{
			name: "Successfully Delete an Existing Article",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Description: "Test Description", Body: "Test Body"}
				db.Create(article)
				return article
			},
			wantErr: false,
		},
		{
			name: "Attempt to Delete a Non-existent Article",
			setup: func(db *gorm.DB) *model.Article {
				return &model.Article{Model: gorm.Model{ID: 9999}}
			},
			wantErr: true,
		},
		{
			name: "Delete an Article with Associated Tags",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Description: "Test Description", Body: "Test Body"}
				tags := []model.Tag{{Name: "tag1"}, {Name: "tag2"}}
				article.Tags = tags
				db.Create(article)
				return article
			},
			wantErr: false,
		},
		{
			name: "Delete an Article with Comments",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Description: "Test Description", Body: "Test Body"}
				comments := []model.Comment{{Body: "Comment 1"}, {Body: "Comment 2"}}
				article.Comments = comments
				db.Create(article)
				return article
			},
			wantErr: false,
		},
		{
			name: "Delete an Article with Favorited Users",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Description: "Test Description", Body: "Test Body"}
				users := []model.User{{Username: "user1"}, {Username: "user2"}}
				article.FavoritedUsers = users
				db.Create(article)
				return article
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error During Deletion",
			setup: func(db *gorm.DB) *model.Article {

				db.Close()
				return &model.Article{Model: gorm.Model{ID: 1}}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			assert.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.Tag{}, &model.User{}, &model.Comment{})

			store := &ArticleStore{db: db}

			article := tt.setup(db)

			err = store.Delete(article)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var count int
				db.Model(&model.Article{}).Where("id = ?", article.ID).Count(&count)
				assert.Equal(t, 0, count)

				db.Model(&model.Tag{}).Where("id IN (?)", article.Tags).Count(&count)
				assert.Equal(t, 0, count)

				db.Model(&model.Comment{}).Where("article_id = ?", article.ID).Count(&count)
				assert.Equal(t, 0, count)

				db.Table("favorite_articles").Where("article_id = ?", article.ID).Count(&count)
				assert.Equal(t, 0, count)
			}
		})
	}
}

func (m *MockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func (m *MockDB) Delete(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

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
			name:    "Delete Comment with Null Comment Pointer",
			comment: nil,
			dbError: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			if tt.comment != nil {
				mockDB.On("Delete", tt.comment).Return(&gorm.DB{Error: tt.dbError})
			}

			store := &ArticleStore{
				db: mockDB,
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
		setupDB         func(*gorm.DB)
		commentID       uint
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name: "Successfully retrieve an existing comment",
			setupDB: func(db *gorm.DB) {
				comment := model.Comment{
					Model:     gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(&comment)
			},
			commentID:     1,
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name: "Attempt to retrieve a non-existent comment",
			setupDB: func(db *gorm.DB) {

			},
			commentID:       999,
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Handle database error",
			setupDB: func(db *gorm.DB) {

				db.AddError(errors.New("database error"))
			},
			commentID:       1,
			expectedError:   errors.New("database error"),
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			assert.NoError(t, err)
			defer db.Close()

			tt.setupDB(db)

			store := &ArticleStore{db: db}

			comment, err := store.GetCommentByID(tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedComment != nil {
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedComment.ID, comment.ID)
				assert.Equal(t, tt.expectedComment.Body, comment.Body)
				assert.Equal(t, tt.expectedComment.UserID, comment.UserID)
				assert.Equal(t, tt.expectedComment.ArticleID, comment.ArticleID)
			} else {
				assert.Nil(t, comment)
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
		name    string
		dbSetup func(*gorm.DB)
		want    []model.Tag
		wantErr bool
	}{
		{
			name: "Successfully retrieve all tags",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag2"})
				db.Create(&model.Tag{Name: "tag3"})
			},
			want: []model.Tag{
				{Name: "tag1"},
				{Name: "tag2"},
				{Name: "tag3"},
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
		{
			name: "Large number of tags",
			dbSetup: func(db *gorm.DB) {
				for i := 0; i < 10000; i++ {
					db.Create(&model.Tag{Name: fmt.Sprintf("tag%d", i)})
				}
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate tag names",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "duplicate"})
				db.Create(&model.Tag{Name: "duplicate"})
				db.Create(&model.Tag{Name: "unique"})
			},
			want: []model.Tag{
				{Name: "duplicate"},
				{Name: "duplicate"},
				{Name: "unique"},
			},
			wantErr: false,
		},
		{
			name: "Tags with special characters",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "特殊字符"})
				db.Create(&model.Tag{Name: "tag with spaces"})
				db.Create(&model.Tag{Name: "tag-with-hyphens"})
				db.Create(&model.Tag{Name: "tag_with_underscores"})
			},
			want: []model.Tag{
				{Name: "特殊字符"},
				{Name: "tag with spaces"},
				{Name: "tag-with-hyphens"},
				{Name: "tag_with_underscores"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("failed to open database: %v", err)
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

			if tt.name == "Large number of tags" {
				if len(got) != 10000 {
					t.Errorf("ArticleStore.GetTags() returned %d tags, want 10000", len(got))
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		mockSetup       func(*MockDB)
		expectedError   error
		expectedArticle *model.Article
	}{
		{
			name: "Successfully retrieve an existing article by ID",
			id:   1,
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Tags").Return(mockDB)
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Find", mock.AnythingOfType("*model.Article"), uint(1)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.Article)
					*arg = model.Article{
						Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Test Article",
						Description: "Test Description",
						Body:        "Test Body",
						Tags:        []model.Tag{{Model: gorm.Model{ID: 1}, Name: "test"}},
						Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
						UserID:      1,
					}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedError: nil,
			expectedArticle: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				Tags:        []model.Tag{{Model: gorm.Model{ID: 1}, Name: "test"}},
				Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:      1,
			},
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
		expectError bool
		validate    func(*testing.T, *gorm.DB, *model.Article)
	}{
		{
			name: "Successfully Update an Existing Article",
			setupDB: func(db *gorm.DB) {
				article := &model.Article{
					Title:       "Original Title",
					Description: "Original Description",
					Body:        "Original Body",
					UserID:      1,
				}
				require.NoError(t, db.Create(article).Error)
			},
			input: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Updated Title",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      1,
			},
			expectError: false,
			validate: func(t *testing.T, db *gorm.DB, input *model.Article) {
				var updatedArticle model.Article
				err := db.First(&updatedArticle, input.ID).Error
				require.NoError(t, err)
				assert.Equal(t, input.Title, updatedArticle.Title)
				assert.Equal(t, input.Description, updatedArticle.Description)
				assert.Equal(t, input.Body, updatedArticle.Body)
			},
		},
		{
			name:    "Attempt to Update a Non-existent Article",
			setupDB: func(db *gorm.DB) {},
			input: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			expectError: true,
			validate: func(t *testing.T, db *gorm.DB, input *model.Article) {
				var article model.Article
				err := db.First(&article, input.ID).Error
				assert.Error(t, err)
				assert.True(t, gorm.IsRecordNotFoundError(err))
			},
		},

		{
			name: "Update Article with No Changes",
			setupDB: func(db *gorm.DB) {
				article := &model.Article{
					Title:       "Unchanged Title",
					Description: "Unchanged Description",
					Body:        "Unchanged Body",
					UserID:      1,
				}
				require.NoError(t, db.Create(article).Error)
			},
			input: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Unchanged Title",
				Description: "Unchanged Description",
				Body:        "Unchanged Body",
				UserID:      1,
			},
			expectError: false,
			validate: func(t *testing.T, db *gorm.DB, input *model.Article) {
				var article model.Article
				err := db.First(&article, input.ID).Error
				require.NoError(t, err)
				assert.Equal(t, input.Title, article.Title)
				assert.Equal(t, input.Description, article.Description)
				assert.Equal(t, input.Body, article.Body)
			},
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

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, db, tt.input)
		})
	}
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
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
			name: "Successful retrieval of comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(1)).Return(mockDB)
				mockDB.On("Find", mock.AnythingOfType("*[]model.Comment"), mock.Anything).Run(func(args mock.Arguments) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			dbWrapper := struct {
				*MockDB
				*gorm.DB
			}{
				MockDB: mockDB,
			}

			store := &ArticleStore{db: &dbWrapper}

			result, err := store.GetComments(tt.article)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func (m *MockDB) Count(value interface{}) *MockDB {
	args := m.Called(value)
	return args.Get(0).(*MockDB)
}

func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Table(name string) *MockDB {
	args := m.Called(name)
	return args.Get(0).(*MockDB)
}

func TestIsFavorited(t *testing.T) {
	tests := []struct {
		name           string
		setupMockDB    func(*MockDB)
		article        *model.Article
		user           *model.User
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Article is favorited by the user",
			setupMockDB: func(db *MockDB) {
				db.On("Table", "favorite_articles").Return(db)
				db.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					*args.Get(0).(*int) = 1
				}).Return(db)
				db.On("Error").Return(nil)
			},
			article:        &model.Article{Model: model.Model{ID: 1}},
			user:           &model.User{Model: model.Model{ID: 1}},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Article is not favorited by the user",
			setupMockDB: func(db *MockDB) {
				db.On("Table", "favorite_articles").Return(db)
				db.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					*args.Get(0).(*int) = 0
				}).Return(db)
				db.On("Error").Return(nil)
			},
			article:        &model.Article{Model: model.Model{ID: 1}},
			user:           &model.User{Model: model.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "Nil Article parameter",
			setupMockDB:    func(db *MockDB) {},
			article:        nil,
			user:           &model.User{Model: model.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "Nil User parameter",
			setupMockDB:    func(db *MockDB) {},
			article:        &model.Article{Model: model.Model{ID: 1}},
			user:           nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Database error",
			setupMockDB: func(db *MockDB) {
				db.On("Table", "favorite_articles").Return(db)
				db.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Return(db)
				db.On("Error").Return(errors.New("database error"))
			},
			article:        &model.Article{Model: model.Model{ID: 1}},
			user:           &model.User{Model: model.Model{ID: 1}},
			expectedResult: false,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Multiple favorites for the same article and user",
			setupMockDB: func(db *MockDB) {
				db.On("Table", "favorite_articles").Return(db)
				db.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					*args.Get(0).(*int) = 3
				}).Return(db)
				db.On("Error").Return(nil)
			},
			article:        &model.Article{Model: model.Model{ID: 1}},
			user:           &model.User{Model: model.Model{ID: 1}},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Zero count",
			setupMockDB: func(db *MockDB) {
				db.On("Table", "favorite_articles").Return(db)
				db.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					*args.Get(0).(*int) = 0
				}).Return(db)
				db.On("Error").Return(nil)
			},
			article:        &model.Article{Model: model.Model{ID: 1}},
			user:           &model.User{Model: model.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			tt.setupMockDB(mockDB)

			store := &ArticleStore{db: mockDB}

			result, err := store.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	callArgs := append([]interface{}{query}, args...)
	returnArgs := m.Called(callArgs...)
	return returnArgs.Get(0).(*MockDB)
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func (m *MockDB) Limit(limit interface{}) *gorm.DB {
	args := m.Called(limit)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Offset(offset interface{}) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}

func TestGetFeedArticles(t *testing.T) {
	tests := []struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(*MockDB)
		expected  []model.Article
		wantErr   bool
	}{
		{
			name:    "Successful Retrieval of Feed Articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m)
				m.On("Where", "user_id in (?)", mock.Anything).Return(m)
				m.On("Offset", int64(0)).Return(m)
				m.On("Limit", int64(10)).Return(m)
				m.On("Find", mock.AnythingOfType("*[]model.Article"), mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Article)
					*arg = []model.Article{
						{ID: 1, UserID: 1, Author: model.User{ID: 1, Username: "user1"}},
						{ID: 2, UserID: 2, Author: model.User{ID: 2, Username: "user2"}},
					}
				}).Return(&gorm.DB{})
			},
			expected: []model.Article{
				{ID: 1, UserID: 1, Author: model.User{ID: 1, Username: "user1"}},
				{ID: 2, UserID: 2, Author: model.User{ID: 2, Username: "user2"}},
			},
			wantErr: false,
		},
		{
			name:    "Empty Result Set",
			userIDs: []uint{99, 100},
			limit:   10,
			offset:  0,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m)
				m.On("Where", "user_id in (?)", mock.Anything).Return(m)
				m.On("Offset", int64(0)).Return(m)
				m.On("Limit", int64(10)).Return(m)
				m.On("Find", mock.AnythingOfType("*[]model.Article"), mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Article)
					*arg = []model.Article{}
				}).Return(&gorm.DB{})
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Database Error Handling",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m)
				m.On("Where", "user_id in (?)", mock.Anything).Return(m)
				m.On("Offset", int64(0)).Return(m)
				m.On("Limit", int64(10)).Return(m)
				m.On("Find", mock.AnythingOfType("*[]model.Article"), mock.Anything).Return(&gorm.DB{Error: errors.New("database error")})
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{
				db: mockDB,
			}

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, got)

			if len(got) > int(tt.limit) {
				t.Errorf("ArticleStore.GetFeedArticles() returned more articles than limit. got = %d, limit = %d", len(got), tt.limit)
			}

			for _, article := range got {
				assert.NotEmpty(t, article.Author, "Article should have a preloaded Author")
				assert.Contains(t, tt.userIDs, article.UserID, "Article UserID should be in the requested userIDs")
			}

			mockDB.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90


 */
func (m *MockAssociation) Append(values ...interface{}) error {
	args := m.Called(values...)
	return args.Error(0)
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

func TestAddFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedCommit bool
	}{
		{
			name: "Successfully Add Favorite",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Append", mock.Anything).Return(nil)
				mockDB.On("Model", mock.AnythingOfType("*model.Article")).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  1,
			expectedCommit: true,
		},
		{
			name: "Database Error on Association Append",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Append", mock.Anything).Return(errors.New("DB error"))
				mockDB.On("Model", mock.AnythingOfType("*model.Article")).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Rollback").Return(mockDB)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("DB error"),
			expectedCount:  0,
			expectedCommit: false,
		},
		{
			name: "Database Error on FavoritesCount Update",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Append", mock.Anything).Return(nil)
				mockDB.On("Model", mock.AnythingOfType("*model.Article")).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.Error = errors.New("Update error")
				mockDB.On("Rollback").Return(mockDB)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("Update error"),
			expectedCount:  0,
			expectedCommit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{db: mockDB}
			err := store.AddFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			if tt.expectedCommit {
				mockDB.AssertCalled(t, "Commit")
			} else {
				mockDB.AssertCalled(t, "Rollback")
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockDB, *MockAssociation)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedCommit bool
	}{
		{
			name: "Successfully Delete a Favorite",
			setupMock: func(db *MockDB, assoc *MockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				db.On("Update", "favorites_count", mock.Anything).Return(db)
				db.On("Commit").Return(db)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockAssoc := new(MockAssociation)
			tt.setupMock(mockDB, mockAssoc)

			dbWrapper := struct {
				*MockDB
			}{mockDB}

			store := &ArticleStore{db: &dbWrapper}
			err := store.DeleteFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			if tt.expectedCommit {
				mockDB.AssertCalled(t, "Commit")
			} else {
				mockDB.AssertCalled(t, "Rollback")
			}

			mockDB.AssertExpectations(t)
			mockAssoc.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestGetArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	require.NoError(t, err)

	s := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(sqlmock.Sqlmock)
		wantCount   int
		wantErr     bool
	}{
		{
			name:        "Retrieve Articles Without Any Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Title 1", "Description 1", "Body 1", 1).
					AddRow(2, time.Now(), time.Now(), nil, "Title 2", "Description 2", "Body 2", 2)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").WillReturnRows(rows)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:        "Filter Articles by Tag Name",
			tagName:     "golang",
			username:    "",
			favoritedBy: nil,
			limit:       5,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Golang Article", "Description", "Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` JOIN article_tags").WillReturnRows(rows)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:        "Filter Articles by Author Username",
			tagName:     "",
			username:    "johndoe",
			favoritedBy: nil,
			limit:       5,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "John's Article", "Description", "Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` JOIN users").WillReturnRows(rows)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:        "Retrieve Favorited Articles",
			tagName:     "",
			username:    "",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       5,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favoriteRows := sqlmock.NewRows([]string{"article_id"}).AddRow(1).AddRow(2)
				mock.ExpectQuery("^SELECT article_id FROM `favorite_articles`").WillReturnRows(favoriteRows)

				articleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Favorite 1", "Description", "Body", 1).
					AddRow(2, time.Now(), time.Now(), nil, "Favorite 2", "Description", "Body", 2)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").WillReturnRows(articleRows)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:        "Error Handling for Database Issues",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").WillReturnError(gorm.ErrRecordNotFound)
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			got, err := s.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantCount)

				switch tt.name {
				case "Filter Articles by Tag Name":

					assert.Equal(t, "Golang Article", got[0].Title)
				case "Filter Articles by Author Username":
					assert.Equal(t, "John's Article", got[0].Title)
				case "Retrieve Favorited Articles":
					assert.Len(t, got, 2)
					assert.Equal(t, "Favorite 1", got[0].Title)
					assert.Equal(t, "Favorite 2", got[1].Title)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

