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

	logger := zerolog.New(nil)

	us := &store.UserStore{}

	as := &store.ArticleStore{}

	tests := []struct {
		name string
		l    *zerolog.Logger
		us   *store.UserStore
		as   *store.ArticleStore
		want *Handler
	}{
		{
			name: "Create Handler with Valid Inputs",
			l:    &logger,
			us:   us,
			as:   as,
			want: &Handler{logger: &logger, us: us, as: as},
		},
		{
			name: "New with Nil Logger",
			l:    nil,
			us:   us,
			as:   as,
			want: &Handler{logger: nil, us: us, as: as},
		},
		{
			name: "New with Nil UserStore",
			l:    &logger,
			us:   nil,
			as:   as,
			want: &Handler{logger: &logger, us: nil, as: as},
		},
		{
			name: "New with Nil ArticleStore",
			l:    &logger,
			us:   us,
			as:   nil,
			want: &Handler{logger: &logger, us: us, as: nil},
		},
		{
			name: "New with All Nil Inputs",
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

