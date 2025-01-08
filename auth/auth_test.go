package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)








/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d

FUNCTION_DEF=func GetUserID(ctx context.Context) (uint, error) 

 */
func TestGetUserID(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedID    uint
		expectedError string
	}{
		{
			name: "Valid Token with Correct User ID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    123,
			expectedError: "",
		},
		{
			name: "Invalid Token Format",
			setupContext: func() context.Context {
				md := metadata.Pairs("authorization", "Token not-a-real-token")
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: it's not even a token",
		},
		{
			name: "Expired Token",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "token expired",
		},
		{
			name: "Missing Authorization in Context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    0,
			expectedError: "Request unauthenticated with Token",
		},
		{
			name: "Token with Invalid Signature",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString([]byte("wrong-secret"))
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Non-Numeric User ID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id":   "not-a-number",
					"ExpiresAt": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: cannot map token to claims",
		},
		{
			name: "Valid Token with Zero User ID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 0,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "",
		},
		{
			name: "Token with Future Not Before Claim",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						NotBefore: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			gotID, err := GetUserID(ctx)

			if gotID != tt.expectedID {
				t.Errorf("GetUserID() gotID = %v, want %v", gotID, tt.expectedID)
			}

			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.expectedError)
			}

			if err != nil && tt.expectedError != "" {
				if !errors.Is(err, errors.New(tt.expectedError)) && err.Error() != tt.expectedError {
					t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.expectedError)
				}
			}
		})
	}
}

