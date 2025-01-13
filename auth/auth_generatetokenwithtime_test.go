package auth

import (
	"testing"
	"time"
)

func TestGenerateTokenWithTime(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		time    time.Time
		wantErr bool
	}{
		{
			name:    "Valid Input",
			id:      1234,
			time:    time.Now(),
			wantErr: false,
		},
		{
			name:    "Zero User ID",
			id:      0,
			time:    time.Now(),
			wantErr: false, // Adjust based on expected behavior
		},
		{
			name:    "Future Time",
			id:      5678,
			time:    time.Now().Add(24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Past Time",
			id:      9012,
			time:    time.Now().Add(-24 * time.Hour),
			wantErr: false, // Adjust based on expected behavior
		},
		{
			name:    "Maximum Uint Value",
			id:      ^uint(0),
			time:    time.Now(),
			wantErr: false,
		},
		{
			name:    "Zero Time",
			id:      3456,
			time:    time.Time{},
			wantErr: false, // Adjust based on expected behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTokenWithTime(tt.id, tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokenWithTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("GenerateTokenWithTime() returned empty token")
			}
		})
	}
}

// TestGenerateMultipleTokens tests that multiple tokens for the same user are unique
func TestGenerateMultipleTokens(t *testing.T) {
	userID := uint(7890)
	time1 := time.Now()
	time2 := time1.Add(1 * time.Hour)

	token1, err1 := GenerateTokenWithTime(userID, time1)
	if err1 != nil {
		t.Fatalf("Failed to generate first token: %v", err1)
	}

	token2, err2 := GenerateTokenWithTime(userID, time2)
	if err2 != nil {
		t.Fatalf("Failed to generate second token: %v", err2)
	}

	if token1 == token2 {
		t.Errorf("Generated tokens are not unique")
	}
}

// TODO: Implement mock for generateToken function if needed
// var generateTokenMock func(id uint, t time.Time) (string, error)
// func generateToken(id uint, t time.Time) (string, error) {
//     if generateTokenMock != nil {
//         return generateTokenMock(id, t)
//     }
//     // Original implementation
// }
