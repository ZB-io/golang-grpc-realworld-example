package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/proto"
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
		want        Article
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
			want: Article{
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
			body:        "New Body",
			want: Article{
				Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "New Title",
				Description: "Original Description",
				Body:        "New Body",
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
			want: Article{
				Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      3,
			},
		},
		{
			name: "Overwrite with Special Characters",
			article: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      4,
			},
			title:       "Title with ñ and é",
			description: "Description with 你好",
			body:        "Body with 🚀 emoji",
			want: Article{
				Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Title with ñ and é",
				Description: "Description with 你好",
				Body:        "Body with 🚀 emoji",
				UserID:      4,
			},
		},
		{
			name: "Overwrite with Very Long Strings",
			article: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
				UserID:      5,
			},
			title:       string(make([]byte, 10000)),
			description: string(make([]byte, 10000)),
			body:        string(make([]byte, 10000)),
			want: Article{
				Model:       gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       string(make([]byte, 10000)),
				Description: string(make([]byte, 10000)),
				Body:        string(make([]byte, 10000)),
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
				UserID:         6,
				FavoritesCount: 10,
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			want: Article{
				Model:          gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:          "New Title",
				Description:    "New Description",
				Body:           "New Body",
				UserID:         6,
				FavoritesCount: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.article.Overwrite(tt.title, tt.description, tt.body)

			if tt.article.Title != tt.want.Title {
				t.Errorf("Title = %v, want %v", tt.article.Title, tt.want.Title)
			}
			if tt.article.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", tt.article.Description, tt.want.Description)
			}
			if tt.article.Body != tt.want.Body {
				t.Errorf("Body = %v, want %v", tt.article.Body, tt.want.Body)
			}
			if tt.article.UserID != tt.want.UserID {
				t.Errorf("UserID = %v, want %v", tt.article.UserID, tt.want.UserID)
			}
			if tt.article.FavoritesCount != tt.want.FavoritesCount {
				t.Errorf("FavoritesCount = %v, want %v", tt.article.FavoritesCount, tt.want.FavoritesCount)
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
		article   Article
		favorited bool
		expected  *pb.Article
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
			favorited: false,
			expected: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				TagList:        []string{"tag1", "tag2"},
				CreatedAt:      now.Format(ISO8601),
				UpdatedAt:      now.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 10,
			},
		},
		{
			name: "Favorited Flag Set to True",
			article: Article{
				Model: gorm.Model{ID: 2},
				Title: "Favorited Article",
			},
			favorited: true,
			expected: &pb.Article{
				Slug:      "2",
				Title:     "Favorited Article",
				Favorited: true,
			},
		},
		{
			name: "Article with No Tags",
			article: Article{
				Model: gorm.Model{ID: 3},
				Title: "No Tags Article",
			},
			favorited: false,
			expected: &pb.Article{
				Slug:    "3",
				Title:   "No Tags Article",
				TagList: []string{},
			},
		},
		{
			name: "Article with Multiple Tags",
			article: Article{
				Model: gorm.Model{ID: 4},
				Title: "Multi-Tag Article",
				Tags:  []Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			},
			favorited: false,
			expected: &pb.Article{
				Slug:    "4",
				Title:   "Multi-Tag Article",
				TagList: []string{"tag1", "tag2", "tag3"},
			},
		},
		{
			name: "Correct Formatting of Timestamps",
			article: Article{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 16, 11, 45, 0, 0, time.UTC),
				},
				Title: "Timestamp Article",
			},
			favorited: false,
			expected: &pb.Article{
				Slug:      "5",
				Title:     "Timestamp Article",
				CreatedAt: "2023-05-15T10:30:00+0000Z",
				UpdatedAt: "2023-05-16T11:45:00+0000Z",
			},
		},
		{
			name: "Article with Zero FavoritesCount",
			article: Article{
				Model:          gorm.Model{ID: 6},
				Title:          "Zero Favorites Article",
				FavoritesCount: 0,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "6",
				Title:          "Zero Favorites Article",
				FavoritesCount: 0,
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
			name: "Article with Valid Fields and Additional Data",
			article: Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:          "Test Title",
				Description:    "Test Description",
				Body:           "Test Body",
				Tags:           []Tag{{Name: "Test Tag"}},
				Author:         User{Model: gorm.Model{ID: 1}},
				UserID:         1,
				FavoritesCount: 10,
			},
			wantErr: false,
		},
		{
			name: "Article with Extremely Long Title",
			article: Article{
				Title: string(make([]byte, 1000)),
				Body:  "Test Body",
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

