package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Overwrite_3d4db6693d
ROOST_METHOD_SIG_HASH=Overwrite_22e8730976
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
			name: "Overwrite Maintaining Other Struct Fields",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				Tags:        []Tag{{Name: "tag1"}, {Name: "tag2"}},
				Author:      User{Model: gorm.Model{ID: 1}, Username: "author"},
				UserID:      1,
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
				Tags:        []Tag{{Name: "tag1"}, {Name: "tag2"}},
				Author:      User{Model: gorm.Model{ID: 1}, Username: "author"},
				UserID:      1,
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
			description: "新描述",
			body:        "新正文",
			expected: Article{
				Title:       "新标题",
				Description: "新描述",
				Body:        "新正文",
			},
		},
		{
			name: "Overwrite Idempotency",
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

			if tt.name == "Overwrite Maintaining Other Struct Fields" {
				if len(article.Tags) != len(tt.expected.Tags) {
					t.Errorf("Tags length = %v, want %v", len(article.Tags), len(tt.expected.Tags))
				}
				if article.Author.ID != tt.expected.Author.ID {
					t.Errorf("Author ID = %v, want %v", article.Author.ID, tt.expected.Author.ID)
				}
				if article.UserID != tt.expected.UserID {
					t.Errorf("UserID = %v, want %v", article.UserID, tt.expected.UserID)
				}
			}

			if tt.name == "Overwrite Idempotency" {
				article.Overwrite(tt.title, tt.description, tt.body)
				if article.Title != tt.expected.Title || article.Description != tt.expected.Description || article.Body != tt.expected.Body {
					t.Errorf("Second overwrite changed values: got %+v, want %+v", article, tt.expected)
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726
*/
func TestArticleProtoArticle(t *testing.T) {
	tests := []struct {
		name      string
		article   Article
		favorited bool
		want      *pb.Article
	}{
		{
			name: "Scenario 1: Convert a valid Article to ProtoArticle with favorited true",
			article: Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC),
				},
				Title:          "Test Article",
				Description:    "This is a test article",
				Body:           "Article body content",
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}},
				FavoritesCount: 10,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "This is a test article",
				Body:           "Article body content",
				TagList:        []string{"tag1", "tag2"},
				CreatedAt:      "2023-05-01T10:00:00+0000Z",
				UpdatedAt:      "2023-05-02T11:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 10,
			},
		},
		{
			name: "Scenario 2: Convert a valid Article to ProtoArticle with favorited false",
			article: Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Date(2023, 5, 3, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 4, 13, 0, 0, 0, time.UTC),
				},
				Title:          "Another Test Article",
				Description:    "This is another test article",
				Body:           "Another article body content",
				Tags:           []Tag{{Name: "tag3"}, {Name: "tag4"}},
				FavoritesCount: 5,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "2",
				Title:          "Another Test Article",
				Description:    "This is another test article",
				Body:           "Another article body content",
				TagList:        []string{"tag3", "tag4"},
				CreatedAt:      "2023-05-03T12:00:00+0000Z",
				UpdatedAt:      "2023-05-04T13:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 5,
			},
		},
		{
			name: "Scenario 3: Convert an Article with no Tags",
			article: Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Date(2023, 5, 5, 14, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 6, 15, 0, 0, 0, time.UTC),
				},
				Title:          "No Tags Article",
				Description:    "This article has no tags",
				Body:           "Article body without tags",
				Tags:           []Tag{},
				FavoritesCount: 2,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "3",
				Title:          "No Tags Article",
				Description:    "This article has no tags",
				Body:           "Article body without tags",
				TagList:        []string{},
				CreatedAt:      "2023-05-05T14:00:00+0000Z",
				UpdatedAt:      "2023-05-06T15:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 2,
			},
		},
		{
			name: "Scenario 4: Verify correct time formatting",
			article: Article{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: time.Date(2023, 5, 7, 16, 30, 45, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 8, 17, 45, 30, 0, time.UTC),
				},
				Title:          "Time Format Article",
				Description:    "This article tests time formatting",
				Body:           "Article body for time format test",
				Tags:           []Tag{{Name: "time"}},
				FavoritesCount: 1,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "4",
				Title:          "Time Format Article",
				Description:    "This article tests time formatting",
				Body:           "Article body for time format test",
				TagList:        []string{"time"},
				CreatedAt:      "2023-05-07T16:30:45+0000Z",
				UpdatedAt:      "2023-05-08T17:45:30+0000Z",
				Favorited:      true,
				FavoritesCount: 1,
			},
		},
		{
			name: "Scenario 5: Verify correct Slug generation",
			article: Article{
				Model: gorm.Model{
					ID:        12345,
					CreatedAt: time.Date(2023, 5, 9, 18, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 10, 19, 0, 0, 0, time.UTC),
				},
				Title:          "Slug Test Article",
				Description:    "This article tests slug generation",
				Body:           "Article body for slug test",
				Tags:           []Tag{{Name: "slug"}},
				FavoritesCount: 3,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "12345",
				Title:          "Slug Test Article",
				Description:    "This article tests slug generation",
				Body:           "Article body for slug test",
				TagList:        []string{"slug"},
				CreatedAt:      "2023-05-09T18:00:00+0000Z",
				UpdatedAt:      "2023-05-10T19:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 3,
			},
		},
		{
			name: "Scenario 6: Convert an Article with maximum values",
			article: Article{
				Model: gorm.Model{
					ID:        9999999,
					CreatedAt: time.Date(2023, 5, 11, 20, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 5, 12, 21, 0, 0, 0, time.UTC),
				},
				Title:          "Very Long Title " + fmt.Sprintf("%0100d", 0),
				Description:    "Very Long Description " + fmt.Sprintf("%0500d", 0),
				Body:           "Very Long Body " + fmt.Sprintf("%01000d", 0),
				Tags:           []Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}, {Name: "tag4"}, {Name: "tag5"}},
				FavoritesCount: 1000000,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "9999999",
				Title:          "Very Long Title " + fmt.Sprintf("%0100d", 0),
				Description:    "Very Long Description " + fmt.Sprintf("%0500d", 0),
				Body:           "Very Long Body " + fmt.Sprintf("%01000d", 0),
				TagList:        []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
				CreatedAt:      "2023-05-11T20:00:00+0000Z",
				UpdatedAt:      "2023-05-12T21:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 1000000,
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
*/
func TestValidate(t *testing.T) {
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
				Tags:  []Tag{{Name: "TestTag"}},
			},
			wantErr: false,
		},
		{
			name: "Article Missing Title",
			article: Article{
				Body: "Test Body",
				Tags: []Tag{{Name: "TestTag"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank.",
		},
		{
			name: "Article Missing Body",
			article: Article{
				Title: "Test Title",
				Tags:  []Tag{{Name: "TestTag"}},
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
			name: "Article with Minimum Valid Data",
			article: Article{
				Title: "T",
				Body:  "B",
				Tags:  []Tag{{Name: "T"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Very Long Title and Body",
			article: Article{
				Title: string(make([]rune, 10000)),
				Body:  string(make([]rune, 10000)),
				Tags:  []Tag{{Name: "TestTag"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Special Characters in Title and Body",
			article: Article{
				Title: "Tést Tîtlé 😊",
				Body:  "Tést Bôdy with špécîal châracters 🚀",
				Tags:  []Tag{{Name: "TestTag"}},
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
				if err == nil {
					t.Errorf("Article.Validate() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Article.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
