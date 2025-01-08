package undefined

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"fmt"
	"time"
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
			name: "Overwrite All Fields",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
			},
		},
		{
			name: "Overwrite Partial Fields",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title",
			description: "",
			body:        "New Body",
			expected: Article{
				Title:       "New Title",
				Description: "Initial Description",
				Body:        "New Body",
			},
		},
		{
			name: "Overwrite with Empty Strings",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "",
			description: "",
			body:        "",
			expected: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
		},
		{
			name: "Overwrite with Very Long Strings",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       string(make([]byte, 10000)),
			description: string(make([]byte, 10000)),
			body:        string(make([]byte, 10000)),
			expected: Article{
				Title:       string(make([]byte, 10000)),
				Description: string(make([]byte, 10000)),
				Body:        string(make([]byte, 10000)),
			},
		},
		{
			name: "Overwrite Maintaining Other Fields",
			initial: Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Initial Title",
				Description:    "Initial Description",
				Body:           "Initial Body",
				UserID:         2,
				FavoritesCount: 10,
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Model:          gorm.Model{ID: 1},
				Title:          "New Title",
				Description:    "New Description",
				Body:           "New Body",
				UserID:         2,
				FavoritesCount: 10,
			},
		},
		{
			name: "Overwrite with Unicode Characters",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "新标题",
			description: "新描述 🌟",
			body:        "新正文 😊",
			expected: Article{
				Title:       "新标题",
				Description: "新描述 🌟",
				Body:        "新正文 😊",
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
			assert.Equal(t, tt.expected.UserID, article.UserID)
			assert.Equal(t, tt.expected.FavoritesCount, article.FavoritesCount)
		})
	}
}


/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726

FUNCTION_DEF=func (a *Article) ProtoArticle(favorited bool) *pb.Article 

 */
func TestArticleProtoArticle(t *testing.T) {
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
					CreatedAt: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC),
				},
				Title:          "Test Article",
				Description:    "This is a test article",
				Body:           "Article body content",
				Tags:           []Tag{{Name: "test"}, {Name: "article"}},
				FavoritesCount: 10,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "This is a test article",
				Body:           "Article body content",
				TagList:        []string{"test", "article"},
				CreatedAt:      "2023-05-01T10:00:00+0000Z",
				UpdatedAt:      "2023-05-02T11:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 10,
			},
		},
		{
			name: "Favorited True",
			article: Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Date(2023, 5, 3, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 4, 13, 0, 0, 0, time.UTC),
				},
				Title:          "Favorited Article",
				Description:    "This is a favorited article",
				Body:           "Favorited article content",
				FavoritesCount: 5,
			},
			favorited: true,
			expected: &pb.Article{
				Slug:           "2",
				Title:          "Favorited Article",
				Description:    "This is a favorited article",
				Body:           "Favorited article content",
				TagList:        []string{},
				CreatedAt:      "2023-05-03T12:00:00+0000Z",
				UpdatedAt:      "2023-05-04T13:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 5,
			},
		},
		{
			name: "No Tags",
			article: Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Date(2023, 5, 5, 14, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 6, 15, 0, 0, 0, time.UTC),
				},
				Title:          "No Tags Article",
				Description:    "This article has no tags",
				Body:           "No tags content",
				FavoritesCount: 0,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "3",
				Title:          "No Tags Article",
				Description:    "This article has no tags",
				Body:           "No tags content",
				TagList:        []string{},
				CreatedAt:      "2023-05-05T14:00:00+0000Z",
				UpdatedAt:      "2023-05-06T15:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 0,
			},
		},
		{
			name: "Maximum Values",
			article: Article{
				Model: gorm.Model{
					ID:        ^uint(0),
					CreatedAt: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
					UpdatedAt: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
				},
				Title:          fmt.Sprintf("%s", make([]byte, 1000)),
				Description:    fmt.Sprintf("%s", make([]byte, 1000)),
				Body:           fmt.Sprintf("%s", make([]byte, 1000)),
				FavoritesCount: 2147483647,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           fmt.Sprintf("%d", ^uint(0)),
				Title:          fmt.Sprintf("%s", make([]byte, 1000)),
				Description:    fmt.Sprintf("%s", make([]byte, 1000)),
				Body:           fmt.Sprintf("%s", make([]byte, 1000)),
				TagList:        []string{},
				CreatedAt:      "9999-12-31T23:59:59+0000Z",
				UpdatedAt:      "9999-12-31T23:59:59+0000Z",
				Favorited:      false,
				FavoritesCount: 2147483647,
			},
		},
		{
			name: "Zero Values",
			article: Article{
				Model: gorm.Model{
					ID:        0,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Title:          "",
				Description:    "",
				Body:           "",
				FavoritesCount: 0,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "0",
				Title:          "",
				Description:    "",
				Body:           "",
				TagList:        []string{},
				CreatedAt:      "0001-01-01T00:00:00+0000Z",
				UpdatedAt:      "0001-01-01T00:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 0,
			},
		},
		{
			name: "Unicode Characters",
			article: Article{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: time.Date(2023, 5, 7, 16, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 8, 17, 0, 0, 0, time.UTC),
				},
				Title:          "Unicode 😊 Title",
				Description:    "Description with 日本語",
				Body:           "Body with Русский текст",
				Tags:           []Tag{{Name: "unicode"}, {Name: "😊"}},
				FavoritesCount: 7,
			},
			favorited: false,
			expected: &pb.Article{
				Slug:           "5",
				Title:          "Unicode 😊 Title",
				Description:    "Description with 日本語",
				Body:           "Body with Русский текст",
				TagList:        []string{"unicode", "😊"},
				CreatedAt:      "2023-05-07T16:00:00+0000Z",
				UpdatedAt:      "2023-05-08T17:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 7,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.article.ProtoArticle(tc.favorited)
			assert.Equal(t, tc.expected, result)
		})
	}
}

