package model

import (
	"testing"
	"unicode/utf8"
	"strings"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Overwrite_3d4db6693d
ROOST_METHOD_SIG_HASH=Overwrite_22e8730976


 */
func TestOverwrite(t *testing.T) {

	tests := []struct {
		name           string
		initialArticle Article
		inputTitle     string
		inputDesc      string
		inputBody      string
		expected       Article
		description    string
	}{
		{
			name: "Scenario 1: Update All Fields with Valid Values",
			initialArticle: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			inputTitle: "New Title",
			inputDesc:  "New Description",
			inputBody:  "New Body",
			expected: Article{
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
			},
			description: "Should update all fields when non-empty values are provided",
		},
		{
			name: "Scenario 2: Update No Fields with Empty Strings",
			initialArticle: Article{
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
			},
			inputTitle: "",
			inputDesc:  "",
			inputBody:  "",
			expected: Article{
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
			},
			description: "Should preserve original values when empty strings are provided",
		},
		{
			name: "Scenario 3: Update Only Title Field",
			initialArticle: Article{
				Title:       "Old Title",
				Description: "Keep Description",
				Body:        "Keep Body",
			},
			inputTitle: "Updated Title",
			inputDesc:  "",
			inputBody:  "",
			expected: Article{
				Title:       "Updated Title",
				Description: "Keep Description",
				Body:        "Keep Body",
			},
			description: "Should update only title while preserving other fields",
		},
		{
			name: "Scenario 4: Update Only Description Field",
			initialArticle: Article{
				Title:       "Keep Title",
				Description: "Old Description",
				Body:        "Keep Body",
			},
			inputTitle: "",
			inputDesc:  "Updated Description",
			inputBody:  "",
			expected: Article{
				Title:       "Keep Title",
				Description: "Updated Description",
				Body:        "Keep Body",
			},
			description: "Should update only description while preserving other fields",
		},
		{
			name: "Scenario 5: Update Only Body Field",
			initialArticle: Article{
				Title:       "Keep Title",
				Description: "Keep Description",
				Body:        "Old Body",
			},
			inputTitle: "",
			inputDesc:  "",
			inputBody:  "Updated Body",
			expected: Article{
				Title:       "Keep Title",
				Description: "Keep Description",
				Body:        "Updated Body",
			},
			description: "Should update only body while preserving other fields",
		},
		{
			name: "Scenario 6: Update with Special Characters",
			initialArticle: Article{
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
			},
			inputTitle: "Title with 特殊文字 and émojis 🎉",
			inputDesc:  "Description with ñ, é, ü characters",
			inputBody:  "Body with symbols: ©®™",
			expected: Article{
				Title:       "Title with 特殊文字 and émojis 🎉",
				Description: "Description with ñ, é, ü characters",
				Body:        "Body with symbols: ©®™",
			},
			description: "Should handle special characters and Unicode content correctly",
		},
		{
			name: "Scenario 7: Update with Maximum Length Values",
			initialArticle: Article{
				Title:       "Short Title",
				Description: "Short Description",
				Body:        "Short Body",
			},
			inputTitle: createLongString(1000),
			inputDesc:  createLongString(2000),
			inputBody:  createLongString(5000),
			expected: Article{
				Title:       createLongString(1000),
				Description: createLongString(2000),
				Body:        createLongString(5000),
			},
			description: "Should handle maximum length content properly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			article := tt.initialArticle

			article.Overwrite(tt.inputTitle, tt.inputDesc, tt.inputBody)

			if article.Title != tt.expected.Title {
				t.Errorf("Title mismatch - got: %v, want: %v", article.Title, tt.expected.Title)
			}
			if article.Description != tt.expected.Description {
				t.Errorf("Description mismatch - got: %v, want: %v", article.Description, tt.expected.Description)
			}
			if article.Body != tt.expected.Body {
				t.Errorf("Body mismatch - got: %v, want: %v", article.Body, tt.expected.Body)
			}

			if tt.name == "Scenario 6: Update with Special Characters" {
				if !utf8.ValidString(article.Title) {
					t.Error("Title contains invalid UTF-8 sequences")
				}
				if !utf8.ValidString(article.Description) {
					t.Error("Description contains invalid UTF-8 sequences")
				}
				if !utf8.ValidString(article.Body) {
					t.Error("Body contains invalid UTF-8 sequences")
				}
			}
		})
	}
}

func createLongString(length int) string {
	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = 'a'
	}
	return string(result)
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
			name: "Valid Article with all required fields",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body Content",
				Tags:  []Tag{{Name: "test-tag"}},
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Missing Title",
			article: Article{
				Body: "Test Body Content",
				Tags: []Tag{{Name: "test-tag"}},
			},
			wantErr: true,
			errMsg:  "title: cannot be blank",
		},
		{
			name: "Missing Body",
			article: Article{
				Title: "Test Title",
				Tags:  []Tag{{Name: "test-tag"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank",
		},
		{
			name: "Empty Tags Array",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body Content",
				Tags:  []Tag{},
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank",
		},
		{
			name: "Nil Tags Array",
			article: Article{
				Title: "Test Title",
				Body:  "Test Body Content",
			},
			wantErr: true,
			errMsg:  "tags: cannot be blank",
		},
		{
			name: "Multiple Validation Errors",
			article: Article{
				Title: "",
				Body:  "",
			},
			wantErr: true,
			errMsg:  "body: cannot be blank; tags: cannot be blank; title: cannot be blank",
		},
		{
			name: "Whitespace Only in Required Fields",
			article: Article{
				Title: "   ",
				Body:  "    ",
				Tags:  []Tag{{Name: "test-tag"}},
			},
			wantErr: true,
			errMsg:  "body: cannot be blank; title: cannot be blank",
		},
		{
			name: "Maximum Field Length Test",
			article: Article{
				Title: strings.Repeat("a", 1000),
				Body:  strings.Repeat("b", 10000),
				Tags:  []Tag{{Name: "test-tag"}},
			},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.name)

			err := tt.article.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
					t.Logf("Expected error received: %v", err)
				}
			} else {
				assert.NoError(t, err)
				t.Log("Validation passed as expected")
			}
		})
	}
}

