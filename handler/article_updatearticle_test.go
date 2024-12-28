package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestHandlerUpdateArticle(t *testing.T) {
	t.Parallel()

	t.Run("Scenario 1: Successful Article Update by Author", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()

		// Mock database operations
		mock.ExpectQuery(`SELECT * FROM "articles" WHERE "id" = ? LIMIT`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "body", "author_id"}).
				AddRow(1, "Old Title", "Old Body", 101))

		mock.ExpectExec(`UPDATE "articles" SET`).WithArgs("New Title", "New Body", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		us := &store.UserStore{DB: db}
		as := &store.ArticleStore{DB: db}
		logger := &zerolog.Logger{}

		handler := Handler{logger: logger, us: us, as: as}

		ctx := context.TODO() // Replace with a more realistic context as needed
		req := &pb.UpdateArticleRequest{Article: &pb.UpdateArticleRequest_Article{
			Slug:        "1",
			Title:       "New Title",
			Body:        "New Body",
			Description: "New Description",
		}}

		resp, err := handler.UpdateArticle(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Article.Title != "New Title" {
			t.Errorf("expected title to be updated to 'New Title', got %s", resp.Article.Title)
		}
		t.Log("UpdateArticle by a valid author passed successfully")
	})

	t.Run("Scenario 2: Unauthenticated User Attempt", func(t *testing.T) {
		t.Parallel()

		handler := &Handler{} // mock dependencies if needed

		ctx := context.TODO()  // This context doesn't have authentication tokens
		req := &pb.UpdateArticleRequest{}

		resp, err := handler.UpdateArticle(ctx, req)
		if err == nil || status.Code(err) != codes.Unauthenticated {
			t.Fatalf("expected Unauthenticated error, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Unauthenticated user attempt handled correctly")
	})

	t.Run("Scenario 3: Attempt to Update Another User's Article", func(t *testing.T) {
		t.Parallel()

		// Create the mock and handler setup as per need
		// TODO: Add setup...

		handler := &Handler{} // mock dependencies if needed

		ctx := context.TODO() // Replace with valid context
		req := &pb.UpdateArticleRequest{Article: &pb.UpdateArticleRequest_Article{Slug: "2"}} // Assuming another user's article

		resp, err := handler.UpdateArticle(ctx, req)
		if err == nil || status.Code(err) != codes.PermissionDenied {
			t.Fatalf("expected PermissionDenied error, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Attempt to update another user's article was blocked as expected")

	})

	t.Run("Scenario 4: Invalid Slug Conversion", func(t *testing.T) {
		t.Parallel()

		// Duplicate setup if needed
		handler := &Handler{} // mock dependencies if needed

		ctx := context.TODO() // Use proper context
		req := &pb.UpdateArticleRequest{Article: &pb.UpdateArticleRequest_Article{Slug: "invalid_slug"}}

		resp, err := handler.UpdateArticle(ctx, req)
		if err == nil || status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument error, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Invalid slug conversion handled correctly")
	})

	t.Run("Scenario 5: Validation Failure on Updated Article Fields", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()

		// Mock setup
		// TODO: Mock store as needed

		us := &store.UserStore{DB: db}
		as := &store.ArticleStore{DB: db}
		logger := &zerolog.Logger{}

		handler := Handler{logger: logger, us: us, as: as}

		ctx := context.TODO() // Replace with suitable context
		req := &pb.UpdateArticleRequest{Article: &pb.UpdateArticleRequest_Article{
			Slug:  "1",
			Title: "", // Invalid: Title is empty
		}}

		resp, err := handler.UpdateArticle(ctx, req)
		if err == nil || status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument error, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Validation failure handled properly")

	})

	t.Run("Scenario 6: Database Failure during Article Update", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()

		// Setup mock expectations
		mock.ExpectQuery(`SELECT * FROM "articles" WHERE "id" = ? LIMIT`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "body", "author_id"}).
				AddRow(1, "Title", "Body", 101))

		mock.ExpectExec(`UPDATE "articles" SET`).WithArgs("New Title", "New Body", 1).
			WillReturnError(errors.New("db error"))

		us := &store.UserStore{DB: db}
		as := &store.ArticleStore{DB: db}
		logger := &zerolog.Logger{}

		handler := Handler{logger: logger, us: us, as: as}

		ctx := context.TODO() // proper context setup
		req := &pb.UpdateArticleRequest{Article: &pb.UpdateArticleRequest_Article{
			Slug:  "1",
			Title: "New Title",
			Body:  "New Body",
		}}

		resp, err := handler.UpdateArticle(ctx, req)
		if err == nil || status.Code(err) != codes.Internal {
			t.Fatalf("expected Internal error due to db failure, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Database failure during update handled correctly")
	})
}
