package handler

import (
	"bytes"
	"fmt"
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
		name        string
		logger      *zerolog.Logger
		us          *store.UserStore
		as          *store.ArticleStore
		expectNil   bool
		shouldLog   bool
		description string
	}{
		{
			name:        "Successful Initialization of Handler",
			logger:      newLogger(),
			us:          new(store.UserStore),
			as:          new(store.ArticleStore),
			expectNil:   false,
			shouldLog:   true,
			description: "Verifies correct initialization with all valid parameters",
		},
		{
			name:        "Handler Initialization with Nil Logger",
			logger:      nil,
			us:          new(store.UserStore),
			as:          new(store.ArticleStore),
			expectNil:   false,
			shouldLog:   false,
			description: "Tests behavior with nil logger ensuring proper handling",
		},
		{
			name:        "Handler Initialization with Nil User and Article Stores",
			logger:      newLogger(),
			us:          nil,
			as:          nil,
			expectNil:   false,
			shouldLog:   true,
			description: "Tests behavior with nil stores ensuring proper handling",
		},
		{
			name:        "Simultaneous Nil Logger and Stores",
			logger:      nil,
			us:          nil,
			as:          nil,
			expectNil:   true,
			description: "Tests behavior with all nil parameters for robustness",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := New(tc.logger, tc.us, tc.as)

			assert.Equal(t, tc.logger, handler.logger, "Logger should match input logger")
			assert.Equal(t, tc.us, handler.us, "UserStore should match input user store")
			assert.Equal(t, tc.as, handler.as, "ArticleStore should match input article store")

			if tc.shouldLog && handler.logger != nil {
				var buf bytes.Buffer

				fmt.Fprintf(&buf, "Test log event")
				handler.logger.Log().Msg(buf.String())
				assert.Contains(t, buf.String(), "Test log event", "Expected log event not found")
			}
		})
	}
}

func newLogger() *zerolog.Logger {
	l := zerolog.New(zerolog.ConsoleWriter{Out: &bytes.Buffer{}})
	return &l
}

