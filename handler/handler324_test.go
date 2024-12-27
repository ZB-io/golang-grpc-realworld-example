package handler

import (
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982
*/
func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		logger        *zerolog.Logger
		userStore     *store.UserStore
		articleStore  *store.ArticleStore
		expectHandler *Handler
	}{
		{
			name:         "Valid Initialization of Handler",
			logger:       zerolog.New(nil),
			userStore:    &store.UserStore{},
			articleStore: &store.ArticleStore{},
			expectHandler: &Handler{
				logger: zerolog.New(nil),
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
		},
		{
			name:         "Handling Nil Logger",
			logger:       nil,
			userStore:    &store.UserStore{},
			articleStore: &store.ArticleStore{},
			expectHandler: &Handler{
				logger: nil,
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
		},
		{
			name:         "Handling Nil UserStore",
			logger:       zerolog.New(nil),
			userStore:    nil,
			articleStore: &store.ArticleStore{},
			expectHandler: &Handler{
				logger: zerolog.New(nil),
				us:     nil,
				as:     &store.ArticleStore{},
			},
		},
		{
			name:         "Handling Nil ArticleStore",
			logger:       zerolog.New(nil),
			userStore:    &store.UserStore{},
			articleStore: nil,
			expectHandler: &Handler{
				logger: zerolog.New(nil),
				us:     &store.UserStore{},
				as:     nil,
			},
		},
		{
			name:         "Complete Nil Inputs",
			logger:       nil,
			userStore:    nil,
			articleStore: nil,
			expectHandler: &Handler{
				logger: nil,
				us:     nil,
				as:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.logger, tt.userStore, tt.articleStore)
			assert.NotNil(t, handler, "Expected handler to be non-nil")

			assert.Equal(t, tt.expectHandler.logger, handler.logger, "Logger should match expected")
			assert.Equal(t, tt.expectHandler.us, handler.us, "UserStore should match expected")
			assert.Equal(t, tt.expectHandler.as, handler.as, "ArticleStore should match expected")

			t.Logf("Successfully created handler for scenario: %s", tt.name)
		})
	}
}
