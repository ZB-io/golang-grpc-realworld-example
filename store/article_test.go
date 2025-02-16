package golang-grpc-realworld-example

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"time"
	"reflect"
)








/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377

FUNCTION_DEF=func (s *ArticleStore) Create(m *model.Article) error 

*/
func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestArticleStoreCreateMultiple(t *testing.T) {
	mockDB := new(MockDB)
	store := &ArticleStore{
		db: mockDB,
	}

	articles := []*model.Article{
		{
			Title:       "Article 1",
			Description: "Description 1",
			Body:        "Body 1",
			UserID:      1,
		},
		{
			Title:       "Article 2",
			Description: "Description 2",
			Body:        "Body 2",
			UserID:      2,
		},
		{
			Title:       "Article 3",
			Description: "Description 3",
			Body:        "Body 3",
			UserID:      3,
		},
	}

	for _, article := range articles {
		mockDB.On("Create", mock.AnythingOfType("*model.Article")).Return(&gorm.DB{Error: nil}).Once()
	}

	for _, article := range articles {
		err := store.Create(article)
		assert.NoError(t, err)
	}

	mockDB.AssertExpectations(t)
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
			dbError: errors.New("Body cannot be empty"),
			wantErr: true,
		},
		{
			name: "Create a Comment with Maximum Allowed Length for Body",
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
			dbError: errors.New("foreign key constraint failed"),
			wantErr: true,
		},
		{
			name: "Create Multiple Comments in Quick Succession",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Duplicate Unique Fields",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			dbError: errors.New("unique constraint failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{createError: tt.dbError}
			store := &ArticleStore{db: mockDB}

			err := store.CreateComment(tt.comment)

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
func (m *mockDB) Delete(value interface{}) *gorm.DB {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.errorOnDelete {
		return &gorm.DB{Error: errors.New("database connection error")}
	}

	article, ok := value.(*model.Article)
	if !ok {
		return &gorm.DB{Error: errors.New("invalid type")}
	}

	if _, exists := m.articles[article.ID]; !exists {
		return &gorm.DB{Error: gorm.ErrRecordNotFound}
	}

	delete(m.articles, article.ID)
	delete(m.favorites, article.ID)

	for _, comment := range m.comments {
		if comment.ArticleID == article.ID {
			delete(m.comments, comment.ID)
		}
	}

	return &gorm.DB{}
}

func TestArticleStoreDeleteConcurrent(t *testing.T) {
	mockDB := newMockDB()
	mockDB.articles[1] = &model.Article{Model: gorm.Model{ID: 1}}

	store := &ArticleStore{db: mockDB}
	article := &model.Article{Model: gorm.Model{ID: 1}}

	var wg sync.WaitGroup
	concurrentCalls := 5
	results := make(chan error, concurrentCalls)

	for i := 0; i < concurrentCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- store.Delete(article)
		}()
	}

	wg.Wait()
	close(results)

	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	if successCount != 1 {
		t.Errorf("Expected exactly one successful deletion, got %d", successCount)
	}

	if _, exists := mockDB.articles[article.ID]; exists {
		t.Errorf("Article was not deleted from the database after concurrent deletion attempts")
	}
}

func newMockDB() *mockDB {
	return &mockDB{
		articles:  make(map[uint]*model.Article),
		tags:      make(map[uint]*model.Tag),
		comments:  make(map[uint]*model.Comment),
		favorites: make(map[uint][]uint),
	}
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12

FUNCTION_DEF=func (s *ArticleStore) DeleteComment(m *model.Comment) error 

*/
func (m *mockDB) Delete(value interface{}) *gorm.DB {
	return m.deleteFunc(value)
}

func TestArticleStoreDeleteComment(t *testing.T) {
	tests := []struct {
		name    string
		comment *model.Comment
		dbFunc  func(interface{}) *gorm.DB
		wantErr error
	}{
		{
			name: "Successfully Delete an Existing Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
			dbFunc: func(value interface{}) *gorm.DB {
				return &gorm.DB{Error: nil}
			},
			wantErr: nil,
		},
		{
			name: "Attempt to Delete a Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{ID: 999},
				Body:  "Non-existent comment",
			},
			dbFunc: func(value interface{}) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 2},
				Body:  "Another test comment",
			},
			dbFunc: func(value interface{}) *gorm.DB {
				return &gorm.DB{Error: errors.New("database connection error")}
			},
			wantErr: errors.New("database connection error"),
		},
		{
			name:    "Delete Comment with Null Comment Pointer",
			comment: nil,
			dbFunc: func(value interface{}) *gorm.DB {
				return &gorm.DB{Error: errors.New("invalid argument")}
			},
			wantErr: errors.New("invalid argument"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{deleteFunc: tt.dbFunc}
			store := &ArticleStore{db: mockDB}

			err := store.DeleteComment(tt.comment)

			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) 

*/
func (m *mockDB) Find(dest interface{}) *gorm.DB {
	return m.findFunc(dest)
}

func (m *mockDB) Joins(query string, args ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Rows() (*sql.Rows, error) {
	return nil, nil
}

func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m
}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52

FUNCTION_DEF=func (s *ArticleStore) GetByID(id uint) (*model.Article, error) 

*/
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b

FUNCTION_DEF=func (s *ArticleStore) GetCommentByID(id uint) (*model.Comment, error) 

*/
func (m *mockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return m.findFunc(dest, conds...)
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs []uint, limit, offset int64) ([]model.Article, error) 

*/
func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Find(out interface{}, where ...interface{}) *MockDB {
	args := m.Called(out, where)
	return args.Get(0).(*MockDB)
}

func (m *MockDB) Limit(limit interface{}) *MockDB {
	m.Called(limit)
	return m
}

func (m *MockDB) Offset(offset interface{}) *MockDB {
	m.Called(offset)
	return m
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *MockDB {
	args := m.Called(column, conditions)
	return args.Get(0).(*MockDB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	m.Called(query, args)
	return m
}


/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0

FUNCTION_DEF=func (s *ArticleStore) GetTags() ([]model.Tag, error) 

*/
func (m *mockDB) Find(out interface{}) *gorm.DB {
	return m.findFunc(out)
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) 

*/
func (m *mockDB) Count(value interface{}) *mockDB {
	*value.(*int) = m.count
	return m
}

func (m *mockDB) Error() error {
	return m.err
}

func (m *mockDB) Table(name string) *mockDB {
	return m
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		dbCount     int
		dbErr       error
		expected    bool
		expectedErr error
	}{
		{
			name:        "Article is favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     1,
			dbErr:       nil,
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "Article is not favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbErr:       nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil Article parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbErr:       nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			dbCount:     0,
			dbErr:       nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Database error",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbErr:       errors.New("database error"),
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:        "Multiple favorites",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     2,
			dbErr:       nil,
			expected:    true,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				count: tt.dbCount,
				err:   tt.dbErr,
			}
			store := &ArticleStore{
				db: mockDB,
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expected {
				t.Errorf("Expected result %v, but got %v", tt.expected, result)
			}

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.expectedErr, err)
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *mockDB {
	return m
}


/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe

FUNCTION_DEF=func (s *ArticleStore) Update(m *model.Article) error 

*/
func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Update(attrs ...interface{}) *gorm.DB {
	args := m.Called(attrs...)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error 

*/
func (m *mockAssociation) Append(values ...interface{}) error {
	args := m.Called(values...)
	return args.Error(0)
}

func (m *mockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

func (m *mockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestArticleStoreAddFavoriteConcurrent(t *testing.T) {

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	store := &ArticleStore{db: db}
	article := &model.Article{Title: "Test Article", FavoritesCount: 0}
	db.Create(article)

	numUsers := 10
	var wg sync.WaitGroup
	wg.Add(numUsers)

	for i := 0; i < numUsers; i++ {
		go func(i int) {
			defer wg.Done()
			user := &model.User{Username: fmt.Sprintf("user%d", i)}
			db.Create(user)
			err := store.AddFavorite(article, user)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	db.First(article, article.ID)

	assert.Equal(t, int32(numUsers), article.FavoritesCount)
	assert.Len(t, article.FavoritedUsers, numUsers)
}

func (m *mockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error 

*/
func (m *mockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

func (m *mockAssociation) Error() error {
	args := m.Called()
	return args.Error(0)
}

func TestArticleStoreDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedCommit bool
	}{
		{
			name: "Successfully Delete a Favorite Article",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", mock.Anything).Return(m)
				m.On("Commit").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
		{
			name: "Attempt to Delete a Non-existent Favorite",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", mock.Anything).Return(m)
				m.On("Commit").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
		{
			name: "Database Error During Association Deletion",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(errors.New("association deletion error"))
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Rollback").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  errors.New("association deletion error"),
			expectedCount:  1,
			expectedCommit: false,
		},
		{
			name: "Database Error During FavoritesCount Update",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{Error: errors.New("update error")})
				m.On("Rollback").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  errors.New("update error"),
			expectedCount:  1,
			expectedCommit: false,
		},
		{
			name: "Delete Favorite When FavoritesCount is Already Zero",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", mock.Anything).Return(m)
				m.On("Commit").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{db: mockDB}

			err := store.DeleteFavorite(tt.article, tt.user)

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

