package handler

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer mockDB.Close()

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

			if tt.logger == nil {
				assert.Nil(t, handler.logger, "Expected logger to be nil")
			} else {
				assert.Equal(t, tt.logger, handler.logger, "Logger should match expected")
			}

			if tt.userStore == nil {
				assert.Nil(t, handler.us, "Expected userStore to be nil")
			} else {
				assert.Equal(t, tt.userStore, handler.us, "UserStore should match expected")
			}

			if tt.articleStore == nil {
				assert.Nil(t, handler.as, "Expected articleStore to be nil")
			} else {
				assert.Equal(t, tt.articleStore, handler.as, "ArticleStore should match expected")
			}

			t.Logf("Successfully created handler for scenario: %s", tt.name)
		})
	}
}
