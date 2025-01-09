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

	type mockStruct struct {
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
		l1, l2 := &zerolog.Logger{}, &zerolog.Logger{}
		us1, us2 := &store.UserStore{}, &store.UserStore{}
		as1, as2 := &store.ArticleStore{}, &store.ArticleStore{}

		h1 := New(l1, us1, as1)
		h2 := New(l2, us2, as2)

		if h1 == h2 {
			t.Error("Expected different Handler instances, got the same")
		}

		if !reflect.DeepEqual(h1.logger, l1) || !reflect.DeepEqual(h1.us, us1) || !reflect.DeepEqual(h1.as, as1) {
			t.Error("Handler 1 fields not set correctly")
		}

		if !reflect.DeepEqual(h2.logger, l2) || !reflect.DeepEqual(h2.us, us2) || !reflect.DeepEqual(h2.as, as2) {
			t.Error("Handler 2 fields not set correctly")
		}
	})
}

