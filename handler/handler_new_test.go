package handler

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)






type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		logger       *zerolog.Logger
		userStore    *store.UserStore
		articleStore *store.ArticleStore
		want         *Handler
	}{
		{
			name:         "Valid Inputs for New Function",
			logger:       &zerolog.Logger{},
			userStore:    &store.UserStore{},
			articleStore: &store.ArticleStore{},
			want:         &Handler{logger: &zerolog.Logger{}, us: &store.UserStore{}, as: &store.ArticleStore{}},
		},
		{
			name:         "Nil Inputs for New Function",
			logger:       nil,
			userStore:    nil,
			articleStore: nil,
			want:         &Handler{logger: nil, us: nil, as: nil},
		},
		{
			name:         "Partial Nil Inputs for New Function",
			logger:       &zerolog.Logger{},
			userStore:    nil,
			articleStore: &store.ArticleStore{},
			want:         &Handler{logger: &zerolog.Logger{}, us: nil, as: &store.ArticleStore{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Scenario: ", tt.name)

			got := New(tt.logger, tt.userStore, tt.articleStore)

			assert.Equal(t, tt.want, got, "they should be equal")

			if tt.logger != nil {
				assert.Equal(t, tt.logger, got.logger, "Logger should be equal to input Logger")
			} else {
				assert.Nil(t, got.logger, "Logger should be nil")
			}

			if tt.userStore != nil {
				assert.Equal(t, tt.userStore, got.us, "UserStore should be equal to input UserStore")
			} else {
				assert.Nil(t, got.us, "UserStore should be nil")
			}

			if tt.articleStore != nil {
				assert.Equal(t, tt.articleStore, got.as, "ArticleStore should be equal to input ArticleStore")
			} else {
				assert.Nil(t, got.as, "ArticleStore should be nil")
			}
		})
	}
}
