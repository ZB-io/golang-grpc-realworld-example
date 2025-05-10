package model

import (
	"fmt"
	"strings"
	"testing"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}



type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
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
			name: "Updating all fields of an article",
			initial: Article{
				Title:       "Old Title",
				Description: "Old Description",
				Body:        "Old Body",
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
			name: "Partially updating fields (only title)",
			initial: Article{
				Title:       "Old Title",
				Description: "Old Description",
				Body:        "Old Body",
			},
			title:       "New Title",
			description: "",
			body:        "",
			expected: Article{
				Title:       "New Title",
				Description: "Old Description",
				Body:        "Old Body",
			},
		},
		{
			name: "No fields update with all empty strings",
			initial: Article{
				Title:       "Old Title",
				Description: "Old Description",
				Body:        "Old Body",
			},
			title:       "",
			description: "",
			body:        "",
			expected: Article{
				Title:       "Old Title",
				Description: "Old Description",
				Body:        "Old Body",
			},
		},
		{
			name:        "Empty initial fields updated with non-empty strings",
			initial:     Article{},
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
			name: "Edge Case: Long string updates",
			initial: Article{
				Title:       "Short Title",
				Description: "Short Description",
				Body:        "Short Body",
			},
			title:       "L" + repeat("ong Title", 500),
			description: "L" + repeat("ong Description", 500),
			body:        "L" + repeat("ong Body", 500),
			expected: Article{
				Title:       "L" + repeat("ong Title", 500),
				Description: "L" + repeat("ong Description", 500),
				Body:        "L" + repeat("ong Body", 500),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running scenario: %s", tt.name)

			tt.initial.Overwrite(tt.title, tt.description, tt.body)

			if tt.initial.Title != tt.expected.Title {
				t.Errorf("Expected Title: %v, got: %v", tt.expected.Title, tt.initial.Title)
			}
			if tt.initial.Description != tt.expected.Description {
				t.Errorf("Expected Description: %v, got: %v", tt.expected.Description, tt.initial.Description)
			}
			if tt.initial.Body != tt.expected.Body {
				t.Errorf("Expected Body: %v, got: %v", tt.expected.Body, tt.initial.Body)
			}

			t.Log("Test passed for scenario:", tt.name)
		})
	}
}
func repeat(str string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(str, count)
}
