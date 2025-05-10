package github

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
				UserID:      1,
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
				UserID:      1,
			},
		},
		{
			name: "Overwrite Partial Fields",
			initial: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				UserID:      2,
			},
			title:       "Updated Title",
			description: "",
			body:        "",
			expected: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Updated Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				UserID:      2,
			},
		},
		{
			name: "Overwrite with Empty Strings",
			initial: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				UserID:      3,
			},
			title:       "",
			description: "",
			body:        "",
			expected: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				UserID:      3,
			},
		},
		{
			name: "Overwrite with Very Long Strings",
			initial: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				UserID:      4,
			},
			title:       string(make([]byte, 10000)),
			description: string(make([]byte, 10000)),
			body:        string(make([]byte, 10000)),
			expected: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       string(make([]byte, 10000)),
				Description: string(make([]byte, 10000)),
				Body:        string(make([]byte, 10000)),
				UserID:      4,
			},
		},
		{
			name: "Overwrite with Special Characters",
			initial: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				UserID:      5,
			},
			title:       "Title with 😊 emoji",
			description: "<script>alert('XSS')</script>",
			body:        "Body with Unicode: こんにちは",
			expected: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Title with 😊 emoji",
				Description: "<script>alert('XSS')</script>",
				Body:        "Body with Unicode: こんにちは",
				UserID:      5,
			},
		},
		{
			name: "Overwrite Maintaining Other Struct Fields",
			initial: Article{
				Model:          gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "Initial Title",
				Description:    "Initial Description",
				Body:           "Initial Body",
				Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "TestTag"}},
				UserID:         6,
				FavoritesCount: 10,
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Model:          gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "New Title",
				Description:    "New Description",
				Body:           "New Body",
				Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "TestTag"}},
				UserID:         6,
				FavoritesCount: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := tt.initial
			article.Overwrite(tt.title, tt.description, tt.body)

			if article.Title != tt.expected.Title {
				t.Errorf("Title mismatch. Got %s, want %s", article.Title, tt.expected.Title)
			}
			if article.Description != tt.expected.Description {
				t.Errorf("Description mismatch. Got %s, want %s", article.Description, tt.expected.Description)
			}
			if article.Body != tt.expected.Body {
				t.Errorf("Body mismatch. Got %s, want %s", article.Body, tt.expected.Body)
			}

			if article.ID != tt.expected.ID {
				t.Errorf("ID changed unexpectedly. Got %d, want %d", article.ID, tt.expected.ID)
			}
			if article.UserID != tt.expected.UserID {
				t.Errorf("UserID changed unexpectedly. Got %d, want %d", article.UserID, tt.expected.UserID)
			}
			if article.FavoritesCount != tt.expected.FavoritesCount {
				t.Errorf("FavoritesCount changed unexpectedly. Got %d, want %d", article.FavoritesCount, tt.expected.FavoritesCount)
			}
			if len(article.Tags) != len(tt.expected.Tags) {
				t.Errorf("Tags changed unexpectedly. Got %v, want %v", article.Tags, tt.expected.Tags)
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
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}},
				FavoritesCount: 5,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
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
				Title:          "No Tags Article",
				Description:    "No Tags Description",
				Body:           "No Tags Body",
				Tags:           []Tag{},
				FavoritesCount: 3,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "2",
				Title:          "No Tags Article",
				Description:    "No Tags Description",
				Body:           "No Tags Body",
				TagList:        []string{},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 3,
			},
		},
		{
			name: "Zero Values Handling",
			article: Article{
				Model: gorm.Model{
					ID:        0,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Zero Values Article",
				Description:    "Zero Values Description",
				Body:           "Zero Values Body",
				Tags:           []Tag{},
				FavoritesCount: 0,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "0",
				Title:          "Zero Values Article",
				Description:    "Zero Values Description",
				Body:           "Zero Values Body",
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
			name: "Valid Article",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
		{
			name: "Empty Title",
			article: Article{
				Body: "Test Body",
				Tags: []Tag{{Name: "Test Tag"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank.",
		},
		{
			name: "Empty Body",
			article: Article{
				Title: "Test Title",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "No Tags",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body",
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank.",
		},
		{
			name:    "All Fields Empty",
			article: Article{},
			wantErr: true,
			errMsg:  "body: cannot be blank; tags: cannot be blank; title: cannot be blank.",
		},
		{
			name: "Very Long Title",
			article: Article{
				Title: string(make([]byte, 1000)),
				Body:  "Test Body",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
		{
			name: "Special Characters in Fields",
			article: Article{
				Title: "Test Title with 🚀 and <html>",
				Body:  "Test Body with 😊 and <script>alert('XSS')</script>",
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

