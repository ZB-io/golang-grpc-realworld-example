package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"time"
	"sync"
	"github.com/stretchr/testify/assert"
)






type mockDB struct {
	createError error
}
type MockDB struct {
	BeginCalled       bool
	CommitCalled      bool
	RollbackCalled    bool
	DeleteAssociation func(interface{}) error
	UpdateModel       func(interface{}, ...interface{}) error
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

			if reflect.TypeOf(got) != reflect.TypeOf(&ArticleStore{}) {
				t.Errorf("NewArticleStore() returned incorrect type: got %T, want *ArticleStore", got)
			}
		})
	}
}

func TestNewArticleStoreMultipleInstances(t *testing.T) {
	db1 := &gorm.DB{}
	db2 := &gorm.DB{}

	store1 := NewArticleStore(db1)
	store2 := NewArticleStore(db2)

	if store1 == store2 {
		t.Error("NewArticleStore created identical instances for different DB connections")
	}

	if store1.db != db1 || store2.db != db2 {
		t.Error("NewArticleStore did not assign the correct DB reference")
	}
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error 

 */
func (m *mockDB) Create(value interface{}) *gorm.DB {
	return &gorm.DB{Error: m.createError}
}

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
			name: "Create a Comment with Maximum Length Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 1000)),
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
			name: "Attempt to Create a Comment When Database Connection Fails",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: errors.New("database connection failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{createError: tt.dbError}
			store := &ArticleStore{db: &gorm.DB{Value: mockDB}}

			err := store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {

				if tt.comment.ID == 0 {
					t.Errorf("Comment was not assigned an ID")
				}
				if tt.comment.CreatedAt.IsZero() {
					t.Errorf("Comment CreatedAt was not set")
				}
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
	type args struct {
		userIDs []uint
		limit   int64
		offset  int64
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
		mockDB  func() *gorm.DB
	}{
		{
			name: "Successful retrieval of feed articles",
			args: args{
				userIDs: []uint{1, 2},
				limit:   10,
				offset:  0,
			},
			want: []model.Article{
				{
					Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:       "Article 1",
					Description: "Description 1",
					Body:        "Body 1",
					UserID:      1,
					Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
				},
				{
					Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:       "Article 2",
					Description: "Description 2",
					Body:        "Body 2",
					UserID:      2,
					Author:      model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
				},
			},
			wantErr: false,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
		},
		{
			name: "Empty result set",
			args: args{
				userIDs: []uint{99, 100},
				limit:   10,
				offset:  0,
			},
			want:    []model.Article{},
			wantErr: false,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
		},
		{
			name: "Database error handling",
			args: args{
				userIDs: []uint{1, 2},
				limit:   10,
				offset:  0,
			},
			want:    nil,
			wantErr: true,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(errors.New("database error"))
				return db
			},
		},
		{
			name: "Limit and offset functionality",
			args: args{
				userIDs: []uint{1, 2, 3},
				limit:   2,
				offset:  1,
			},
			want: []model.Article{
				{
					Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:       "Article 2",
					Description: "Description 2",
					Body:        "Body 2",
					UserID:      2,
					Author:      model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
				},
				{
					Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:       "Article 3",
					Description: "Description 3",
					Body:        "Body 3",
					UserID:      3,
					Author:      model.User{Model: gorm.Model{ID: 3}, Username: "user3"},
				},
			},
			wantErr: false,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
		},
		{
			name: "Large number of user IDs",
			args: args{
				userIDs: func() []uint {
					ids := make([]uint, 1000)
					for i := range ids {
						ids[i] = uint(i + 1)
					}
					return ids
				}(),
				limit:  50,
				offset: 0,
			},
			want: func() []model.Article {
				articles := make([]model.Article, 50)
				for i := range articles {
					articles[i] = model.Article{
						Model:       gorm.Model{ID: uint(i + 1), CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Article",
						Description: "Description",
						Body:        "Body",
						UserID:      uint(i + 1),
						Author:      model.User{Model: gorm.Model{ID: uint(i + 1)}, Username: "user"},
					}
				}
				return articles
			}(),
			wantErr: false,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}
			got, err := s.GetFeedArticles(tt.args.userIDs, tt.args.limit, tt.args.offset)
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


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error 

 */
func (m *MockDB) Association(column string) *MockDB {
	return m
}

func (m *MockDB) Begin() *MockDB {
	m.BeginCalled = true
	return m
}

func (m *MockDB) Commit() *MockDB {
	m.CommitCalled = true
	return m
}

func (m *MockDB) Delete(value interface{}) error {
	return m.DeleteAssociation(value)
}

func (s *MockArticleStore) DeleteFavorite(a *model.Article, u *model.User) error {
	tx := s.db.Begin()

	err := tx.Model(a).Association("FavoritedUsers").Delete(u)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(a).Update("favorites_count", gorm.Expr("favorites_count - ?", 1))
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	a.FavoritesCount--

	return nil
}

func (m *MockDB) Model(value interface{}) *MockDB {
	return m
}

func (m *MockDB) Rollback() *MockDB {
	m.RollbackCalled = true
	return m
}

func TestArticleStoreDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		user           *model.User
		mockDB         *MockDB
		expectedError  error
		expectedCount  int32
		expectedCommit bool
	}{
		{
			name: "Successfully Delete a Favorite Article",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
				FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
			},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			mockDB:         &MockDB{},
			expectedError:  nil,
			expectedCount:  1,
			expectedCommit: true,
		},
		{
			name: "Attempt to Delete a Non-existent Favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 1,
			},
			user:           &model.User{Model: gorm.Model{ID: 2}},
			mockDB:         &MockDB{},
			expectedError:  nil,
			expectedCount:  1,
			expectedCommit: true,
		},
		{
			name: "Database Error During Association Deletion",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: &MockDB{
				DeleteAssociation: func(interface{}) error {
					return errors.New("association deletion error")
				},
			},
			expectedError:  errors.New("association deletion error"),
			expectedCount:  2,
			expectedCommit: false,
		},
		{
			name: "Database Error During FavoritesCount Update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: &MockDB{
				DeleteAssociation: func(interface{}) error { return nil },
				UpdateModel: func(interface{}, ...interface{}) error {
					return errors.New("update error")
				},
			},
			expectedError:  errors.New("update error"),
			expectedCount:  2,
			expectedCommit: false,
		},
		{
			name: "Delete Favorite for Article with Zero FavoritesCount",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			mockDB:         &MockDB{},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockDB.DeleteAssociation == nil {
				tt.mockDB.DeleteAssociation = func(interface{}) error { return nil }
			}
			if tt.mockDB.UpdateModel == nil {
				tt.mockDB.UpdateModel = func(interface{}, ...interface{}) error { return nil }
			}

			store := &MockArticleStore{db: tt.mockDB}
			err := store.DeleteFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
			assert.Equal(t, tt.expectedCommit, tt.mockDB.CommitCalled)
			assert.Equal(t, !tt.expectedCommit, tt.mockDB.RollbackCalled)
		})
	}
}

func TestArticleStoreDeleteFavoriteConcurrent(t *testing.T) {
	article := &model.Article{
		Model:          gorm.Model{ID: 1},
		FavoritesCount: 5,
		FavoritedUsers: []model.User{
			{Model: gorm.Model{ID: 1}},
			{Model: gorm.Model{ID: 2}},
			{Model: gorm.Model{ID: 3}},
			{Model: gorm.Model{ID: 4}},
			{Model: gorm.Model{ID: 5}},
		},
	}

	mockDB := &MockDB{
		DeleteAssociation: func(interface{}) error { return nil },
		UpdateModel:       func(interface{}, ...interface{}) error { return nil },
	}

	store := &MockArticleStore{db: mockDB}

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(userID uint) {
			defer wg.Done()
			user := &model.User{Model: gorm.Model{ID: userID}}
			err := store.DeleteFavorite(article, user)
			assert.NoError(t, err)
		}(uint(i))
	}

	wg.Wait()

	assert.Equal(t, int32(0), article.FavoritesCount)
}

func (m *MockDB) Update(column string, value interface{}) error {
	return m.UpdateModel(column, value)
}

