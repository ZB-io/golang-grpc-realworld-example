package handler

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
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
func TestNew(t *testing.T) {

	_, userDbMock, _ := sqlmock.New()
	_, articleDbMock, _ := sqlmock.New()

	validUserStore := &store.UserStore{db: userDbMock}
	validArticleStore := &store.ArticleStore{db: articleDbMock}

	validLogger := zerolog.New(nil).Level(zerolog.InfoLevel)

	tests := []struct {
		name               string
		logger             *zerolog.Logger
		userStore          *store.UserStore
		articleStore       *store.ArticleStore
		expectNilHandler   bool
		expectLogger       bool
		expectUserStore    bool
		expectArticleStore bool
	}{
		{
			name:               "Valid Logger, UserStore, and ArticleStore",
			logger:             &validLogger,
			userStore:          validUserStore,
			articleStore:       validArticleStore,
			expectNilHandler:   false,
			expectLogger:       true,
			expectUserStore:    true,
			expectArticleStore: true,
		},
		{
			name:               "Null Logger",
			logger:             nil,
			userStore:          validUserStore,
			articleStore:       validArticleStore,
			expectNilHandler:   false,
			expectLogger:       false,
			expectUserStore:    true,
			expectArticleStore: true,
		},
		{
			name:               "Null UserStore",
			logger:             &validLogger,
			userStore:          nil,
			articleStore:       validArticleStore,
			expectNilHandler:   false,
			expectLogger:       true,
			expectUserStore:    false,
			expectArticleStore: true,
		},
		{
			name:               "Null ArticleStore",
			logger:             &validLogger,
			userStore:          validUserStore,
			articleStore:       nil,
			expectNilHandler:   false,
			expectLogger:       true,
			expectUserStore:    true,
			expectArticleStore: false,
		},
		{
			name:               "All Null Arguments",
			logger:             nil,
			userStore:          nil,
			articleStore:       nil,
			expectNilHandler:   false,
			expectLogger:       false,
			expectUserStore:    false,
			expectArticleStore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.logger, tt.userStore, tt.articleStore)

			if tt.expectNilHandler {
				assert.Nil(t, handler, "Expected handler to be nil")
			} else {
				assert.NotNil(t, handler, "Expected handler to be non-nil")
			}

			if tt.expectLogger {
				assert.NotNil(t, handler.logger, "Expected logger to be non-nil")
			} else {
				assert.Nil(t, handler.logger, "Expected logger to be nil")
			}

			if tt.expectUserStore {
				assert.NotNil(t, handler.us, "Expected UserStore to be non-nil")
			} else {
				assert.Nil(t, handler.us, "Expected UserStore to be nil")
			}

			if tt.expectArticleStore {
				assert.NotNil(t, handler.as, "Expected ArticleStore to be non-nil")
			} else {
				assert.Nil(t, handler.as, "Expected ArticleStore to be nil")
			}
		})
	}
}
