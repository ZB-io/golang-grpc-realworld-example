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
		article     Article
		title       string
		description string
		body        string
		expected    Article
	}{
		{
			name: "Overwrite All Fields",
			article: Article{
				Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
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
			article: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      2,
			},
			title:       "New Title",
			description: "",
			body:        "",
			expected: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      2,
			},
		},
		{
			name: "Overwrite with Empty Strings",
			article: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      3,
			},
			title:       "",
			description: "",
			body:        "",
			expected: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      3,
			},
		},
		{
			name: "Overwrite with Very Long Strings",
			article: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
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
			article: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      5,
			},
			title:       "New Title 😊 <script>alert('XSS')</script>",
			description: "New Description 🚀 &lt;html&gt;",
			body:        "New Body 💡 こんにちは",
			expected: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title 😊 <script>alert('XSS')</script>",
				Description: "New Description 🚀 &lt;html&gt;",
				Body:        "New Body 💡 こんにちは",
				UserID:      5,
			},
		},
		{
			name: "Overwrite Maintaining Other Struct Fields",
			article: Article{
				Model:          gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "Original Title",
				Description:    "Original Description",
				Body:           "Original Body",
				Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "Tag1"}},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "author"},
				UserID:         6,
				FavoritesCount: 10,
				FavoritedUsers: []User{{Model: gorm.Model{ID: 2}, Username: "user1"}},
				Comments:       []Comment{{Model: gorm.Model{ID: 1}, Body: "Comment1"}},
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Model:          gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "New Title",
				Description:    "New Description",
				Body:           "New Body",
				Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "Tag1"}},
				Author:         User{Model: gorm.Model{ID: 1}, Username: "author"},
				UserID:         6,
				FavoritesCount: 10,
				FavoritedUsers: []User{{Model: gorm.Model{ID: 2}, Username: "user1"}},
				Comments:       []Comment{{Model: gorm.Model{ID: 1}, Body: "Comment1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.article.Overwrite(tt.title, tt.description, tt.body)

			if tt.article.Title != tt.expected.Title {
				t.Errorf("Title mismatch. Got: %s, Want: %s", tt.article.Title, tt.expected.Title)
			}
			if tt.article.Description != tt.expected.Description {
				t.Errorf("Description mismatch. Got: %s, Want: %s", tt.article.Description, tt.expected.Description)
			}
			if tt.article.Body != tt.expected.Body {
				t.Errorf("Body mismatch. Got: %s, Want: %s", tt.article.Body, tt.expected.Body)
			}

			if tt.article.UserID != tt.expected.UserID {
				t.Errorf("UserID changed unexpectedly. Got: %d, Want: %d", tt.article.UserID, tt.expected.UserID)
			}
			if tt.article.FavoritesCount != tt.expected.FavoritesCount {
				t.Errorf("FavoritesCount changed unexpectedly. Got: %d, Want: %d", tt.article.FavoritesCount, tt.expected.FavoritesCount)
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
	testCases := []struct {
		name      string
		article   *Article
		favorited bool
		expected  *pb.Article
	}{
		{
			name: "Basic Article Conversion",
			article: &Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}},
				Author:         User{Username: "testuser"},
				UserID:         1,
				FavoritesCount: 10,
			},
			favorited: true,
			expected: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				TagList:        []string{"tag1", "tag2"},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      true,
				FavoritesCount: 10,
			},
		},
		{
			name: "Article Conversion with Empty Tags",
			article: &Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "No Tags Article",
				Description:    "No Tags Description",
				Body:           "No Tags Body",
				Tags:           []Tag{},
				Author:         User{Username: "testuser"},
				UserID:         1,
				FavoritesCount: 5,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "2",
				Title:          "No Tags Article",
				Description:    "No Tags Description",
				Body:           "No Tags Body",
				TagList:        []string{},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 5,
			},
		},
		{
			name: "Zero FavoritesCount",
			article: &Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Zero Favorites",
				Description:    "Zero Favorites Description",
				Body:           "Zero Favorites Body",
				Tags:           []Tag{{Name: "tag"}},
				Author:         User{Username: "testuser"},
				UserID:         1,
				FavoritesCount: 0,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "3",
				Title:          "Zero Favorites",
				Description:    "Zero Favorites Description",
				Body:           "Zero Favorites Body",
				TagList:        []string{"tag"},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 0,
			},
		},
		{
			name: "Large FavoritesCount",
			article: &Article{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:          "Large Favorites",
				Description:    "Large Favorites Description",
				Body:           "Large Favorites Body",
				Tags:           []Tag{{Name: "popular"}},
				Author:         User{Username: "testuser"},
				UserID:         1,
				FavoritesCount: 2147483647,
			},
			favorited: true,
			expected: &pb.Article{
				Slug:           "4",
				Title:          "Large Favorites",
				Description:    "Large Favorites Description",
				Body:           "Large Favorites Body",
				TagList:        []string{"popular"},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      true,
				FavoritesCount: 2147483647,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.article.ProtoArticle(tc.favorited)
			assert.Equal(t, tc.expected.Slug, result.Slug)
			assert.Equal(t, tc.expected.Title, result.Title)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.Body, result.Body)
			assert.Equal(t, tc.expected.TagList, result.TagList)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
			assert.Equal(t, tc.expected.Favorited, result.Favorited)
			assert.Equal(t, tc.expected.FavoritesCount, result.FavoritesCount)
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
			name: "Article Missing Title",
			article: Article{
				Body: "Test Body",
				Tags: []Tag{{Name: "Test Tag"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank.",
		},
		{
			name: "Article Missing Body",
			article: Article{
				Title: "Test Title",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Article with Empty Tags Slice",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body",
				Tags:  []Tag{},
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank.",
		},
		{
			name:    "Article with All Fields Empty",
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
			name: "Article with Special Characters in Fields",
			article: Article{
				Title: "Test Title 🚀 with ñ and 你好",
				Body:  "Test Body with ñ and 你好 🌍",
				Tags:  []Tag{{Name: "Test Tag"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Minimum Valid Data",
			article: Article{
				Title: "T",
				Body:  "B",
				Tags:  []Tag{{Name: "T"}},
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

