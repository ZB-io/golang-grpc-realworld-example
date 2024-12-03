package model

import (
	"testing"
	"unicode/utf8"
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
			name: "Scenario 1: Update All Fields with Valid Values",
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
			name: "Scenario 2: Update Only Title",
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
			name: "Scenario 3: Update Only Description",
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
			name: "Scenario 4: Update Only Body",
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
			name: "Scenario 5: No Updates with Empty Strings",
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
			name: "Scenario 6: Update with Special Characters",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       "Title with 特殊文字 and émojis 🎉",
			description: "Description with ñ, é, ü characters",
			body:        "Body with symbols @#$%^&*()",
			expected: Article{
				Title:       "Title with 特殊文字 and émojis 🎉",
				Description: "Description with ñ, é, ü characters",
				Body:        "Body with symbols @#$%^&*()",
			},
		},
		{
			name: "Scenario 7: Update with Maximum Length Values",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			title:       string(make([]byte, 255)),
			description: string(make([]byte, 1000)),
			body:        string(make([]byte, 5000)),
			expected: Article{
				Title:       string(make([]byte, 255)),
				Description: string(make([]byte, 1000)),
				Body:        string(make([]byte, 5000)),
			},
		},
		{
			name: "Scenario 8: Preserve Related Data During Update",
			initial: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
				Tags: []Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
				UserID:         1,
				FavoritesCount: 5,
			},
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
			expected: Article{
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
				Tags: []Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
				UserID:         1,
				FavoritesCount: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.name)

			article := tt.initial

			article.Overwrite(tt.title, tt.description, tt.body)

			if article.Title != tt.expected.Title {
				t.Errorf("Title mismatch\nwant: %v\ngot:  %v", tt.expected.Title, article.Title)
			}

			if article.Description != tt.expected.Description {
				t.Errorf("Description mismatch\nwant: %v\ngot:  %v", tt.expected.Description, article.Description)
			}

			if article.Body != tt.expected.Body {
				t.Errorf("Body mismatch\nwant: %v\ngot:  %v", tt.expected.Body, article.Body)
			}

			if tt.name == "Scenario 6: Update with Special Characters" {
				if !utf8.ValidString(article.Title) || !utf8.ValidString(article.Description) || !utf8.ValidString(article.Body) {
					t.Error("Invalid UTF-8 encoding in special characters test")
				}
			}

			if tt.name == "Scenario 8: Preserve Related Data During Update" {
				if len(article.Tags) != len(tt.expected.Tags) {
					t.Error("Tags were not preserved during update")
				}
				if article.UserID != tt.expected.UserID {
					t.Error("UserID was not preserved during update")
				}
				if article.FavoritesCount != tt.expected.FavoritesCount {
					t.Error("FavoritesCount was not preserved during update")
				}
			}

			t.Log("Test completed successfully:", tt.name)
		})
	}
}

