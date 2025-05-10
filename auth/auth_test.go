package github

import (
	"os"
	"testing"
	"github.com/dgrijalva/jwt-go"
	"time"
)









/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3

FUNCTION_DEF=func GenerateToken(id uint) (string, error) 

 */
func TestConcurrentTokenGeneration(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	userIDs := []uint{1, 2, 3, 4, 5}
	results := make(chan struct {
		token string
		err   error
	}, len(userIDs))

	for _, id := range userIDs {
		go func(id uint) {
			token, err := GenerateToken(id)
			results <- struct {
				token string
				err   error
			}{token, err}
		}(id)
	}

	for i := 0; i < len(userIDs); i++ {
		result := <-results
		if result.err != nil {
			t.Errorf("Unexpected error in concurrent generation: %v", result.err)
		}
		if result.token == "" {
			t.Errorf("Expected a non-empty token in concurrent generation, but got an empty string")
		}
	}
}

func TestGenerateToken(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tests := []struct {
		name        string
		userID      uint
		setupEnv    func()
		expectError bool
	}{
		{
			name:        "Successful Token Generation",
			userID:      123,
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectError: false,
		},
		{
			name:        "Token Generation with Zero User ID",
			userID:      0,
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectError: true,
		},
		{
			name:        "Token Generation with Maximum uint Value",
			userID:      ^uint(0),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectError: false,
		},
		{
			name:        "Token Generation with Missing JWT Secret",
			userID:      123,
			setupEnv:    func() { os.Unsetenv("JWT_SECRET") },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			token, err := GenerateToken(tt.userID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("Expected a non-empty token, but got an empty string")
				}

				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}

				if claims.UserID != tt.userID {
					t.Errorf("Expected UserID %d, but got %d", tt.userID, claims.UserID)
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

 */
func TestGenerateTokenWithTime(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test_secret")

	tests := []struct {
		name        string
		id          uint
		time        time.Time
		expectError bool
	}{
		{
			name:        "Successful Token Generation",
			id:          1,
			time:        time.Now(),
			expectError: false,
		},
		{
			name:        "Token Generation with Zero User ID",
			id:          0,
			time:        time.Now(),
			expectError: true,
		},
		{
			name:        "Token Generation with Far Future Time",
			id:          1,
			time:        time.Now().AddDate(100, 0, 0),
			expectError: false,
		},
		{
			name:        "Token Generation with Past Time",
			id:          1,
			time:        time.Now().AddDate(0, 0, -1),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateTokenWithTime(tt.id, tt.time)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("Expected a non-empty token, but got an empty string")
				}

				claims := &claims{}
				parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}

				if !parsedToken.Valid {
					t.Errorf("Token is not valid")
				}

				if claims.UserID != tt.id {
					t.Errorf("Expected UserID %d, but got %d", tt.id, claims.UserID)
				}

				expectedExp := tt.time.Add(time.Hour * 24).Unix()
				if claims.ExpiresAt != expectedExp {
					t.Errorf("Expected expiration time %v, but got %v", expectedExp, claims.ExpiresAt)
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8

FUNCTION_DEF=func generateToken(id uint, now time.Time) (string, error) 

 */
func TestGenerateToken(t *testing.T) {

	jwtSecret = []byte("test_secret")

	tests := []struct {
		name    string
		id      uint
		now     time.Time
		wantErr bool
	}{
		{
			name:    "Successfully Generate Token for Valid User ID",
			id:      1,
			now:     time.Now(),
			wantErr: false,
		},
		{
			name:    "Handle Zero User ID",
			id:      0,
			now:     time.Now(),
			wantErr: false,
		},
		{
			name:    "Token Generation with Maximum uint Value",
			id:      ^uint(0),
			now:     time.Now(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateToken(tt.id, tt.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("generateToken() returned empty token")
			}

			if !tt.wantErr {
				token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrSignatureInvalid
					}
					return jwtSecret, nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if uint(claims["user_id"].(float64)) != tt.id {
						t.Errorf("Token user_id = %v, want %v", claims["user_id"], tt.id)
					}
					if int64(claims["exp"].(float64)) != tt.now.Add(time.Hour*72).Unix() {
						t.Errorf("Token expiration time is incorrect")
					}
				} else {
					t.Errorf("Token claims are invalid")
				}
			}
		})
	}
}

func TestGenerateTokenConsistency(t *testing.T) {
	jwtSecret = []byte("test_secret")
	id := uint(1)
	now := time.Now()

	token1, err := generateToken(id, now)
	if err != nil {
		t.Fatalf("Failed to generate first token: %v", err)
	}

	token2, err := generateToken(id, now)
	if err != nil {
		t.Fatalf("Failed to generate second token: %v", err)
	}

	if token1 != token2 {
		t.Errorf("Tokens are not consistent for the same input")
	}
}

func TestGenerateTokenInvalidSecret(t *testing.T) {

	originalSecret := jwtSecret
	jwtSecret = []byte{}
	defer func() { jwtSecret = originalSecret }()

	_, err := generateToken(1, time.Now())
	if err == nil {
		t.Errorf("Expected error with invalid JWT secret, got nil")
	}
}

func TestGenerateTokenPerformance(t *testing.T) {
	jwtSecret = []byte("test_secret")
	now := time.Now()
	iterations := 1000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := generateToken(uint(i), now)
		if err != nil {
			t.Fatalf("Failed to generate token on iteration %d: %v", i, err)
		}
	}
	duration := time.Since(start)

	t.Logf("Generated %d tokens in %v", iterations, duration)
	if duration > time.Second*5 {
		t.Errorf("Token generation took too long: %v", duration)
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	jwtSecret = []byte("test_secret")
	now := time.Now()

	token1, err := generateToken(1, now)
	if err != nil {
		t.Fatalf("Failed to generate first token: %v", err)
	}

	token2, err := generateToken(2, now)
	if err != nil {
		t.Fatalf("Failed to generate second token: %v", err)
	}

	if token1 == token2 {
		t.Errorf("Tokens for different user IDs are not unique")
	}
}

