package handler

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext
}
/*
ROOST_METHOD_HASH=New_5541bf24ba
ROOST_METHOD_SIG_HASH=New_7d9b4d5982


 */
func TestNew(t *testing.T) {
	logger := zerolog.Nop()
	userStore := store.UserStore{}
	articleStore := store.ArticleStore{}

	h := handler.New(&logger, &userStore, &articleStore)

	if h.Logger != &logger {
		t.Errorf("Logger does not match")
	}

	if h.Us != &userStore {
		t.Errorf("UserStore does not match")
	}

	if h.As != &articleStore {
		t.Errorf("ArticleStore does not match")
	}
}

