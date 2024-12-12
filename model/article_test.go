package model

import (
	"errors"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"

	// "fmt"
	"math"
	"time"

	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

// const ISO8601 = "2006-01-02T15:04:05-0700Z"
/*
ROOST_METHOD_HASH=Overwrite_3d4db6693d
ROOST_METHOD_SIG_HASH=Overwrite_22e8730976


*/
func TestArticleOverwrite(t *testing.T) {
	type args struct {
		title       string
		description string
		body        string
	}
	tests := []struct {
		name     string
		initial  Article
		args     args
		expected Article
	}{
		{
			name: "Scenario 1: Overwrite All Fields with New Non-Empty Values",
			initial: Article{
				Title:       "Old Title",
				Description: "Old Description",
				Body:        "Old Body",
			},
			args: args{
				title:       "New Title",
				description: "New Description",
				body:        "New Body",
			},
			expected: Article{
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
			},
		},
		{
			name: "Scenario 2: Keep Original Values When New Values Are Empty",
			initial: Article{
				Title:       "Existing Title",
				Description: "Existing Description",
				Body:        "Existing Body",
			},
			args: args{
				title:       "",
				description: "",
				body:        "",
			},
			expected: Article{
				Title:       "Existing Title",
				Description: "Existing Description",
				Body:        "Existing Body",
			},
		},
		{
			name: "Scenario 3: Update Only the Title Field",
			initial: Article{
				Title:       "Final Title",
				Description: "Unchanged Description",
				Body:        "Unchanged Body",
			},
			args: args{
				title:       "New Final Title",
				description: "",
				body:        "",
			},
			expected: Article{
				Title:       "New Final Title",
				Description: "Unchanged Description",
				Body:        "Unchanged Body",
			},
		},
		{
			name: "Scenario 4: Simultaneous Partial and Full Field Updates",
			initial: Article{
				Title:       "Mixed Title",
				Description: "Mixed Description",
				Body:        "Mixed Body",
			},
			args: args{
				title:       "Updated Title",
				description: "",
				body:        "Updated Body",
			},
			expected: Article{
				Title:       "Updated Title",
				Description: "Mixed Description",
				Body:        "Updated Body",
			},
		},
		{
			name: "Scenario 5: No Operation with All Empty Arguments",
			initial: Article{
				Title:       "No Change Title",
				Description: "No Change Description",
				Body:        "No Change Body",
			},
			args: args{
				title:       "",
				description: "",
				body:        "",
			},
			expected: Article{
				Title:       "No Change Title",
				Description: "No Change Description",
				Body:        "No Change Body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := tt.initial
			article.Overwrite(tt.args.title, tt.args.description, tt.args.body)

			if article.Title != tt.expected.Title || article.Description != tt.expected.Description || article.Body != tt.expected.Body {
				t.Errorf("Test: %s - Overwrite() error: got [%v], expected [%v]", tt.name, article, tt.expected)
			} else {
				t.Logf("Test: %s - Passed. Article updated as expected.", tt.name)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Validate_f6d09c3ac5
ROOST_METHOD_SIG_HASH=Validate_99e41aac91


 */
func TestArticleValidate(t *testing.T) {
	type test struct {
		name     string
		article  Article
		expError bool
		errMsg   string
	}

	tests := []test{
		{
			name: "Validation Passes for a Fully Populated and Correct Article",
			article: Article{
				Title: "Valid Title",
				Body:  "Valid Body",
				Tags:  []Tag{{Name: "ValidTag"}},
			},
			expError: false,
			errMsg:   "expected no error, got error",
		},
		{
			name: "Validation Fails for Missing Title",
			article: Article{
				Body: "Valid Body",
				Tags: []Tag{{Name: "ValidTag"}},
			},
			expError: true,
			errMsg:   "expected error for missing Title, got nil",
		},
		{
			name: "Validation Fails for Missing Body",
			article: Article{
				Title: "Valid Title",
				Tags:  []Tag{{Name: "ValidTag"}},
			},
			expError: true,
			errMsg:   "expected error for missing Body, got nil",
		},
		{
			name: "Validation Fails for Missing Tags",
			article: Article{
				Title: "Valid Title",
				Body:  "Valid Body",
			},
			expError: true,
			errMsg:   "expected error for missing Tags, got nil",
		},
		{
			name: "Validation Fails for Empty Tags",
			article: Article{
				Title: "Valid Title",
				Body:  "Valid Body",
				Tags:  []Tag{},
			},
			expError: true,
			errMsg:   "expected error for empty Tags, got nil",
		},
		{
			name: "Comprehensive Validation Covering Multiple Missing Fields",
			article: Article{
				Tags: nil,
			},
			expError: true,
			errMsg:   "expected error for multiple missing fields (Title, Body, Tags), got nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()

			if (err != nil) != tt.expError {
				t.Errorf("%s: %s", tt.name, tt.errMsg)
				if err != nil {
					t.Logf("unexpected error: %v", err)
				}
				return
			}

			if err != nil {
				var errs validation.Errors
				if errors.As(err, &errs) {
					missingFields := []string{}
					for fieldName := range errs {
						missingFields = append(missingFields, fieldName)
					}
					t.Logf("validation errors: missing fields = [%s]", strings.Join(missingFields, ", "))
				} else {
					t.Errorf("unexpected error type, expected validation.Errors, got: %T", err)
				}
			} else {
				t.Logf("validation passed with no errors for %s", tt.name)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726


 */
func TestProtoArticle(t *testing.T) {
	type args struct {
		article    Article
		favorited  bool
		expectedPb pb.Article
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal Operation with Favorited True",
			args: args{
				article: Article{
					Model:          gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:          "Sample Title",
					Description:    "Sample Description",
					Body:           "Sample Body",
					FavoritesCount: 10,
					Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "tag1"}, {Model: gorm.Model{ID: 2}, Name: "tag2"}},
				},
				favorited: true,
				expectedPb: pb.Article{
					Slug:           "1",
					Title:          "Sample Title",
					Description:    "Sample Description",
					Body:           "Sample Body",
					FavoritesCount: 10,
					Favorited:      true,
					TagList:        []string{"tag1", "tag2"},
				},
			},
		},
		{
			name: "Normal Operation with Favorited False",
			args: args{
				article: Article{
					Model:          gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:          "Sample Title 2",
					Description:    "Sample Description 2",
					Body:           "Sample Body 2",
					FavoritesCount: 5,
					Tags:           []Tag{{Model: gorm.Model{ID: 1}, Name: "tag1"}},
				},
				favorited: false,
				expectedPb: pb.Article{
					Slug:           "2",
					Title:          "Sample Title 2",
					Description:    "Sample Description 2",
					Body:           "Sample Body 2",
					FavoritesCount: 5,
					Favorited:      false,
					TagList:        []string{"tag1"},
				},
			},
		},
		{
			name: "Edge Case - Article with No Tags",
			args: args{
				article: Article{
					Model:          gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:          "Tagless Article",
					Description:    "Description",
					Body:           "Body",
					FavoritesCount: 0,
				},
				favorited: false,
				expectedPb: pb.Article{
					Slug:           "3",
					Title:          "Tagless Article",
					Description:    "Description",
					Body:           "Body",
					FavoritesCount: 0,
					Favorited:      false,
					TagList:        nil,
				},
			},
		},
		{
			name: "Edge Case - Article with Maximum Integer Favorites Count",
			args: args{
				article: Article{
					Model:          gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:          "Max Count Article",
					Description:    "Description",
					Body:           "Body",
					FavoritesCount: math.MaxInt32,
					Tags:           []Tag{},
				},
				favorited: false,
				expectedPb: pb.Article{
					Slug:           "4",
					Title:          "Max Count Article",
					Description:    "Description",
					Body:           "Body",
					FavoritesCount: math.MaxInt32,
					Favorited:      false,
					TagList:        nil,
				},
			},
		},
		{
			name: "Edge Case - Article with Timestamps and Date Formats",
			args: args{
				article: Article{
					Model:          gorm.Model{ID: 5, CreatedAt: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2021, time.February, 1, 0, 0, 0, 0, time.UTC)},
					Title:          "Article with Timestamps",
					Description:    "Description",
					Body:           "Body",
					FavoritesCount: 0,
					Tags:           []Tag{},
				},
				favorited: false,
				expectedPb: pb.Article{
					Slug:           "5",
					Title:          "Article with Timestamps",
					Description:    "Description",
					Body:           "Body",
					FavoritesCount: 0,
					Favorited:      false,
					TagList:        nil,
					CreatedAt:      "2021-01-01T00:00:00-0700Z",
					UpdatedAt:      "2021-02-01T00:00:00-0700Z",
				},
			},
		},
		{
			name: "Empty Article Content",
			args: args{
				article: Article{
					Model:          gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Title:          "",
					Description:    "",
					Body:           "",
					FavoritesCount: 0,
					Tags:           []Tag{},
				},
				favorited: false,
				expectedPb: pb.Article{
					Slug:           "6",
					Title:          "",
					Description:    "",
					Body:           "",
					FavoritesCount: 0,
					Favorited:      false,
					TagList:        nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.article.ProtoArticle(tt.args.favorited)

			t.Logf("Running test case: %s", tt.name)

			if got.Slug != tt.args.expectedPb.Slug {
				t.Errorf("expected Slug %v, got %v", tt.args.expectedPb.Slug, got.Slug)
			}
			if got.Title != tt.args.expectedPb.Title {
				t.Errorf("expected Title %v, got %v", tt.args.expectedPb.Title, got.Title)
			}
			if got.Description != tt.args.expectedPb.Description {
				t.Errorf("expected Description %v, got %v", tt.args.expectedPb.Description, got.Description)
			}
			if got.Body != tt.args.expectedPb.Body {
				t.Errorf("expected Body %v, got %v", tt.args.expectedPb.Body, got.Body)
			}
			if got.FavoritesCount != tt.args.expectedPb.FavoritesCount {
				t.Errorf("expected FavoritesCount %v, got %v", tt.args.expectedPb.FavoritesCount, got.FavoritesCount)
			}
			if got.Favorited != tt.args.expectedPb.Favorited {
				t.Errorf("expected Favorited %v, got %v", tt.args.expectedPb.Favorited, got.Favorited)
			}
			if len(got.TagList) != len(tt.args.expectedPb.TagList) {
				t.Errorf("expected TagList length %v, got %v", len(tt.args.expectedPb.TagList), len(got.TagList))
			}
			for i := range got.TagList {
				if got.TagList[i] != tt.args.expectedPb.TagList[i] {
					t.Errorf("expected TagList[%d] %v, got %v", i, tt.args.expectedPb.TagList[i], got.TagList[i])
				}
			}
		})
	}
}

