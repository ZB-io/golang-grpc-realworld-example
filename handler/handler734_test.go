package handler_test

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982


 */
func TestNew(t *testing.T) {

	var tests = []struct {
		name          string
		expectHandler bool
		logger        *zerolog.Logger
		userStore     *store.UserStore
		articleStore  *store.ArticleStore
	}{
		{"Validate Successful Creation of Handler", true, mockLogger, mockUserStore, mockArticleStore},
		{"Validate actions when Null Logger Input is given", false, nil, mockUserStore, mockArticleStore},
		{"Validate actions when Null UserStore Input is given", false, mockLogger, nil, mockArticleStore},
		{"Validate actions when Null ArticleStore Input is given", false, mockLogger, mockUserStore, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			h := handler.New(tt.logger, tt.userStore, tt.articleStore)

			if tt.expectHandler {
				assert.NotNil(t, h, "Handler should not be nil")
			} else {
				assert.Nil(t, h, "Handler should be nil")
			}
		})
	}
}

