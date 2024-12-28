package model

import (
	"strings"
	"testing"
	validation "github.com/go-ozzo/ozzo-validation"
)



func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		comment   Comment
		expectErr bool
	}{
		{
			name: "Valid Comment Body",
			comment: Comment{
				Body: "This is a valid comment",
			},
			expectErr: false,
		},
		{
			name: "Empty Comment Body",
			comment: Comment{
				Body: "",
			},
			expectErr: true,
		},
		{
			name: "Whitespace-only Comment Body",
			comment: Comment{
				Body: "    ",
			},
			expectErr: true,
		},
		{
			name: "Large Comment Body",
			comment: Comment{
				Body: strings.Repeat("a", 10*1024),
			},
			expectErr: false,
		},
		{
			name: "Comment Body with Special Characters",
			comment: Comment{
				Body: "!@#%&*()",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected an error, but got none for %s", tt.name)
				} else {
					t.Logf("success: expected an error and received: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s: %v", tt.name, err)
				} else {
					t.Logf("success: no error as expected for %s", tt.name)
				}
			}
		})
	}
}


