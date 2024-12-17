package handler

import (
	"testing"
	"github.com/rs/zerolog"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982


 */
func TestNew(t *testing.T) {

	testCases := []struct {
		name         string
		logger       zerolog.Logger
		userStore    store.UserStore
		articleStore store.ArticleStore
		expectNil    bool
	}{
		{
			name:         "Basic Initialization",
			logger:       zerolog.New(nil),
			userStore:    mockUserStore{},
			articleStore: mockArticleStore{},
			expectNil:    false,
		},
		{
			name:         "Nil Logger",
			logger:       zerolog.Logger{},
			userStore:    mockUserStore{},
			articleStore: mockArticleStore{},
			expectNil:    false,
		},
		{
			name:         "Empty UserStore",
			logger:       zerolog.New(nil),
			userStore:    store.UserStore{},
			articleStore: mockArticleStore{},
			expectNil:    false,
		},
		{
			name:         "Empty ArticleStore",
			logger:       zerolog.New(nil),
			userStore:    mockUserStore{},
			articleStore: store.ArticleStore{},
			expectNil:    false,
		},
		{
			name:         "All Dependencies Nil",
			logger:       zerolog.Logger{},
			userStore:    store.UserStore{},
			articleStore: store.ArticleStore{},
			expectNil:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := New(&tc.logger, &tc.userStore, &tc.articleStore)
			if tc.expectNil {
				assert.Nil(t, handler, "Expected handler to be nil but it wasn't.")
			} else {
				assert.NotNil(t, handler, "Expected handler to be non-nil but it was nil.")
			}
			t.Logf("Test %s: Handler creation returned %v handler.", tc.name, handler)
		})
	}
}

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	dialector := postgres.New(postgres.Config{
		DriverName: "postgres",
		DSN:        "sqlmock_db_0",
		Conn:       db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

