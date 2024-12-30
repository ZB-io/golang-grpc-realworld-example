package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestUserProtoProfile remains unchanged

// TestUserProtoUser remains unchanged

// TestUserCheckPassword remains unchanged

// TestUserHashPassword remains unchanged

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		// Test cases remain unchanged
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("User.Validate() expected error, got nil")
					return
				}
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}
