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

	type mockLogger struct {
		zerolog.Logger
	}
	type mockUserStore struct {
		store.UserStore
	}
	type mockArticleStore struct {
		store.ArticleStore
	}

	tests := []struct {
		name string
		l    *zerolog.Logger
		us   *store.UserStore
		as   *store.ArticleStore
		want *Handler
	}{
		{
			name: "Create Handler with Valid Inputs",
			l:    &zerolog.Logger{},
			us:   &store.UserStore{},
			as:   &store.ArticleStore{},
			want: &Handler{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
		},
		{
			name: "Create Handler with Nil Logger",
			l:    nil,
			us:   &store.UserStore{},
			as:   &store.ArticleStore{},
			want: &Handler{
				logger: nil,
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
		},
		{
			name: "Create Handler with Nil UserStore",
			l:    &zerolog.Logger{},
			us:   nil,
			as:   &store.ArticleStore{},
			want: &Handler{
				logger: &zerolog.Logger{},
				us:     nil,
				as:     &store.ArticleStore{},
			},
		},
		{
			name: "Create Handler with Nil ArticleStore",
			l:    &zerolog.Logger{},
			us:   &store.UserStore{},
			as:   nil,
			want: &Handler{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
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

	t.Run("Create Multiple Handlers", func(t *testing.T) {
		l1 := &zerolog.Logger{}
		us1 := &store.UserStore{}
		as1 := &store.ArticleStore{}
		h1 := New(l1, us1, as1)

		l2 := &zerolog.Logger{}
		us2 := &store.UserStore{}
		as2 := &store.ArticleStore{}
		h2 := New(l2, us2, as2)

		if h1 == h2 {
			t.Errorf("Expected different instances, got the same instance")
		}

		if !reflect.DeepEqual(h1, &Handler{logger: l1, us: us1, as: as1}) {
			t.Errorf("First handler not initialized correctly")
		}

		if !reflect.DeepEqual(h2, &Handler{logger: l2, us: us2, as: as2}) {
			t.Errorf("Second handler not initialized correctly")
		}
	})
}

