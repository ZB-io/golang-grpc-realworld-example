package handler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerDeleteComment(t *testing.T) {
	// Setup mocks and handler
	db, mock, _ := sqlmock.New() // Assume sqlmock is configured properly
	defer db.Close()

	userStore := &store.UserStore{db}
	articleStore := &store.ArticleStore{db}
	logger := zerolog.New(os.Stdout)
	handler := &Handler{
		logger: &logger,
		us:     userStore,
		as:     articleStore,
	}

	type args struct {
		context  context.Context
		request  *pb.DeleteCommentRequest
		userID   uint
		comment  *model.Comment
		mockFunc func()
	}

	tests := []struct {
		name      string
		args      args
		wantError error
	}{
		{
			name: "Successful Deletion of a Comment",
			args: args{
				context: context.WithValue(context.Background(), "userID", uint(1)),
				request: &pb.DeleteCommentRequest{
					Slug: "1",
					Id:   "1",
				},
				userID: 1,
				comment: &model.Comment{
					ID:        1,
					ArticleID: 1,
					UserID:    1,
				},
				mockFunc: func() {
					mock.ExpectQuery("SELECT (.+) FROM `users`").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectQuery("SELECT (.+) FROM `comments`").
						WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "user_id"}).AddRow(1, 1, 1))
					mock.ExpectExec("DELETE FROM `comments`").
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
		},
		{
			name: "Unauthenticated User Attempting to Delete a Comment",
			args: args{
				context: context.Background(),
				request: &pb.DeleteCommentRequest{
					Slug: "1",
					Id:   "1",
				},
				mockFunc: func() {
					// no mock required
				},
			},
			wantError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Current User Not Found in the System",
			args: args{
				context: context.WithValue(context.Background(), "userID", uint(1)),
				request: &pb.DeleteCommentRequest{
					Slug: "1",
					Id:   "2",
				},
				mockFunc: func() {
					mock.ExpectQuery("SELECT (.+) FROM `users`").
						WillReturnError(fmt.Errorf("user not found"))
				},
			},
			wantError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Invalid Comment ID Format",
			args: args{
				context: context.WithValue(context.Background(), "userID", uint(1)),
				request: &pb.DeleteCommentRequest{
					Slug: "1",
					Id:   "abc",
				},
				mockFunc: func(){
					// no mock required
				},
			},
			wantError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Comment Not Found for Given ID",
			args: args{
				context: context.WithValue(context.Background(), "userID", uint(1)),
				request: &pb.DeleteCommentRequest{
					Slug: "1",
					Id:   "123",
				},
				mockFunc: func() {
					mock.ExpectQuery("SELECT (.+) FROM `users`").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectQuery("SELECT (.+) FROM `comments`").
						WillReturnError(fmt.Errorf("record not found"))
				},
			},
			wantError: status.Error(codes.InvalidArgument, "failed to get comment"),
		},
		{
			name: "Slug Mismatch with Comment's Article ID",
			args: args{
				context: context.WithValue(context.Background(), "userID", uint(1)),
				request: &pb.DeleteCommentRequest{
					Slug: "2",
					Id:   "1",
				},
				comment: &model.Comment{
					ID:        1,
					ArticleID: 1,
					UserID:    1,
				},
				mockFunc: func() {
					mock.ExpectQuery("SELECT (.+) FROM `users`").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectQuery("SELECT (.+) FROM `comments`").
						WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "user_id"}).AddRow(1, 1, 1))
				},
			},
			wantError: status.Error(codes.InvalidArgument, "the comment is not in the article"),
		},
		{
			name: "Unauthorized User Attempts to Delete Comment",
			args: args{
				context: context.WithValue(context.Background(), "userID", uint(1)),
				request: &pb.DeleteCommentRequest{
					Slug: "1",
					Id:   "1",
				},
				comment: &model.Comment{
					ID:        1,
					ArticleID: 1,
					UserID:    2, // different user
				},
				mockFunc: func() {
					mock.ExpectQuery("SELECT (.+) FROM `users`").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectQuery("SELECT (.+) FROM `comments`").
						WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "user_id"}).AddRow(1, 1, 2))
				},
			},
			wantError: status.Error(codes.InvalidArgument, "forbidden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mockFunc()
			_, err := handler.DeleteComment(tt.args.context, tt.args.request)

			if tt.wantError != nil && err != nil {
				if status.Code(err) != status.Code(tt.wantError) {
					t.Fatalf("expected error code %v, got %v", status.Code(tt.wantError), status.Code(err))
				}
			} else if tt.wantError != nil && err == nil {
				t.Fatalf("expected error %v, got nil", tt.wantError)
			} else if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
