package handler

import (
	debug "runtime/debug"
	testing "testing"

	store "github.com/raahii/golang-grpc-realworld-example/store"
	zerolog "github.com/rs/zerolog"
)

/*
ROOST_METHOD_HASH=New_437eff3b29
ROOST_METHOD_SIG_HASH=New_6e92a7c68a

FUNCTION_DEF=func New(l *zerolog.Logger, us *store.UserStore, as *store.ArticleStore) *Handler // New returns a new handler with logger and database
*/
func TestNew(t *testing.T) {
	type args struct {
		logger *zerolog.Logger
		us     *store.UserStore
		as     *store.ArticleStore
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "Normal Operation",
			args: args{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
			want: &Handler{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
		},
		{
			name: "Nil Logger",
			args: args{
				logger: nil,
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
			want: &Handler{
				logger: nil,
				us:     &store.UserStore{},
				as:     &store.ArticleStore{},
			},
		},
		{
			name: "Nil UserStore",
			args: args{
				logger: &zerolog.Logger{},
				us:     nil,
				as:     &store.ArticleStore{},
			},
			want: &Handler{
				logger: &zerolog.Logger{},
				us:     nil,
				as:     &store.ArticleStore{},
			},
		},
		{
			name: "Nil ArticleStore",
			args: args{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
				as:     nil,
			},
			want: &Handler{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
				as:     nil,
			},
		},
		{
			name: "All Nil Parameters",
			args: args{
				logger: nil,
				us:     nil,
				as:     nil,
			},
			want: &Handler{
				logger: nil,
				us:     nil,
				as:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()
			got := New(tt.args.logger, tt.args.us, tt.args.as)
			if got == nil {
				t.Errorf("New() got = nil, want %v", tt.want)
			}
			if got.logger != tt.want.logger {
				t.Errorf("New().logger = %v, want %v", got.logger, tt.want.logger)
			}
			if got.us != tt.want.us {
				t.Errorf("New().us = %v, want %v", got.us, tt.want.us)
			}
			if got.as != tt.want.as {
				t.Errorf("New().as = %v, want %v", got.as, tt.want.as)
			}
			t.Logf("Test %s passed", tt.name)
		})
	}
}
