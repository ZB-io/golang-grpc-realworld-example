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
			want: &Handler{logger: mockLogger, us: mockUserStore, as: mockArticleStore},
		},
		{
			name: "New with Nil Logger",
			l:    nil,
			us:   mockUserStore,
			as:   mockArticleStore,
			want: &Handler{logger: nil, us: mockUserStore, as: mockArticleStore},
		},
		{
			name: "New with Nil UserStore",
			l:    mockLogger,
			us:   nil,
			as:   mockArticleStore,
			want: &Handler{logger: mockLogger, us: nil, as: mockArticleStore},
		},
		{
			name: "New with Nil ArticleStore",
			l:    mockLogger,
			us:   mockUserStore,
			as:   nil,
			want: &Handler{logger: mockLogger, us: mockUserStore, as: nil},
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

func TestNewConsistency(t *testing.T) {
	mockLogger := &zerolog.Logger{}
	mockUserStore := &store.UserStore{}
	mockArticleStore := &store.ArticleStore{}

	handler1 := New(mockLogger, mockUserStore, mockArticleStore)
	handler2 := New(mockLogger, mockUserStore, mockArticleStore)

	if !reflect.DeepEqual(handler1, handler2) {
		t.Errorf("Multiple calls to New() produced inconsistent results")
	}
}
