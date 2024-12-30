package handler

import (
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"reflect"
	"sync"
	"testing"
)




type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestNew(t *testing.T) {
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
			name: "Handle Nil Logger Input",
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
			name: "Handle Nil UserStore Input",
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
			name: "Handle Nil ArticleStore Input",
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

	t.Run("Concurrent Handler Creation", func(t *testing.T) {
		var wg sync.WaitGroup
		handlerChan := make(chan *Handler, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				l := &zerolog.Logger{}
				us := &store.UserStore{}
				as := &store.ArticleStore{}
				h := New(l, us, as)
				handlerChan <- h
			}()
		}

		wg.Wait()
		close(handlerChan)

		handlers := make([]*Handler, 0, 10)
		for h := range handlerChan {
			handlers = append(handlers, h)
		}

		if len(handlers) != 10 {
			t.Errorf("Expected 10 handlers, got %d", len(handlers))
		}

		for i := 0; i < len(handlers); i++ {
			if handlers[i] == nil {
				t.Errorf("Handler at index %d is nil", i)
			}
			for j := i + 1; j < len(handlers); j++ {
				if handlers[i] == handlers[j] {
					t.Errorf("Handlers at index %d and %d are the same instance", i, j)
				}
			}
		}
	})
}
