package model

import (
	"testing"
	"time"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)



type Tag struct {
	Name string
}



func TestArticleProtoArticle(t *testing.T) {
	type fields struct {
		ID             uint
		Title          string
		Description    string
		Body           string
		FavoritesCount int
		CreatedAt      time.Time
		UpdatedAt      time.Time
		Tags           []Tag
	}
	type testCase struct {
		name            string
		fields          fields
		favorited       bool
		expectedArticle pb.Article
	}

	tests := []testCase{
		{
			name: "Convert Article with No Tags to ProtoArticle",
			fields: fields{
				ID:          1,
				Title:       "Test Article",
				Description: "Description of test article",
				Body:        "This is the body of the test article",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Tags:        []Tag{},
			},
			favorited: false,
			expectedArticle: pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "Description of test article",
				Body:           "This is the body of the test article",
				FavoritesCount: 0,
				Favorited:      false,
				TagList:        []string{},
			},
		},
		{
			name: "Verify Favorited Field Transformation",
			fields: fields{
				ID:          2,
				Title:       "Favorite Article",
				Description: "An article to check favorited feature",
				Body:        "Body of favorite article",
				Tags:        []Tag{},
			},
			favorited: true,
			expectedArticle: pb.Article{
				Slug:      "2",
				Title:     "Favorite Article",
				Favorited: true,
				TagList:   []string{},
			},
		},
		{
			name: "Handle Article with Multiple Tags",
			fields: fields{
				ID:          3,
				Title:       "Multi-tag Article",
				Description: "Article with multiple tags",
				Body:        "Body of a multi-tag article",
				Tags: []Tag{
					{Name: "Go"},
					{Name: "Golang"},
				},
			},
			favorited: false,
			expectedArticle: pb.Article{
				Slug:    "3",
				Title:   "Multi-tag Article",
				TagList: []string{"Go", "Golang"},
			},
		},
		{
			name: "Validate Date Formatting",
			fields: fields{
				ID:        4,
				Title:     "Date Formatted Article",
				CreatedAt: time.Date(2023, time.January, 1, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, time.January, 2, 11, 0, 0, 0, time.UTC),
			},
			expectedArticle: pb.Article{
				Slug:      "4",
				Title:     "Date Formatted Article",
				CreatedAt: "2023-01-01T10:00:00+0000Z",
				UpdatedAt: "2023-01-02T11:00:00+0000Z",
			},
		},
		{
			name: "Verify Transformation with Maximum Field Length",
			fields: fields{
				ID:          5,
				Title:       string(make([]byte, 1000)),
				Description: string(make([]byte, 5000)),
				Body:        string(make([]byte, 10000)),
			},
			expectedArticle: pb.Article{
				Slug:        "5",
				Title:       string(make([]byte, 1000)),
				Description: string(make([]byte, 5000)),
				Body:        string(make([]byte, 10000)),
			},
		},
		{
			name: "Confirm Zero Favorites Count Handling",
			fields: fields{
				ID:             6,
				Title:          "Zero Favorites Article",
				FavoritesCount: 0,
			},
			expectedArticle: pb.Article{
				Slug:           "6",
				Title:          "Zero Favorites Article",
				FavoritesCount: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			article := &Article{
				ID:             tc.fields.ID,
				Title:          tc.fields.Title,
				Description:    tc.fields.Description,
				Body:           tc.fields.Body,
				FavoritesCount: tc.fields.FavoritesCount,
				CreatedAt:      tc.fields.CreatedAt,
				UpdatedAt:      tc.fields.UpdatedAt,
				Tags:           tc.fields.Tags,
			}

			result := article.ProtoArticle(tc.favorited)

			assert.Equal(t, tc.expectedArticle.Slug, result.Slug)
			assert.Equal(t, tc.expectedArticle.Title, result.Title)
			assert.Equal(t, tc.expectedArticle.Description, result.Description)
			assert.Equal(t, tc.expectedArticle.Body, result.Body)
			assert.Equal(t, tc.expectedArticle.FavoritesCount, result.FavoritesCount)
			assert.Equal(t, tc.expectedArticle.Favorited, result.Favorited)
			assert.Equal(t, tc.expectedArticle.TagList, result.TagList)

			if tc.expectedArticle.CreatedAt != "" {
				assert.Equal(t, tc.expectedArticle.CreatedAt, result.CreatedAt)
			}
			if tc.expectedArticle.UpdatedAt != "" {
				assert.Equal(t, tc.expectedArticle.UpdatedAt, result.UpdatedAt)
			}
		})
	}
}


