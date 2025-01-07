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

FUNCTION_DEF=func New(l *zerolog.Logger, us *store.UserStore, as *store.ArticleStore) *Handler 

 */
func TestNew(t *testing.T) {

	logger := zerolog.New(nil)

	userStore := &store.UserStore{}
	articleStore := &store.ArticleStore{}

	tests := []struct {
		name           string
		logger         *zerolog.Logger
		userStore      *store.UserStore
		articleStore   *store.ArticleStore
		expectedNil    bool
		validateFields bool
	}{
		{
			name:           "Successfully Create New Handler with Valid Parameters",
			logger:         &logger,
			userStore:      userStore,
			articleStore:   articleStore,
			expectedNil:    false,
			validateFields: true,
		},
		{
			name:           "Create Handler with Nil Logger",
			logger:         nil,
			userStore:      userStore,
			articleStore:   articleStore,
			expectedNil:    false,
			validateFields: true,
		},
		{
			name:           "Create Handler with Nil UserStore",
			logger:         &logger,
			userStore:      nil,
			articleStore:   articleStore,
			expectedNil:    false,
			validateFields: true,
		},
		{
			name:           "Create Handler with Nil ArticleStore",
			logger:         &logger,
			userStore:      userStore,
			articleStore:   nil,
			expectedNil:    false,
			validateFields: true,
		},
		{
			name:           "Create Handler with All Nil Parameters",
			logger:         nil,
			userStore:      nil,
			articleStore:   nil,
			expectedNil:    false,
			validateFields: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing:", tt.name)

			handler := New(tt.logger, tt.userStore, tt.articleStore)

			if tt.expectedNil {
				assert.Nil(t, handler, "Handler should be nil")
				return
			}

			assert.NotNil(t, handler, "Handler should not be nil")

			if tt.validateFields {

				assert.Equal(t, tt.logger, handler.logger, "Logger field mismatch")
				assert.Equal(t, tt.userStore, handler.us, "UserStore field mismatch")
				assert.Equal(t, tt.articleStore, handler.as, "ArticleStore field mismatch")

				t.Log("Successfully validated all fields")
			}
		})
	}
}

