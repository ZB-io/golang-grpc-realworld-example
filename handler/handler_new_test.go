package handler

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)


func TestNew(t *testing.T) {

	tests := []struct {
		name         string
		logger       *zerolog.Logger
		userStore    *store.UserStore
		articleStore *store.ArticleStore
	}{
		{
			name:         "Valid Inputs for New Function",
			logger:       &zerolog.Logger{},
			userStore:    &store.UserStore{},
			articleStore: &store.ArticleStore{},
		},
		{
			name:         "Nil Inputs for New Function",
			logger:       nil,
			userStore:    nil,
			articleStore: nil,
		},
		{
			name:         "Partial Nil Inputs for New Function",
			logger:       nil,
			userStore:    &store.UserStore{},
			articleStore: &store.ArticleStore{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.New(tt.logger, tt.userStore, tt.articleStore)

			assert.Equal(t, tt.logger, h.GetLogger())
			assert.Equal(t, tt.userStore, h.GetUserStore())
			assert.Equal(t, tt.articleStore, h.GetArticleStore())
			if tt.logger == nil {
				t.Log("Logger is nil")
			}
			if tt.userStore == nil {
				t.Log("UserStore is nil")
			}
			if tt.articleStore == nil {
				t.Log("ArticleStore is nil")
			}
		})
	}
}
