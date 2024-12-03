package handler

import (
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"reflect"
	"testing"
)

/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982


 */
func TestNew(t *testing.T) {

	testLogger := zerolog.New(nil)

	type args struct {
		logger       *zerolog.Logger
		userStore    *store.UserStore
		articleStore *store.ArticleStore
	}

	tests := []struct {
		name     string
		args     args
		want     *Handler
		validate func(*testing.T, *Handler)
	}{
		{
			name: "Successfully create new handler with valid parameters",
			args: args{
				logger:       &testLogger,
				userStore:    &store.UserStore{},
				articleStore: &store.ArticleStore{},
			},
			validate: func(t *testing.T, h *Handler) {
				t.Log("Validating handler with all valid parameters")
				if h == nil {
					t.Error("Expected non-nil handler")
					return
				}
				if h.logger == nil {
					t.Error("Expected non-nil logger")
				}
				if h.us == nil {
					t.Error("Expected non-nil UserStore")
				}
				if h.as == nil {
					t.Error("Expected non-nil ArticleStore")
				}
			},
		},
		{
			name: "Create handler with nil logger",
			args: args{
				logger:       nil,
				userStore:    &store.UserStore{},
				articleStore: &store.ArticleStore{},
			},
			validate: func(t *testing.T, h *Handler) {
				t.Log("Validating handler with nil logger")
				if h == nil {
					t.Error("Expected non-nil handler even with nil logger")
					return
				}
				if h.logger != nil {
					t.Error("Expected nil logger")
				}
			},
		},
		{
			name: "Create handler with nil UserStore",
			args: args{
				logger:       &testLogger,
				userStore:    nil,
				articleStore: &store.ArticleStore{},
			},
			validate: func(t *testing.T, h *Handler) {
				t.Log("Validating handler with nil UserStore")
				if h == nil {
					t.Error("Expected non-nil handler even with nil UserStore")
					return
				}
				if h.us != nil {
					t.Error("Expected nil UserStore")
				}
			},
		},
		{
			name: "Create handler with nil ArticleStore",
			args: args{
				logger:       &testLogger,
				userStore:    &store.UserStore{},
				articleStore: nil,
			},
			validate: func(t *testing.T, h *Handler) {
				t.Log("Validating handler with nil ArticleStore")
				if h == nil {
					t.Error("Expected non-nil handler even with nil ArticleStore")
					return
				}
				if h.as != nil {
					t.Error("Expected nil ArticleStore")
				}
			},
		},
		{
			name: "Create handler with all nil parameters",
			args: args{
				logger:       nil,
				userStore:    nil,
				articleStore: nil,
			},
			validate: func(t *testing.T, h *Handler) {
				t.Log("Validating handler with all nil parameters")
				if h == nil {
					t.Error("Expected non-nil handler even with all nil parameters")
					return
				}
				if !reflect.DeepEqual(h, &Handler{logger: nil, us: nil, as: nil}) {
					t.Error("Expected handler with all nil fields")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.logger, tt.args.userStore, tt.args.articleStore)
			tt.validate(t, got)
		})
	}
}

