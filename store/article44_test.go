package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"math"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"fmt"
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
	b.ResetTimer()
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
			name: "Create ArticleStore with a Valid Database Connection",
			db:   &gorm.DB{},
			want: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name: "Create ArticleStore with a Nil Database Connection",
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
}

func TestNewArticleStoreFields(t *testing.T) {
	db := &gorm.DB{}
	store := NewArticleStore(db)

	storeType := reflect.TypeOf(*store)
	if storeType.NumField() != 1 {
		t.Errorf("ArticleStore has %d fields, want 1", storeType.NumField())
	}

	dbField, ok := storeType.FieldByName("db")
	if !ok {
		t.Error("ArticleStore does not have a 'db' field")
	}

	if dbField.Type != reflect.TypeOf(&gorm.DB{}) {
		t.Errorf("ArticleStore 'db' field is of type %v, want *gorm.DB", dbField.Type)
	}
}

func TestNewArticleStoreInstanceIndependence(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewArticleStore(db)
	store2 := NewArticleStore(db)

	if store1 == store2 {
		t.Error("NewArticleStore() returned the same instance for multiple calls")
	}

	if store1.db != store2.db {
		t.Error("NewArticleStore() returned instances with different db references")
	}
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b

FUNCTION_DEF=func (s *ArticleStore) GetCommentByID(id uint) (*model.Comment, error) 

 */
func TestArticleStoreGetCommentById(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Comment
		wantErr error
	}{
		{
			name: "Successfully Retrieve an Existing Comment",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{id: 1},
			want: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
			wantErr: nil,
		},
		{
			name: "Attempt to Retrieve a Non-existent Comment",
			fields: fields{
				db: &gorm.DB{},
			},
			args:    args{id: 999},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Handle Database Connection Error",
			fields: fields{
				db: &gorm.DB{},
			},
			args:    args{id: 1},
			want:    nil,
			wantErr: errors.New("database connection error"),
		},
		{
			name: "Retrieve Comment with Associated Data",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{id: 2},
			want: &model.Comment{
				Model: gorm.Model{ID: 2},
				Body:  "Comment with associations",
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
				},
				Article: model.Article{
					Model: gorm.Model{ID: 1},
					Title: "Test Article",
				},
			},
			wantErr: nil,
		},
		{
			name: "Performance Test with Large ID",
			fields: fields{
				db: &gorm.DB{},
			},
			args:    args{id: math.MaxUint32},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.fields.db,
			}
			got, err := s.GetCommentByID(tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArticleStore.GetCommentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetCommentByID() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0

FUNCTION_DEF=func (s *ArticleStore) GetTags() ([]model.Tag, error) 

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
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
				{Model: gorm.Model{ID: 3}, Name: "tag3"},
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
			want:    []model.Tag{},
			wantErr: true,
		},
		{
			name: "Large Number of Tags",
			dbSetup: func(db *gorm.DB) {
				for i := 1; i <= 1000; i++ {
					db.Create(&model.Tag{Name: fmt.Sprintf("tag%d", i)})
				}
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate Tag Names",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "duplicate"})
				db.Create(&model.Tag{Name: "duplicate"})
				db.Create(&model.Tag{Name: "unique"})
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "duplicate"},
				{Model: gorm.Model{ID: 2}, Name: "duplicate"},
				{Model: gorm.Model{ID: 3}, Name: "unique"},
			},
			wantErr: false,
		},
		{
			name: "Deleted Tags Handling",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "active1"})
				db.Create(&model.Tag{Name: "deleted"}).Delete(&model.Tag{})
				db.Create(&model.Tag{Name: "active2"})
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "active1"},
				{Model: gorm.Model{ID: 3}, Name: "active2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := &gorm.DB{}
			tt.dbSetup(db)

			s := &ArticleStore{db: db}
			got, err := s.GetTags()

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "Large Number of Tags" {
				if len(got) != 1000 {
					t.Errorf("ArticleStore.GetTags() returned %d tags, want 1000", len(got))
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error 

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
		setupMock      func(*MockDB, *MockAssociation)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedAppend bool
	}{
		{
			name: "Successfully Add Favorite",
			setupMock: func(db *MockDB, assoc *MockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Append", mock.Anything).Return(nil)
				db.On("Update", "favorites_count", mock.Anything).Return(db)
				db.On("Commit").Return(db)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  1,
			expectedAppend: true,
		},
		{
			name: "Database Error During Association",
			setupMock: func(db *MockDB, assoc *MockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Append", mock.Anything).Return(errors.New("database error"))
				db.On("Rollback").Return(db)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("database error"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name: "Database Error During FavoritesCount Update",
			setupMock: func(db *MockDB, assoc *MockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Append", mock.Anything).Return(nil)
				db.On("Update", "favorites_count", mock.Anything).Return(db).Return(errors.New("update error"))
				db.On("Rollback").Return(db)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("update error"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name:           "Add Favorite with Nil Article",
			setupMock:      func(db *MockDB, assoc *MockAssociation) {},
			article:        nil,
			user:           &model.User{},
			expectedError:  errors.New("invalid input: article is nil"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name:           "Add Favorite with Nil User",
			setupMock:      func(db *MockDB, assoc *MockAssociation) {},
			article:        &model.Article{FavoritesCount: 0},
			user:           nil,
			expectedError:  errors.New("invalid input: user is nil"),
			expectedCount:  0,
			expectedAppend: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockAssoc := new(MockAssociation)
			tt.setupMock(mockDB, mockAssoc)

			db := &gorm.DB{
				Value: mockDB,
			}

			store := &ArticleStore{db: db}
			err := store.AddFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			if tt.article != nil {
				assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
			}

			mockDB.AssertExpectations(t)
			mockAssoc.AssertExpectations(t)
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
				return db.Preload("Author").Offset(0).Limit(10)
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
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Preload("Author").
					Joins("join article_tags on articles.id = article_tags.article_id join tags on tags.id = article_tags.tag_id").
					Where("tags.name = ?", "golang").
					Offset(0).Limit(10)
			},
			expected: []model.Article{
				{Title: "Golang Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Filter Articles by Author Username",
			tagName:     "",
			username:    "john",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Preload("Author").
					Joins("join users on articles.user_id = users.id").
					Where("users.username = ?", "john").
					Offset(0).Limit(10)
			},
			expected: []model.Article{
				{Title: "John's Article"},
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
				return db.Preload("Author").
					Where("id in (?)", []uint{1, 2}).
					Offset(0).Limit(10)
			},
			expected: []model.Article{
				{Title: "Favorited Article 1"},
				{Title: "Favorited Article 2"},
			},
			expectedErr: nil,
		},
		{
			name:        "Test Pagination",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       5,
			offset:      10,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Preload("Author").Offset(10).Limit(5)
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
			username:    "john",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.Preload("Author").
					Joins("join users on articles.user_id = users.id").
					Where("users.username = ?", "john").
					Joins("join article_tags on articles.id = article_tags.article_id join tags on tags.id = article_tags.tag_id").
					Where("tags.name = ?", "golang").
					Where("id in (?)", []uint{1}).
					Offset(0).Limit(10)
			},
			expected: []model.Article{
				{Title: "John's Golang Article"},
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
				return db.Preload("Author").
					Joins("join article_tags on articles.id = article_tags.article_id join tags on tags.id = article_tags.tag_id").
					Where("tags.name = ?", "nonexistent").
					Offset(0).Limit(10)
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
				return db.Preload("Author").Offset(0).Limit(10)
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &ArticleStore{
				db: tt.mockDB(),
			}

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expected, articles)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

