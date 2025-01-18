package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"time"
	"github.com/raahii/golang-grpc-realworld-example/model"
)









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

