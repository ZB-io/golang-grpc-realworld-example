package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
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
func TestGetUserID(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedUserID uint
		expectedError  string
	}{
		{
			name: "Successfully Retrieve User ID from Valid Token",
			setupContext: func() context.Context {
				claims := claims{
					UserID:    123,
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 123,
			expectedError:  "",
		},
		{
			name: "Handle Missing Authorization Metadata",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedUserID: 0,
			expectedError:  "Request unauthenticated with Token",
		},
		{
			name: "Handle Invalid Token Error",
			setupContext: func() context.Context {
				md := metadata.Pairs("authorization", "Token invalidtokenstring")
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: it's not even a token",
		},
		{
			name: "Handle Expired Token",
			setupContext: func() context.Context {
				claims := claims{
					UserID:    123,
					ExpiresAt: time.Now().Add(-time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "token expired",
		},
		{
			name: "Handle Token with Invalid Claims",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"invalid": "data",
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: cannot map token to claims",
		},
		{
			name: "Token with Invalid Signing Method",
			setupContext: func() context.Context {
				claims := claims{
					UserID:    123,
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "signing method HS384 is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			userID, err := GetUserID(ctx)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedUserID, userID)

			t.Logf("Test '%s': expected userID=%d, got=%d; expected error='%s', got='%v'", tt.name, tt.expectedUserID, userID, tt.expectedError, err)
		})
	}
}
