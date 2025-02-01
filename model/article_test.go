package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)








/*
ROOST_METHOD_HASH=Overwrite_3d4db6693d
ROOST_METHOD_SIG_HASH=Overwrite_22e8730976

FUNCTION_DEF=func (a *Article) Overwrite(title, description, body string) 

*/
func TestArticleOverwrite(t *testing.T) {
	tests := []struct {
		name        string
		initial     Article
		title       string
		description string
		body        string
		expected    Article
	}{
		{
			name: "Overwrite All Fields",
			initial: Article{
				Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
			},
		},
		{
			name: "Overwrite Partial Fields",
			initial: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title",
			description: "",
			body:        "New Body",
			expected: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title",
				Description: "Initial Description",
				Body:        "New Body",
			},
		},
		{
			name: "Overwrite with Empty Strings",
			initial: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "",
			description: "",
			body:        "",
			expected: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
		},
		{
			name: "Overwrite with Very Long Strings",
			initial: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       string(make([]byte, 10000)),
			description: string(make([]byte, 10000)),
			body:        string(make([]byte, 10000)),
			expected: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       string(make([]byte, 10000)),
				Description: string(make([]byte, 10000)),
				Body:        string(make([]byte, 10000)),
			},
		},
		{
			name: "Overwrite Maintaining Other Fields",
			initial: Article{
				Model:          gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "Initial Title",
				Description:    "Initial Description",
				Body:           "Initial Body",
				Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "Tag1"}},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:         1,
				FavoritesCount: 10,
				FavoritedUsers: []User{{Model: gorm.Model{ID: 2}, Username: "fan"}},
				Comments:       []Comment{{Model: gorm.Model{ID: 1}, Body: "Great article!"}},
			},
			title:       "New Title",
			description: "",
			body:        "",
			expected: Article{
				Model:          gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "New Title",
				Description:    "Initial Description",
				Body:           "Initial Body",
				Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "Tag1"}},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:         1,
				FavoritesCount: 10,
				FavoritedUsers: []User{{Model: gorm.Model{ID: 2}, Username: "fan"}},
				Comments:       []Comment{{Model: gorm.Model{ID: 1}, Body: "Great article!"}},
			},
		},
		{
			name: "Overwrite with Unicode Characters",
			initial: Article{
				Model:       gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "新标题",
			description: "新描述",
			body:        "新正文",
			expected: Article{
				Model:       gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "新标题",
				Description: "新描述",
				Body:        "新正文",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := tt.initial
			article.Overwrite(tt.title, tt.description, tt.body)

			if article.Title != tt.expected.Title {
				t.Errorf("Title = %v, want %v", article.Title, tt.expected.Title)
			}
			if article.Description != tt.expected.Description {
				t.Errorf("Description = %v, want %v", article.Description, tt.expected.Description)
			}
			if article.Body != tt.expected.Body {
				t.Errorf("Body = %v, want %v", article.Body, tt.expected.Body)
			}

			if tt.name == "Overwrite Maintaining Other Fields" {
				if len(article.Tags) != len(tt.expected.Tags) {
					t.Errorf("Tags length = %v, want %v", len(article.Tags), len(tt.expected.Tags))
				}
				if article.Author.Username != tt.expected.Author.Username {
					t.Errorf("Author Username = %v, want %v", article.Author.Username, tt.expected.Author.Username)
				}
				if article.UserID != tt.expected.UserID {
					t.Errorf("UserID = %v, want %v", article.UserID, tt.expected.UserID)
				}
				if article.FavoritesCount != tt.expected.FavoritesCount {
					t.Errorf("FavoritesCount = %v, want %v", article.FavoritesCount, tt.expected.FavoritesCount)
				}
				if len(article.FavoritedUsers) != len(tt.expected.FavoritedUsers) {
					t.Errorf("FavoritedUsers length = %v, want %v", len(article.FavoritedUsers), len(tt.expected.FavoritedUsers))
				}
				if len(article.Comments) != len(tt.expected.Comments) {
					t.Errorf("Comments length = %v, want %v", len(article.Comments), len(tt.expected.Comments))
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726

FUNCTION_DEF=func (a *Article) ProtoArticle(favorited bool) *pb.Article 

*/
func TestArticleProtoArticle(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		article   Article
		favorited bool
		want      *pb.Article
	}{
		{
			name: "Basic Article Conversion",
			article: Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Test Title",
				Description:    "Test Description",
				Body:           "Test Body",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:         1,
				FavoritesCount: 5,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "1",
				Title:          "Test Title",
				Description:    "Test Description",
				Body:           "Test Body",
				TagList:        []string{"tag1", "tag2"},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      true,
				FavoritesCount: 5,
			},
		},
		{
			name: "Article Conversion with Empty Tags",
			article: Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "No Tags",
				Description:    "Article without tags",
				Body:           "Body text",
				Tags:           []Tag{},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:         1,
				FavoritesCount: 0,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "2",
				Title:          "No Tags",
				Description:    "Article without tags",
				Body:           "Body text",
				TagList:        []string{},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 0,
			},
		},
		{
			name: "Article with Zero Values",
			article: Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "",
				Description:    "",
				Body:           "",
				Tags:           []Tag{},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:         1,
				FavoritesCount: 0,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "3",
				Title:          "",
				Description:    "",
				Body:           "",
				TagList:        []string{},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.article.ProtoArticle(tt.favorited)
			assert.Equal(t, tt.want.Slug, got.Slug)
			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.Body, got.Body)
			assert.Equal(t, tt.want.TagList, got.TagList)
			assert.Equal(t, tt.want.CreatedAt, got.CreatedAt)
			assert.Equal(t, tt.want.UpdatedAt, got.UpdatedAt)
			assert.Equal(t, tt.want.Favorited, got.Favorited)
			assert.Equal(t, tt.want.FavoritesCount, got.FavoritesCount)
		})
	}
}


/*
ROOST_METHOD_HASH=Validate_f6d09c3ac5
ROOST_METHOD_SIG_HASH=Validate_99e41aac91

FUNCTION_DEF=func (a Article) Validate() error 

*/
func TestArticleValidate(t *testing.T) {
	tests := []struct {
		name    string
		article Article
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Article with All Required Fields",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Missing Title",
			article: Article{
				Body: "Test Body",
				Tags: []Tag{{Name: "Test Tag"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank.",
		},
		{
			name: "Article with Missing Body",
			article: Article{
				Title: "Test Title",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Article with No Tags",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body",
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank.",
		},
		{
			name:    "Article with All Fields Missing",
			article: Article{},
			wantErr: true,
			errMsg:  "body: cannot be blank; tags: cannot be blank; title: cannot be blank.",
		},
		{
			name: "Article with Valid Data and Extra Fields",
			article: Article{
				Model:          gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "Test Title",
				Description:    "Test Description",
				Body:           "Test Body",
				Tags:           []Tag{{Name: "Test Tag"}},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "testuser"},
				UserID:         1,
				FavoritesCount: 10,
				FavoritedUsers: []User{{Model: gorm.Model{ID: 2}, Username: "fan"}},
				Comments:       []Comment{{Body: "Great article!"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Very Long Title",
			article: Article{
				Title: string(make([]byte, 1000)),
				Body:  "Test Body",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Special Characters in Fields",
			article: Article{
				Title: "Test Title !@#$%^&*()",
				Body:  "Test Body 你好世界",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

