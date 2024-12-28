package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
)





func TestGetUserID(t *testing.T) {
	type testCase struct {
		desc           string
		token          string
		expectedUserID uint
		expectedError  error
		setupToken     func() string
	}

	validUserID := uint(123)

	testCases := []testCase{
		{
			desc: "Successfully Retrieve User ID from Valid Token",
			setupToken: func() string {
				claims := &claims{
					UserID: validUserID,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			expectedUserID: validUserID,
			expectedError:  nil,
		},
		{
			desc: "Handle Missing Authorization Metadata",
			setupToken: func() string {
				return ""
			},
			expectedUserID: 0,
			expectedError:  grpc_auth.ErrUnauthenticated,
		},
		{
			desc: "Handle Invalid Token Error",
			setupToken: func() string {
				return "invalid.token.string"
			},
			expectedUserID: 0,
			expectedError:  errors.New("invalid token: it's not even a token"),
		},
		{
			desc: "Handle Expired Token",
			setupToken: func() string {
				claims := &claims{
					UserID: validUserID,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			expectedUserID: 0,
			expectedError:  errors.New("token expired"),
		},
		{
			desc: "Handle Token with Invalid Claims",
			setupToken: func() string {

				type invalidClaims struct {
					Name string `json:"name"`
					jwt.StandardClaims
				}
				claims := &invalidClaims{
					Name: "UnauthorizedAccess",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			expectedUserID: 0,
			expectedError:  errors.New("invalid token: cannot map token to claims"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenString := tc.setupToken()
			if tokenString != "" {
				md := map[string]string{
					"authorization": "Token " + tokenString,
				}
				ctx := context.WithValue(context.Background(), grpc_auth.AuthHeaderPayloadKey{}, md)
				userID, err := GetUserID(ctx)

				assert.Equal(t, tc.expectedUserID, userID)
				if err != nil {
					assert.EqualError(t, err, tc.expectedError.Error())
				} else {
					assert.Nil(t, tc.expectedError)
				}
			} else {
				ctx := context.Background()
				userID, err := GetUserID(ctx)
				assert.Equal(t, tc.expectedUserID, userID)
				if err != nil {
					assert.EqualError(t, err, tc.expectedError.Error())
				} else {
					assert.Nil(t, tc.expectedError)
				}
			}
		})
	}
}




func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "defaultSecretKey"
	}
	return []byte(secret)
}


