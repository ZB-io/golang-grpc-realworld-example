package handler

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type Handler struct {
	logger *zerolog.Logger
	us     *store.UserStore
	as     *store.ArticleStore
}

/* Rest of struct declarations... */

func TestNew(t *testing.T) {
	logger := zerolog.New(log.Logger)
	userStore := &store.UserStore{}
	articleStore := &store.ArticleStore{}

	tests := []struct {
		name string
		l    *zerolog.Logger
		us   *store.UserStore
		as   *store.ArticleStore
		want *Handler
	}{
		//...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Scenario:", tt.name)
			h := New(tt.l, tt.us, tt.as)
			//...
		})
	}
}
