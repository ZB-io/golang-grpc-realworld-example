package store

import (
		"reflect"
		"sync"
		"testing"
		"time"
		"github.com/jinzhu/gorm"
		"errors"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/mock"
		"github.com/stretchr/testify/require"
		"github.com/DATA-DOG/go-sqlmock"
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
type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc *Location
}
type MockDB struct {
	mock.Mock
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
type Tag struct {
	gorm.Model
	Name string `gorm:"not null"`
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
type DBInterface interface {
	Delete(value interface{}) *gorm.DB
}
type DBInterface interface {
	Preload(column string, conditions ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
}
type mockDB struct {
	mock.Mock
}
type MockDB struct {
	gorm.DB
	mock.Mock
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
type Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
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
				t.Fatal("NewArticleStore returned nil")
			}
			if !reflect.DeepEqual(got.db, tt.want.db) {
				t.Errorf("NewArticleStore() = %v, want %v", got.db, tt.want.db)
			}
		})
	}

	t.Run("Verify ArticleStore Immutability", func(t *testing.T) {
		db := &gorm.DB{}
		store1 := NewArticleStore(db)
		store2 := NewArticleStore(db)
		if store1 == store2 {
			t.Error("NewArticleStore returned the same instance for multiple calls")
		}
		if store1.db != store2.db {
			t.Error("NewArticleStore did not use the same DB reference for multiple calls")
		}
	})

	t.Run("Verify DB Reference Integrity", func(t *testing.T) {
		db := &gorm.DB{Value: "unique_identifier"}
		store := NewArticleStore(db)
		if store.db != db {
			t.Error("NewArticleStore did not maintain DB reference integrity")
		}
	})

	t.Run("Performance Test for NewArticleStore", func(t *testing.T) {
		db := &gorm.DB{}
		iterations := 1000
		start := time.Now()
		for i := 0; i < iterations; i++ {
			NewArticleStore(db)
		}
		duration := time.Since(start)
		t.Logf("Time taken for %d iterations: %v", iterations, duration)

	})

	t.Run("Concurrent Access Safety", func(t *testing.T) {
		db := &gorm.DB{}
		var wg sync.WaitGroup
		concurrency := 100
		stores := make([]*ArticleStore, concurrency)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				stores[index] = NewArticleStore(db)
			}(i)
		}

		wg.Wait()

		for _, store := range stores {
			if store == nil {
				t.Error("Concurrent call to NewArticleStore resulted in nil ArticleStore")
			}
			if store.db != db {
				t.Error("Concurrent call to NewArticleStore did not maintain DB reference integrity")
			}
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

func TestArticleStoreCreate(t *testing.T) {
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
			name: "Create an Article with Database Connection Error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			dbError: errors.New("database connection error"),
			wantErr: true,
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

			mockDB.AssertExpectations(t)
		})
	}

	t.Run("Create Multiple Articles in Succession", func(t *testing.T) {
		mockDB := new(MockDB)
		mockDB.On("Create", mock.AnythingOfType("*model.Article")).Return(&gorm.DB{Error: nil}).Times(3)

		store := &ArticleStore{
			db: mockDB,
		}

		articles := []*model.Article{
			{Title: "Article 1", Description: "Desc 1", Body: "Body 1", UserID: 1},
			{Title: "Article 2", Description: "Desc 2", Body: "Body 2", UserID: 2},
			{Title: "Article 3", Description: "Desc 3", Body: "Body 3", UserID: 3},
		}

		for _, article := range articles {
			err := store.Create(article)
			assert.NoError(t, err)
		}

		mockDB.AssertExpectations(t)
	})
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestArticleStoreCreateComment(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*gorm.DB)
		comment *model.Comment
		wantErr bool
	}{
		{
			name: "Successfully Create a Comment",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{Model: gorm.Model{ID: 1}})
				db.Create(&model.Article{Model: gorm.Model{ID: 1}})
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
				db.Create(&model.User{Model: gorm.Model{ID: 1}})
				db.Create(&model.Article{Model: gorm.Model{ID: 1}})
			},
			comment: &model.Comment{

				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
		},
		{
			name: "Create a Comment with Maximum Length Body",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{Model: gorm.Model{ID: 1}})
				db.Create(&model.Article{Model: gorm.Model{ID: 1}})
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
				db.Create(&model.User{Model: gorm.Model{ID: 1}})
			},
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 9999,
			},
			wantErr: true,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent User",
			setup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}})
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

			s := &ArticleStore{db: db}

			err = s.CreateComment(tt.comment)

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
func TestArticleStoreDelete(t *testing.T) {
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
				tag := &model.Tag{Name: "TestTag"}
				db.Create(tag)
				article.Tags = []model.Tag{*tag}
				db.Create(article)
				return article
			},
			wantErr: false,
		},
		{
			name: "Delete an Article with Comments",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Description: "Test Description", Body: "Test Body"}
				db.Create(article)
				comment := &model.Comment{Body: "Test Comment", ArticleID: article.ID}
				db.Create(comment)
				return article
			},
			wantErr: false,
		},
		{
			name: "Delete an Article with Favorited Users",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Description: "Test Description", Body: "Test Body"}
				user := &model.User{Username: "testuser"}
				db.Create(user)
				article.FavoritedUsers = []model.User{*user}
				db.Create(article)
				return article
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error During Delete",
			setup: func(db *gorm.DB) *model.Article {

				db.Close()
				return &model.Article{Model: gorm.Model{ID: 1}}
			},
			wantErr: true,
		},
		{
			name: "Delete an Article with Maximum Integer ID",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{Model: gorm.Model{ID: ^uint(0)}, Title: "Max ID Article", Description: "Test Description", Body: "Test Body"}
				db.Create(article)
				return article
			},
			wantErr: false,
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

				var count int64
				db.Model(&model.Article{}).Where("id = ?", article.ID).Count(&count)
				assert.Equal(t, int64(0), count)

				if tt.name == "Delete an Article with Associated Tags" {
					var tagCount int64
					db.Model(&model.Tag{}).Count(&tagCount)
					assert.Equal(t, int64(1), tagCount)
					var articleTagCount int64
					db.Table("article_tags").Where("article_id = ?", article.ID).Count(&articleTagCount)
					assert.Equal(t, int64(0), articleTagCount)
				}

				if tt.name == "Delete an Article with Comments" {
					var commentCount int64
					db.Model(&model.Comment{}).Where("article_id = ?", article.ID).Count(&commentCount)
					assert.Equal(t, int64(0), commentCount)
				}

				if tt.name == "Delete an Article with Favorited Users" {
					var favoriteCount int64
					db.Table("favorite_articles").Where("article_id = ?", article.ID).Count(&favoriteCount)
					assert.Equal(t, int64(0), favoriteCount)
				}
			}
		})
	}
}

func (m *MockDB) Delete(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func TestArticleStoreDeleteComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		dbSetup func(*MockDB)
		wantErr bool
	}{
		{
			name: "Successfully Delete an Existing Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
			dbSetup: func(mockDB *MockDB) {
				mockDB.On("Delete", mock.AnythingOfType("*model.Comment")).Return(&gorm.DB{Error: nil})
			},
			wantErr: false,
		},
		{
			name: "Attempt to Delete a Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 999},
				Body:  "Non-existent comment",
			},
			dbSetup: func(mockDB *MockDB) {
				mockDB.On("Delete", mock.AnythingOfType("*model.Comment")).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			wantErr: true,
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 2},
				Body:  "Test comment",
			},
			dbSetup: func(mockDB *MockDB) {
				mockDB.On("Delete", mock.AnythingOfType("*model.Comment")).Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			wantErr: true,
		},
		{
			name:    "Delete Comment with Null Comment Pointer",
			comment: nil,
			dbSetup: func(mockDB *MockDB) {

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

			err := s.DeleteComment(tt.comment)

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
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {
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
				comment := &model.Comment{
					Model:     gorm.Model{ID: 1},
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(comment)
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
			name:            "Attempt to retrieve a non-existent comment",
			setupDB:         func(db *gorm.DB) {},
			commentID:       999,
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Retrieve a comment with associated data",
			setupDB: func(db *gorm.DB) {
				user := &model.User{Model: gorm.Model{ID: 1}, Username: "testuser"}
				article := &model.Article{Model: gorm.Model{ID: 1}, Title: "Test Article"}
				comment := &model.Comment{
					Model:     gorm.Model{ID: 2},
					Body:      "Comment with associations",
					UserID:    1,
					Author:    *user,
					ArticleID: 1,
					Article:   *article,
				}
				db.Create(user)
				db.Create(article)
				db.Create(comment)
			},
			commentID:     2,
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 2},
				Body:      "Comment with associations",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name: "Retrieve a soft-deleted comment",
			setupDB: func(db *gorm.DB) {
				comment := &model.Comment{
					Model:     gorm.Model{ID: 3},
					Body:      "Soft-deleted comment",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(comment)
				db.Delete(comment)
			},
			commentID:       3,
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.User{}, &model.Article{}, &model.Comment{})

			tt.setupDB(db)

			store := &ArticleStore{db: db}

			comment, err := store.GetCommentByID(tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedComment.ID, comment.ID)
				assert.Equal(t, tt.expectedComment.Body, comment.Body)
				assert.Equal(t, tt.expectedComment.UserID, comment.UserID)
				assert.Equal(t, tt.expectedComment.ArticleID, comment.ArticleID)
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
			name: "Database error handling",
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Large number of tags",
			dbSetup: func(db *gorm.DB) {
				for i := 0; i < 100; i++ {
					db.Create(&model.Tag{Name: "tag"})
				}
			},
			want:    make([]model.Tag, 100),
			wantErr: false,
		},
		{
			name: "Duplicate tag handling",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag2"})
				db.Create(&model.Tag{Name: "tag1"})
			},
			want: []model.Tag{
				{Name: "tag1"},
				{Name: "tag2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			assert.NoError(t, err, "failed to open database")
			defer db.Close()

			err = db.AutoMigrate(&model.Tag{}).Error
			assert.NoError(t, err, "failed to migrate schema")

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			got, err := s.GetTags()

			if tt.wantErr {
				assert.Error(t, err, "ArticleStore.GetTags() error expected")
			} else {
				assert.NoError(t, err, "ArticleStore.GetTags() unexpected error")
			}

			if tt.wantErr {
				assert.Nil(t, got, "Expected nil result when error occurs")
			} else {
				assert.Equal(t, len(tt.want), len(got), "Unexpected number of tags")
				for i, tag := range tt.want {
					assert.Equal(t, tag.Name, got[i].Name, "Tag name mismatch")
				}
			}
		})
	}
}

func TestArticleStoreGetTagsConcurrent(t *testing.T) {

}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (s *ArticleStore) GetByID(id uint) (*model.Article, error) {
	var m model.Article
	err := s.db.Preload("Tags").Preload("Author").Find(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}

func TestArticleStoreGetByID(t *testing.T) {
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
						Model:       gorm.Model{ID: 1},
						Title:       "Test Article",
						Description: "Test Description",
						Body:        "Test Body",
						Tags:        []model.Tag{{Name: "test"}},
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
				Tags:        []model.Tag{{Name: "test"}},
				Author:      model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
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
func (m *mockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockDB) Model(value interface{}) *mockDB {
	args := m.Called(value)
	return args.Get(0).(*mockDB)
}

func TestArticleStoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbError error
		wantErr bool
	}{
		{
			name: "Successfully Update an Existing Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Updated Title",
				Body:  "Updated Body",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Update a Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			dbError: gorm.ErrRecordNotFound,
			wantErr: true,
		},
		{
			name: "Update Article with Invalid Data",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
				Title: "",
			},
			dbError: errors.New("invalid data"),
			wantErr: true,
		},
		{
			name: "Update Article Tags",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
				Title: "Article with Tags",
				Tags:  []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Update Article Favorites Count",
			article: &model.Article{
				Model:          gorm.Model{ID: 4},
				Title:          "Popular Article",
				FavoritesCount: 10,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Update Article with Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
				Title: "Connection Error Article",
			},
			dbError: errors.New("database connection error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			store := &ArticleStore{db: mockDB}

			mockDB.On("Model", mock.AnythingOfType("*model.Article")).Return(mockDB)
			mockDB.On("Update", mock.AnythingOfType("*model.Article")).Return(mockDB)
			mockDB.On("Error").Return(tt.dbError)

			err := store.Update(tt.article)

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

func (m *mockDB) Update(attrs ...interface{}) *mockDB {
	args := m.Called(attrs...)
	return args.Get(0).(*mockDB)
}

func (m *MockDB) Update(attrs ...interface{}) *gorm.DB {
	args := m.Called(attrs...)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestArticleStoreGetComments(t *testing.T) {
	tests := []struct {
		name            string
		setupFunc       func(*gorm.DB) *model.Article
		expectedCount   int
		expectedError   error
		validateFunc    func(*testing.T, []model.Comment)
		performanceTest bool
	}{
		{
			name: "Successful retrieval of comments",
			setupFunc: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Test Article", Body: "Test Body"}
				db.Create(article)

				for i := 0; i < 3; i++ {
					user := &model.User{Username: "user" + string(rune(i+'0'))}
					db.Create(user)
					comment := &model.Comment{Body: "Comment " + string(rune(i+'0')), UserID: user.ID, ArticleID: article.ID}
					db.Create(comment)
				}

				return article
			},
			expectedCount: 3,
			expectedError: nil,
			validateFunc: func(t *testing.T, comments []model.Comment) {
				for _, comment := range comments {
					assert.NotEmpty(t, comment.Author)
					assert.NotEmpty(t, comment.Body)
				}
			},
		},
		{
			name: "Article with no comments",
			setupFunc: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "No Comments Article", Body: "Test Body"}
				db.Create(article)
				return article
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name: "Database error",
			setupFunc: func(db *gorm.DB) *model.Article {

				db.Close()
				return &model.Article{Model: gorm.Model{ID: 1}}
			},
			expectedCount: 0,
			expectedError: errors.New("database error"),
		},
		{
			name: "Comments with deleted authors",
			setupFunc: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Deleted Authors Article", Body: "Test Body"}
				db.Create(article)

				user := &model.User{Username: "deleted_user"}
				db.Create(user)
				comment := &model.Comment{Body: "Orphan Comment", UserID: user.ID, ArticleID: article.ID}
				db.Create(comment)

				db.Delete(user)

				return article
			},
			expectedCount: 1,
			expectedError: nil,
			validateFunc: func(t *testing.T, comments []model.Comment) {
				assert.Equal(t, uint(0), comments[0].Author.ID)
			},
		},
		{
			name: "Performance with large number of comments",
			setupFunc: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Performance Test Article", Body: "Test Body"}
				db.Create(article)

				for i := 0; i < 1000; i++ {
					user := &model.User{Username: "user" + string(rune(i+'0'))}
					db.Create(user)
					comment := &model.Comment{Body: "Comment " + string(rune(i+'0')), UserID: user.ID, ArticleID: article.ID}
					db.Create(comment)
				}

				return article
			},
			expectedCount:   1000,
			expectedError:   nil,
			performanceTest: true,
		},
		{
			name: "Consistency of comment order",
			setupFunc: func(db *gorm.DB) *model.Article {
				article := &model.Article{Title: "Order Test Article", Body: "Test Body"}
				db.Create(article)

				for i := 0; i < 5; i++ {
					user := &model.User{Username: "user" + string(rune(i+'0'))}
					db.Create(user)
					comment := &model.Comment{Body: "Comment " + string(rune(i+'0')), UserID: user.ID, ArticleID: article.ID}
					db.Create(comment)
				}

				return article
			},
			expectedCount: 5,
			expectedError: nil,
			validateFunc: func(t *testing.T, comments []model.Comment) {
				firstCallComments := comments

				store := &ArticleStore{db: db}
				secondCallComments, err := store.GetComments(&model.Article{Model: gorm.Model{ID: comments[0].ArticleID}})
				require.NoError(t, err)

				for i := range firstCallComments {
					assert.Equal(t, firstCallComments[i].ID, secondCallComments[i].ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.User{}, &model.Comment{})

			article := tt.setupFunc(db)

			store := &ArticleStore{db: db}

			start := time.Now()
			comments, err := store.GetComments(article)
			duration := time.Since(start)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tt.expectedCount)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, comments)
			}

			if tt.performanceTest {
				assert.Less(t, duration, 1*time.Second, "GetComments took too long for large number of comments")
			}
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

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(*MockDB)
		expectedFav   bool
		expectedError error
	}{
		{
			name:    "Article is favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Table", "favorite_articles").Return(mockDB)
				mockDB.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					*args.Get(0).(*int) = 1
				}).Return(mockDB)
			},
			expectedFav:   true,
			expectedError: nil,
		},
		{
			name:    "Article is not favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Table", "favorite_articles").Return(mockDB)
				mockDB.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					*args.Get(0).(*int) = 0
				}).Return(mockDB)
			},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Null Article parameter",
			article:       nil,
			user:          &model.User{Model: gorm.Model{ID: 1}},
			setupMock:     func(mockDB *MockDB) {},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:          "Null User parameter",
			article:       &model.Article{Model: gorm.Model{ID: 1}},
			user:          nil,
			setupMock:     func(mockDB *MockDB) {},
			expectedFav:   false,
			expectedError: nil,
		},
		{
			name:    "Database error occurs",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Table", "favorite_articles").Return(mockDB)
				mockDB.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("Count", mock.AnythingOfType("*int")).Return(&gorm.DB{Error: errors.New("database error")})
			},
			expectedFav:   false,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{db: mockDB}

			isFavorited, err := store.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedFav, isFavorited)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func TestArticleStoreGetFeedArticles(t *testing.T) {
	tests := []struct {
		name    string
		userIDs []uint
		limit   int64
		offset  int64
		mockDB  func() (*gorm.DB, sqlmock.Sqlmock)
		want    []model.Article
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Successful retrieval of feed articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "title", "user_id"}).
					AddRow(1, "Article 1", 1).
					AddRow(2, "Article 2", 2)

				mock.ExpectQuery("SELECT (.+) FROM `articles`").
					WithArgs(1, 2).
					WillReturnRows(rows)

				return gormDB, mock
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 2},
			},
			wantErr: false,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{99, 100},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "title", "user_id"})

				mock.ExpectQuery("SELECT (.+) FROM `articles`").
					WithArgs(99, 100).
					WillReturnRows(rows)

				return gormDB, mock
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Database error handling",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				mock.ExpectQuery("SELECT (.+) FROM `articles`").
					WithArgs(1).
					WillReturnError(errors.New("database error"))

				return gormDB, mock
			},
			want:    nil,
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock := tt.mockDB()
			s := &ArticleStore{
				db: mockDB,
			}
			got, err := s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.errMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.want)
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

func TestArticleStoreAddFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedAppend bool
	}{
		{
			name: "Successfully Add Favorite",
			setupMock: func(mockDB *MockDB) {
				tx := &gorm.DB{}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Append", mock.Anything).Return(nil)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  1,
			expectedAppend: true,
		},
		{
			name: "Database Error During Association Append",
			setupMock: func(mockDB *MockDB) {
				tx := &gorm.DB{}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Append", mock.Anything).Return(errors.New("DB error"))
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Rollback").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("DB error"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name: "Database Error During FavoritesCount Update",
			setupMock: func(mockDB *MockDB) {
				tx := &gorm.DB{}
				mockDB.On("Begin").Return(tx)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Append", mock.Anything).Return(nil)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{Error: errors.New("Update error")})
				mockDB.On("Rollback").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("Update error"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name:           "Add Favorite with Nil Article",
			setupMock:      func(mockDB *MockDB) {},
			article:        nil,
			user:           &model.User{},
			expectedError:  errors.New("invalid argument: nil article"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name:           "Add Favorite with Nil User",
			setupMock:      func(mockDB *MockDB) {},
			article:        &model.Article{FavoritesCount: 0},
			user:           nil,
			expectedError:  errors.New("invalid argument: nil user"),
			expectedCount:  0,
			expectedAppend: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			db := &gorm.DB{
				Value: mockDB,
			}

			store := &ArticleStore{db: db}
			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.article != nil {
				assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMockDB    func() *MockDB
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedCommit bool
	}{
		{
			name: "Successfully Delete a Favorite",
			setupMockDB: func() *MockDB {
				mockDB := new(MockDB)
				mockDB.On("Begin").Return(&gorm.DB{})
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{})
				mockDB.On("Commit").Return(&gorm.DB{})
				return mockDB
			},
			article: &model.Article{
				FavoritesCount: 1,
				FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
			},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
		{
			name: "Attempt to Delete a Non-existent Favorite",
			setupMockDB: func() *MockDB {
				mockDB := new(MockDB)
				mockDB.On("Begin").Return(&gorm.DB{})
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{})
				mockDB.On("Commit").Return(&gorm.DB{})
				return mockDB
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
		{
			name: "Database Error During Association Deletion",
			setupMockDB: func() *MockDB {
				mockDB := new(MockDB)
				mockDB.On("Begin").Return(&gorm.DB{})
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{Error: errors.New("association deletion error")})
				mockDB.On("Rollback").Return(&gorm.DB{})
				return mockDB
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedError:  errors.New("association deletion error"),
			expectedCount:  1,
			expectedCommit: false,
		},
		{
			name: "Database Error During Favorites Count Update",
			setupMockDB: func() *MockDB {
				mockDB := new(MockDB)
				mockDB.On("Begin").Return(&gorm.DB{})
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{Error: errors.New("update error")})
				mockDB.On("Rollback").Return(&gorm.DB{})
				return mockDB
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedError:  errors.New("update error"),
			expectedCount:  1,
			expectedCommit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setupMockDB()
			store := &ArticleStore{db: mockDB}

			err := store.DeleteFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			if tt.expectedCommit {
				mockDB.AssertCalled(t, "Commit")
				mockDB.AssertNotCalled(t, "Rollback")
			} else {
				mockDB.AssertCalled(t, "Rollback")
				mockDB.AssertNotCalled(t, "Commit")
			}
		})
	}
}

func TestArticleStoreDeleteFavoriteConcurrent(t *testing.T) {
	mockDB := new(MockDB)
	mockDB.On("Begin").Return(&gorm.DB{})
	mockDB.On("Model", mock.Anything).Return(mockDB)
	mockDB.On("Association", "FavoritedUsers").Return(&gorm.Association{})
	mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{})
	mockDB.On("Commit").Return(&gorm.DB{})

	store := &ArticleStore{db: mockDB}

	article := &model.Article{
		FavoritesCount: 5,
		FavoritedUsers: []model.User{
			{Model: gorm.Model{ID: 1}},
			{Model: gorm.Model{ID: 2}},
			{Model: gorm.Model{ID: 3}},
			{Model: gorm.Model{ID: 4}},
			{Model: gorm.Model{ID: 5}},
		},
	}

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(userID uint) {
			defer wg.Done()
			user := &model.User{Model: gorm.Model{ID: userID}}
			_ = store.DeleteFavorite(article, user)
		}(uint(i))
	}

	wg.Wait()

	assert.Equal(t, int32(0), article.FavoritesCount)
	assert.Empty(t, article.FavoritedUsers)

	mockDB.AssertNumberOfCalls(t, "Begin", 5)
	mockDB.AssertNumberOfCalls(t, "Model", 10)
	mockDB.AssertNumberOfCalls(t, "Association", 5)
	mockDB.AssertNumberOfCalls(t, "Update", 5)
	mockDB.AssertNumberOfCalls(t, "Commit", 5)
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
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

func (m *MockDB) Rows() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *MockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		wantCount   int
		wantErr     bool
	}{
		{
			name:        "No filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantCount:   10,
			wantErr:     false,
		},
		{
			name:        "Filter by tag",
			tagName:     "test-tag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantCount:   5,
			wantErr:     false,
		},
		{
			name:        "Filter by author",
			tagName:     "",
			username:    "test-author",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantCount:   3,
			wantErr:     false,
		},
		{
			name:        "Filter by favorited",
			tagName:     "",
			username:    "",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			wantCount:   2,
			wantErr:     false,
		},
		{
			name:        "Pagination",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       5,
			offset:      5,
			wantCount:   5,
			wantErr:     false,
		},
		{
			name:        "Combined filters",
			tagName:     "test-tag",
			username:    "test-author",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			wantCount:   1,
			wantErr:     false,
		},
		{
			name:        "Empty result",
			tagName:     "non-existent-tag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			wantCount:   0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			s := &ArticleStore{
				db: mockDB,
			}

			mockDB.On("Preload", "Author").Return(mockDB)
			if tt.username != "" {
				mockDB.On("Joins", "join users on articles.user_id = users.id").Return(mockDB)
				mockDB.On("Where", "users.username = ?", tt.username).Return(mockDB)
			}
			if tt.tagName != "" {
				mockDB.On("Joins", "join article_tags on articles.id = article_tags.article_id join tags on tags.id = article_tags.tag_id").Return(mockDB)
				mockDB.On("Where", "tags.name = ?", tt.tagName).Return(mockDB)
			}
			if tt.favoritedBy != nil {
				mockDB.On("Select", "article_id").Return(mockDB)
				mockDB.On("Table", "favorite_articles").Return(mockDB)
				mockDB.On("Where", "user_id = ?", tt.favoritedBy.ID).Return(mockDB)
				mockDB.On("Offset", tt.offset).Return(mockDB)
				mockDB.On("Limit", tt.limit).Return(mockDB)
				mockDB.On("Rows").Return(mockDB, nil)
				mockDB.On("Where", "id in (?)", mock.Anything).Return(mockDB)
			}
			mockDB.On("Offset", tt.offset).Return(mockDB)
			mockDB.On("Limit", tt.limit).Return(mockDB)
			mockDB.On("Find", mock.AnythingOfType("*[]model.Article")).Return(mockDB)

			got, err := s.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantCount)

				if tt.tagName != "" {
					for _, article := range got {
						assert.Contains(t, article.Tags, model.Tag{Name: tt.tagName})
					}
				}
				if tt.username != "" {
					for _, article := range got {
						assert.Equal(t, tt.username, article.Author.Username)
					}
				}
				if tt.favoritedBy != nil {
					for _, article := range got {
						assert.Contains(t, article.FavoritedUsers, *tt.favoritedBy)
					}
				}
			}

			mockDB.AssertExpectations(t)
		})
	}

	t.Run("Database error", func(t *testing.T) {
		mockDB := new(MockDB)
		s := &ArticleStore{
			db: mockDB,
		}

		mockDB.On("Preload", "Author").Return(mockDB)
		mockDB.On("Offset", int64(0)).Return(mockDB)
		mockDB.On("Limit", int64(10)).Return(mockDB)
		mockDB.On("Find", mock.AnythingOfType("*[]model.Article")).Return(mockDB).Run(func(args mock.Arguments) {
			db := args.Get(0).(*gorm.DB)
			db.Error = errors.New("database error")
		})

		_, err := s.GetArticles("", "", nil, 10, 0)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		mockDB.AssertExpectations(t)
	})
}

