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
	t.Run("Scenario 1: Update only title of the Article", func(t *testing.T) {
		article := Article{Title: "Old Title", Description: "Initial Description", Body: "Initial Body"}
		newTitle := "New Title"
		expected := Article{Title: newTitle, Description: "Initial Description", Body: "Initial Body"}

		article.Overwrite(newTitle, "", "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated the title only while preserving other fields.")
		}
	})

	t.Run("Scenario 2: Update only description of the Article", func(t *testing.T) {
		article := Article{Title: "Initial Title", Description: "Old Description", Body: "Initial Body"}
		newDescription := "New Description"
		expected := Article{Title: "Initial Title", Description: newDescription, Body: "Initial Body"}

		article.Overwrite("", newDescription, "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated the description only while preserving other fields.")
		}
	})

	t.Run("Scenario 3: Update only body of the Article", func(t *testing.T) {
		article := Article{Title: "Initial Title", Description: "Initial Description", Body: "Old Body"}
		newBody := "New Body"
		expected := Article{Title: "Initial Title", Description: "Initial Description", Body: newBody}

		article.Overwrite("", "", newBody)

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated the body only while preserving other fields.")
		}
	})

	t.Run("Scenario 4: Update all fields of the Article", func(t *testing.T) {
		article := Article{Title: "Old Title", Description: "Old Description", Body: "Old Body"}
		newTitle := "New Title"
		newDescription := "New Description"
		newBody := "New Body"
		expected := Article{Title: newTitle, Description: newDescription, Body: newBody}

		article.Overwrite(newTitle, newDescription, newBody)

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated all fields of the article.")
		}
	})

	t.Run("Scenario 5: No update when all parameters are empty", func(t *testing.T) {
		article := Article{Title: "Initial Title", Description: "Initial Description", Body: "Initial Body"}
		expected := Article{Title: "Initial Title", Description: "Initial Description", Body: "Initial Body"}

		article.Overwrite("", "", "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("No changes to the article when all update parameters are empty.")
		}
	})

	t.Run("Scenario 6: Partial update of non-empty fields only", func(t *testing.T) {
		article := Article{Title: "Old Title", Description: "Old Description", Body: "Initial Body"}
		newTitle := "New Title"
		expected := Article{Title: newTitle, Description: "Old Description", Body: "Initial Body"}

		article.Overwrite(newTitle, "", "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated non-empty fields only, preserving remaining fields.")
		}
	})
}

/*
ROOST_METHOD_HASH=ProtoArticle_4b12477d53
ROOST_METHOD_SIG_HASH=ProtoArticle_31d9b4d726
*/
func TestArticleProtoArticle(t *testing.T) {
	type fields struct {
		ID             uint
		Title          string
		Description    string
		Body           string
		FavoritesCount int32
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
				Model:          gorm.Model{ID: tc.fields.ID},
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

/*
ROOST_METHOD_HASH=Validate_f6d09c3ac5
ROOST_METHOD_SIG_HASH=Validate_99e41aac91
*/
func TestArticleValidate(t *testing.T) {
	tests := []struct {
		name    string
		article Article
		wantErr bool
	}{
		{
			name: "Scenario 1: Successful Validation of a Complete Article",
			article: Article{
				Title: "Valid Title",
				Body:  "This is a valid body content.",
				Tags:  []Tag{{Name: "Tag1"}, {Name: "Tag2"}},
			},
			wantErr: false,
		},
		{
			name: "Scenario 2: Validation Fails for Missing Title",
			article: Article{
				Title: "",
				Body:  "Has Body",
				Tags:  []Tag{{Name: "Tag"}},
			},
			wantErr: true,
		},
		{
			name: "Scenario 3: Validation Fails for Missing Body",
			article: Article{
				Title: "Has Title",
				Body:  "",
				Tags:  []Tag{{Name: "Tag"}},
			},
			wantErr: true,
		},
		{
			name: "Scenario 4: Validation Fails for Missing Tags",
			article: Article{
				Title: "Has Title",
				Body:  "Has Body",
				Tags:  []Tag{},
			},
			wantErr: true,
		},
		{
			name: "Scenario 5: Validation Failure with All Fields Empty",
			article: Article{
				Title: "",
				Body:  "",
				Tags:  []Tag{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test case: %s", tt.name)
			err := tt.article.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation error status = %v, wantErr = %v", err != nil, tt.wantErr)
			}
			if err != nil {
				t.Logf("Expected error: %v", err)
			} else {
				t.Logf("Validation passed with no errors")
			}
		})
	}
}
