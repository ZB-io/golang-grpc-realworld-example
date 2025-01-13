package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"math"
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
			title:       string(make([]byte, 10000)),
			description: string(make([]byte, 10000)),
			body:        string(make([]byte, 10000)),
			expected: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       string(make([]byte, 10000)),
				Description: string(make([]byte, 10000)),
				Body:        string(make([]byte, 10000)),
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
			title:       "New Title 😊",
			description: "New Description 你好",
			body:        "New Body こんにちは",
			expected: Article{
				Model:       gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title 😊",
				Description: "New Description 你好",
				Body:        "New Body こんにちは",
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
		})
	}
}


/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726

FUNCTION_DEF=func (a *Article) ProtoArticle(favorited bool) *pb.Article 

 */
func TestArticleProtoArticle(t *testing.T) {
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
					CreatedAt: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC),
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
				CreatedAt:      "2023-05-01T10:00:00+0000Z",
				UpdatedAt:      "2023-05-02T11:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 10,
			},
		},
		{
			name: "Empty Tags",
			article: Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Date(2023, 5, 3, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 4, 13, 0, 0, 0, time.UTC),
				},
				Title:          "No Tags Article",
				Description:    "Article without tags",
				Body:           "Body content",
				Tags:           []Tag{},
				FavoritesCount: 5,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "2",
				Title:          "No Tags Article",
				Description:    "Article without tags",
				Body:           "Body content",
				TagList:        []string{},
				CreatedAt:      "2023-05-03T12:00:00+0000Z",
				UpdatedAt:      "2023-05-04T13:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 5,
			},
		},
		{
			name: "Multiple Tags",
			article: Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Date(2023, 5, 5, 14, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 6, 15, 0, 0, 0, time.UTC),
				},
				Title:          "Multi-Tag Article",
				Description:    "Article with multiple tags",
				Body:           "Content with tags",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
				FavoritesCount: 15,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "3",
				Title:          "Multi-Tag Article",
				Description:    "Article with multiple tags",
				Body:           "Content with tags",
				TagList:        []string{"tag1", "tag2", "tag3"},
				CreatedAt:      "2023-05-05T14:00:00+0000Z",
				UpdatedAt:      "2023-05-06T15:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 15,
			},
		},
		{
			name: "Max FavoritesCount",
			article: Article{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: time.Date(2023, 5, 7, 16, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 8, 17, 0, 0, 0, time.UTC),
				},
				Title:          "Popular Article",
				Description:    "Article with max favorites",
				Body:           "Very popular content",
				Tags:           []Tag{{Name: "popular"}},
				FavoritesCount: math.MaxInt32,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "4",
				Title:          "Popular Article",
				Description:    "Article with max favorites",
				Body:           "Very popular content",
				TagList:        []string{"popular"},
				CreatedAt:      "2023-05-07T16:00:00+0000Z",
				UpdatedAt:      "2023-05-08T17:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: math.MaxInt32,
			},
		},
		{
			name: "Empty Strings",
			article: Article{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: time.Date(2023, 5, 9, 18, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 10, 19, 0, 0, 0, time.UTC),
				},
				Title:          "",
				Description:    "",
				Body:           "",
				Tags:           []Tag{{Name: "empty"}},
				FavoritesCount: 0,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "5",
				Title:          "",
				Description:    "",
				Body:           "",
				TagList:        []string{"empty"},
				CreatedAt:      "2023-05-09T18:00:00+0000Z",
				UpdatedAt:      "2023-05-10T19:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.article.ProtoArticle(tt.favorited)
			assert.Equal(t, tt.want, got)
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
				Title: "Test Article",
				Body:  "This is a test article body",
				Tags:  []Tag{{Name: "test"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Missing Title",
			article: Article{
				Body: "This is a test article body",
				Tags: []Tag{{Name: "test"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank.",
		},
		{
			name: "Article with Missing Body",
			article: Article{
				Title: "Test Article",
				Tags:  []Tag{{Name: "test"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Article with No Tags",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
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
				Title: string(make([]rune, 1000)),
				Body:  "This is a test article body",
				Tags:  []Tag{{Name: "test"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Special Characters in Title and Body",
			article: Article{
				Title: "Test Article with Special Characters: !@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
				Body:  "This is a test article body with emojis: 😀🎉🚀",
				Tags:  []Tag{{Name: "test"}},
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

