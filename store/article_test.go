package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"time"
	"github.com/stretchr/testify/assert"
)






type mockDB struct {
	err    error
	result []model.Comment
}
type mockDB struct {
	db *gorm.DB
}


/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore 

 */
func TestNewArticleStore(t *testing.T) {
	tests := []struct {
		name      string
		db        *gorm.DB
		wantNil   bool
		wantDBNil bool
	}{
		{
			name:      "Valid DB connection",
			db:        &gorm.DB{},
			wantNil:   false,
			wantDBNil: false,
		},
		{
			name:      "Nil DB connection",
			db:        nil,
			wantNil:   false,
			wantDBNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewArticleStore() returned nil: %v, want nil: %v", got == nil, tt.wantNil)
			}

			if got != nil {
				if (got.db == nil) != tt.wantDBNil {
					t.Errorf("NewArticleStore().db is nil: %v, want nil: %v", got.db == nil, tt.wantDBNil)
				}

				if !tt.wantDBNil && !reflect.DeepEqual(got.db, tt.db) {
					t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.db)
				}
			}
		})
	}
}

func TestNewArticleStoreImmutability(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewArticleStore(db)
	store2 := NewArticleStore(db)

	if store1 == store2 {
		t.Error("NewArticleStore() returned the same instance for different calls")
	}

	if store1.db != store2.db {
		t.Error("NewArticleStore() did not use the same DB instance for different calls")
	}
}


/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1

FUNCTION_DEF=func (s *ArticleStore) Delete(m *model.Article) error 

 */
func TestArticleStoreDelete(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		dbSetup func(*gorm.DB)
		wantErr error
	}{
		{
			name: "Successfully Delete an Existing Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Test Article"})
			},
			wantErr: nil,
		},
		{
			name: "Attempt to Delete a Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Delete an Article with Associated Tags",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
				Title: "Article with Tags",
				Tags:  []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			dbSetup: func(db *gorm.DB) {
				article := &model.Article{Model: gorm.Model{ID: 2}, Title: "Article with Tags"}
				db.Create(article)
				db.Model(article).Association("Tags").Append([]model.Tag{{Name: "tag1"}, {Name: "tag2"}})
			},
			wantErr: nil,
		},
		{
			name: "Delete an Article with Comments",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
				Title: "Article with Comments",
			},
			dbSetup: func(db *gorm.DB) {
				article := &model.Article{Model: gorm.Model{ID: 3}, Title: "Article with Comments"}
				db.Create(article)
				db.Create(&model.Comment{Body: "Comment 1", ArticleID: 3})
				db.Create(&model.Comment{Body: "Comment 2", ArticleID: 3})
			},
			wantErr: nil,
		},
		{
			name: "Delete an Article with Favorites",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
				Title: "Favorited Article",
			},
			dbSetup: func(db *gorm.DB) {
				article := &model.Article{Model: gorm.Model{ID: 4}, Title: "Favorited Article"}
				db.Create(article)
				user1 := &model.User{Model: gorm.Model{ID: 1}, Username: "user1"}
				user2 := &model.User{Model: gorm.Model{ID: 2}, Username: "user2"}
				db.Create(user1)
				db.Create(user2)
				db.Model(user1).Association("FavoriteArticles").Append(article)
				db.Model(user2).Association("FavoriteArticles").Append(article)
			},
			wantErr: nil,
		},
		{
			name: "Handle Database Connection Error During Deletion",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
				Title: "Error Article",
			},
			dbSetup: func(db *gorm.DB) {

				db.AddError(errors.New("database connection error"))
			},
			wantErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, _ := gorm.Open("sqlite3", ":memory:")
			defer db.Close()

			tt.dbSetup(db)

			store := &ArticleStore{db: db}

			err := store.Delete(tt.article)

			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {

				var count int64
				db.Model(&model.Article{}).Where("id = ?", tt.article.ID).Count(&count)
				if count != 0 {
					t.Errorf("Article was not deleted")
				}

				if tt.name == "Delete an Article with Associated Tags" {
					var tagCount int64
					db.Model(&model.Tag{}).Where("id IN (?)", db.Table("article_tags").Select("tag_id").Where("article_id = ?", tt.article.ID).QueryExpr()).Count(&tagCount)
					if tagCount != 0 {
						t.Errorf("Associated tags were not deleted")
					}
				}

				if tt.name == "Delete an Article with Comments" {
					var commentCount int64
					db.Model(&model.Comment{}).Where("article_id = ?", tt.article.ID).Count(&commentCount)
					if commentCount != 0 {
						t.Errorf("Associated comments were not deleted")
					}
				}

				if tt.name == "Delete an Article with Favorites" {
					var favoriteCount int64
					db.Table("favorite_articles").Where("article_id = ?", tt.article.ID).Count(&favoriteCount)
					if favoriteCount != 0 {
						t.Errorf("Article was not removed from favorites")
					}
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12

FUNCTION_DEF=func (s *ArticleStore) DeleteComment(m *model.Comment) error 

 */
func TestArticleStoreDeleteComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		dbSetup func(*gorm.DB)
		wantErr bool
	}{
		{
			name: "Successfully Delete an Existing Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Comment{Model: gorm.Model{ID: 1}, Body: "Test comment"})
			},
			wantErr: false,
		},
		{
			name: "Attempt to Delete a Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 999},
				Body:  "Non-existent comment",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
		},
		{
			name: "Delete Comment with Associated Relationships",
			comment: &model.Comment{
				Model:     gorm.Model{ID: 2},
				Body:      "Comment with relationship",
				ArticleID: 1,
			},
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Test Article"})
				db.Create(&model.Comment{Model: gorm.Model{ID: 2}, Body: "Comment with relationship", ArticleID: 1})
			},
			wantErr: false,
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 3},
				Body:  "Error comment",
			},
			dbSetup: func(db *gorm.DB) {

				db.AddError(errors.New("database connection error"))
			},
			wantErr: true,
		},
		{
			name:    "Delete Comment with Null Comment Pointer",
			comment: nil,
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("failed to open database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Comment{}, &model.Article{})

			tt.dbSetup(db)

			store := &ArticleStore{db: db}

			err = store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.comment != nil {
				var count int
				db.Model(&model.Comment{}).Where("id = ?", tt.comment.ID).Count(&count)
				if count != 0 {
					t.Errorf("Comment was not deleted from the database")
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e

FUNCTION_DEF=func (s *ArticleStore) GetComments(m *model.Article) ([]model.Comment, error) 

 */
func (m *mockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	if m.err != nil {
		return &gorm.DB{Error: m.err}
	}
	*(out.(*[]model.Comment)) = m.result
	return &gorm.DB{Value: m}
}

func (m *mockDB) Preload(column string) *gorm.DB {
	return &gorm.DB{Value: m}
}

func TestArticleStoreGetComments(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		mockDBResult   []model.Comment
		mockDBError    error
		expectedResult []model.Comment
		expectedError  error
	}{
		{
			name: "Successfully retrieve comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockDBResult: []model.Comment{
				{
					Model:     gorm.Model{ID: 1},
					Body:      "Comment 1",
					UserID:    1,
					ArticleID: 1,
					Author:    model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
				},
				{
					Model:     gorm.Model{ID: 2},
					Body:      "Comment 2",
					UserID:    2,
					ArticleID: 1,
					Author:    model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
				},
			},
			mockDBError: nil,
			expectedResult: []model.Comment{
				{
					Model:     gorm.Model{ID: 1},
					Body:      "Comment 1",
					UserID:    1,
					ArticleID: 1,
					Author:    model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
				},
				{
					Model:     gorm.Model{ID: 2},
					Body:      "Comment 2",
					UserID:    2,
					ArticleID: 1,
					Author:    model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			mockDBResult:   []model.Comment{},
			mockDBError:    nil,
			expectedResult: []model.Comment{},
			expectedError:  nil,
		},
		{
			name: "Database error",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			mockDBResult:   nil,
			mockDBError:    errors.New("database error"),
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Comments with deleted authors",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
			},
			mockDBResult: []model.Comment{
				{
					Model:     gorm.Model{ID: 3},
					Body:      "Comment 3",
					UserID:    3,
					ArticleID: 4,
					Author:    model.User{},
				},
				{
					Model:     gorm.Model{ID: 4},
					Body:      "Comment 4",
					UserID:    4,
					ArticleID: 4,
					Author:    model.User{Model: gorm.Model{ID: 4}, Username: "user4"},
				},
			},
			mockDBError: nil,
			expectedResult: []model.Comment{
				{
					Model:     gorm.Model{ID: 3},
					Body:      "Comment 3",
					UserID:    3,
					ArticleID: 4,
					Author:    model.User{},
				},
				{
					Model:     gorm.Model{ID: 4},
					Body:      "Comment 4",
					UserID:    4,
					ArticleID: 4,
					Author:    model.User{Model: gorm.Model{ID: 4}, Username: "user4"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Verify correct ordering of comments",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
			},
			mockDBResult: []model.Comment{
				{
					Model:     gorm.Model{ID: 5, CreatedAt: time.Now().Add(-2 * time.Hour)},
					Body:      "Older Comment",
					UserID:    5,
					ArticleID: 5,
					Author:    model.User{Model: gorm.Model{ID: 5}, Username: "user5"},
				},
				{
					Model:     gorm.Model{ID: 6, CreatedAt: time.Now().Add(-1 * time.Hour)},
					Body:      "Newer Comment",
					UserID:    6,
					ArticleID: 5,
					Author:    model.User{Model: gorm.Model{ID: 6}, Username: "user6"},
				},
			},
			mockDBError: nil,
			expectedResult: []model.Comment{
				{
					Model:     gorm.Model{ID: 5, CreatedAt: time.Now().Add(-2 * time.Hour)},
					Body:      "Older Comment",
					UserID:    5,
					ArticleID: 5,
					Author:    model.User{Model: gorm.Model{ID: 5}, Username: "user5"},
				},
				{
					Model:     gorm.Model{ID: 6, CreatedAt: time.Now().Add(-1 * time.Hour)},
					Body:      "Newer Comment",
					UserID:    6,
					ArticleID: 5,
					Author:    model.User{Model: gorm.Model{ID: 6}, Username: "user6"},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				err:    tt.mockDBError,
				result: tt.mockDBResult,
			}

			store := &ArticleStore{
				db: &gorm.DB{Value: mockDB},
			}

			result, err := store.GetComments(tt.article)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)

			if tt.name == "Verify correct ordering of comments" {
				assert.True(t, result[0].CreatedAt.Before(result[1].CreatedAt), "Comments should be ordered by creation time")
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Value: m}
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs []uint, limit, offset int64) ([]model.Article, error) 

 */
func (m *mockDB) Find(out interface{}) *gorm.DB {
	return m.db.Find(out)
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	return m.db.Limit(limit)
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m.db.Offset(offset)
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m.db.Preload(column, conditions...)
}

func TestArticleStoreGetFeedArticles(t *testing.T) {
	tests := []struct {
		name    string
		userIDs []uint
		limit   int64
		offset  int64
		mockDB  func() *gorm.DB
		want    []model.Article
		wantErr bool
	}{
		{
			name:    "Successful Retrieval of Feed Articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:auto_preload", true)
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
				{Model: gorm.Model{ID: 2}, UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}}},
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
				return db.InstantSet("gorm:auto_preload", true)
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Pagination with Offset",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:auto_preload", true)
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 2}, UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
				{Model: gorm.Model{ID: 3}, UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}}},
			},
			wantErr: false,
		},
		{
			name:    "Handling Database Errors",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{Error: errors.New("database error")}
				return db
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Limit Exceeds Available Articles",
			userIDs: []uint{1, 2},
			limit:   100,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:auto_preload", true)
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
				{Model: gorm.Model{ID: 2}, UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}}},
			},
			wantErr: false,
		},
		{
			name:    "Zero Limit Handling",
			userIDs: []uint{1, 2},
			limit:   0,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:auto_preload", true)
			},
			want:    []model.Article{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}
			got, err := s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.db.Where(query, args...)
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) 

 */
func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockDB      func() *gorm.DB
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
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1"},
				{Model: gorm.Model{ID: 2}, Title: "Article 2"},
			},
			expectedErr: nil,
		},
		{
			name:        "Filter Articles by Tag Name",
			tagName:     "technology",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Tech Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Filter Articles by Author Username",
			tagName:     "",
			username:    "john_doe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "John's Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Retrieve Favorited Articles",
			tagName:     "",
			username:    "",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Favorited Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Test Pagination with Limit and Offset",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      20,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 21}, Title: "Article 21"},
				{Model: gorm.Model{ID: 22}, Title: "Article 22"},
			},
			expectedErr: nil,
		},
		{
			name:        "Combine Multiple Filters",
			tagName:     "technology",
			username:    "john_doe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "John's Tech Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Handle Empty Result Set",
			tagName:     "nonexistent",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
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
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(errors.New("database error"))
				return db
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.mockDB()
			store := &ArticleStore{db: mockDB}

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expected, articles)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

