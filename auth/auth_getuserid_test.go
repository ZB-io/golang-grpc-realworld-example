package auth

import (
	"context"
	"errors"
	"testing"
	"time"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

var jwtSecretTest = []byte("testsecret")

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
	t.Run("Scenario 1: Successfully Retrieve User ID from Valid Token", func(t *testing.T) {
		expirationTime := time.Now().Add(1 * time.Hour).Unix()
		validClaims := &claims{
			UserID: 12345,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime,
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
		tokenString, err := token.SignedString(jwtSecretTest)
		if err != nil {
			t.Fatalf("unexpected error signing token: %v", err)
		}

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Token "+tokenString))
		userID, err := GetUserID(ctx)
		if err != nil || userID != validClaims.UserID {
			t.Fatalf("expected userID: %v, got: %v, error: %v", validClaims.UserID, userID, err)
		}
		t.Log("Successfully retrieved User ID from valid token")
	})

	t.Run("Scenario 2: Handle Missing Authorization Metadata", func(t *testing.T) {
		ctx := context.Background()
		userID, err := GetUserID(ctx)
		if userID != 0 || err == nil {
			t.Fatalf("expected userID: 0, got: %v, expected error due to missing metadata", userID)
		}
		t.Log("Properly handled missing authorization metadata")
	})

	t.Run("Scenario 3: Handle Invalid Token Error", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Token "+invalidToken))
		userID, err := GetUserID(ctx)
		if userID != 0 || err == nil || err.Error() != "invalid token: it's not even a token" {
			t.Fatalf("expected userID: 0, error about invalid token, got: %v, error: %v", userID, err)
		}
		t.Log("Detected invalid token structure correctly")
	})

	t.Run("Scenario 4: Handle Expired Token", func(t *testing.T) {
		expiredClaims := &claims{
			UserID: 12345,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		tokenString, err := token.SignedString(jwtSecretTest)
		if err != nil {
			t.Fatalf("unexpected error signing token: %v", err)
		}

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Token "+tokenString))
		userID, err := GetUserID(ctx)
		if userID != 0 || err == nil || err.Error() != "token expired" {
			t.Fatalf("expected userID: 0, error about expired token, got: %v, error: %v", userID, err)
		}
		t.Log("Correctly handled expired token")
	})

	t.Run("Scenario 5: Handle Token with Invalid Claims", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "not-a-uint",
		})
		tokenString, err := token.SignedString(jwtSecretTest)
		if err != nil {
			t.Fatalf("unexpected error signing token: %v", err)
		}

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Token "+tokenString))
		userID, err := GetUserID(ctx)
		if userID != 0 || err == nil || err.Error() != "invalid token: cannot map token to claims" {
			t.Fatalf("expected userID: 0, error about claim mapping, got: %v, error: %v", userID, err)
		}
		t.Log("Correctly handled invalid claim token")
	})

	t.Run("Scenario 6: Token with Invalid Signing Method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims{
			UserID: 12345,
		})
		tokenString, err := token.SignedString(jwtSecretTest)
		if err != nil {
			t.Fatalf("unexpected error signing token: %v", err)
		}

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Token "+tokenString))
		userID, err := GetUserID(ctx)
		if userID != 0 || err == nil {
			t.Fatalf("expected userID: 0, got: %v, expected error due to invalid signing method", userID)
		}
		t.Log("Handled invalid signing method")
	})
}
