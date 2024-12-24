package handler

import (
	"testing"
	"os"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type Handler struct {
	logger *zerolog.Logger
	us     *store.UserStore
	as     *store.ArticleStore
}

type ArticleStore struct {
	db *gorm.DB
}

type UserStore struct {
	db *gorm.DB
}

type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext
}T struct {
	common
	isEnvSet bool
	context  *testContext
}
/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982


 */
func NewTestArticleStore() *store.ArticleStore {
	return &store.ArticleStore{}
}

func NewTestLogger() *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	return &logger
}

func NewTestUserStore() *store.UserStore {
	return &store.UserStore{}
}

func TestNew(t *testing.T) {
	tt := []struct {
		name           string
		loggerIn       *zerolog.Logger
		userStoreIn    *store.UserStore
		articleStoreIn *store.ArticleStore
	}{
		{
			name:           "Successful creation of Handler",
			loggerIn:       NewTestLogger(),
			userStoreIn:    NewTestUserStore(),
			articleStoreIn: NewTestArticleStore(),
		},
		{
			name:           "New Function is provided a nil Logger",
			loggerIn:       nil,
			userStoreIn:    NewTestUserStore(),
			articleStoreIn: NewTestArticleStore(),
		},
		{
			name:           "New Function is provided a nil UserStore",
			loggerIn:       NewTestLogger(),
			userStoreIn:    nil,
			articleStoreIn: NewTestArticleStore(),
		},
		{
			name:           "New Function is provided a nil ArticleStore",
			loggerIn:       NewTestLogger(),
			userStoreIn:    NewTestUserStore(),
			articleStoreIn: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := New(tc.loggerIn, tc.userStoreIn, tc.articleStoreIn)

			assert.Equal(t, tc.loggerIn, h.logger, "They should be equal")
			assert.Equal(t, tc.userStoreIn, h.us, "They should be equal")
			assert.Equal(t, tc.articleStoreIn, h.as, "They should be equal")
		})
	}
}

