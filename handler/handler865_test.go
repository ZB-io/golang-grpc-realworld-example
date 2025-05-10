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
			name: "Create Handler with valid inputs",
			l:    mockLogger,
			us:   mockUserStore,
			as:   mockArticleStore,
			want: &Handler{logger: mockLogger, us: mockUserStore, as: mockArticleStore},
		},
		{
			name: "Create Handler with nil logger",
			l:    nil,
			us:   mockUserStore,
			as:   mockArticleStore,
			want: &Handler{logger: nil, us: mockUserStore, as: mockArticleStore},
		},
		{
			name: "Create Handler with nil UserStore",
			l:    mockLogger,
			us:   nil,
			as:   mockArticleStore,
			want: &Handler{logger: mockLogger, us: nil, as: mockArticleStore},
		},
		{
			name: "Create Handler with nil ArticleStore",
			l:    mockLogger,
			us:   mockUserStore,
			as:   nil,
			want: &Handler{logger: mockLogger, us: mockUserStore, as: nil},
		},
		{
			name: "Create Handler with all nil inputs",
			l:    nil,
			us:   nil,
			as:   nil,
			want: &Handler{logger: nil, us: nil, as: nil},
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

