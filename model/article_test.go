package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Overwrite_3d4db6693d
ROOST_METHOD_SIG_HASH=Overwrite_22e8730976
*/
func TestOverwrite(t *testing.T) {
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
			name: "Overwrite No Fields",
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
			name: "Overwrite Only Title",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title",
			description: "",
			body:        "",
			expected: Article{
				Title:       "New Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
		},
		{
			name: "Overwrite Only Description",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "",
			description: "New Description",
			body:        "",
			expected: Article{
				Title:       "Initial Title",
				Description: "New Description",
				Body:        "Initial Body",
			},
		},
		{
			name: "Overwrite Only Body",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "",
			description: "",
			body:        "New Body",
			expected: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "New Body",
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
			name: "Overwrite with Special Characters",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "New Title 🚀",
			description: "New Description 你好",
			body:        "New Body ñ",
			expected: Article{
				Title:       "New Title 🚀",
				Description: "New Description 你好",
				Body:        "New Body ñ",
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
		})
	}
}

/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726
*/
func TestProtoArticle(t *testing.T) {
	tests := []struct {
		name      string
		article   Article
		favorited bool
		want      *pb.Article
	}{
		{
			name: "Valid Article to ProtoArticle with favorited true",
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
			favorited: true,
			want: &pb.Article{
				Slug:           "1",
				Title:          "Test Article",
				Description:    "This is a test article",
				Body:           "Article body content",
				TagList:        []string{"test", "article"},
				CreatedAt:      "2023-05-01T10:00:00+0000Z",
				UpdatedAt:      "2023-05-02T11:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 10,
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
			name: "Article with Nil Tags Slice",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body",
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank.",
		},
		{
			name: "Article with Very Long Title",
			article: Article{
				Title: string(make([]byte, 1000)),
				Body:  "Test Body",
				Tags:  []Tag{{Name: "TestTag"}},
			},
			wantErr: false,
		},
		{
			name: "Article with Special Characters in Title and Body",
			article: Article{
				Title: "Test Title 😊 with 特殊字符",
				Body:  "Test Body with 特殊字符 and 😊",
				Tags:  []Tag{{Name: "TestTag"}},
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
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
