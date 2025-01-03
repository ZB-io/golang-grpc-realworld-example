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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticleStore() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Verify ArticleStore Immutability", func(t *testing.T) {
		db := &gorm.DB{}
		store1 := NewArticleStore(db)
		store2 := NewArticleStore(db)
		if store1 == store2 {
			t.Errorf("NewArticleStore() returned the same instance for multiple calls")
		}
	})

	t.Run("Check DB Field Accessibility", func(t *testing.T) {
		db := &gorm.DB{Value: "test"}
		store := NewArticleStore(db)
		if !reflect.DeepEqual(store.db, db) {
			t.Errorf("NewArticleStore().db = %v, want %v", store.db, db)
		}
	})

	t.Run("Performance Test for Multiple Instantiations", func(t *testing.T) {
		db := &gorm.DB{}
		start := time.Now()
		for i := 0; i < 1000; i++ {
			NewArticleStore(db)
		}
		duration := time.Since(start)
		if duration > time.Second {
			t.Errorf("NewArticleStore() took too long for 1000 instantiations: %v", duration)
		}
	})

	t.Run("Concurrent Access Safety", func(t *testing.T) {
		db := &gorm.DB{}
		var wg sync.WaitGroup
		storesChan := make(chan *ArticleStore, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				store := NewArticleStore(db)
				storesChan <- store
			}()
		}

		wg.Wait()
		close(storesChan)

		stores := make([]*ArticleStore, 0, 100)
		for store := range storesChan {
			stores = append(stores, store)
		}

		if len(stores) != 100 {
			t.Errorf("Expected 100 ArticleStore instances, got %d", len(stores))
		}

		for _, store := range stores {
			if store == nil || store.db != db {
				t.Errorf("Invalid ArticleStore instance created in concurrent execution")
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


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func (m *MockDB) Delete(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

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
			name: "Successfully Retrieve All Tags",
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
			name:    "Empty Tag List",
			dbSetup: func(db *gorm.DB) {},
			want:    []model.Tag{},
			wantErr: false,
		},
		{
			name: "Database Connection Error",
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Large Number of Tags",
			dbSetup: func(db *gorm.DB) {
				for i := 0; i < 1000; i++ {
					db.Create(&model.Tag{Name: "tag"})
				}
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate Tags in Database",
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
			assert.NoError(t, err)
			defer db.Close()

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			got, err := s.GetTags()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.name == "Large Number of Tags" {
				assert.Equal(t, 1000, len(got))
			} else if !tt.wantErr {
				assert.Equal(t, tt.want, got)
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


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	arguments := m.Called(out, where)
	return arguments.Get(0).(*gorm.DB)
}

func TestArticleStoreGetComments(t *testing.T) {
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
				comments := []model.Comment{
					{Model: gorm.Model{ID: 1}, ArticleID: 1, Body: "Comment 1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
					{Model: gorm.Model{ID: 2}, ArticleID: 1, Body: "Comment 2", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
				}
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(1)).Return(mockDB)
				mockDB.On("Find", &[]model.Comment{}, []interface{}(nil)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Comment)
					*arg = comments
				}).Return(mockDB)
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 1}, ArticleID: 1, Body: "Comment 1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, ArticleID: 1, Body: "Comment 2", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			expectedError: nil,
		},
		{
			name: "Retrieve comments for an article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(2)).Return(mockDB)
				mockDB.On("Find", &[]model.Comment{}, []interface{}(nil)).Return(mockDB)
			},
			expectedResult: []model.Comment{},
			expectedError:  nil,
		},
		{
			name: "Handle database error when retrieving comments",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(3)).Return(mockDB)
				mockDB.On("Find", &[]model.Comment{}, []interface{}(nil)).Return(&gorm.DB{Error: errors.New("database error")})
			},
			expectedResult: []model.Comment{},
			expectedError:  errors.New("database error"),
		},
		{
			name: "Retrieve comments for a non-existent article",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
			},
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(4)).Return(mockDB)
				mockDB.On("Find", &[]model.Comment{}, []interface{}(nil)).Return(mockDB)
			},
			expectedResult: []model.Comment{},
			expectedError:  nil,
		},
		{
			name: "Verify correct ordering of retrieved comments",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
			},
			mockSetup: func(mockDB *MockDB) {
				comments := []model.Comment{
					{Model: gorm.Model{ID: 3, CreatedAt: time.Now().Add(-1 * time.Hour)}, ArticleID: 5, Body: "Old Comment", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
					{Model: gorm.Model{ID: 4, CreatedAt: time.Now()}, ArticleID: 5, Body: "New Comment", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
				}
				mockDB.On("Preload", "Author").Return(mockDB)
				mockDB.On("Where", "article_id = ?", uint(5)).Return(mockDB)
				mockDB.On("Find", &[]model.Comment{}, []interface{}(nil)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Comment)
					*arg = comments
				}).Return(mockDB)
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 3, CreatedAt: time.Now().Add(-1 * time.Hour)}, ArticleID: 5, Body: "Old Comment", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 4, CreatedAt: time.Now()}, ArticleID: 5, Body: "New Comment", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB}

			result, err := store.GetComments(tt.article)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

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
		name     string
		userIDs  []uint
		limit    int64
		offset   int64
		mockDB   func() *gorm.DB
		expected []model.Article
		wantErr  bool
	}{
		{
			name:    "Successful Retrieval of Feed Articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
		},
		{
			name:    "Empty Result Set",
			userIDs: []uint{3, 4},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Pagination with Offset and Limit",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 3}, Title: "Article 3", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
		},
		{
			name:    "Database Error Handling",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(errors.New("database error"))
				return db
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Large Number of User IDs",
			userIDs: func() []uint {
				ids := make([]uint, 1000)
				for i := range ids {
					ids[i] = uint(i + 1)
				}
				return ids
			}(),
			limit:  50,
			offset: 0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: func() []model.Article {
				articles := make([]model.Article, 50)
				for i := range articles {
					articles[i] = model.Article{Model: gorm.Model{ID: uint(i + 1)}, Title: "Article", UserID: uint(i + 1)}
				}
				return articles
			}(),
			wantErr: false,
		},
		{
			name:    "Zero Limit and Offset",
			userIDs: []uint{1, 2},
			limit:   0,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			store := &ArticleStore{
				db: tt.mockDB(),
			}

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, got)
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

func (m *MockDB) Update(attrs ...interface{}) *gorm.DB {
	args := m.Called(attrs...)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestArticleStoreGetArticles(t *testing.T) {

	mockDB := &gorm.DB{}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(*gorm.DB)
		expected    []model.Article
		expectedErr error
	}{
		{
			name:        "Retrieve Articles Without Any Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Preload("Author").Find(&[]model.Article{
					{Title: "Article 1"},
					{Title: "Article 2"},
				})
			},
			expected: []model.Article{
				{Title: "Article 1"},
				{Title: "Article 2"},
			},
			expectedErr: nil,
		},
		{
			name:        "Filter Articles by Tag Name",
			tagName:     "golang",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Preload("Author").Find(&[]model.Article{
					{Title: "Golang Article"},
				})
			},
			expected: []model.Article{
				{Title: "Golang Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Filter Articles by Author Username",
			tagName:     "",
			username:    "johndoe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Preload("Author").Find(&[]model.Article{
					{Title: "John's Article", Author: model.User{Username: "johndoe"}},
				})
			},
			expected: []model.Article{
				{Title: "John's Article", Author: model.User{Username: "johndoe"}},
			},
			expectedErr: nil,
		},
		{
			name:     "Retrieve Favorited Articles",
			tagName:  "",
			username: "",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Preload("Author").Find(&[]model.Article{
					{Title: "Favorited Article"},
				})
			},
			expected: []model.Article{
				{Title: "Favorited Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Test Pagination with Limit and Offset",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       5,
			offset:      10,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Preload("Author").Offset(10).Limit(5).Find(&[]model.Article{
					{Title: "Paginated Article 1"},
					{Title: "Paginated Article 2"},
				})
			},
			expected: []model.Article{
				{Title: "Paginated Article 1"},
				{Title: "Paginated Article 2"},
			},
			expectedErr: nil,
		},
		{
			name:        "Combine Multiple Filters",
			tagName:     "golang",
			username:    "johndoe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Preload("Author").Find(&[]model.Article{
					{Title: "John's Golang Article", Author: model.User{Username: "johndoe"}},
				})
			},
			expected: []model.Article{
				{Title: "John's Golang Article", Author: model.User{Username: "johndoe"}},
			},
			expectedErr: nil,
		},
		{
			name:        "Handle Empty Result Set",
			tagName:     "nonexistenttag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

				db.Error = gorm.ErrRecordNotFound
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Error Handling for Database Issues",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB = &gorm.DB{}
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB}
			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expected, articles)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

