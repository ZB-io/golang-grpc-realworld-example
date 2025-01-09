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
		name string
		l    *zerolog.Logger
		us   *store.UserStore
		as   *store.ArticleStore
		want *Handler
	}{
		{
			name: "Create Handler with Valid Inputs",
			l:    mockLogger,
			us:   mockUserStore,
			as:   mockArticleStore,
			want: &Handler{
				logger: mockLogger,
				us:     mockUserStore,
				as:     mockArticleStore,
			},
		},
		{
			name: "Create Handler with Nil Logger",
			l:    nil,
			us:   mockUserStore,
			as:   mockArticleStore,
			want: &Handler{
				logger: nil,
				us:     mockUserStore,
				as:     mockArticleStore,
			},
		},
		{
			name: "Create Handler with Nil UserStore",
			l:    mockLogger,
			us:   nil,
			as:   mockArticleStore,
			want: &Handler{
				logger: mockLogger,
				us:     nil,
				as:     mockArticleStore,
			},
		},
		{
			name: "Create Handler with Nil ArticleStore",
			l:    mockLogger,
			us:   mockUserStore,
			as:   nil,
			want: &Handler{
				logger: mockLogger,
				us:     mockUserStore,
				as:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.l, tt.us, tt.as)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMultipleHandlers(t *testing.T) {
	mockLogger1 := &zerolog.Logger{}
	mockUserStore1 := &store.UserStore{}
	mockArticleStore1 := &store.ArticleStore{}

	mockLogger2 := &zerolog.Logger{}
	mockUserStore2 := &store.UserStore{}
	mockArticleStore2 := &store.ArticleStore{}

	handler1 := New(mockLogger1, mockUserStore1, mockArticleStore1)
	handler2 := New(mockLogger2, mockUserStore2, mockArticleStore2)

	if reflect.DeepEqual(handler1, handler2) {
		t.Errorf("Handlers should be independent, but they are equal")
	}

	if handler1.logger == handler2.logger {
		t.Errorf("Handlers should have different loggers")
	}

	if handler1.us == handler2.us {
		t.Errorf("Handlers should have different UserStores")
	}

	if handler1.as == handler2.as {
		t.Errorf("Handlers should have different ArticleStores")
	}
}

