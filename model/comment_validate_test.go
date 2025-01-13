package model

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Comment",
			comment: Comment{
				Body: "This is a valid comment",
			},
			wantErr: false,
		},
		{
			name: "Empty Body",
			comment: Comment{
				Body: "",
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Whitespace-only Body",
			comment: Comment{
				Body: "   \t   ",
			},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
		{
			name: "Very Long Body",
			comment: Comment{
				Body: string(make([]rune, 10000)),
			},
			wantErr: false,
		},
		{
			name: "Special Characters in Body",
			comment: Comment{
				Body: "!@#$%^&*()_+",
			},
			wantErr: false,
		},
		{
			name: "Unicode Characters in Body",
			comment: Comment{
				Body: "こんにちは世界 - Здравствуй, мир - مرحبا بالعالم",
			},
			wantErr: false,
		},
		{
			name:    "Null Body",
			comment: Comment{},
			wantErr: true,
			errMsg:  "body: cannot be blank.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("Comment.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TODO: Implement the following mock struct if needed
// type (ROOST_MOCK_STRUCT) struct {}
