package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"fmt"
	"reflect"
	"github.com/raahii/golang-grpc-realworld-example/proto"
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
			name: "Overwrite all fields with non-empty values",
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
			name: "Overwrite only the title field",
			initial: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title",
			description: "",
			body:        "",
			expected: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
		},
		{
			name: "Overwrite with all empty strings",
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
			name: "Overwrite description and body, leaving title unchanged",
			initial: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "New Description",
				Body:        "New Body",
			},
		},
		{
			name: "Overwrite with very long strings",
			initial: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       string(make([]byte, 1000)),
			description: string(make([]byte, 1000)),
			body:        string(make([]byte, 1000)),
			expected: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       string(make([]byte, 1000)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 1000)),
			},
		},
		{
			name: "Overwrite with special characters",
			initial: Article{
				Model:       gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title 😊 <script>alert('XSS')</script>",
			description: "New Description 🌈 &lt;html&gt;",
			body:        "New Body ☀️ <p>HTML content</p>",
			expected: Article{
				Model:       gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title 😊 <script>alert('XSS')</script>",
				Description: "New Description 🌈 &lt;html&gt;",
				Body:        "New Body ☀️ <p>HTML content</p>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := tt.initial
			article.Overwrite(tt.title, tt.description, tt.body)

			assert.Equal(t, tt.expected.Title, article.Title)
			assert.Equal(t, tt.expected.Description, article.Description)
			assert.Equal(t, tt.expected.Body, article.Body)
			assert.Equal(t, tt.expected.ID, article.ID)
			assert.Equal(t, tt.expected.CreatedAt.Unix(), article.CreatedAt.Unix())
			assert.Equal(t, tt.expected.UpdatedAt.Unix(), article.UpdatedAt.Unix())
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
	iso8601Format := "2006-01-02T15:04:05-0700Z"

	tests := []struct {
		name      string
		article   Article
		favorited bool
		want      *pb.Article
	}{
		{
			name: "Basic Conversion",
			article: Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}},
				FavoritesCount: 10,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				TagList:        []string{"tag1", "tag2"},
				CreatedAt:      now.Format(iso8601Format),
				UpdatedAt:      now.Format(iso8601Format),
				Favorited:      true,
				FavoritesCount: 10,
			},
		},
		{
			name: "Empty Tags",
			article: Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "No Tags Article",
				Description:    "Article without tags",
				Body:           "Body text",
				Tags:           []Tag{},
				FavoritesCount: 5,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "2",
				Title:          "No Tags Article",
				Description:    "Article without tags",
				Body:           "Body text",
				TagList:        []string{},
				CreatedAt:      now.Format(iso8601Format),
				UpdatedAt:      now.Format(iso8601Format),
				Favorited:      false,
				FavoritesCount: 5,
			},
		},
		{
			name: "Multiple Tags",
			article: Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Multiple Tags",
				Description:    "Article with multiple tags",
				Body:           "Content",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
				FavoritesCount: 15,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "3",
				Title:          "Multiple Tags",
				Description:    "Article with multiple tags",
				Body:           "Content",
				TagList:        []string{"tag1", "tag2", "tag3"},
				CreatedAt:      now.Format(iso8601Format),
				UpdatedAt:      now.Format(iso8601Format),
				Favorited:      true,
				FavoritesCount: 15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.article.ProtoArticle(tt.favorited)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoArticle() = %v, want %v", got, tt.want)
			}

			if got.Slug != fmt.Sprintf("%d", tt.article.ID) {
				t.Errorf("ProtoArticle() Slug = %v, want %v", got.Slug, tt.article.ID)
			}

			if got.CreatedAt != tt.article.CreatedAt.Format(iso8601Format) {
				t.Errorf("ProtoArticle() CreatedAt = %v, want %v", got.CreatedAt, tt.article.CreatedAt.Format(iso8601Format))
			}

			if got.UpdatedAt != tt.article.UpdatedAt.Format(iso8601Format) {
				t.Errorf("ProtoArticle() UpdatedAt = %v, want %v", got.UpdatedAt, tt.article.UpdatedAt.Format(iso8601Format))
			}

			if got.Favorited != tt.favorited {
				t.Errorf("ProtoArticle() Favorited = %v, want %v", got.Favorited, tt.favorited)
			}

			if len(got.TagList) != len(tt.article.Tags) {
				t.Errorf("ProtoArticle() TagList length = %v, want %v", len(got.TagList), len(tt.article.Tags))
			}
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
			name: "Article with Very Long Title",
			article: Article{
				Title: string(make([]byte, 1000)),
				Body:  "Test Body",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Special Characters in Title and Body",
			article: Article{
				Title: "Special 🚀 Tïtlé",
				Body:  "Special 🌟 Bödy",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Article.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("Article.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

