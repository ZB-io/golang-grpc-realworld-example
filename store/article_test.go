package store

import (
		"reflect"
		"testing"
		"time"
		"github.com/jinzhu/gorm"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/require"
		"errors"
		"github.com/stretchr/testify/mock"
		"fmt"
		"github.com/DATA-DOG/go-sqlmock"
		"sync"
)


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
			db: &gorm.DB{
				Value: "test_db",
			},
			want: &ArticleStore{
				db: &gorm.DB{
					Value: "test_db",
				},
			},
		},
		{
			name: "Create ArticleStore with Nil DB Connection",
			db:   nil,
			want: &ArticleStore{
				db: nil,
			},
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

func TestNewArticleStoreDBFieldAccessibility(t *testing.T) {
	db := &gorm.DB{
		Value: "test_db",
	}

	store := NewArticleStore(db)

	if !reflect.DeepEqual(store.db, db) {
		t.Errorf("NewArticleStore() db field does not match the provided DB instance")
	}
}

func TestNewArticleStoreImmutability(t *testing.T) {
	db := &gorm.DB{
		Value: "test_db",
	}

	store1 := NewArticleStore(db)
	store2 := NewArticleStore(db)

	if store1 == store2 {
		t.Errorf("NewArticleStore() returned the same instance for multiple calls")
	}

	if store1.db != store2.db {
		t.Errorf("NewArticleStore() did not use the same DB instance for multiple calls")
	}
}

func TestNewArticleStorePerformance(t *testing.T) {
	db := &gorm.DB{
		Value: "test_db",
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		NewArticleStore(db)
	}
	duration := time.Since(start)

	if duration > time.Millisecond*100 {
		t.Errorf("NewArticleStore() took too long for 1000 calls: %v", duration)
	}
}


/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		article *model.Article
		wantErr bool
	}{
		{
			name: "Successfully Create a New Article",
			article: &model.Article{
				Title:       "Test Article",
				Description: "This is a test article",
				Body:        "This is the body of the test article",
				UserID:      1,
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create an Article with Missing Required Fields",
			article: &model.Article{

				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "Create an Article with Maximum Length Content",
			article: &model.Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 10000)),
				UserID:      1,
			},
			wantErr: false,
		},
		{
			name: "Create an Article with Associated Tags and Author",
			article: &model.Article{
				Title:       "Article with Tags and Author",
				Description: "This article has tags and an author",
				Body:        "Body of the article with tags and author",
				UserID:      1,
				Tags: []model.Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
				},
			},
			wantErr: false,
		},

		{
			name: "Attempt to Create an Article with Duplicate Title",
			article: &model.Article{
				Title:       "Test Article",
				Description: "This is another test article",
				Body:        "This is the body of another test article",
				UserID:      2,
			},
			wantErr: true,
		},
	}

	db, err := gorm.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db.AutoMigrate(&model.Article{}, &model.Tag{}, &model.User{}).Error
	require.NoError(t, err)

	s := &ArticleStore{db: db}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := db.Exec("DELETE FROM articles").Error
			require.NoError(t, err)

			err = s.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var savedArticle model.Article
				err = db.Preload("Tags").Preload("Author").First(&savedArticle, "title = ?", tt.article.Title).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.article.Title, savedArticle.Title)
				assert.Equal(t, tt.article.Description, savedArticle.Description)
				assert.Equal(t, tt.article.Body, savedArticle.Body)
				assert.Equal(t, tt.article.UserID, savedArticle.UserID)

				assert.Len(t, savedArticle.Tags, len(tt.article.Tags))
				for i, tag := range tt.article.Tags {
					assert.Equal(t, tag.Name, savedArticle.Tags[i].Name)
				}

				if tt.article.Author.ID != 0 {
					assert.Equal(t, tt.article.Author.ID, savedArticle.Author.ID)
					assert.Equal(t, tt.article.Author.Username, savedArticle.Author.Username)
				}
			}
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
		setupDB func(*gorm.DB)
		wantErr bool
	}{
		{
			name: "Successfully Create a Comment",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			wantErr: false,
		},
		{
			name: "Attempt to Create a Comment with Missing Required Fields",
			comment: &model.Comment{
				UserID:    1,
				ArticleID: 1,
			},
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			wantErr: true,
		},
		{
			name: "Create a Comment for a Non-existent Article",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 999,
			},
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
			},
			wantErr: true,
		},
		{
			name: "Create Multiple Comments in Succession",
			comment: &model.Comment{
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
				db.Create(&model.Comment{Body: "Existing comment", UserID: 1, ArticleID: 1})
			},
			wantErr: false,
		},
		{
			name: "Create a Comment with Maximum Allowed Length for Body",
			comment: &model.Comment{
				Body:      string(make([]byte, 65535)),
				UserID:    1,
				ArticleID: 1,
			},
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{ID: 1})
				db.Create(&model.Article{ID: 1})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.Comment{}, &model.User{}, &model.Article{})

			tt.setupDB(db)

			store := &ArticleStore{db: db}

			err = store.CreateComment(tt.comment)

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
		{
			name: "Delete an Article with High Favorites Count",
			setup: func(db *gorm.DB) *model.Article {
				article := &model.Article{
					Title:          "Popular Article",
					Description:    "Very popular",
					Body:           "Everyone likes this",
					FavoritesCount: 2147483647,
				}
				db.Create(article)
				return article
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
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

				if len(article.Tags) > 0 {
					var tagCount int
					db.Model(&model.Tag{}).Where("id IN (?)", article.Tags).Count(&tagCount)
					assert.Equal(t, 0, tagCount)
				}

				if len(article.Comments) > 0 {
					var commentCount int
					db.Model(&model.Comment{}).Where("article_id = ?", article.ID).Count(&commentCount)
					assert.Equal(t, 0, commentCount)
				}

				if len(article.FavoritedUsers) > 0 {
					var favoriteCount int
					db.Table("favorite_articles").Where("article_id = ?", article.ID).Count(&favoriteCount)
					assert.Equal(t, 0, favoriteCount)
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
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
			name: "Handle database connection error",
			setupDB: func(db *gorm.DB) {

				db.Close()
			},
			commentID:       1,
			expectedError:   errors.New("sql: database is closed"),
			expectedComment: nil,
		},
		{
			name: "Retrieve a comment with associated user and article data",
			setupDB: func(db *gorm.DB) {
				user := &model.User{Model: gorm.Model{ID: 1}, Username: "testuser"}
				article := &model.Article{Model: gorm.Model{ID: 1}, Title: "Test Article"}
				comment := &model.Comment{
					Model:     gorm.Model{ID: 1},
					Body:      "Test comment with associations",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(user)
				db.Create(article)
				db.Create(comment)
			},
			commentID:     1,
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment with associations",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name: "Retrieve a soft-deleted comment",
			setupDB: func(db *gorm.DB) {
				deletedAt := time.Now()
				comment := &model.Comment{
					Model:     gorm.Model{ID: 1, DeletedAt: &deletedAt},
					Body:      "Soft-deleted comment",
					UserID:    1,
					ArticleID: 1,
				}
				db.Create(comment)
			},
			commentID:       1,
			expectedError:   gorm.ErrRecordNotFound,
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

			assert.Equal(t, tt.expectedComment, comment)
		})
	}

	t.Run("Performance test with a large number of comments", func(t *testing.T) {
		db, err := gorm.Open("sqlite3", ":memory:")
		assert.NoError(t, err)
		defer db.Close()

		for i := 1; i <= 10000; i++ {
			comment := &model.Comment{
				Model:     gorm.Model{ID: uint(i)},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			}
			db.Create(comment)
		}

		store := &ArticleStore{db: db}

		start := time.Now()
		comment, err := store.GetCommentByID(9999)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, uint(9999), comment.ID)
		assert.Less(t, duration, 50*time.Millisecond, "Query took too long")
	})
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
				for i := 0; i < 1000; i++ {
					db.Create(&model.Tag{Name: fmt.Sprintf("tag%d", i)})
				}
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate tags in database",
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
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}
			defer db.Close()

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			got, err := s.GetTags()

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "Large number of tags" {
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
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func TestGetByID(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		mockSetup       func(*MockDB)
		expectedError   error
		expectedArticle *model.Article
	}{}

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
		name         string
		setupDB      func(*gorm.DB)
		inputArticle *model.Article
		expectedErr  error
		validate     func(*testing.T, *gorm.DB, *model.Article)
	}{

		{
			name: "Update All Fields of an Article",
			setupDB: func(db *gorm.DB) {
				db.Create(&model.Article{
					Model:       gorm.Model{ID: 6},
					Title:       "Original Title",
					Description: "Original Description",
					Body:        "Original Body",
					UserID:      1,
				})
			},
			inputArticle: &model.Article{
				Model:       gorm.Model{ID: 6},
				Title:       "Updated Title",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      2,
			},
			expectedErr: nil,
			validate: func(t *testing.T, db *gorm.DB, a *model.Article) {
				var updatedArticle model.Article
				db.First(&updatedArticle, a.ID)
				assert.Equal(t, "Updated Title", updatedArticle.Title)
				assert.Equal(t, "Updated Description", updatedArticle.Description)
				assert.Equal(t, "Updated Body", updatedArticle.Body)
				assert.Equal(t, uint(2), updatedArticle.UserID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			assert.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.Tag{}, &model.User{})

			tt.setupDB(db)

			store := &ArticleStore{db: db}

			err = store.Update(tt.inputArticle)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, db, tt.inputArticle)
		})
	}
}


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestGetComments(t *testing.T) {
	tests := []struct {
		name             string
		setupMockDB      func() *gorm.DB
		article          *model.Article
		expectedComments []model.Comment
		expectedError    error
	}{
		{
			name: "Successful Retrieval of Comments",
			setupMockDB: func() *gorm.DB {
				mockDB := &gorm.DB{}
				mockDB.AddError(nil)
				return mockDB
			},
			article: &model.Article{Model: gorm.Model{ID: 1}},
			expectedComments: []model.Comment{
				{Model: gorm.Model{ID: 1}, ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, ArticleID: 1, Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			expectedError: nil,
		},
		{
			name: "Article with No Comments",
			setupMockDB: func() *gorm.DB {
				mockDB := &gorm.DB{}
				mockDB.AddError(nil)
				return mockDB
			},
			article:          &model.Article{Model: gorm.Model{ID: 2}},
			expectedComments: []model.Comment{},
			expectedError:    nil,
		},
		{
			name: "Non-existent Article",
			setupMockDB: func() *gorm.DB {
				mockDB := &gorm.DB{}
				mockDB.AddError(nil)
				return mockDB
			},
			article:          &model.Article{Model: gorm.Model{ID: 999}},
			expectedComments: []model.Comment{},
			expectedError:    nil,
		},
		{
			name: "Database Error",
			setupMockDB: func() *gorm.DB {
				mockDB := &gorm.DB{}
				mockDB.AddError(errors.New("database error"))
				return mockDB
			},
			article:          &model.Article{Model: gorm.Model{ID: 1}},
			expectedComments: []model.Comment{},
			expectedError:    errors.New("database error"),
		},
		{
			name: "Large Number of Comments",
			setupMockDB: func() *gorm.DB {
				mockDB := &gorm.DB{}
				mockDB.AddError(nil)
				return mockDB
			},
			article: &model.Article{Model: gorm.Model{ID: 3}},
			expectedComments: func() []model.Comment {
				comments := make([]model.Comment, 10000)
				for i := range comments {
					comments[i] = model.Comment{
						Model:     gorm.Model{ID: uint(i + 1)},
						ArticleID: 3,
						Author:    model.User{Model: gorm.Model{ID: uint(i % 100)}, Username: "user"},
					}
				}
				return comments
			}(),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := tt.setupMockDB()
			store := &ArticleStore{db: mockDB}

			comments, err := store.GetComments(tt.article)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, len(tt.expectedComments), len(comments))

			if len(tt.expectedComments) > 0 {
				for i, comment := range comments {
					assert.Equal(t, tt.expectedComments[i].ID, comment.ID)
					assert.Equal(t, tt.expectedComments[i].ArticleID, comment.ArticleID)
					assert.Equal(t, tt.expectedComments[i].Author.ID, comment.Author.ID)
					assert.Equal(t, tt.expectedComments[i].Author.Username, comment.Author.Username)
				}
			}

			if tt.name == "Large Number of Comments" {
				assert.Equal(t, 10000, len(comments))
			}
		})
	}
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func TestIsFavorited(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		user           *model.User
		setupMock      func(*MockDB)
		expectedResult bool
		expectedError  error
	}{
		{
			name:    "Article is favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(mockDB *MockDB) {
				mockDB.On("Table", "favorite_articles").Return(mockDB)
				mockDB.On("Where", "article_id = ? AND user_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*int)
					*arg = 1
				}).Return(mockDB)
			},
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: &MockGormDB{MockDB: mockDB},
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockDB.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func TestGetFeedArticles(t *testing.T) {
	tests := []struct {
		name     string
		userIDs  []uint
		limit    int64
		offset   int64
		mockDB   func() (*gorm.DB, sqlmock.Sqlmock)
		expected []model.Article
		wantErr  bool
	}{
		{
			name:    "Successful Retrieval of Feed Articles",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "author_id", "author_username"}).
					AddRow(1, 1, "Article 1", 1, "user1").
					AddRow(2, 2, "Article 2", 2, "user2")

				mock.ExpectQuery("SELECT .*").
					WithArgs(1, 2).
					WillReturnRows(rows)

				return gormDB, mock
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 1}, UserID: 1, Title: "Article 1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, UserID: 2, Title: "Article 2", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
		},
		{
			name:    "Empty Result Set",
			userIDs: []uint{99, 100},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "author_id", "author_username"})

				mock.ExpectQuery("SELECT .*").
					WithArgs(99, 100).
					WillReturnRows(rows)

				return gormDB, mock
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Pagination with Offset",
			userIDs: []uint{1, 2, 3},
			limit:   2,
			offset:  2,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "author_id", "author_username"}).
					AddRow(3, 2, "Article 3", 2, "user2").
					AddRow(4, 3, "Article 4", 3, "user3")

				mock.ExpectQuery("SELECT .*").
					WithArgs(1, 2, 3).
					WillReturnRows(rows)

				return gormDB, mock
			},
			expected: []model.Article{
				{Model: gorm.Model{ID: 3}, UserID: 2, Title: "Article 3", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
				{Model: gorm.Model{ID: 4}, UserID: 3, Title: "Article 4", Author: model.User{Model: gorm.Model{ID: 3}, Username: "user3"}},
			},
			wantErr: false,
		},
		{
			name:    "Error Handling - Database Error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockDB: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("mysql", db)

				mock.ExpectQuery("SELECT .*").
					WithArgs(1).
					WillReturnError(errors.New("database error"))

				return gormDB, mock
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := tt.mockDB()
			s := &ArticleStore{
				db: gormDB,
			}

			got, err := s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.expected)
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
func TestAddFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedUsers  []model.User
		concurrentCalls int
	}{
	
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

		
			store := &ArticleStore{db: mockDB}

			if tt.concurrentCalls > 0 {
				var wg sync.WaitGroup
				wg.Add(tt.concurrentCalls)
				for i := 0; i < tt.concurrentCalls; i++ {
					go func(i int) {
						defer wg.Done()
						user := &model.User{Model: gorm.Model{ID: uint(i + 1)}}
						err := store.AddFavorite(tt.article, user)
						assert.NoError(t, err)
					}(i)
				}
				wg.Wait()
			} else {
				err := store.AddFavorite(tt.article, tt.user)
				assert.Equal(t, tt.expectedError, err)
			}

			if tt.article != nil {
				assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
				assert.ElementsMatch(t, tt.expectedUsers, tt.article.FavoritedUsers)
			}

			mockDB.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestDeleteFavorite(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(*MockDB, *model.Article, *model.User)
		article         *model.Article
		user            *model.User
		expectedError   error
		expectedCount   int32
		concurrentCalls int
	}{
		{
			name: "Successfully Delete a Favorite",
			setupMock: func(mockDB *MockDB, a *model.Article, u *model.User) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", a).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", u).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Attempt to Delete a Non-existent Favorite",
			setupMock: func(mockDB *MockDB, a *model.Article, u *model.User) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", a).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", u).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
			},
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Database Error During Association Deletion",
			setupMock: func(mockDB *MockDB, a *model.Article, u *model.User) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", a).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", u).Return(&gorm.Association{Error: errors.New("DB error")})
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Rollback").Return(mockDB)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: errors.New("DB error"),
			expectedCount: 1,
		},
		{
			name: "Database Error During FavoritesCount Update",
			setupMock: func(mockDB *MockDB, a *model.Article, u *model.User) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", a).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", u).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(&gorm.DB{Error: errors.New("DB error")})
				mockDB.On("Rollback").Return(mockDB)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: errors.New("DB error"),
			expectedCount: 1,
		},
		{
			name: "Concurrent Deletion of Favorites",
			setupMock: func(mockDB *MockDB, a *model.Article, u *model.User) {
				mockDB.On("Begin").Return(mockDB)
				mockDB.On("Model", a).Return(mockDB)
				mockAssoc := &MockAssociation{}
				mockAssoc.On("Delete", u).Return(mockAssoc)
				mockDB.On("Association", "FavoritedUsers").Return(mockAssoc)
				mockDB.On("Update", "favorites_count", mock.Anything).Return(mockDB)
				mockDB.On("Commit").Return(mockDB)
			},
			article:         &model.Article{FavoritesCount: 5},
			user:            &model.User{},
			expectedError:   nil,
			expectedCount:   0,
			concurrentCalls: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB, tt.article, tt.user)

		
			store := &ArticleStore{db: mockDB}

			if tt.concurrentCalls > 0 {
				var wg sync.WaitGroup
				wg.Add(tt.concurrentCalls)
				for i := 0; i < tt.concurrentCalls; i++ {
					go func() {
						defer wg.Done()
						err := store.DeleteFavorite(tt.article, tt.user)
						assert.Equal(t, tt.expectedError, err)
					}()
				}
				wg.Wait()
			} else {
				err := store.DeleteFavorite(tt.article, tt.user)
				assert.Equal(t, tt.expectedError, err)
			}

			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
			mockDB.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_6382a4fe7a
ROOST_METHOD_SIG_HASH=GetArticles_1a0b3b0e8b


 */
func TestGetArticles(t *testing.T) {

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()


	store := &ArticleStore{db: db}


	setupTestData(db)

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
			name:      "Retrieve Articles Without Filters",
			limit:     10,
			offset:    0,
			wantCount: 3,
		},
		{
			name:      "Filter Articles by Tag Name",
			tagName:   "golang",
			limit:     10,
			offset:    0,
			wantCount: 2,
		},
		{
			name:      "Filter Articles by Author Username",
			username:  "johndoe",
			limit:     10,
			offset:    0,
			wantCount: 2,
		},
		{
			name:        "Retrieve Favorited Articles",
			favoritedBy: &model.User{Model: gorm.Model{ID: 2}},
			limit:       10,
			offset:      0,
			wantCount:   1,
		},
		{
			name:      "Test Pagination",
			limit:     1,
			offset:    1,
			wantCount: 1,
		},
		{
			name:      "Combine Multiple Filters",
			tagName:   "golang",
			username:  "johndoe",
			limit:     10,
			offset:    0,
			wantCount: 1,
		},
		{
			name:      "Handle Empty Result Set",
			tagName:   "nonexistent",
			limit:     10,
			offset:    0,
			wantCount: 0,
		},
		{
			name:    "Error Handling for Database Issues",
			wantErr: true,
		
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
			
				db.Close()
				store.db = db
			}

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, articles, tt.wantCount)

			if tt.tagName != "" {
				for _, article := range articles {
					hasTag := false
					for _, tag := range article.Tags {
						if tag.Name == tt.tagName {
							hasTag = true
							break
						}
					}
					assert.True(t, hasTag)
				}
			}

			if tt.username != "" {
				for _, article := range articles {
					assert.Equal(t, tt.username, article.Author.Username)
				}
			}

			if tt.favoritedBy != nil {
				for _, article := range articles {
					assert.Contains(t, article.FavoritedUsers, *tt.favoritedBy)
				}
			}
		})
	}
}

func setupTestData(db *gorm.DB) {

	users := []model.User{
		{Model: gorm.Model{ID: 1}, Username: "johndoe", Email: "john@example.com"},
		{Model: gorm.Model{ID: 2}, Username: "janedoe", Email: "jane@example.com"},
	}
	for _, user := range users {
		db.Create(&user)
	}


	tags := []model.Tag{
		{Model: gorm.Model{ID: 1}, Name: "golang"},
		{Model: gorm.Model{ID: 2}, Name: "testing"},
	}
	for _, tag := range tags {
		db.Create(&tag)
	}


	articles := []model.Article{
		{Model: gorm.Model{ID: 1}, Title: "Article 1", Description: "Description 1", Body: "Content 1", UserID: 1},
		{Model: gorm.Model{ID: 2}, Title: "Article 2", Description: "Description 2", Body: "Content 2", UserID: 1},
		{Model: gorm.Model{ID: 3}, Title: "Article 3", Description: "Description 3", Body: "Content 3", UserID: 2},
	}
	for _, article := range articles {
		db.Create(&article)
	}


	db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", 1, 1)
	db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", 2, 1)
	db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", 2, 2)
	db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", 3, 2)


	db.Exec("INSERT INTO favorite_articles (user_id, article_id) VALUES (?, ?)", 2, 1)
}

