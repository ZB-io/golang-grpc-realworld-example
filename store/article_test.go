package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"time"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)






type MockDB struct {
	mock.Mock
}


/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore 

 */
func BenchmarkNewArticleStore(b *testing.B) {
	db := &gorm.DB{}
	for i := 0; i < b.N; i++ {
		NewArticleStore(db)
	}
}

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
				t.Error("NewArticleStore returned nil")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticleStore() = %v, want %v", got, tt.want)
			}

			if got.db != tt.db {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.db)
			}

			if _, ok := interface{}(got).(*ArticleStore); !ok {
				t.Errorf("NewArticleStore() did not return *ArticleStore")
			}
		})
	}
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error 

 */
func TestArticleStoreCreateComment(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		m *model.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successfully Create a Valid Comment",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{
				m: &model.Comment{
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{
				m: &model.Comment{
					Body: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Create a Comment with Maximum Length Body",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{
				m: &model.Comment{
					Body:      string(make([]byte, 1000)),
					UserID:    1,
					ArticleID: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent Article",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{
				m: &model.Comment{
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 9999,
				},
			},
			wantErr: true,
		},
		{
			name: "Attempt to Create a Comment When Database Connection Fails",
			fields: fields{
				db: &gorm.DB{Error: errors.New("database connection failed")},
			},
			args: args{
				m: &model.Comment{
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.fields.db,
			}
			err := s.CreateComment(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
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
		wantErr bool
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
			wantErr: false,
		},
		{
			name: "Attempt to Delete a Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
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
			wantErr: false,
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
				db.Create(&model.Comment{ArticleID: 3, Body: "Test Comment"})
			},
			wantErr: false,
		},
		{
			name: "Delete an Article That's Been Favorited",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
				Title: "Favorited Article",
			},
			dbSetup: func(db *gorm.DB) {
				article := &model.Article{Model: gorm.Model{ID: 4}, Title: "Favorited Article"}
				db.Create(article)
				user := &model.User{Model: gorm.Model{ID: 1}, Username: "testuser"}
				db.Create(user)
				db.Model(user).Association("FavoriteArticles").Append(article)
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error During Deletion",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
				Title: "Error Article",
			},
			dbSetup: func(db *gorm.DB) {

				db.AddError(errors.New("database connection error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open mock database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.Tag{}, &model.Comment{}, &model.User{})

			tt.dbSetup(db)

			store := &ArticleStore{db: db}

			err = store.Delete(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var count int64
				db.Model(&model.Article{}).Where("id = ?", tt.article.ID).Count(&count)
				if count != 0 {
					t.Errorf("ArticleStore.Delete() failed to delete article, count = %d, want 0", count)
				}

				if tt.name == "Delete an Article with Associated Tags" {
					var tagCount int64
					db.Model(&model.Tag{}).Count(&tagCount)
					if tagCount != 0 {
						t.Errorf("ArticleStore.Delete() failed to delete associated tags, count = %d, want 0", tagCount)
					}
				}

				if tt.name == "Delete an Article with Comments" {
					var commentCount int64
					db.Model(&model.Comment{}).Where("article_id = ?", tt.article.ID).Count(&commentCount)
					if commentCount != 0 {
						t.Errorf("ArticleStore.Delete() failed to delete associated comments, count = %d, want 0", commentCount)
					}
				}

				if tt.name == "Delete an Article That's Been Favorited" {
					var favoriteCount int64
					db.Table("favorite_articles").Where("article_id = ?", tt.article.ID).Count(&favoriteCount)
					if favoriteCount != 0 {
						t.Errorf("ArticleStore.Delete() failed to remove favorite associations, count = %d, want 0", favoriteCount)
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
				Model: gorm.Model{ID: 2},
				Body:  "Comment with relationships",
			},
			dbSetup: func(db *gorm.DB) {
				comment := &model.Comment{Model: gorm.Model{ID: 2}, Body: "Comment with relationships"}
				db.Create(comment)

			},
			wantErr: false,
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 3},
				Body:  "Comment for connection error",
			},
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("simulated connection error"))
			},
			wantErr: true,
		},
		{
			name:    "Delete Comment with Null Input",
			comment: nil,
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open mock database: %v", err)
			}
			defer db.Close()

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			err = s.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.comment != nil {
				var count int64
				db.Model(&model.Comment{}).Where("id = ?", tt.comment.ID).Count(&count)
				if count != 0 {
					t.Errorf("Comment was not deleted from the database")
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) 

 */
func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		user    *model.User
		dbSetup func(*gorm.DB) error
		want    bool
		wantErr bool
	}{
		{
			name:    "Article is favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			dbSetup: func(db *gorm.DB) error {
				return db.Table("favorite_articles").Create(map[string]interface{}{
					"article_id": 1,
					"user_id":    1,
				}).Error
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Article is not favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 2}},
			user:    &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(db *gorm.DB) error { return nil },
			want:    false,
			wantErr: false,
		},
		{
			name:    "Null article parameter",
			article: nil,
			user:    &model.User{Model: gorm.Model{ID: 3}},
			dbSetup: func(db *gorm.DB) error { return nil },
			want:    false,
			wantErr: false,
		},
		{
			name:    "Null user parameter",
			article: &model.Article{Model: gorm.Model{ID: 4}},
			user:    nil,
			dbSetup: func(db *gorm.DB) error { return nil },
			want:    false,
			wantErr: false,
		},
		{
			name:    "Database error",
			article: &model.Article{Model: gorm.Model{ID: 5}},
			user:    &model.User{Model: gorm.Model{ID: 5}},
			dbSetup: func(db *gorm.DB) error {
				return errors.New("database error")
			},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Multiple favorites for the same article-user pair",
			article: &model.Article{Model: gorm.Model{ID: 6}},
			user:    &model.User{Model: gorm.Model{ID: 6}},
			dbSetup: func(db *gorm.DB) error {
				err := db.Table("favorite_articles").Create(map[string]interface{}{
					"article_id": 6,
					"user_id":    6,
				}).Error
				if err != nil {
					return err
				}
				return db.Table("favorite_articles").Create(map[string]interface{}{
					"article_id": 6,
					"user_id":    6,
				}).Error
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}
			defer db.Close()

			if err := tt.dbSetup(db); err != nil {
				t.Fatalf("Failed to setup database: %v", err)
			}

			s := &ArticleStore{db: db}

			got, err := s.IsFavorited(tt.article, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs []uint, limit, offset int64) ([]model.Article, error) 

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
			name:    "Successful retrieval of feed articles",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Scopes(func(d *gorm.DB) *gorm.DB {
					return d.Where("user_id in (?)", []uint{1, 2}).Offset(0).Limit(2)
				})
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{99, 100},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Scopes(func(d *gorm.DB) *gorm.DB {
					return d.Where("user_id in (?)", []uint{99, 100}).Offset(0).Limit(10)
				})
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Database error handling",
			userIDs: []uint{1, 2},
			limit:   5,
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
			name:    "Limit and offset functionality",
			userIDs: []uint{1, 2, 3},
			limit:   2,
			offset:  1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Scopes(func(d *gorm.DB) *gorm.DB {
					return d.Where("user_id in (?)", []uint{1, 2, 3}).Offset(1).Limit(2)
				})
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
				{Model: gorm.Model{ID: 3}, Title: "Article 3", UserID: 3, Author: model.User{Model: gorm.Model{ID: 3}, Username: "user3"}},
			},
			wantErr: false,
		},
		{
			name:    "Preloading of Author information",
			userIDs: []uint{1},
			limit:   1,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Scopes(func(d *gorm.DB) *gorm.DB {
					return d.Where("user_id in (?)", []uint{1}).Offset(0).Limit(1)
				})
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1", Email: "user1@example.com"}},
			},
			wantErr: false,
		},
		{
			name: "Large number of user IDs",
			userIDs: func() []uint {
				ids := make([]uint, 1000)
				for i := range ids {
					ids[i] = uint(i + 1)
				}
				return ids
			}(),
			limit:  10,
			offset: 0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Scopes(func(d *gorm.DB) *gorm.DB {
					return d.Where("user_id in (?)", func() []uint {
						ids := make([]uint, 1000)
						for i := range ids {
							ids[i] = uint(i + 1)
						}
						return ids
					}()).Offset(0).Limit(10)
				})
			},
			expected: func() []model.Article {
				articles := make([]model.Article, 10)
				for i := range articles {
					articles[i] = model.Article{Model: gorm.Model{ID: uint(i + 1)}, Title: "Article", UserID: uint(i + 1), Author: model.User{Model: gorm.Model{ID: uint(i + 1)}, Username: "user"}}
				}
				return articles
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}

			got, err := s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

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
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error 

 */
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

func (m *MockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestArticleStoreDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockDB, *model.Article, *model.User)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedCommit bool
	}{
		{
			name: "Successfully Delete a Favorite",
			setupMock: func(mockDB *MockDB, article *model.Article, user *model.User) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", article).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", user).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
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
			tt.setupMock(mockDB, tt.article, tt.user)

			db := &gorm.DB{
				Value: mockDB,
			}

			store := &ArticleStore{db: db}
			err := store.DeleteFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			if tt.expectedCommit {
				mockDB.AssertCalled(t, "Commit")
			} else {
				mockDB.AssertCalled(t, "Rollback")
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
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) 

 */
func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name         string
		tagName      string
		username     string
		favoritedBy  *model.User
		limit        int64
		offset       int64
		mockDB       func() *gorm.DB
		wantArticles []model.Article
		wantErr      bool
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
			wantArticles: []model.Article{
				{Title: "Article 1"},
				{Title: "Article 2"},
			},
			wantErr: false,
		},
		{
			name:        "Filter Articles by Tag Name",
			tagName:     "golang",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			wantArticles: []model.Article{
				{Title: "Golang Article"},
			},
			wantErr: false,
		},
		{
			name:        "Filter Articles by Author Username",
			tagName:     "",
			username:    "johndoe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			wantArticles: []model.Article{
				{Title: "John's Article"},
			},
			wantErr: false,
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
			wantArticles: []model.Article{
				{Title: "Favorited Article"},
			},
			wantErr: false,
		},
		{
			name:        "Combine Multiple Filters",
			tagName:     "golang",
			username:    "johndoe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			wantArticles: []model.Article{
				{Title: "John's Golang Article"},
			},
			wantErr: false,
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
			wantArticles: []model.Article{},
			wantErr:      false,
		},
		{
			name:        "Test Pagination with Large Offset",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      1000,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			wantArticles: []model.Article{},
			wantErr:      false,
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
			wantArticles: []model.Article{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}

			gotArticles, err := s.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantArticles, gotArticles)
		})
	}
}

