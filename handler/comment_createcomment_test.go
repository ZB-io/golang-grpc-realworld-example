package handler

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler struct to simulate real-world example; assume correct imports
type Handler struct {
	// Assume necessary fields exist in your actual implementation
	us *UserService
	as *ArticleService

	logger Logger // Assume logger is defined in your actual package
}

// Assume necessary service and logger interfaces are defined in your package
type UserService interface {
	GetByID(userID uint) (*model.User, error)
}

type ArticleService interface {
	GetByID(articleID uint) (*model.Article, error)
	CreateComment(comment *model.Comment) error
}

type Logger interface {
	Info() *LogEntry
	Error() *LogEntry
}

type LogEntry interface {
	Msgf(format string, args ...interface{})
	Err(error) *LogEntry
	Msg(msg string)
}

// Function the tests are targeting
func (h *Handler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CommentResponse, error) {
	h.logger.Info().Msgf("Create comment | req: %+v", req)

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	currentUser, err := h.us.GetByID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	articleID, err := strconv.Atoi(req.GetSlug())
	if err != nil {
		msg := fmt.Sprintf("cannot convert slug (%s) into integer", req.GetSlug())
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.InvalidArgument, "invalid article id")
	}

	article, err := h.as.GetByID(uint(articleID))
	if err != nil {
		msg := fmt.Sprintf("requested article (slug=%d) not found", articleID)
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Errorf(codes.InvalidArgument, "invalid article id")
	}

	comment := model.Comment{
		Body:      req.GetComment().GetBody(),
		Author:    *currentUser,
		ArticleID: article.ID,
	}

	err = comment.Validate()
	if err != nil {
		err = fmt.Errorf("validation error: %w", err)
		h.logger.Error().Err(err).Msg("validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.as.CreateComment(&comment)
	if err != nil {
		msg := "failed to create comment."
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, msg)
	}

	pc := comment.ProtoComment()
	pc.Author = currentUser.ProtoProfile(false)

	return &pb.CommentResponse{Comment: pc}, nil
}

func TestHandlerCreateComment(t *testing.T) {
	mockUserService := &MockUserService{}
	mockArticleService := &MockArticleService{}
	mockLogger := &MockLogger{}

	handle := &Handler{
		us:     mockUserService,
		as:     mockArticleService,
		logger: mockLogger,
	}

	tests := []struct {
		name   string
		setup  func()
		input  *pb.CreateCommentRequest
		verify func(t *testing.T, resp *pb.CommentResponse, err error)
	}{
		{
			name: "Scenario 1: Unauthorized User Cannot Create Comment",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, status.Error(codes.Unauthenticated, "unauthenticated")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This is a comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.Unauthenticated, status.Code(err))
				t.Log("Unauthorized user was correctly denied comment creation")
			},
		},
		{
			name: "Scenario 2: Non-existent User Cannot Create Comment",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return nil, status.Error(codes.NotFound, "user not found")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This is a comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.NotFound, status.Code(err))
				t.Log("Correctly identified non-existent user")
			},
		},
		{
			name: "Scenario 3: Invalid Article Slug Cannot Lead to Comment Creation",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "NaN",
				Comment: &pb.NewComment{
					Body: "This is a comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				t.Log("Invalid article slug was correctly rejected")
			},
		},
		{
			name: "Scenario 4: Non-existent Article Cannot Have a Comment",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return nil, status.Error(codes.InvalidArgument, "invalid article id")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "99",
				Comment: &pb.NewComment{
					Body: "Trying to comment on a non-existent article",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				t.Log("Non-existent article correctly identified")
			},
		},
		{
			name: "Scenario 5: Valid Comment on Existing Article",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
				mockArticleService.CreateCommentMock = func(comment *model.Comment) error {
					return nil
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This is a valid comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				t.Log("Comment successfully created")
			},
		},
		{
			name: "Scenario 6: Comment Fails Validation",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				t.Log("Correctly failed on comment validation")
			},
		},
		{
			name: "Scenario 7: Database Error on Comment Creation",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
				mockArticleService.CreateCommentMock = func(comment *model.Comment) error {
					return status.Error(codes.Aborted, "database error")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This comment should test db error",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.Aborted, status.Code(err))
				t.Log("Database error correctly handled")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := handle.CreateComment(context.Background(), tt.input)
			tt.verify(t, resp, err)
		})
	}
}

// Mock implementations for testing
type MockUserService struct {
	GetByIDMock func(userID uint) (*model.User, error)
}

func (mus *MockUserService) GetByID(userID uint) (*model.User, error) {
	return mus.GetByIDMock(userID)
}

type MockArticleService struct {
	GetByIDMock       func(articleID uint) (*model.Article, error)
	CreateCommentMock func(comment *model.Comment) error
}

func (mas *MockArticleService) GetByID(articleID uint) (*model.Article, error) {
	return mas.GetByIDMock(articleID)
}

func (mas *MockArticleService) CreateComment(comment *model.Comment) error {
	return mas.CreateCommentMock(comment)
}

type MockLogger struct{}

func (ml *MockLogger) Info() *MockLogEntry {
	return &MockLogEntry{}
}

func (ml *MockLogger) Error() *MockLogEntry {
	return &MockLogEntry{}
}

type MockLogEntry struct{}

func (mle *MockLogEntry) Msgf(format string, args ...interface{}) {}

func (mle *MockLogEntry) Err(err error) *MockLogEntry {
	return mle
}

func (mle *MockLogEntry) Msg(msg string) {}
