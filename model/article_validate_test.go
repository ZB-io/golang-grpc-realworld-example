package model

import (
	"testing"
	validation "github.com/go-ozzo/ozzo-validation"
)



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
				Tags:  []string{"Tag1", "Tag2"},
			},
			wantErr: false,
		},
		{
			name: "Scenario 2: Validation Fails for Missing Title",
			article: Article{
				Title: "",
				Body:  "Has Body",
				Tags:  []string{"Tag"},
			},
			wantErr: true,
		},
		{
			name: "Scenario 3: Validation Fails for Missing Body",
			article: Article{
				Title: "Has Title",
				Body:  "",
				Tags:  []string{"Tag"},
			},
			wantErr: true,
		},
		{
			name: "Scenario 4: Validation Fails for Missing Tags",
			article: Article{
				Title: "Has Title",
				Body:  "Has Body",
				Tags:  []string{},
			},
			wantErr: true,
		},
		{
			name: "Scenario 5: Validation Failure with All Fields Empty",
			article: Article{
				Title: "",
				Body:  "",
				Tags:  []string{},
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

