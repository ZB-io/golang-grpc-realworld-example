package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"reflect"
	"time"
	"fmt"
	"github.com/stretchr/testify/assert"
)









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
			name:    "Successfully Delete an Existing Article",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			dbSetup: func(db *gorm.DB) {
				db.AddError(nil)
			},
			wantErr: false,
		},
		{
			name:    "Attempt to Delete a Non-existent Article",
			article: &model.Article{Model: gorm.Model{ID: 999}},
			dbSetup: func(db *gorm.DB) {
				db.AddError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:    "Database Connection Error During Deletion",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &gorm.DB{}
			tt.dbSetup(mockDB)

			s := &ArticleStore{
				db: mockDB,
			}

			err := s.Delete(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
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
			name: "Successfully retrieve an existing comment",
			fields: fields{
				db: mockDBWithComment(1, "Test comment", 1, 1),
			},
			args: args{id: 1},
			want: &model.Comment{
				Model:     gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: nil,
		},
		{
			name: "Attempt to retrieve a non-existent comment",
			fields: fields{
				db: mockDBWithNoComments(),
			},
			args:    args{id: 999},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Handle database connection error",
			fields: fields{
				db: mockDBWithError(errors.New("database connection error")),
			},
			args:    args{id: 1},
			want:    nil,
			wantErr: errors.New("database connection error"),
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

func mockDBWithComment(id uint, body string, userID, articleID uint) *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.AutoMigrate(&model.Comment{})
	comment := model.Comment{
		Model:     gorm.Model{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Body:      body,
		UserID:    userID,
		ArticleID: articleID,
	}
	db.Create(&comment)
	return db
}

func mockDBWithError(err error) *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.AddError(err)
	return db
}

func mockDBWithNoComments() *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.AutoMigrate(&model.Comment{})
	return db
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
			want: func() []model.Tag {
				var tags []model.Tag
				for i := 1; i <= 1000; i++ {
					tags = append(tags, model.Tag{Model: gorm.Model{ID: uint(i)}, Name: fmt.Sprintf("tag%d", i)})
				}
				return tags
			}(),
			wantErr: false,
		},
		{
			name: "Duplicate Tag Names",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag2"})
				db.Create(&model.Tag{Name: "tag1"})
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
				{Model: gorm.Model{ID: 3}, Name: "tag1"},
			},
			wantErr: false,
		},
		{
			name: "Deleted Tags Handling",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag2"})
				db.Create(&model.Tag{Name: "tag3"})
				db.Delete(&model.Tag{}, "name = ?", "tag2")
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 3}, Name: "tag3"},
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
		want    *model.Article
		wantErr error
		dbSetup func(*gorm.DB)
	}{
		{
			name: "Successful Retrieval of an Existing Article",
			id:   1,
			want: &model.Article{
				Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title: "Test Article",
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "tag1"},
					{Model: gorm.Model{ID: 2}, Name: "tag2"},
				},
				Author: model.User{Model: gorm.Model{ID: 1}, Username: "testuser"},
			},
			wantErr: nil,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{
					Model:  gorm.Model{ID: 1},
					Title:  "Test Article",
					UserID: 1,
					Tags: []model.Tag{
						{Model: gorm.Model{ID: 1}, Name: "tag1"},
						{Model: gorm.Model{ID: 2}, Name: "tag2"},
					},
				})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "testuser"})
			},
		},
		{
			name:    "Attempt to Retrieve a Non-existent Article",
			id:      9999,
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
			dbSetup: func(db *gorm.DB) {},
		},
		{
			name:    "Database Connection Error",
			id:      1,
			want:    nil,
			wantErr: errors.New("database connection error"),
			dbSetup: func(db *gorm.DB) {

				db.AddError(errors.New("database connection error"))
			},
		},
		{
			name: "Retrieval of Article with No Tags",
			id:   2,
			want: &model.Article{
				Model:  gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:  "Tagless Article",
				Tags:   []model.Tag{},
				Author: model.User{Model: gorm.Model{ID: 2}, Username: "taglessuser"},
			},
			wantErr: nil,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{
					Model:  gorm.Model{ID: 2},
					Title:  "Tagless Article",
					UserID: 2,
				})
				db.Create(&model.User{Model: gorm.Model{ID: 2}, Username: "taglessuser"})
			},
		},
		{
			name: "Retrieval of Article with Multiple Tags",
			id:   3,
			want: &model.Article{
				Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title: "Multi-tagged Article",
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 3}, Name: "tag3"},
					{Model: gorm.Model{ID: 4}, Name: "tag4"},
					{Model: gorm.Model{ID: 5}, Name: "tag5"},
				},
				Author: model.User{Model: gorm.Model{ID: 3}, Username: "multitaguser"},
			},
			wantErr: nil,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{
					Model:  gorm.Model{ID: 3},
					Title:  "Multi-tagged Article",
					UserID: 3,
					Tags: []model.Tag{
						{Model: gorm.Model{ID: 3}, Name: "tag3"},
						{Model: gorm.Model{ID: 4}, Name: "tag4"},
						{Model: gorm.Model{ID: 5}, Name: "tag5"},
					},
				})
				db.Create(&model.User{Model: gorm.Model{ID: 3}, Username: "multitaguser"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.Tag{}, &model.User{})

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			got, err := s.GetByID(tt.id)

			if !errors.Is(err, tt.wantErr) {
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
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) 

 */
func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name           string
		setupMockDB    func() *gorm.DB
		article        *model.Article
		user           *model.User
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Article is favorited by the user",
			setupMockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Exec("CREATE TABLE favorite_articles (article_id int, user_id int)")
				db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (1, 1)")
				return db
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Article is not favorited by the user",
			setupMockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Exec("CREATE TABLE favorite_articles (article_id int, user_id int)")
				return db
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Nil article parameter",
			setupMockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db
			},
			article:        nil,
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Nil user parameter",
			setupMockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Database error",
			setupMockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.AddError(errors.New("database error"))
				return db
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: false,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Multiple favorites for the same article and user",
			setupMockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Exec("CREATE TABLE favorite_articles (article_id int, user_id int)")
				db.Exec("INSERT INTO favorite_articles (article_id, user_id) VALUES (1, 1), (1, 1)")
				return db
			},
			article:        &model.Article{Model: gorm.Model{ID: 1}},
			user:           &model.User{Model: gorm.Model{ID: 1}},
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setupMockDB()
			store := &ArticleStore{db: mockDB}

			result, err := store.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
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
			name:    "Successful Retrieval of Feed Articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.InstantSet("gorm:auto_preload", true)
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}}},
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
			name:    "Pagination with Offset",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.InstantSet("gorm:auto_preload", true)
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
				{Model: gorm.Model{ID: 3}, Title: "Article 3", UserID: 2, Author: model.User{Model: gorm.Model{ID: 2}}},
			},
			wantErr: false,
		},
		{
			name:    "Error Handling for Database Issues",
			userIDs: []uint{1},
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
			name:    "Limit Exceeds Available Articles",
			userIDs: []uint{1},
			limit:   100,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.InstantSet("gorm:auto_preload", true)
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
			},
			wantErr: false,
		},
		{
			name:    "Single User ID",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db.InstantSet("gorm:auto_preload", true)
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 1, Author: model.User{Model: gorm.Model{ID: 1}}},
			},
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

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.expected)
			}

			if len(got) > int(tt.limit) {
				t.Errorf("ArticleStore.GetFeedArticles() returned more articles than the limit: got %d, limit %d", len(got), tt.limit)
			}

			for _, article := range got {
				found := false
				for _, userID := range tt.userIDs {
					if article.UserID == userID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ArticleStore.GetFeedArticles() returned an article with UserID %d, which is not in the requested userIDs", article.UserID)
				}

				if article.Author.ID == 0 {
					t.Errorf("ArticleStore.GetFeedArticles() returned an article with unloaded Author: %+v", article)
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) 

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

			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1"},
				{Model: gorm.Model{ID: 2}, Title: "Article 2"},
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

			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Golang Article"},
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

			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 3}, Title: "John's Article"},
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
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 4}, Title: "Favorited Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Combine Multiple Filters",
			tagName:     "tech",
			username:    "janedoe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 5}, Title: "Jane's Tech Article"},
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
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Test Pagination with Large Offset",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      1000,
			mockSetup: func(db *gorm.DB) {
				db.AddError(nil)

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
				db.AddError(errors.New("database error"))
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database error"),
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

