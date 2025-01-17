package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"fmt"
	"time"
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
}

func TestNewArticleStoreDifferentConnections(t *testing.T) {
	db1 := &gorm.DB{}
	db2 := &gorm.DB{}
	store1 := NewArticleStore(db1)
	store2 := NewArticleStore(db2)

	if store1.db == store2.db {
		t.Error("ArticleStore instances should have different DB references")
	}
}

func TestNewArticleStoreImmutability(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewArticleStore(db)
	store2 := NewArticleStore(db)

	if store1 == store2 {
		t.Error("NewArticleStore should return different instances")
	}
	if store1.db != store2.db {
		t.Error("ArticleStore instances should have the same DB reference")
	}
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error 

 */
func TestArticleStoreCreateComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		dbError error
		wantErr bool
	}{
		{
			name: "Successfully Create a New Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			comment: &model.Comment{

				UserID:    1,
				ArticleID: 1,
			},
			dbError: errors.New("required field missing"),
			wantErr: true,
		},
		{
			name: "Create Comment with Very Long Body Text",
			comment: &model.Comment{
				Body:      string(make([]byte, 10000)),
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment for a Non-existent Article",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 9999,
			},
			dbError: errors.New("foreign key constraint violation"),
			wantErr: true,
		},
		{
			name: "Create Comment with Special Characters in the Body",
			comment: &model.Comment{
				Body:      "Test comment with special characters: !@#$%^&*()_+ and emojis: 😊🎉",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &gorm.DB{}

			mockDB.AddError(tt.dbError)

			s := &ArticleStore{
				db: mockDB,
			}

			err := s.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 2},
				Body:  "Error comment",
			},
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			wantErr: true,
		},
		{
			name: "Delete Comment with Null Fields",
			comment: &model.Comment{
				Model: gorm.Model{ID: 3},
				Body:  "",
			},
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Comment{Model: gorm.Model{ID: 3}, Body: ""})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open mock database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Comment{})

			tt.dbSetup(db)

			store := &ArticleStore{db: db}

			err = store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
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
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b

FUNCTION_DEF=func (s *ArticleStore) GetCommentByID(id uint) (*model.Comment, error) 

 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (s *MockArticleStore) GetCommentByID(id uint) (*model.Comment, error) {
	var m model.Comment
	err := s.db.Find(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func TestArticleStoreGetCommentById(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		mockSetup       func(*MockDB)
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name: "Successfully retrieve an existing comment",
			id:   1,
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"), uint(1)).Return(&gorm.DB{Error: nil}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.Comment)
					*arg = model.Comment{
						Model: gorm.Model{ID: 1},
						Body:  "Test comment",
					}
				})
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
		},
		{
			name: "Attempt to retrieve a non-existent comment",
			id:   999,
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"), uint(999)).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Handle database connection error",
			id:   2,
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Find", mock.AnythingOfType("*model.Comment"), uint(2)).Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			expectedError:   errors.New("database connection error"),
			expectedComment: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			store := &MockArticleStore{db: mockDB}

			comment, err := store.GetCommentByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedComment, comment)

			mockDB.AssertExpectations(t)
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
				for i := 1; i <= 10000; i++ {
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

			if tt.name == "Large Number of Tags" {
				if len(got) != 10000 {
					t.Errorf("ArticleStore.GetTags() returned %d tags, want 10000", len(got))
				}
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

FUNCTION_DEF=func (s *ArticleStore) GetByID(id uint) (*model.Article, error) 

 */
func TestArticleStoreGetById(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mockDB  func() *gorm.DB
		want    *model.Article
		wantErr error
	}{
		{
			name: "Successful Retrieval of an Existing Article",
			id:   1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:mock_article", &model.Article{
					Model: gorm.Model{ID: 1},
					Title: "Test Article",
					Tags:  []model.Tag{{Model: gorm.Model{ID: 1}, Name: "test"}},
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
				})
			},
			want: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
				Tags:  []model.Tag{{Model: gorm.Model{ID: 1}, Name: "test"}},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
				},
			},
			wantErr: nil,
		},
		{
			name: "Attempt to Retrieve a Non-existent Article",
			id:   999,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = gorm.ErrRecordNotFound
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			id:   1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = errors.New("database connection error")
				return db
			},
			want:    nil,
			wantErr: errors.New("database connection error"),
		},
		{
			name: "Retrieval of Article with No Tags",
			id:   2,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:mock_article", &model.Article{
					Model: gorm.Model{ID: 2},
					Title: "Tagless Article",
					Tags:  []model.Tag{},
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
				})
			},
			want: &model.Article{
				Model: gorm.Model{ID: 2},
				Title: "Tagless Article",
				Tags:  []model.Tag{},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
				},
			},
			wantErr: nil,
		},
		{
			name: "Retrieval of Article with Multiple Tags",
			id:   3,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:mock_article", &model.Article{
					Model: gorm.Model{ID: 3},
					Title: "Multi-tagged Article",
					Tags: []model.Tag{
						{Model: gorm.Model{ID: 1}, Name: "tag1"},
						{Model: gorm.Model{ID: 2}, Name: "tag2"},
						{Model: gorm.Model{ID: 3}, Name: "tag3"},
					},
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
				})
			},
			want: &model.Article{
				Model: gorm.Model{ID: 3},
				Title: "Multi-tagged Article",
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "tag1"},
					{Model: gorm.Model{ID: 2}, Name: "tag2"},
					{Model: gorm.Model{ID: 3}, Name: "tag3"},
				},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}
			got, err := s.GetByID(tt.id)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ArticleStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe

FUNCTION_DEF=func (s *ArticleStore) Update(m *model.Article) error 

 */
func TestArticleStoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		mockDB  func() *gorm.DB
		wantErr bool
	}{
		{
			name: "Successfully Update an Existing Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Updated Title",
				Body:  "Updated Body",
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.
					Scopes(func(d *gorm.DB) *gorm.DB {
						d.Error = nil
						d.RowsAffected = 1
						return d
					})
			},
			wantErr: false,
		},
		{
			name: "Attempt to Update a Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.
					Scopes(func(d *gorm.DB) *gorm.DB {
						d.Error = gorm.ErrRecordNotFound
						return d
					})
			},
			wantErr: true,
		},
		{
			name: "Update Article with Invalid Data",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
				Title: "",
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.
					Scopes(func(d *gorm.DB) *gorm.DB {
						d.Error = errors.New("validation error: Title cannot be empty")
						return d
					})
			},
			wantErr: true,
		},
		{
			name: "Handle Database Connection Error During Update",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
				Title: "Connection Error Test",
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.
					Scopes(func(d *gorm.DB) *gorm.DB {
						d.Error = errors.New("database connection error")
						return d
					})
			},
			wantErr: true,
		},
		{
			name: "Update Article with New Tags",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
				Title: "Article with New Tags",
				Tags:  []model.Tag{{Name: "NewTag1"}, {Name: "NewTag2"}},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.
					Scopes(func(d *gorm.DB) *gorm.DB {
						d.Error = nil
						d.RowsAffected = 1
						return d
					})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}

			err := s.Update(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e

FUNCTION_DEF=func (s *ArticleStore) GetComments(m *model.Article) ([]model.Comment, error) 

 */
func TestArticleStoreGetComments(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		mockDB         func() *gorm.DB
		expectedResult []model.Comment
		expectedError  error
	}{
		{
			name: "Successfully retrieve comments for an article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Callback().Query().Register("mock_query", func(scope *gorm.Scope) {
					if scope.Search.Where("article_id = ?", 1) != nil {
						comments := []model.Comment{
							{Model: gorm.Model{ID: 1}, Body: "Comment 1", UserID: 1, ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
							{Model: gorm.Model{ID: 2}, Body: "Comment 2", UserID: 2, ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
						}
						scope.DB().AddError(scope.DB().Find(scope.Value, comments).Error)
					}
				})
				return db
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 1}, Body: "Comment 1", UserID: 1, ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Body: "Comment 2", UserID: 2, ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			expectedError: nil,
		},
		{
			name: "Retrieve comments for an article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Callback().Query().Register("mock_query", func(scope *gorm.Scope) {
					if scope.Search.Where("article_id = ?", 2) != nil {
						scope.DB().AddError(scope.DB().Find(scope.Value, []model.Comment{}).Error)
					}
				})
				return db
			},
			expectedResult: []model.Comment{},
			expectedError:  nil,
		},
		{
			name: "Handle database error when retrieving comments",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Callback().Query().Register("mock_query", func(scope *gorm.Scope) {
					if scope.Search.Where("article_id = ?", 3) != nil {
						scope.DB().AddError(errors.New("database error"))
					}
				})
				return db
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Verify correct query construction",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Callback().Query().Register("mock_query", func(scope *gorm.Scope) {
					if scope.Search.Where("article_id = ?", 4) != nil && scope.Search.Preload("Author") != nil {
						comments := []model.Comment{
							{Model: gorm.Model{ID: 1}, Body: "Comment 1", UserID: 1, ArticleID: 4, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
						}
						scope.DB().AddError(scope.DB().Find(scope.Value, comments).Error)
					} else {
						scope.DB().AddError(errors.New("incorrect query construction"))
					}
				})
				return db
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 1}, Body: "Comment 1", UserID: 1, ArticleID: 4, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
			},
			expectedError: nil,
		},
		{
			name: "Handle large number of comments",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Callback().Query().Register("mock_query", func(scope *gorm.Scope) {
					if scope.Search.Where("article_id = ?", 5) != nil {
						comments := make([]model.Comment, 1000)
						for i := 0; i < 1000; i++ {
							comments[i] = model.Comment{
								Model:     gorm.Model{ID: uint(i + 1)},
								Body:      "Comment",
								UserID:    1,
								ArticleID: 5,
								Author:    model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
							}
						}
						scope.DB().AddError(scope.DB().Find(scope.Value, comments).Error)
					}
				})
				return db
			},
			expectedResult: func() []model.Comment {
				comments := make([]model.Comment, 1000)
				for i := 0; i < 1000; i++ {
					comments[i] = model.Comment{
						Model:     gorm.Model{ID: uint(i + 1)},
						Body:      "Comment",
						UserID:    1,
						ArticleID: 5,
						Author:    model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
					}
				}
				return comments
			}(),
			expectedError: nil,
		},
		{
			name: "Verify comment order",
			article: &model.Article{
				Model: gorm.Model{ID: 6},
			},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Callback().Query().Register("mock_query", func(scope *gorm.Scope) {
					if scope.Search.Where("article_id = ?", 6) != nil {
						comments := []model.Comment{
							{Model: gorm.Model{ID: 3, CreatedAt: time.Now()}, Body: "Comment 3", UserID: 3, ArticleID: 6, Author: model.User{Model: gorm.Model{ID: 3}, Username: "user3"}},
							{Model: gorm.Model{ID: 2, CreatedAt: time.Now().Add(-1 * time.Hour)}, Body: "Comment 2", UserID: 2, ArticleID: 6, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
							{Model: gorm.Model{ID: 1, CreatedAt: time.Now().Add(-2 * time.Hour)}, Body: "Comment 1", UserID: 1, ArticleID: 6, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
						}
						scope.DB().AddError(scope.DB().Order("created_at desc").Find(scope.Value, comments).Error)
					}
				})
				return db
			},
			expectedResult: []model.Comment{
				{Model: gorm.Model{ID: 3}, Body: "Comment 3", UserID: 3, ArticleID: 6, Author: model.User{Model: gorm.Model{ID: 3}, Username: "user3"}},
				{Model: gorm.Model{ID: 2}, Body: "Comment 2", UserID: 2, ArticleID: 6, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
				{Model: gorm.Model{ID: 1}, Body: "Comment 1", UserID: 1, ArticleID: 6, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &ArticleStore{
				db: tt.mockDB(),
			}

			result, err := store.GetComments(tt.article)

			if !reflect.DeepEqual(err, tt.expectedError) {
				t.Errorf("GetComments() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("GetComments() = %v, expected %v", result, tt.expectedResult)
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
		name           string
		setupMock      func(*gorm.DB)
		article        *model.Article
		user           *model.User
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Article is favorited by the user",
			setupMock: func(db *gorm.DB) {
				db.AddError(nil)
				db.RowsAffected = 1
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Article is not favorited by the user",
			setupMock: func(db *gorm.DB) {
				db.AddError(nil)
				db.RowsAffected = 0
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "Nil article parameter",
			setupMock:      func(db *gorm.DB) {},
			article:        nil,
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "Nil user parameter",
			setupMock:      func(db *gorm.DB) {},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Database error",
			setupMock: func(db *gorm.DB) {
				db.AddError(errors.New("database error"))
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: false,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Multiple favorites for the same article and user",
			setupMock: func(db *gorm.DB) {
				db.AddError(nil)
				db.RowsAffected = 3
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &gorm.DB{}
			tt.setupMock(mockDB)

			store := &ArticleStore{db: mockDB}
			result, err := store.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
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

func (m *MockDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(value, where)
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

func TestArticleStoreDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectRollback bool
		expectCommit   bool
	}{
		{
			name: "Successfully Delete a Favorite",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", mock.Anything).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 0,
			expectCommit:  true,
		},
		{
			name: "Association Deletion Failure",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", mock.Anything).Return(&gorm.Association{Error: errors.New("deletion failed")})
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Rollback").Return(mockDB)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  errors.New("deletion failed"),
			expectedCount:  1,
			expectRollback: true,
		},
		{
			name: "Update Failure",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", mock.Anything).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{Error: errors.New("update failed")})
				mockDB.On("Rollback").Return(mockDB)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  errors.New("update failed"),
			expectedCount:  1,
			expectRollback: true,
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
			err := store.DeleteFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			if tt.expectRollback {
				mockDB.AssertCalled(t, "Rollback")
			}
			if tt.expectCommit {
				mockDB.AssertCalled(t, "Commit")
			}
		})
	}
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}

