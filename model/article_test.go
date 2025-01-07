package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"strings"
)








/*
ROOST_METHOD_HASH=Overwrite_3d4db6693d
ROOST_METHOD_SIG_HASH=Overwrite_22e8730976

FUNCTION_DEF=func (a *Article) Overwrite(title, description, body string) 

 */
func TestArticleOverwrite(t *testing.T) {

	type testCase struct {
		name           string
		initialArticle Article
		inputTitle     string
		inputDesc      string
		inputBody      string
		expected       Article
	}

	baseTime := time.Now()

	baseArticle := Article{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: baseTime,
			UpdatedAt: baseTime,
		},
		Title:          "Original Title",
		Description:    "Original Description",
		Body:           "Original Body",
		UserID:         1,
		FavoritesCount: 0,
	}

	tests := []testCase{
		{
			name:           "Complete Article Overwrite",
			initialArticle: baseArticle,
			inputTitle:     "New Title",
			inputDesc:      "New Description",
			inputBody:      "New Body",
			expected: Article{
				Model:          baseArticle.Model,
				Title:          "New Title",
				Description:    "New Description",
				Body:           "New Body",
				UserID:         baseArticle.UserID,
				FavoritesCount: baseArticle.FavoritesCount,
			},
		},
		{
			name:           "Partial Article Overwrite - Title Only",
			initialArticle: baseArticle,
			inputTitle:     "Updated Title",
			inputDesc:      "",
			inputBody:      "",
			expected: Article{
				Model:          baseArticle.Model,
				Title:          "Updated Title",
				Description:    baseArticle.Description,
				Body:           baseArticle.Body,
				UserID:         baseArticle.UserID,
				FavoritesCount: baseArticle.FavoritesCount,
			},
		},
		{
			name:           "No Changes When All Empty Strings",
			initialArticle: baseArticle,
			inputTitle:     "",
			inputDesc:      "",
			inputBody:      "",
			expected:       baseArticle,
		},
		{
			name:           "Mixed Update with Empty and Non-Empty Values",
			initialArticle: baseArticle,
			inputTitle:     "",
			inputDesc:      "Only Description Updated",
			inputBody:      "",
			expected: Article{
				Model:          baseArticle.Model,
				Title:          baseArticle.Title,
				Description:    "Only Description Updated",
				Body:           baseArticle.Body,
				UserID:         baseArticle.UserID,
				FavoritesCount: baseArticle.FavoritesCount,
			},
		},
		{
			name:           "Update with Special Characters",
			initialArticle: baseArticle,
			inputTitle:     "Title with 特殊文字 and émojis 🎉",
			inputDesc:      "Description with $p€c!@l characters",
			inputBody:      "Body with üñîçødé characters",
			expected: Article{
				Model:          baseArticle.Model,
				Title:          "Title with 特殊文字 and émojis 🎉",
				Description:    "Description with $p€c!@l characters",
				Body:           "Body with üñîçødé characters",
				UserID:         baseArticle.UserID,
				FavoritesCount: baseArticle.FavoritesCount,
			},
		},
		{
			name:           "Update with Maximum Length Content",
			initialArticle: baseArticle,
			inputTitle:     string(make([]byte, 255)),
			inputDesc:      string(make([]byte, 255)),
			inputBody:      string(make([]byte, 255)),
			expected: Article{
				Model:          baseArticle.Model,
				Title:          string(make([]byte, 255)),
				Description:    string(make([]byte, 255)),
				Body:           string(make([]byte, 255)),
				UserID:         baseArticle.UserID,
				FavoritesCount: baseArticle.FavoritesCount,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			article := tc.initialArticle

			t.Logf("Testing scenario: %s", tc.name)
			t.Logf("Initial state: Title=%s, Description=%s, Body=%s",
				article.Title, article.Description, article.Body)

			article.Overwrite(tc.inputTitle, tc.inputDesc, tc.inputBody)

			if article.Title != tc.expected.Title {
				t.Errorf("Title mismatch - got: %s, want: %s", article.Title, tc.expected.Title)
			}
			if article.Description != tc.expected.Description {
				t.Errorf("Description mismatch - got: %s, want: %s",
					article.Description, tc.expected.Description)
			}
			if article.Body != tc.expected.Body {
				t.Errorf("Body mismatch - got: %s, want: %s", article.Body, tc.expected.Body)
			}

			t.Logf("Test completed successfully - Final state: Title=%s, Description=%s, Body=%s",
				article.Title, article.Description, article.Body)
		})
	}
}


/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726

FUNCTION_DEF=func (a *Article) ProtoArticle(favorited bool) *pb.Article 

 */
func TestArticleProtoArticle(t *testing.T) {

	type testCase struct {
		name      string
		article   *Article
		favorited bool
		want      *pb.Article
	}

	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []testCase{
		{
			name: "Scenario 1: Basic Article Conversion with Minimal Data",
			article: &Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: baseTime,
					UpdatedAt: baseTime,
				},
				Title:          "Test Title",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 0,
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "1",
				Title:          "Test Title",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 0,
				Favorited:      false,
				CreatedAt:      baseTime.Format(ISO8601),
				UpdatedAt:      baseTime.Format(ISO8601),
				TagList:        []string{},
			},
		},
		{
			name: "Scenario 2: Article Conversion with Tags",
			article: &Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: baseTime,
					UpdatedAt: baseTime,
				},
				Title:       "Tagged Article",
				Description: "Article with tags",
				Body:        "Content with tags",
				Tags: []Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "2",
				Title:          "Tagged Article",
				Description:    "Article with tags",
				Body:           "Content with tags",
				TagList:        []string{"tag1", "tag2"},
				CreatedAt:      baseTime.Format(ISO8601),
				UpdatedAt:      baseTime.Format(ISO8601),
				Favorited:      false,
				FavoritesCount: 0,
			},
		},
		{
			name: "Scenario 3: Article Conversion with Favorited Status",
			article: &Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: baseTime,
					UpdatedAt: baseTime,
				},
				Title:          "Favorited Article",
				Description:    "Popular article",
				Body:           "Favorited content",
				FavoritesCount: 10,
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "3",
				Title:          "Favorited Article",
				Description:    "Popular article",
				Body:           "Favorited content",
				FavoritesCount: 10,
				Favorited:      true,
				CreatedAt:      baseTime.Format(ISO8601),
				UpdatedAt:      baseTime.Format(ISO8601),
				TagList:        []string{},
			},
		},
		{
			name: "Scenario 4: Article Conversion with Zero Values",
			article: &Article{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: baseTime,
					UpdatedAt: baseTime,
				},
				Title:       "",
				Description: "",
				Body:        "",
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "4",
				Title:          "",
				Description:    "",
				Body:           "",
				FavoritesCount: 0,
				Favorited:      false,
				CreatedAt:      baseTime.Format(ISO8601),
				UpdatedAt:      baseTime.Format(ISO8601),
				TagList:        []string{},
			},
		},
		{
			name: "Scenario 5: Article Conversion with Maximum Values",
			article: &Article{
				Model: gorm.Model{
					ID:        ^uint(0),
					CreatedAt: baseTime,
					UpdatedAt: baseTime,
				},
				Title:          "Max Values",
				Description:    "Testing maximum values",
				Body:           "Content",
				FavoritesCount: int32(^uint32(0) >> 1),
			},
			favorited: true,
			want: &pb.Article{
				Slug:           "18446744073709551615",
				Title:          "Max Values",
				Description:    "Testing maximum values",
				Body:           "Content",
				FavoritesCount: int32(^uint32(0) >> 1),
				Favorited:      true,
				CreatedAt:      baseTime.Format(ISO8601),
				UpdatedAt:      baseTime.Format(ISO8601),
				TagList:        []string{},
			},
		},
		{
			name: "Scenario 6: Article Conversion with Special Characters",
			article: &Article{
				Model: gorm.Model{
					ID:        6,
					CreatedAt: baseTime,
					UpdatedAt: baseTime,
				},
				Title:       "Special 特殊 Characters !@#$%^&*()",
				Description: "Description with émojis 🎉",
				Body:        "Body with\nmultiple\nlines and unicode ★",
			},
			favorited: false,
			want: &pb.Article{
				Slug:           "6",
				Title:          "Special 特殊 Characters !@#$%^&*()",
				Description:    "Description with émojis 🎉",
				Body:           "Body with\nmultiple\nlines and unicode ★",
				FavoritesCount: 0,
				Favorited:      false,
				CreatedAt:      baseTime.Format(ISO8601),
				UpdatedAt:      baseTime.Format(ISO8601),
				TagList:        []string{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Testing:", tc.name)

			got := tc.article.ProtoArticle(tc.favorited)

			assert.Equal(t, tc.want.Slug, got.Slug, "Slug mismatch")
			assert.Equal(t, tc.want.Title, got.Title, "Title mismatch")
			assert.Equal(t, tc.want.Description, got.Description, "Description mismatch")
			assert.Equal(t, tc.want.Body, got.Body, "Body mismatch")
			assert.Equal(t, tc.want.FavoritesCount, got.FavoritesCount, "FavoritesCount mismatch")
			assert.Equal(t, tc.want.Favorited, got.Favorited, "Favorited status mismatch")
			assert.Equal(t, tc.want.CreatedAt, got.CreatedAt, "CreatedAt timestamp mismatch")
			assert.Equal(t, tc.want.UpdatedAt, got.UpdatedAt, "UpdatedAt timestamp mismatch")
			assert.Equal(t, tc.want.TagList, got.TagList, "TagList mismatch")

			t.Log("Test completed successfully")
		})
	}
}


/*
ROOST_METHOD_HASH=Validate_f6d09c3ac5
ROOST_METHOD_SIG_HASH=Validate_99e41aac91

FUNCTION_DEF=func (a Article) Validate() error 

 */
func TestArticleValidate() {

	tests := []struct {
		name    string
		article Article
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Article with all required fields",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
				Tags:  []Tag{{Name: "test"}},
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Missing Title",
			article: Article{
				Body: "This is a test article body",
				Tags: []Tag{{Name: "test"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank",
		},
		{
			name: "Missing Body",
			article: Article{
				Title: "Test Article",
				Tags:  []Tag{{Name: "test"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank",
		},
		{
			name: "Empty Tags Slice",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
				Tags:  []Tag{},
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank",
		},
		{
			name: "Nil Tags Slice",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank",
		},
		{
			name: "Multiple Validation Errors",
			article: Article{
				Description: "Test Description",
			},
			wantErr: true,
			errMsg:  "title: cannot be blank; body: cannot be blank; tags: cannot be blank",
		},
		{
			name: "Whitespace Only in Required Fields",
			article: Article{
				Title: "   ",
				Body:  "\t\n",
				Tags:  []Tag{{Name: "test"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank; body: cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.name)

			err := tt.article.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Article.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
					return
				}

				errStr := err.Error()

				actualMsgs := strings.Split(errStr, "; ")
				expectedMsgs := strings.Split(tt.errMsg, "; ")
				sort.Strings(actualMsgs)
				sort.Strings(expectedMsgs)

				if !reflect.DeepEqual(actualMsgs, expectedMsgs) {
					t.Errorf("Expected error message '%s', but got '%s'", tt.errMsg, errStr)
				}
			}

			if err == nil {
				t.Log("Validation passed successfully")
			} else {
				t.Logf("Validation failed as expected with error: %v", err)
			}
		})
	}
}

