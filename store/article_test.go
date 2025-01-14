package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"reflect"
	"time"
)






type mockDB struct {
	countResult int
	countError  error
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
		errMsg  string
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
			errMsg:  "record not found",
		},
		{
			name: "Delete Comment with Database Connection Error",
			comment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
			},
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:    "Delete Comment with Null Comment Pointer",
			comment: nil,
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
			errMsg:  "invalid argument",
		},
		{
			name: "Delete Comment with Cascading Effects",
			comment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				ArticleID: 1,
			},
			dbSetup: func(db *gorm.DB) {
				article := &model.Article{Model: gorm.Model{ID: 1}, Title: "Test Article"}
				db.Create(article)
				db.Create(&model.Comment{Model: gorm.Model{ID: 1}, Body: "Test comment", ArticleID: 1})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open mock database: %v", err)
			}
			defer mockDB.Close()

			mockDB.AutoMigrate(&model.Comment{}, &model.Article{})
			tt.dbSetup(mockDB)

			store := &ArticleStore{db: mockDB}

			err = store.DeleteComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("DeleteComment() error message = %v, want %v", err.Error(), tt.errMsg)
				return
			}

			if !tt.wantErr && tt.comment != nil {
				var count int64
				mockDB.Model(&model.Comment{}).Where("id = ?", tt.comment.ID).Count(&count)
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
func (m *mockDB) Count(value interface{}) *gorm.DB {
	*(value.(*int)) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Value: m}
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name            string
		article         *model.Article
		user            *model.User
		mockCountResult int
		mockCountError  error
		expectedResult  bool
		expectedError   error
	}{
		{
			name:            "Article is favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 1,
			mockCountError:  nil,
			expectedResult:  true,
			expectedError:   nil,
		},
		{
			name:            "Article is not favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil article parameter",
			article:         nil,
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil user parameter",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Database error",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  errors.New("database error"),
			expectedResult:  false,
			expectedError:   errors.New("database error"),
		},
		{
			name:            "Article with no favorites",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "User with no favorite articles",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				countResult: tt.mockCountResult,
				countError:  tt.mockCountError,
			}

			store := &ArticleStore{
				db: &gorm.DB{Value: mockDB},
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expectedResult {
				t.Errorf("Expected result %v, but got %v", tt.expectedResult, result)
			}

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
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
func TestArticleStoreGetFeedArticles(t *testing.T) {
	tests := []struct {
		name    string
		userIDs []uint
		limit   int64
		offset  int64
		mockDB  func() *gorm.DB
		want    []model.Article
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Successful retrieval of feed articles",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
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
		},
		{
			name:    "Empty result set",
			userIDs: []uint{99, 100},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Database error handling",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(errors.New("database error"))
				return db
			},
			want:    nil,
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:    "Limit and offset behavior",
			userIDs: []uint{1, 2, 3},
			limit:   2,
			offset:  1,
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
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
				return db
			},
			want: func() []model.Article {
				articles := make([]model.Article, 10)
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

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ArticleStore.GetFeedArticles() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.want)
			}

			for _, article := range got {
				if article.Author.ID == 0 {
					t.Errorf("ArticleStore.GetFeedArticles() Author not preloaded for article ID %v", article.ID)
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
	type args struct {
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
	}

	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
		dbSetup func(*gorm.DB)
	}{
		{
			name: "Retrieve Articles Without Any Filters",
			args: args{
				tagName:     "",
				username:    "",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Article 1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 2}, Title: "Article 2", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 2})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
				db.Create(&model.User{Model: gorm.Model{ID: 2}, Username: "user2"})
			},
		},
		{
			name: "Filter Articles by Tag Name",
			args: args{
				tagName:     "golang",
				username:    "",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Golang Article", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
			},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Golang Article", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 2}, Title: "Python Article", UserID: 2})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
				db.Create(&model.User{Model: gorm.Model{ID: 2}, Username: "user2"})
				db.Exec("INSERT INTO tags (name) VALUES ('golang')")
				db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (1, 1)")
			},
		},
		{
			name: "Filter Articles by Author Username",
			args: args{
				tagName:     "",
				username:    "user1",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "User1 Article", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
			},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "User1 Article", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 2}, Title: "User2 Article", UserID: 2})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
				db.Create(&model.User{Model: gorm.Model{ID: 2}, Username: "user2"})
			},
		},
		{
			name: "Retrieve Favorited Articles",
			args: args{
				tagName:     "",
				username:    "",
				favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
				limit:       10,
				offset:      0,
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Favorited Article", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
			},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Favorited Article", UserID: 2})
				db.Create(&model.Article{Model: gorm.Model{ID: 2}, Title: "Unfavorited Article", UserID: 2})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
				db.Create(&model.User{Model: gorm.Model{ID: 2}, Username: "user2"})
				db.Exec("INSERT INTO favorite_articles (user_id, article_id) VALUES (1, 1)")
			},
		},
		{
			name: "Test Pagination with Limit and Offset",
			args: args{
				tagName:     "",
				username:    "",
				favoritedBy: nil,
				limit:       2,
				offset:      1,
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 2}, Title: "Article 2", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
				{Model: gorm.Model{ID: 3}, Title: "Article 3", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
			},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 2}, Title: "Article 2", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 3}, Title: "Article 3", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 4}, Title: "Article 4", UserID: 1})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
			},
		},
		{
			name: "Combine Multiple Filters",
			args: args{
				tagName:     "golang",
				username:    "user1",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want: []model.Article{
				{Model: gorm.Model{ID: 1}, Title: "Golang Article by User1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
			},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Golang Article by User1", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 2}, Title: "Python Article by User1", UserID: 1})
				db.Create(&model.Article{Model: gorm.Model{ID: 3}, Title: "Golang Article by User2", UserID: 2})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
				db.Create(&model.User{Model: gorm.Model{ID: 2}, Username: "user2"})
				db.Exec("INSERT INTO tags (name) VALUES ('golang')")
				db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (1, 1), (3, 1)")
			},
		},
		{
			name: "Handle Non-Existent Tag or User",
			args: args{
				tagName:     "nonexistent",
				username:    "",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want:    []model.Article{},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Article{Model: gorm.Model{ID: 1}, Title: "Article 1", UserID: 1})
				db.Create(&model.User{Model: gorm.Model{ID: 1}, Username: "user1"})
			},
		},
		{
			name: "Test with Empty Database",
			args: args{
				tagName:     "",
				username:    "",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want:    []model.Article{},
			wantErr: false,
			dbSetup: func(db *gorm.DB) {},
		},
		{
			name: "Handle Database Connection Errors",
			args: args{
				tagName:     "",
				username:    "",
				favoritedBy: nil,
				limit:       10,
				offset:      0,
			},
			want:    []model.Article{},
			wantErr: true,
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Article{}, &model.User{}, &model.Tag{})

			tt.dbSetup(db)

			s := &ArticleStore{
				db: db,
			}

			got, err := s.GetArticles(tt.args.tagName, tt.args.username, tt.args.favoritedBy, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i := range got {
				got[i].CreatedAt = time.Time{}
				got[i].UpdatedAt = time.Time{}
				got[i].DeletedAt = nil
				got[i].Author.CreatedAt = time.Time{}
				got[i].Author.UpdatedAt = time.Time{}
				got[i].Author.DeletedAt = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}

