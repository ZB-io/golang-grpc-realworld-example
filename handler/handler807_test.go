package handler

import (
	"reflect"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
)








/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982

FUNCTION_DEF=func New(l *zerolog.Logger, us *store.UserStore, as *store.ArticleStore) *Handler 

 */
func TestNew(t *testing.T) {

	mockLogger := &zerolog.Logger{}
	mockUserStore := &store.UserStore{}
	mockArticleStore := &store.ArticleStore{}

	tests := []struct {
		name     string
		logger   *zerolog.Logger
		us       *store.UserStore
		as       *store.ArticleStore
		expected *Handler
	}{
		{
			name:   "Create a new Handler with valid inputs",
			logger: mockLogger,
			us:     mockUserStore,
			as:     mockArticleStore,
			expected: &Handler{
				logger: mockLogger,
				us:     mockUserStore,
				as:     mockArticleStore,
			},
		},
		{
			name:   "Create Handler with nil logger",
			logger: nil,
			us:     mockUserStore,
			as:     mockArticleStore,
			expected: &Handler{
				logger: nil,
				us:     mockUserStore,
				as:     mockArticleStore,
			},
		},
		{
			name:   "Create Handler with nil UserStore",
			logger: mockLogger,
			us:     nil,
			as:     mockArticleStore,
			expected: &Handler{
				logger: mockLogger,
				us:     nil,
				as:     mockArticleStore,
			},
		},
		{
			name:   "Create Handler with nil ArticleStore",
			logger: mockLogger,
			us:     mockUserStore,
			as:     nil,
			expected: &Handler{
				logger: mockLogger,
				us:     mockUserStore,
				as:     nil,
			},
		},
		{
			name:   "Create Handler with all nil inputs",
			logger: nil,
			us:     nil,
			as:     nil,
			expected: &Handler{
				logger: nil,
				us:     nil,
				as:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.logger, tt.us, tt.as)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("New() = %v, want %v", got, tt.expected)
			}

			if got.logger != tt.logger {
				t.Errorf("New().logger = %v, want %v", got.logger, tt.logger)
			}
			if got.us != tt.us {
				t.Errorf("New().us = %v, want %v", got.us, tt.us)
			}
			if got.as != tt.as {
				t.Errorf("New().as = %v, want %v", got.as, tt.as)
			}
		})
	}
}

