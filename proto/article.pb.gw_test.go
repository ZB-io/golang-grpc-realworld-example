package proto

import (
	"context"
	"net/http"
	"testing"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/metadata"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"strings"
	"time"
	"github.com/yourrepo/mocks"
	"net"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"github.com/DATA-DOG/go-sqlmock"
	"yourpackage/proto"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"github.com/your_module/proto"
	proto_package "example.com/proto"
	"io"
	"net/url"
	"sync"
)

type ArticlesServer interface {
	GetTags(ctx context.Context, req *Empty) (*Tags, error)
}
type Empty struct{}
type MockArticlesServer struct{}
type Tags struct {
	Tags []string
}
type StreamDesc struct {
	StreamName	string
	Handler		StreamHandler
	ServerStreams	bool
	ClientStreams	bool
}// At least one of these is true.


type mockClientConn struct {
	grpc.ClientConnInterface
}
type DeleteArticleRequest struct {
	Slug string
}
type MockArticlesServer struct{}
type mockArticlesServer struct{}
type mockArticlesServer struct{}
type MockArticlesServer struct {
	mock.Mock
}
type TestResponse struct {
}
type GetCommentsRequest struct {
	Slug string
}
type MockArticlesServer struct{}
type ServerMetadata struct {
	HeaderMD	metadata.MD
	TrailerMD	metadata.MD
}

type AlreadyUnfavoritedResponse struct{}
type MockArticlesServer struct{}
type SuccessResponse struct{}
type Empty struct{}
type MockArticlesClient struct{}
type TagsResponse struct {
	Tags []string
}
type CreateAritcleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}
type MockArticlesServer struct{}
type MockArticlesServerMockRecorder struct {
	mock *MockArticlesServer
}
type SomeResponse struct {
	Message string `json:"message"`
}
type Controller struct {
	T		TestHelper
	mu		sync.Mutex
	expectedCalls	*callSet
	finished	bool
}// T should only be called within a generated mock. It is not intended to
// be used in user code and may be changed in future versions. T is the
// TestReporter passed in when creating the Controller via NewController.
// If the TestReporter does not implement a TestHelper it will be wrapped
// with a nopTestHelper.


type ArticlesClient interface {
	DeleteArticle(ctx context.Context, req *DeleteArticleRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type DeleteArticleRequest struct {
	Slug string
}
type DeleteArticleResponse struct{}
type mockArticlesClient struct{}
type MockArticlesClient struct{}
type mockArticlesClient struct {
	GetArticleFunc func(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error)
}
type ArticlesClient interface {
	GetComments(ctx context.Context, req *GetCommentsRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type Comment struct{}
type CommentsResponse struct {
	Comments []*Comment
}
type GetCommentsRequest struct {
	Slug string
}
type MockArticlesClient struct{}
type MockArticlesClient struct{}
type mockArticlesServer struct {
	CommentResponse proto.Message
	Err             error
}
type MockArticlesServer struct {
	mock.Mock
}
type mockArticlesServer struct{}
type MockArticlesClient struct{}
type MockArticlesClientMockRecorder struct {
	mock *MockArticlesClient
}
type Controller struct {
	T		TestHelper
	mu		sync.Mutex
	expectedCalls	*callSet
	finished	bool
}// T should only be called within a generated mock. It is not intended to
// be used in user code and may be changed in future versions. T is the
// TestReporter passed in when creating the Controller via NewController.
// If the TestReporter does not implement a TestHelper it will be wrapped
// with a nopTestHelper.


type GetFeedArticlesRequest struct {
}
type GetFeedArticlesResponse struct {
}
type MockArticlesClient struct{}
type mockArticlesClient struct{}
type ArticlesClient interface {
	FavoriteArticle(ctx context.Context, req *FavoriteArticleRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type FavoriteArticleRequest struct {
	Slug string
}
type MockArticlesClient struct{}
type mockArticlesClient struct{}
type CreateAritcleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}
type Empty struct{}
type mockArticlesClient struct{}
type CreateAritcleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}
type DeleteArticleRequest struct {
	Slug string
}
type Empty struct{}
type FavoriteArticleRequest struct {
	Slug string
}
type GetCommentsRequest struct {
	Slug string
}
type GetFeedArticlesRequest struct {
}
type MockArticlesServer struct{}
type ServeMux struct {
	handlers	map // handlers maps HTTP method to a list of handlers.
	[string][]handler
	forwardResponseOptions		[]func(context.Context, http.ResponseWriter, proto.Message) error
	marshalers			marshalerRegistry
	incomingHeaderMatcher		HeaderMatcherFunc
	outgoingHeaderMatcher		HeaderMatcherFunc
	metadataAnnotators		[]func(context.Context, *http.Request) metadata.MD
	streamErrorHandler		StreamErrorHandlerFunc
	protoErrorHandler		ProtoErrorHandlerFunc
	disablePathLengthFallback	bool
	lastMatchWins			bool
}

/*
ROOST_METHOD_HASH=local_request_Articles_GetTags_0_8ac54ca3ab
ROOST_METHOD_SIG_HASH=local_request_Articles_GetTags_0_8687d9e2c2


 */
func (m *MockArticlesServer) GetTags(ctx context.Context, req *Empty) (*Tags, error) {
	if m.ReturnError != nil {
		return nil, m.ReturnError
	}
	return m.Tags, nil
}

func Testlocal_request_Articles_GetTags_0(t *testing.T) {
	type args struct {
		ctx        context.Context
		marshaler  runtime.Marshaler
		server     ArticlesServer
		req        *http.Request
		pathParams map[string]string
	}

	tests := []struct {
		name             string
		args             args
		wantTags         *Tags
		wantErr          bool
		expectedErrorMsg string
	}{
		{
			name: "Scenario 1: Successful Retrieval of Tags",
			args: args{
				ctx: context.Background(),
				server: &MockArticlesServer{
					Tags: &Tags{
						Tags: []string{"Go", "Golang", "Programming"},
					},
				},
				req: nil,
			},
			wantTags: &Tags{
				Tags: []string{"Go", "Golang", "Programming"},
			},
			wantErr: false,
		},
		{
			name: "Scenario 2: Server Method Error",
			args: args{
				ctx: context.Background(),
				server: &MockArticlesServer{
					ReturnError: status.Error(codes.Internal, "internal error"),
				},
				req: nil,
			},
			wantTags:         nil,
			wantErr:          true,
			expectedErrorMsg: "internal error",
		},
		{
			name: "Scenario 3: Context Cancellation",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				server: &MockArticlesServer{
					Tags: &Tags{
						Tags: []string{"Go", "Golang"},
					},
				},
				req: nil,
			},
			wantTags:         nil,
			wantErr:          true,
			expectedErrorMsg: "context canceled",
		},
		{
			name: "Scenario 4: Empty Response from Server",
			args: args{
				ctx: context.Background(),
				server: &MockArticlesServer{
					Tags: &Tags{},
				},
				req: nil,
			},
			wantTags: &Tags{},
			wantErr:  false,
		},
		{
			name: "Scenario 5: Nil Request Parameter",
			args: args{
				ctx: context.Background(),
				server: &MockArticlesServer{
					Tags: &Tags{
						Tags: []string{"Go", "Golang"},
					},
				},
				req: nil,
			},
			wantTags: &Tags{
				Tags: []string{"Go", "Golang"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := local_request_Articles_GetTags_0(tt.args.ctx, nil, tt.args.server, tt.args.req, tt.args.pathParams)

			if (err != nil) != tt.wantErr {
				t.Errorf("local_request_Articles_GetTags_0() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
				t.Errorf("local_request_Articles_GetTags_0() error message = %v, expectedErrorMsg %v", err.Error(), tt.expectedErrorMsg)
			}

			if !tt.wantErr && !proto.Equal(got, tt.wantTags) {
				t.Errorf("local_request_Articles_GetTags_0() = %v, want %v", got, tt.wantTags)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=RegisterArticlesHandler_9b17f012c0
ROOST_METHOD_SIG_HASH=RegisterArticlesHandler_018d6724b4


 */
func (m *mockClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return nil
}

func (m *mockClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func TestRegisterArticlesHandler(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		mux           *runtime.ServeMux
		conn          grpc.ClientConnInterface
		expectedError error
	}{
		{
			name: "Scenario 1: Successful Registration of Articles Handler",
			ctx:  context.Background(),
			mux:  runtime.NewServeMux(),
			conn: &mockClientConn{},
		},
		{
			name:          "Scenario 2: Registration with Nil Context",
			ctx:           nil,
			mux:           runtime.NewServeMux(),
			conn:          &mockClientConn{},
			expectedError: status.Error(codes.InvalidArgument, "context is nil"),
		},
		{
			name:          "Scenario 3: Invalid gRPC Client Connection",
			ctx:           context.Background(),
			mux:           runtime.NewServeMux(),
			conn:          nil,
			expectedError: status.Error(codes.Internal, "grpc connection is invalid"),
		},
		{
			name: "Scenario 4: Validate Integration with ArticlesClient",
			ctx:  context.Background(),
			mux:  runtime.NewServeMux(),
			conn: &mockClientConn{},
		},
		{
			name: "Scenario 5: Empty ServeMux",
			ctx:  context.Background(),
			mux:  runtime.NewServeMux(),
			conn: &mockClientConn{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterArticlesHandler(tt.ctx, tt.mux, tt.conn.(*grpc.ClientConn))

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
			t.Log("Test case:", tt.name, "completed successfully")
		})
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_DeleteArticle_0_ff2c7120b3
ROOST_METHOD_SIG_HASH=local_request_Articles_DeleteArticle_0_29f03e64ed


 */
func (m *MockArticlesServer) DeleteArticle(ctx context.Context, req *DeleteArticleRequest) (proto.Message, error) {
	if req.Slug == "existing-slug" {
		return &DeleteArticleResponse{Message: "Deleted"}, nil
	}
	if req.Slug == "error-on-delete" {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return nil, status.Errorf(codes.NotFound, "article not found")
}

func Testlocal_request_Articles_DeleteArticle_0(t *testing.T) {
	mockServer := &MockArticlesServer{}

	tests := []struct {
		name        string
		pathParams  map[string]string
		expectedMsg proto.Message
		expectedErr error
	}{
		{
			name: "Successful Article Deletion",
			pathParams: map[string]string{
				"slug": "existing-slug",
			},
			expectedMsg: &DeleteArticleResponse{Message: "Deleted"},
			expectedErr: nil,
		},
		{
			name:        "Missing Slug Parameter",
			pathParams:  map[string]string{},
			expectedMsg: nil,
			expectedErr: status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug"),
		},
		{
			name: "Invalid Slug Type Mismatch",
			pathParams: map[string]string{
				"slug": "%%%",
			},
			expectedMsg: nil,
			expectedErr: status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", errors.New("invalid type")),
		},
		{
			name: "Server-Side Deletion Fails",
			pathParams: map[string]string{
				"slug": "error-on-delete",
			},
			expectedMsg: nil,
			expectedErr: status.Errorf(codes.Internal, "internal error"),
		},
		{
			name: "Edge Case with Special Characters in Slug",
			pathParams: map[string]string{
				"slug": "special@!#",
			},
			expectedMsg: nil,
			expectedErr: status.Errorf(codes.NotFound, "article not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := &http.Request{}
			marshaler := &runtime.JSONPb{}

			msg, _, err := local_request_Articles_DeleteArticle_0(ctx, marshaler, mockServer, req, tt.pathParams)

			if want, got := tt.expectedErr, err; !errors.Is(got, want) {
				t.Errorf("expected error %v, got %v", want, got)
			}

			if want, got := tt.expectedMsg, msg; !proto.Equal(want, got) {
				t.Errorf("expected message %v, got %v", want, got)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_DeleteComment_0_00619a3bf5
ROOST_METHOD_SIG_HASH=local_request_Articles_DeleteComment_0_722319db7d


 */
func Testlocal_request_Articles_DeleteComment_0(t *testing.T) {
	tests := []struct {
		name           string
		pathParams     map[string]string
		serverResponse proto.Message
		serverError    error
		expectedError  error
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := newMockArticlesServer(t)
			defer mockServer.finish()

			if tc.serverError == nil {
				mockServer.expectDeleteComment(context.TODO(), &DeleteCommentRequest{
					Slug: tc.pathParams["slug"],
					Id:   tc.pathParams["id"],
				}, tc.serverResponse, tc.serverError)
			}

			resp, _, err := local_request_Articles_DeleteComment_0(context.TODO(), nil, mockServer.mockInterface, &http.Request{}, tc.pathParams)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if !proto.Equal(resp, tc.serverResponse) {
					t.Errorf("expected response: %v, got: %v", tc.serverResponse, resp)
				}
			}
		})
	}
}

func (m *mockArticlesServer) expectDeleteComment(ctx context.Context, req *DeleteCommentRequest, resp proto.Message, err error) {
	m.mockInterface.EXPECT().DeleteComment(ctx, req).Return(resp, err).Times(1)
}

func (m *mockArticlesServer) finish() {
	m.mockCtrl.Finish()
}

func newMockArticlesServer(t *testing.T) *mockArticlesServer {
	ctrl := gomock.NewController(t)
	mock := NewMockArticlesServer(ctrl)
	return &mockArticlesServer{mockCtrl: ctrl, mockInterface: mock}
}

/*
ROOST_METHOD_HASH=local_request_Articles_GetArticle_0_05eebc8ecf
ROOST_METHOD_SIG_HASH=local_request_Articles_GetArticle_0_46be1e7697


 */
func (m *mockArticlesServer) CreateArticle(ctx context.Context, req *CreateArticleRequest) (proto.Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method not implemented")
}

func (m *mockArticlesServer) GetArticle(ctx context.Context, req *GetArticleRequest) (proto.Message, error) {
	if req.Slug == "valid-slug" {
		return &ArticleResponse{Title: "Test Article"}, nil
	} else {
		return nil, status.Errorf(codes.NotFound, "article not found")
	}
}

func Testlocal_request_Articles_GetArticle_0(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		expectErr  bool
		errCode    codes.Code
		expectData bool
	}{
		{
			name:       "Valid Slug",
			pathParams: map[string]string{"slug": "valid-slug"},
			expectErr:  false,
			expectData: true,
		},
		{
			name:       "Missing Slug Parameter",
			pathParams: map[string]string{},
			expectErr:  true,
			errCode:    codes.InvalidArgument,
			expectData: false,
		},
		{
			name:       "Simulated Article Retrieval Failure",
			pathParams: map[string]string{"slug": "invalid-slug"},
			expectErr:  true,
			errCode:    codes.NotFound,
			expectData: false,
		},
	}

	server := &mockArticlesServer{}
	ctx := context.TODO()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, md, err := local_request_Articles_GetArticle_0(ctx, &runtime.JSONPb{}, server, &http.Request{}, tt.pathParams)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				st, ok := status.FromError(err)
				if !ok || st.Code() != tt.errCode {
					t.Fatalf("expected error code %v, got %v", tt.errCode, st.Code())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if tt.expectData {
				if resp == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
			} else {
				if resp != nil {
					t.Fatalf("expected nil response, got non-nil")
				}
			}

			if md.HeaderMD == nil {
				t.Fatalf("expected non-nil metadata, got nil")
			}
		})
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_GetArticles_0_e3af9abfe7
ROOST_METHOD_SIG_HASH=local_request_Articles_GetArticles_0_34a360ba7c


 */
func (m *MockArticlesServer) GetArticles(ctx context.Context, req *GetArticlesRequest) (proto.Message, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(proto.Message), args.Error(1)
}

func (t *TestResponse) ProtoMessage() {}

func (t *TestResponse) Reset() {}

func (t *TestResponse) String() string { return "" }

func Testlocal_request_Articles_GetArticles_0(t *testing.T) {
	type testCase struct {
		name         string
		setupMock    func(*MockArticlesServer)
		request      *http.Request
		expectedMsg  proto.Message
		expectedCode codes.Code
	}

	tests := []testCase{
		{
			name: "Scenario 1: Successful Execution with Valid Query Parameters",
			setupMock: func(m *MockArticlesServer) {
				m.On("GetArticles", mock.Anything, mock.Anything).
					Return(&TestResponse{}, nil)
			},
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles?validParam=value", nil)
				return req
			}(),
			expectedMsg:  &TestResponse{},
			expectedCode: codes.OK,
		},
		{
			name: "Scenario 2: Error Handling on Invalid Query Parameters",
			setupMock: func(m *MockArticlesServer) {

			},
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles?invalidParam=value", nil)
				return req
			}(),
			expectedMsg:  nil,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Scenario 3: Null Request Handling",
			setupMock: func(m *MockArticlesServer) {

			},
			request:      nil,
			expectedMsg:  nil,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Scenario 4: Handling Server Method Failure",
			setupMock: func(m *MockArticlesServer) {
				m.On("GetArticles", mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.Internal, "internal error"))
			},
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles", nil)
				return req
			}(),
			expectedMsg:  nil,
			expectedCode: codes.Internal,
		},
		{
			name: "Scenario 5: Successful Execution with Empty Query Parameters",
			setupMock: func(m *MockArticlesServer) {
				m.On("GetArticles", mock.Anything, mock.Anything).
					Return(&TestResponse{}, nil)
			},
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles", nil)
				return req
			}(),
			expectedMsg:  &TestResponse{},
			expectedCode: codes.OK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := new(MockArticlesServer)
			tc.setupMock(mockServer)

			msg, _, err := local_request_Articles_GetArticles_0(context.Background(), &runtime.JSONPb{}, mockServer, tc.request, nil)

			if tc.expectedCode == codes.OK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.expectedCode, st.Code())
			}
			assert.Equal(t, tc.expectedMsg, msg)
		})
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_GetComments_0_5da1ac4759
ROOST_METHOD_SIG_HASH=local_request_Articles_GetComments_0_935f99d11e


 */
func (m *MockArticlesServer) CreateArticle(ctx context.Context, req *CreateArticleRequest) (*CreateArticleResponse, error) {
	return nil, nil
}

func (m *MockArticlesServer) GetComments(ctx context.Context, req *GetCommentsRequest) (proto.Message, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(proto.Message), args.Error(1)
}

func Testlocal_request_Articles_GetComments_0(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		setupServer  func() *MockArticlesServer
		reqPath      string
		expectedErr  codes.Code
		expectedBody string
	}

	tests := []testCase{
		{
			name: "Successful Retrieval",
			setupServer: func() *MockArticlesServer {
				mockServer := new(MockArticlesServer)
				mockServer.On("GetComments", mock.Anything, &GetCommentsRequest{Slug: "valid-slug"}).Return(&CommentsResponse{}, nil)
				return mockServer
			},
			reqPath:      "/articles/valid-slug/comments",
			expectedErr:  codes.OK,
			expectedBody: "successful response body",
		},
		{
			name: "Missing Slug Parameter",
			setupServer: func() *MockArticlesServer {
				return new(MockArticlesServer)
			},
			reqPath:      "/articles//comments",
			expectedErr:  codes.InvalidArgument,
			expectedBody: "",
		},
		{
			name: "Type Mismatch for Slug Parameter",
			setupServer: func() *MockArticlesServer {
				return new(MockArticlesServer)
			},
			reqPath:      "/articles/12345/comments",
			expectedErr:  codes.InvalidArgument,
			expectedBody: "",
		},
		{
			name: "Server Returns an Error",
			setupServer: func() *MockArticlesServer {
				mockServer := new(MockArticlesServer)
				mockServer.On("GetComments", mock.Anything, &GetCommentsRequest{Slug: "error-slug"}).Return(nil, status.Error(codes.Internal, "internal server error"))
				return mockServer
			},
			reqPath:      "/articles/error-slug/comments",
			expectedErr:  codes.Internal,
			expectedBody: "",
		},
		{
			name: "Server Returns No Comments for a Valid Slug",
			setupServer: func() *MockArticlesServer {
				mockServer := new(MockArticlesServer)
				mockServer.On("GetComments", mock.Anything, &GetCommentsRequest{Slug: "no-comments-slug"}).Return(&CommentsResponse{}, nil)
				return mockServer
			},
			reqPath:      "/articles/no-comments-slug/comments",
			expectedErr:  codes.OK,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticlesServer := tt.setupServer()

			ctx := context.Background()
			req := httptest.NewRequest(http.MethodGet, tt.reqPath, nil)
			resp, _, err := local_request_Articles_GetComments_0(ctx, &runtime.JSONPb{}, mockArticlesServer, req, map[string]string{"slug": "slug-value"})

			if tt.expectedErr == codes.OK && err != nil {
				t.Fatalf("expected no error, got %v", err)
			} else if status.Code(err) != tt.expectedErr {
				t.Fatalf("expected error code %v, got %v", tt.expectedErr, status.Code(err))
			}

			if resp != nil {
				body := resp.String()
				if !strings.Contains(body, tt.expectedBody) {
					t.Logf("Response body: %v", body)
					t.Fatalf("expected body to contain %v", tt.expectedBody)
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_GetFeedArticles_0_6e1296622f
ROOST_METHOD_SIG_HASH=local_request_Articles_GetFeedArticles_0_646bc2f91c


 */
func Testlocal_request_Articles_GetFeedArticles_0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name        string
		serverSetup func(*mocks.MockArticlesServer)
		request     func() *http.Request
		ctxSetup    func() (context.Context, context.CancelFunc)
		expectedMsg proto.Message
		expectedErr error
		expectedMD  runtime.ServerMetadata
	}

	tests := []testCase{
		{
			name: "Successfully Retrieve Feed Articles",
			serverSetup: func(s *mocks.MockArticlesServer) {
				s.EXPECT().GetFeedArticles(gomock.Any(), gomock.Any()).Return(&GetFeedArticlesResponse{Articles: []*Article{{Id: "1", Title: "Article 1"}}}, nil)
			},
			request: func() *http.Request {
				r, _ := http.NewRequest("GET", "/articles/feed", nil)
				return r
			},
			ctxSetup: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
			expectedMsg: &GetFeedArticlesResponse{Articles: []*Article{{Id: "1", Title: "Article 1"}}},
			expectedErr: nil,
		},
		{
			name:        "Invalid Query Parameters",
			serverSetup: func(s *mocks.MockArticlesServer) {},
			request: func() *http.Request {
				r, _ := http.NewRequest("GET", "/articles/feed?invalid=param", nil)
				return r
			},
			ctxSetup: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
			expectedMsg: nil,
			expectedErr: status.Errorf(codes.InvalidArgument, "invalid query parameter"),
		},
		{
			name: "Server Error on Fetching Feed Articles",
			serverSetup: func(s *mocks.MockArticlesServer) {
				s.EXPECT().GetFeedArticles(gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.Internal, "internal server error"))
			},
			request: func() *http.Request {
				r, _ := http.NewRequest("GET", "/articles/feed", nil)
				return r
			},
			ctxSetup: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
			expectedMsg: nil,
			expectedErr: status.Errorf(codes.Internal, "internal server error"),
		},
		{
			name: "No Articles Returned",
			serverSetup: func(s *mocks.MockArticlesServer) {
				s.EXPECT().GetFeedArticles(gomock.Any(), gomock.Any()).Return(&GetFeedArticlesResponse{Articles: []*Article{}}, nil)
			},
			request: func() *http.Request {
				r, _ := http.NewRequest("GET", "/articles/feed", nil)
				return r
			},
			ctxSetup: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
			expectedMsg: &GetFeedArticlesResponse{Articles: []*Article{}},
			expectedErr: nil,
		},
		{
			name: "Context Deadline Exceeded",
			serverSetup: func(s *mocks.MockArticlesServer) {

				s.EXPECT().GetFeedArticles(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, req *GetFeedArticlesRequest) (*GetFeedArticlesResponse, error) {
						time.Sleep(3 * time.Second)
						return nil, nil
					},
				)
			},
			request: func() *http.Request {
				r, _ := http.NewRequest("GET", "/articles/feed", nil)
				return r
			},
			ctxSetup: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 1*time.Second)
			},
			expectedMsg: nil,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name: "Network Failure Simulation",
			serverSetup: func(s *mocks.MockArticlesServer) {
				s.EXPECT().GetFeedArticles(gomock.Any(), gomock.Any()).Return(nil, errors.New("network error"))
			},
			request: func() *http.Request {
				r, _ := http.NewRequest("GET", "/articles/feed", nil)
				return r
			},
			ctxSetup: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
			expectedMsg: nil,
			expectedErr: errors.New("network error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctxSetup()
			defer cancel()

			server := mocks.NewMockArticlesServer(ctrl)
			tt.serverSetup(server)

			request := tt.request()
			marshaler := &runtime.JSONPb{}

			msg, md, err := local_request_Articles_GetFeedArticles_0(ctx, marshaler, server, request, nil)

			if tt.expectedErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !proto.Equal(msg, tt.expectedMsg) {
					t.Errorf("expected message %v, got %v", tt.expectedMsg, msg)
				}
			} else {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			}

			if !messagesAreEqual(md, tt.expectedMD) {
				t.Errorf("expected metadata %v, got %v", tt.expectedMD, md)
			}
		})
	}
}

func messagesAreEqual(md1, md2 runtime.ServerMetadata) bool {

	return md1 == md2
}

/*
ROOST_METHOD_HASH=local_request_Articles_UnfavoriteArticle_0_ba003c030c
ROOST_METHOD_SIG_HASH=local_request_Articles_UnfavoriteArticle_0_f4798c29ac


 */
func (s *SuccessResponse) ProtoMessage() {}

func (s *SuccessResponse) Reset() {}

func (s *SuccessResponse) String() string { return "success" }

func Testlocal_request_Articles_UnfavoriteArticle_0(t *testing.T) {
	tests := []struct {
		name         string
		slug         string
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name:         "Successful Unfavorite Article",
			slug:         "valid-slug",
			expectedCode: codes.OK,
			expectedMsg:  "success",
		},
		{
			name:         "Missing Slug in Path Parameters",
			slug:         "",
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "missing parameter slug",
		},
		{
			name:         "Server Error upon Unfavoriting",
			slug:         "server-error",
			expectedCode: codes.Internal,
			expectedMsg:  "internal server error",
		},
		{
			name:         "Valid Slug but Article Already Unfavorited",
			slug:         "already-unfavorited",
			expectedCode: codes.OK,
			expectedMsg:  "already unfavorited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := &MockArticlesServer{}

			req := httptest.NewRequest(http.MethodDelete, "/articles/"+tt.slug+"/favorite", nil)
			pathParams := map[string]string{"slug": tt.slug}

			resp, _, err := local_request_Articles_UnfavoriteArticle_0(context.Background(), &runtime.JSONPb{}, mockServer, req, pathParams)

			if tt.slug == "" {
				if err == nil || status.Code(err) != tt.expectedCode {
					t.Errorf("expected error code %v, got %v", tt.expectedCode, err)
					return
				}
				if status.Convert(err).Message() != tt.expectedMsg {
					t.Errorf("expected message %v, got %v", tt.expectedMsg, err)
				}
				return
			}

			if status.Convert(err).Code() != tt.expectedCode {
				t.Errorf("expected code %v, got %v", tt.expectedCode, status.Convert(err).Code())
			}

			if err == nil {
				if resp.String() != tt.expectedMsg {
					t.Errorf("expected message %v, got %v", tt.expectedMsg, resp.String())
				}
			} else {
				if status.Convert(err).Message() != tt.expectedMsg {
					t.Errorf("expected message %v, got %v", tt.expectedMsg, status.Convert(err).Message())
				}
			}

			t.Logf("Test scenario %s passed", tt.name)
		})
	}
}

func (m *MockArticlesServer) UnfavoriteArticle(ctx context.Context, req *UnfavoriteArticleRequest) (proto.Message, error) {
	switch req.Slug {
	case "valid-slug":
		return &SuccessResponse{}, nil
	case "server-error":
		return nil, status.Errorf(codes.Internal, "internal server error")
	case "already-unfavorited":
		return &AlreadyUnfavoritedResponse{}, nil
	default:
		return nil, status.Errorf(codes.NotFound, "article not found")
	}
}

/*
ROOST_METHOD_HASH=request_Articles_GetTags_0_c7f91452b0
ROOST_METHOD_SIG_HASH=request_Articles_GetTags_0_30370e8d01


 */
func (m *MockArticlesClient) GetTags(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error) {
	return m.GetTagsFunc(ctx, in, opts...)
}

func Testrequest_Articles_GetTags_0(t *testing.T) {
	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error)
		expectMsg proto.Message
		expectErr error
	}{
		{
			name: "Successful Retrieval of Tags",
			mockFunc: func(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error) {
				return &TagsResponse{Tags: []string{"go", "programming"}}, nil
			},
			expectMsg: &TagsResponse{Tags: []string{"go", "programming"}},
			expectErr: nil,
		},
		{
			name: "Handle Client Error Gracefully",
			mockFunc: func(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error) {
				return nil, status.Error(codes.Internal, "client error")
			},
			expectMsg: nil,
			expectErr: status.Error(codes.Internal, "client error"),
		},
		{
			name: "Empty Tags Response",
			mockFunc: func(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error) {
				return &TagsResponse{Tags: []string{}}, nil
			},
			expectMsg: &TagsResponse{Tags: []string{}},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockArticlesClient{
				GetTagsFunc: tt.mockFunc,
			}

			ctx := context.Background()
			mockMetadata := runtime.ServerMetadata{
				HeaderMD: metadata.MD{"key": []string{"value"}},
			}

			msg, metadata, err := request_Articles_GetTags_0(ctx, nil, client, nil, nil)
			if tt.expectErr != nil {
				assert.NotNil(t, err, "Expected error")
				assert.Equal(t, tt.expectErr.Error(), err.Error(), "Error did not match")
			} else {
				assert.Nil(t, err, "Unexpected error occurred")
			}

			assert.Equal(t, tt.expectMsg, msg, "Message did not match")
			assert.Equal(t, mockMetadata.HeaderMD, metadata.HeaderMD, "Metadata header did not match")

			t.Logf("Test %s passed successfully", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=RegisterArticlesHandlerFromEndpoint_a35e501413
ROOST_METHOD_SIG_HASH=RegisterArticlesHandlerFromEndpoint_0d4fca2f6c


 */
func TestRegisterArticlesHandlerFromEndpoint(t *testing.T) {
	type args struct {
		ctx      context.Context
		mux      *runtime.ServeMux
		endpoint string
		opts     []grpc.DialOption
	}

	tests := []struct {
		name        string
		args        args
		expectError bool
	}{
		{
			name: "Successful Registration with Valid Endpoint and Options",
			args: args{
				ctx:      context.Background(),
				mux:      runtime.NewServeMux(),
				endpoint: "127.0.0.1:9090",
				opts:     []grpc.DialOption{grpc.WithInsecure()},
			},
			expectError: false,
		},
		{
			name: "Error Handling for Invalid Endpoint",
			args: args{
				ctx:      context.Background(),
				mux:      runtime.NewServeMux(),
				endpoint: "invalid-endpoint",
				opts:     []grpc.DialOption{grpc.WithInsecure()},
			},
			expectError: true,
		},
		{
			name: "Error Handling for Invalid gRPC Dial Options",
			args: args{
				ctx:      context.Background(),
				mux:      runtime.NewServeMux(),
				endpoint: "127.0.0.1:9090",
				opts:     nil,
			},
			expectError: true,
		},
		{
			name: "Context Cancellation During Registration Process",
			args: func() args {
				ctx, cancel := context.WithCancel(context.Background())
				addr, listener := getAvailableListener()

				defer listener.Close()

				return args{
					ctx:      ctx,
					mux:      runtime.NewServeMux(),
					endpoint: addr,
					opts:     []grpc.DialOption{grpc.WithInsecure()},
				}
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterArticlesHandlerFromEndpoint(tt.args.ctx, tt.args.mux, tt.args.endpoint, tt.args.opts)
			if (err != nil) != tt.expectError {
				t.Errorf("Test %s failed: expected error = %v, got = %v", tt.name, tt.expectError, err)
			}

			if tt.name == "Context Cancellation During Registration Process" {
				cancel := func() {
					tt.args.ctx.Done()
				}

				cancel()

				select {
				case <-tt.args.ctx.Done():
					t.Log("Context successfully canceled, resources cleaned.")
				default:
					t.Errorf("Context should have been canceled, but it wasn't.")
				}
			}
		})
	}
}

func getAvailableListener() (string, net.Listener) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	return listener.Addr().String(), listener
}

/*
ROOST_METHOD_HASH=local_request_Articles_CreateArticle_0_bd3daf4d54
ROOST_METHOD_SIG_HASH=local_request_Articles_CreateArticle_0_d8bf3dbbbc


 */
func (m *MockArticlesServer) CreateArticle(ctx context.Context, req *CreateAritcleRequest) (proto.Message, error) {
	ret := m.ctrl.Call(m, "CreateArticle", ctx, req)
	ret0, _ := ret[0].(proto.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func NewMockArticlesServer(ctrl *gomock.Controller) *MockArticlesServer {
	mock := &MockArticlesServer{ctrl: ctrl}
	mock.recorder = &MockArticlesServerMockRecorder{mock}
	return mock
}

func Testlocal_request_Articles_CreateArticle_0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := NewMockArticlesServer(ctrl)

	tests := []struct {
		name           string
		requestBody    string
		expectedError  error
		expectedResult proto.Message
		mockSetup      func()
	}{
		{
			name:          "Valid Article Creation",
			requestBody:   `{"title": "Article", "description": "Description", "body": "Content"}`,
			expectedError: nil,
			expectedResult: &SomeResponse{
				Message: "Article created",
			},
			mockSetup: func() {
				mockServer.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(&SomeResponse{
					Message: "Article created",
				}, nil).Times(1)
			},
		},
		{
			name:           "Invalid JSON Request Body",
			requestBody:    `{title: "Article"}`,
			expectedError:  status.Errorf(codes.InvalidArgument, "invalid character 't' looking for beginning of object key string"),
			expectedResult: nil,
			mockSetup:      func() {},
		},
		{
			name:           "Empty Request Body",
			requestBody:    ``,
			expectedError:  status.Errorf(codes.InvalidArgument, "EOF"),
			expectedResult: nil,
			mockSetup:      func() {},
		},
		{
			name:           "Server-Side Error During Article Creation",
			requestBody:    `{"title": "Article", "description": "Description", "body": "Content"}`,
			expectedError:  status.Errorf(codes.Internal, "Server error"),
			expectedResult: nil,
			mockSetup: func() {
				mockServer.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.Internal, "Server error")).Times(1)
			},
		},
		{
			name:           "Handle Decoding Error Gracefully",
			requestBody:    `{"unexpected_field": "value"}`,
			expectedError:  status.Errorf(codes.InvalidArgument, "proto: syntax error (line 1:2): unexpected \"value\""),
			expectedResult: nil,
			mockSetup:      func() {},
		},
		{
			name:           "Incomplete Article Data",
			requestBody:    `{"title": "Article"}`,
			expectedError:  status.Errorf(codes.InvalidArgument, "field missing"),
			expectedResult: nil,
			mockSetup: func() {
				mockServer.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.InvalidArgument, "field missing")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(tt.requestBody)),
			}
			marshaler := &runtime.JSONPB{}

			resp, metadata, err := local_request_Articles_CreateArticle_0(context.Background(), marshaler, mockServer, req, nil)
			if err != nil && tt.expectedError == nil || err == nil && tt.expectedError != nil || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			if !eq(resp, tt.expectedResult) {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, resp)
			}

			if (metadata != runtime.ServerMetadata{}) {
				t.Errorf("Unexpected metadata, got %v", metadata)
			}
		})
	}
}

func eq(a, b proto.Message) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	ab, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return bytes.Equal(ab, bb)
}

/*
ROOST_METHOD_HASH=request_Articles_DeleteArticle_0_f61d5cef63
ROOST_METHOD_SIG_HASH=request_Articles_DeleteArticle_0_77d0697b5f


 */
func (m *mockArticlesClient) DeleteArticle(ctx context.Context, req *DeleteArticleRequest, opts ...grpc.CallOption) (proto.Message, error) {
	if req.Slug == "existing-slug" {
		return &DeleteArticleResponse{}, nil
	} else if req.Slug == "nonexistent-slug" {
		return nil, status.Errorf(codes.NotFound, "article not found")
	} else if req.Slug == "network-error" {
		return nil, status.Errorf(codes.Internal, "network error")
	}
	return nil, status.Errorf(codes.InvalidArgument, "invalid argument")
}

func Testrequest_Articles_DeleteArticle_0(t *testing.T) {
	tests := []struct {
		name         string
		pathParams   map[string]string
		expectedCode codes.Code
		errMessage   string
	}{
		{
			name:         "Valid Request with Existing Article Slug",
			pathParams:   map[string]string{"slug": "existing-slug"},
			expectedCode: codes.OK,
			errMessage:   "",
		},
		{
			name:         "Missing Article Slug Parameter",
			pathParams:   map[string]string{},
			expectedCode: codes.InvalidArgument,
			errMessage:   "missing parameter slug",
		},
		{
			name:         "Invalid Slug Type Conversion",
			pathParams:   map[string]string{"slug": "123"},
			expectedCode: codes.InvalidArgument,
			errMessage:   "type mismatch",
		},
		{
			name:         "Deletion Failure Due to Nonexistent Article",
			pathParams:   map[string]string{"slug": "nonexistent-slug"},
			expectedCode: codes.NotFound,
			errMessage:   "article not found",
		},
		{
			name:         "Network or System Error During Deletion",
			pathParams:   map[string]string{"slug": "network-error"},
			expectedCode: codes.Internal,
			errMessage:   "network error",
		},
	}

	ctx := context.Background()
	client := &mockArticlesClient{}
	marshaler := &runtime.JSONPb{}
	req := &http.Request{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, _, err := request_Articles_DeleteArticle_0(ctx, marshaler, client, req, tt.pathParams)

			if tt.expectedCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				t.Logf("Request succeeded for slug: %s", tt.pathParams["slug"])
			} else {
				assert.Error(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tt.expectedCode, st.Code())
				assert.Contains(t, st.Message(), tt.errMessage)
				t.Logf("Expected error for slug %s: %s", tt.pathParams["slug"], st.Message())
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_DeleteComment_0_c6c6d3d1a0
ROOST_METHOD_SIG_HASH=request_Articles_DeleteComment_0_06bcae45df


 */
func (m *MockArticlesClient) DeleteComment(ctx context.Context, req *DeleteCommentRequest, opts ...grpc.CallOption) (proto.Message, error) {
	if req.Slug == "validSlug" && req.Id == "validId" {
		return &DeleteCommentResponse{Success: true}, nil
	}
	return nil, status.Errorf(codes.NotFound, "comment not found")
}

func Testrequest_Articles_DeleteComment_0(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		wantError  bool
		errorCode  codes.Code
	}{
		{
			name:       "Successfully Delete Comment",
			pathParams: map[string]string{"slug": "validSlug", "id": "validId"},
			wantError:  false,
		},
		{
			name:       "Missing `slug` Parameter",
			pathParams: map[string]string{"id": "validId"},
			wantError:  true,
			errorCode:  codes.InvalidArgument,
		},
		{
			name:       "Missing `id` Parameter",
			pathParams: map[string]string{"slug": "validSlug"},
			wantError:  true,
			errorCode:  codes.InvalidArgument,
		},
		{
			name:       "Invalid `slug` Type",
			pathParams: map[string]string{"slug": "123", "id": "validId"},
			wantError:  true,
			errorCode:  codes.InvalidArgument,
		},
		{
			name:       "Invalid `id` Type",
			pathParams: map[string]string{"slug": "validSlug", "id": "abc"},
			wantError:  true,
			errorCode:  codes.InvalidArgument,
		},
		{
			name:       "Client Returns Error",
			pathParams: map[string]string{"slug": "invalidSlug", "id": "invalidId"},
			wantError:  true,
			errorCode:  codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockArticlesClient{}
			ctx := context.Background()
			req := &http.Request{}

			response, _, err := request_Articles_DeleteComment_0(ctx, &runtime.JSONPb{}, client, req, tt.pathParams)

			if tt.wantError {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errorCode, st.Code())
				t.Logf("Test '%s' correctly returned error: %v", tt.name, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				t.Logf("Test '%s' succeeded with response: %v", tt.name, response)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_GetArticle_0_735c1e04a5
ROOST_METHOD_SIG_HASH=request_Articles_GetArticle_0_2f5f69725a


 */
func (m *mockArticlesClient) GetArticle(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error) {
	return m.GetArticleFunc(ctx, req, opts...)
}

func Testrequest_Articles_GetArticle_0(t *testing.T) {
	tests := []struct {
		name        string
		pathParams  map[string]string
		getArticle  func(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error)
		expectedErr error
	}{
		{
			name:       "Missing Slug Parameter",
			pathParams: map[string]string{},
			getArticle: func(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error) {
				return nil, nil
			},
			expectedErr: status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug"),
		},
		{
			name: "Slug Parameter Type Mismatch",
			pathParams: map[string]string{
				"slug": "invalid!slug",
			},
			getArticle: func(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error) {
				return nil, nil
			},
			expectedErr: status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", errors.New("invalid slug format")),
		},
		{
			name: "Successful Article Retrieval",
			pathParams: map[string]string{
				"slug": "valid_slug",
			},
			getArticle: func(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error) {
				return &proto.Article{Title: "Test Article"}, nil
			},
			expectedErr: nil,
		},
		{
			name: "Client GetArticle Error Handling",
			pathParams: map[string]string{
				"slug": "valid_slug",
			},
			getArticle: func(ctx context.Context, req *proto.GetArticleRequest, opts ...grpc.CallOption) (*proto.Article, error) {
				return nil, status.Errorf(codes.Internal, "internal error")
			},
			expectedErr: status.Errorf(codes.Internal, "internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			marshaler := &runtime.JSONPb{}
			client := &mockArticlesClient{
				GetArticleFunc: tt.getArticle,
			}
			req, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			msg, _, err := proto.request_Articles_GetArticle_0(ctx, marshaler, client, req, tt.pathParams)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_GetComments_0_1a4ff90d5d
ROOST_METHOD_SIG_HASH=request_Articles_GetComments_0_e547ca5bba


 */
func (m *MockArticlesClient) GetComments(ctx context.Context, req *GetCommentsRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return m.mockGetComments(ctx, req, opts...)
}

func Testrequest_Articles_GetComments_0(t *testing.T) {
	testCases := []struct {
		name            string
		pathParams      map[string]string
		mockGetComments func(ctx context.Context, req *GetCommentsRequest, opts ...grpc.CallOption) (proto.Message, error)
		expectedError   error
		expectedMessage proto.Message
	}{
		{
			name: "Valid Request with Existing Slug",
			pathParams: map[string]string{
				"slug": "valid-slug",
			},
			mockGetComments: func(ctx context.Context, req *GetCommentsRequest, opts ...grpc.CallOption) (proto.Message, error) {
				return &CommentsResponse{}, nil
			},
			expectedError:   nil,
			expectedMessage: &CommentsResponse{},
		},
		{
			name:            "Missing Slug in Path Parameters",
			pathParams:      map[string]string{},
			mockGetComments: nil,
			expectedError:   status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug"),
		},
		{
			name: "Slug Type Mismatch Error",
			pathParams: map[string]string{
				"slug": "",
			},
			mockGetComments: nil,
			expectedError:   status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", errors.New("invalid slug")),
		},
		{
			name: "Error from GetComments Method",
			pathParams: map[string]string{
				"slug": "error-slug",
			},
			mockGetComments: func(ctx context.Context, req *GetCommentsRequest, opts ...grpc.CallOption) (proto.Message, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expectedError: status.Error(codes.Internal, "internal error"),
		},
		{
			name: "Valid Request with No Comments",
			pathParams: map[string]string{
				"slug": "no-comments-slug",
			},
			mockGetComments: func(ctx context.Context, req *GetCommentsRequest, opts ...grpc.CallOption) (proto.Message, error) {
				return &CommentsResponse{Comments: []*Comment{}}, nil
			},
			expectedError:   nil,
			expectedMessage: &CommentsResponse{Comments: []*Comment{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockArticlesClient{
				mockGetComments: tc.mockGetComments,
			}

			ctx := context.Background()

			marshaler := &runtime.JSONPb{}
			req := &http.Request{}

			msg, _, err := request_Articles_GetComments_0(ctx, marshaler, mockClient, req, tc.pathParams)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMessage, msg)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_UnfavoriteArticle_0_072cd5d64b
ROOST_METHOD_SIG_HASH=request_Articles_UnfavoriteArticle_0_f4f0b9e771


 */
func Testrequest_Articles_UnfavoriteArticle_0(t *testing.T) {
	type testCase struct {
		description string
		setup       func() (ArticlesClient, *http.Request, map[string]string)
		assert      func(t *testing.T, msg proto.Message, md runtime.ServerMetadata, err error)
	}

	tests := []testCase{
		{
			description: "Successfully Unfavorite an Article",
			setup: func() (ArticlesClient, *http.Request, map[string]string) {
				client := new(MockArticlesClient)
				client.On("UnfavoriteArticle", mock.Anything, &UnfavoriteArticleRequest{Slug: "test-slug"}).Return(&UnfavoriteArticleResponse{}, nil)

				pathParams := map[string]string{"slug": "test-slug"}
				req := &http.Request{}

				return client, req, pathParams
			},
			assert: func(t *testing.T, msg proto.Message, md runtime.ServerMetadata, err error) {
				assert.NotNil(t, msg)
				assert.NoError(t, err)
			},
		},
		{
			description: "Missing Slug Parameter",
			setup: func() (ArticlesClient, *http.Request, map[string]string) {
				client := new(MockArticlesClient)
				pathParams := map[string]string{}
				req := &http.Request{}

				return client, req, pathParams
			},
			assert: func(t *testing.T, msg proto.Message, md runtime.ServerMetadata, err error) {
				assert.Nil(t, msg)
				assert.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			description: "Invalid Slug Type Conversion",
			setup: func() (ArticlesClient, *http.Request, map[string]string) {
				client := new(MockArticlesClient)
				pathParams := map[string]string{"slug": "123"}
				req := &http.Request{}

				return client, req, pathParams
			},
			assert: func(t *testing.T, msg proto.Message, md runtime.ServerMetadata, err error) {
				assert.Nil(t, msg)
				assert.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			description: "GRPC Client Error Handling",
			setup: func() (ArticlesClient, *http.Request, map[string]string) {
				client := new(MockArticlesClient)
				mockError := status.Errorf(codes.Internal, "internal error")
				client.On("UnfavoriteArticle", mock.Anything, &UnfavoriteArticleRequest{Slug: "test-slug"}).Return(nil, mockError)

				pathParams := map[string]string{"slug": "test-slug"}
				req := &http.Request{}

				return client, req, pathParams
			},
			assert: func(t *testing.T, msg proto.Message, md runtime.ServerMetadata, err error) {
				assert.Nil(t, msg)
				assert.Error(t, err)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			description: "Metadata Handling Verification",
			setup: func() (ArticlesClient, *http.Request, map[string]string) {
				client := new(MockArticlesClient)
				mdHeader := metadata.New(map[string]string{"key": "value"})
				mockResponse := &UnfavoriteArticleResponse{}
				client.On("UnfavoriteArticle", mock.Anything, &UnfavoriteArticleRequest{Slug: "test-slug"}, grpc.Header(&mdHeader), grpc.Trailer(nil)).Return(mockResponse, nil)

				pathParams := map[string]string{"slug": "test-slug"}
				req := &http.Request{}

				return client, req, pathParams
			},
			assert: func(t *testing.T, msg proto.Message, md runtime.ServerMetadata, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				assert.Equal(t, "value", md.HeaderMD["key"][0])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			client, req, pathParams := tc.setup()
			msg, md, err := request_Articles_UnfavoriteArticle_0(context.Background(), nil, client, req, pathParams)
			tc.assert(t, msg, md, err)
		})
	}
}

func (m *MockArticlesClient) UnfavoriteArticle(ctx context.Context, in *UnfavoriteArticleRequest, opts ...grpc.CallOption) (proto.Message, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(proto.Message), args.Error(1)
}

/*
ROOST_METHOD_HASH=local_request_Articles_CreateComment_0_85228c065a
ROOST_METHOD_SIG_HASH=local_request_Articles_CreateComment_0_ff06be8431


 */
func (m *mockArticlesServer) CreateComment(ctx context.Context, req *proto.CreateCommentRequest) (proto.Message, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.CommentResponse, nil
}

func Testlocal_request_Articles_CreateComment_0(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		slug           string
		expectError    bool
		expectedErrMsg string
		expectedCode   codes.Code
		setupMock      func() *mockArticlesServer
	}{
		{
			name: "Valid Comment is Created Successfully",
			body: `{"content": "This is a comment"}`,
			slug: "valid-slug",
			setupMock: func() *mockArticlesServer {
				return &mockArticlesServer{
					CommentResponse: &proto.Comment{},
				}
			},
		},
		{
			name:           "Missing 'slug' Parameter",
			body:           `{"content": "This is a comment"}`,
			setupMock:      func() *mockArticlesServer { return &mockArticlesServer{} },
			expectError:    true,
			expectedErrMsg: "missing parameter slug",
			expectedCode:   codes.InvalidArgument,
		},
		{
			name: "Invalid JSON Body",
			body: `{"content": "This is a comment"`,
			slug: "valid-slug",
			setupMock: func() *mockArticlesServer {
				return &mockArticlesServer{}
			},
			expectError:    true,
			expectedErrMsg: "invalid character",
			expectedCode:   codes.InvalidArgument,
		},
		{
			name:           "Invalid 'slug' Type",
			body:           `{"content": "This is a comment"}`,
			slug:           `{"malformed_slug"}`,
			setupMock:      func() *mockArticlesServer { return &mockArticlesServer{} },
			expectError:    true,
			expectedErrMsg: "type mismatch",
			expectedCode:   codes.InvalidArgument,
		},
		{
			name: "Server-Side Error on Comment Creation",
			body: `{"content": "This is a comment"}`,
			slug: "valid-slug",
			setupMock: func() *mockArticlesServer {
				return &mockArticlesServer{
					Err: status.Errorf(codes.Internal, "internal error"),
				}
			},
			expectError:    true,
			expectedErrMsg: "internal error",
			expectedCode:   codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := tt.setupMock()

			req, err := http.NewRequest("POST", "/articles/:slug/comments", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			pathParams := map[string]string{
				"slug": tt.slug,
			}

			marshaler := &runtime.JSONPb{}

			_, _, err = proto.Local_request_Articles_CreateComment_0(context.Background(), marshaler, mockServer, req, pathParams)

			if tt.expectError {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error with message containing '%v', got '%v'", tt.expectedErrMsg, err)
				}
				st, ok := status.FromError(err)
				if !ok || st.Code() != tt.expectedCode {
					t.Errorf("Expected code %v, got %v", tt.expectedCode, err)
				}
			} else {
				if err != nil {
					t.Errorf("Didn't expect an error, but got '%v'", err)
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_FavoriteArticle_0_ee68322d37
ROOST_METHOD_SIG_HASH=local_request_Articles_FavoriteArticle_0_bd14de1048


 */
func (m *MockArticlesServer) FavoriteArticle(ctx context.Context, req *proto_package.FavoriteArticleRequest) (proto.Message, error) {
	if req.Slug == "errSlug" {
		return nil, errors.New("internal server error")
	}
	return &proto_package.FavoriteArticleResponse{}, nil
}

func Testlocal_request_Articles_FavoriteArticle_0(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		pathParams     map[string]string
		expectedErr    error
		serverResponse proto.Message
	}{
		{
			name:           "Successful Article Favoriting",
			body:           `{"userId": "123"}`,
			pathParams:     map[string]string{"slug": "some-valid-slug"},
			expectedErr:    nil,
			serverResponse: &proto_package.FavoriteArticleResponse{},
		},
		{
			name:        "Missing slug parameter",
			body:        `{"userId": "123"}`,
			pathParams:  map[string]string{},
			expectedErr: status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug"),
		},
		{
			name:        "Invalid slug parameter type",
			body:        `{"userId": "123"}`,
			pathParams:  map[string]string{"slug": "bad_slug!"},
			expectedErr: status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", "invalid slug parameter"),
		},
		{
			name:        "Decoding error in request body",
			body:        `{"userId": 123}`,
			pathParams:  map[string]string{"slug": "validslug"},
			expectedErr: status.Errorf(codes.InvalidArgument, "%v", errors.New("EOF")),
		},
		{
			name:           "FavoriteArticle method returns an error",
			body:           `{"userId": "123"}`,
			pathParams:     map[string]string{"slug": "errSlug"},
			expectedErr:    errors.New("internal server error"),
			serverResponse: nil,
		},
		{
			name:        "Proper handling of request reader error",
			body:        "",
			pathParams:  map[string]string{"slug": "validslug"},
			expectedErr: status.Errorf(codes.InvalidArgument, "%v", errors.New("EOF")),
		},
	}

	mockServer := &MockArticlesServer{}
	marshaler := &runtime.JSONPb{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := createHTTPRequest(tc.body)
			resp, _, err := local_request_Articles_FavoriteArticle_0(context.Background(), marshaler, mockServer, req, tc.pathParams)

			if tc.expectedErr != nil && !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.expectedErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tc.serverResponse != nil && !proto.Equal(resp, tc.serverResponse) {
				t.Errorf("expected response %v, got %v", tc.serverResponse, resp)
			}

			t.Logf("Test %s completed successfully", tc.name)
		})
	}
}

func createHTTPRequest(body string) *http.Request {
	return &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(body)),
	}
}

/*
ROOST_METHOD_HASH=local_request_Articles_UpdateArticle_0_b3c1e567fb
ROOST_METHOD_SIG_HASH=local_request_Articles_UpdateArticle_0_d339fbf2dc


 */
func Testlocal_request_Articles_UpdateArticle_0(t *testing.T) {

	marshaler := &runtime.JSONPb{
		EmitDefaults: true,
		OrigName:     true,
	}

	tests := []struct {
		name           string
		body           io.Reader
		pathParams     map[string]string
		mockResponse   proto.Message
		mockError      error
		expectedError  error
		expectedOutput proto.Message
	}{
		{
			name:           "Scenario 1: Test Normal Operation with Valid Parameters",
			body:           strings.NewReader(`{"slug": "test-slug", "title": "Test Title", "content": "Test Content"}`),
			pathParams:     map[string]string{"article.slug": "test-slug"},
			mockResponse:   &UpdateArticleResponse{},
			expectedOutput: &UpdateArticleResponse{},
		},
		{
			name:          "Scenario 2: Missing Path Parameter",
			body:          strings.NewReader(`{"title": "Test Title", "content": "Test Content"}`),
			pathParams:    map[string]string{},
			expectedError: status.Errorf(codes.InvalidArgument, "missing parameter %s", "article.slug"),
		},
		{
			name:          "Scenario 3: Malformed Request Body",
			body:          strings.NewReader(`malformed body`),
			pathParams:    map[string]string{"article.slug": "test-slug"},
			expectedError: status.Errorf(codes.InvalidArgument, "unknown wire type %d", 4),
		},
		{
			name:          "Scenario 4: Type Mismatch in Path Parameter",
			body:          strings.NewReader(`{"title": "Test Title", "content": "Test Content"}`),
			pathParams:    map[string]string{"article.slug": "123"},
			expectedError: status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "article.slug", errors.New("expected string")),
		},
		{
			name:          "Scenario 5: Error Propagation from Server Update",
			body:          strings.NewReader(`{"slug": "test-slug", "title": "Test Title", "content": "Test Content"}`),
			pathParams:    map[string]string{"article.slug": "test-slug"},
			mockError:     errors.New("update failed"),
			expectedError: errors.New("update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := &mockArticlesServer{
				mockResponse: tt.mockResponse,
				mockError:    tt.mockError,
			}

			req, err := http.NewRequest("POST", "/articles", tt.body)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, _, err := local_request_Articles_UpdateArticle_0(context.Background(), marshaler, server, req, tt.pathParams)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !proto.Equal(tt.expectedOutput, resp) {
				t.Errorf("Expected response: %v, got: %v", tt.expectedOutput, resp)
			}
		})
	}
}

func (m *mockArticlesServer) UpdateArticle(ctx context.Context, req *UpdateArticleRequest) (proto.Message, error) {
	return m.mockResponse, m.mockError
}

/*
ROOST_METHOD_HASH=request_Articles_GetArticles_0_41358b839b
ROOST_METHOD_SIG_HASH=request_Articles_GetArticles_0_0078253459


 */
func (m *MockArticlesClient) EXPECT() *MockArticlesClientMockRecorder {
	return m.recorder
}

func (m *MockArticlesClient) GetArticles(ctx context.Context, in *GetArticlesRequest, opts ...grpc.CallOption) (*GetArticlesResponse, error) {
	m.ctrl.T.Helper()
	args := m.ctrl.Call(m, "GetArticles", ctx, in, opts...)
	return args[0].(*GetArticlesResponse), args[1].(error)
}

func NewMockArticlesClient(ctrl *gomock.Controller) *MockArticlesClient {
	mock := &MockArticlesClient{ctrl: ctrl}
	mock.recorder = &MockArticlesClientMockRecorder{mock}
	return mock
}

func Testrequest_Articles_GetArticles_0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name             string
		queryParams      url.Values
		setupClientMock  func(*MockArticlesClient)
		expectedErr      error
		expectedResponse *GetArticlesResponse
	}{
		{
			name:        "Scenario 1: Valid Request with No Query Parameters",
			queryParams: url.Values{},
			setupClientMock: func(client *MockArticlesClient) {
				client.EXPECT().GetArticles(gomock.Any(), &GetArticlesRequest{}, gomock.Any()).Return(&GetArticlesResponse{}, nil)
			},
			expectedErr:      nil,
			expectedResponse: &GetArticlesResponse{},
		},
		{
			name:        "Scenario 2: Valid Request with Query Parameters",
			queryParams: url.Values{"author": []string{"JohnDoe"}, "tag": []string{"Go"}},
			setupClientMock: func(client *MockArticlesClient) {
				client.EXPECT().GetArticles(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *GetArticlesRequest, opts ...grpc.CallOption) (*GetArticlesResponse, error) {

						return &GetArticlesResponse{}, nil
					})
			},
			expectedErr:      nil,
			expectedResponse: &GetArticlesResponse{},
		},
		{
			name:            "Scenario 3: Malformed Request (Invalid Query Parameters)",
			queryParams:     url.Values{"invalidParam": []string{"%%%"}},
			setupClientMock: func(client *MockArticlesClient) {},
			expectedErr:     status.Errorf(codes.InvalidArgument, "%v", "invalid argument error"),
		},
		{
			name:        "Scenario 4: Error Handling When Client Fails",
			queryParams: url.Values{},
			setupClientMock: func(client *MockArticlesClient) {
				client.EXPECT().GetArticles(gomock.Any(), &GetArticlesRequest{}, gomock.Any()).Return(nil, status.Errorf(codes.Internal, "Internal error"))
			},
			expectedErr: status.Errorf(codes.Internal, "Internal error"),
		},
		{
			name:        "Scenario 5: Complex Query with Multiple Parameters",
			queryParams: url.Values{"author": []string{"JohnDoe"}, "tag": []string{"Go"}, "limit": []string{"10"}},
			setupClientMock: func(client *MockArticlesClient) {
				client.EXPECT().GetArticles(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *GetArticlesRequest, opts ...grpc.CallOption) (*GetArticlesResponse, error) {

						return &GetArticlesResponse{}, nil
					})
			},
			expectedErr:      nil,
			expectedResponse: &GetArticlesResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockArticlesClient(ctrl)
			tt.setupClientMock(client)

			req := httptest.NewRequest("GET", "http://example.com/articles", nil)
			req.URL.RawQuery = tt.queryParams.Encode()

			res, _, err := request_Articles_GetArticles_0(context.Background(), nil, client, req, nil)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResponse, res)
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_GetFeedArticles_0_4ca93aa5f8
ROOST_METHOD_SIG_HASH=request_Articles_GetFeedArticles_0_8f9fb504a6


 */
func (m *MockArticlesClient) GetFeedArticles(ctx context.Context, in *GetFeedArticlesRequest, opts ...grpc.CallOption) (*GetFeedArticlesResponse, error) {

	return nil, status.Error(codes.Internal, "internal error")
}

func Testrequest_Articles_GetFeedArticles_0(t *testing.T) {
	type testCase struct {
		name        string
		context     context.Context
		request     *http.Request
		mockClient  func() *MockArticlesClient
		expectedErr error
		expectedRes proto.Message
	}

	tests := []testCase{
		{
			name:    "Scenario 1: Successfully retrieve feed articles",
			context: context.Background(),
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles/feed", nil)
				req.Form = url.Values{}
				return req
			}(),
			mockClient: func() *MockArticlesClient {
				mock := &MockArticlesClient{}

				return mock
			},
			expectedErr: nil,
			expectedRes: &GetFeedArticlesResponse{},
		},
		{
			name:    "Scenario 2: Handle malformed HTTP request",
			context: context.Background(),
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles/feed", nil)

				req.Body = ioutil.NopCloser(bytes.NewBufferString("invalid form data"))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			}(),
			mockClient: func() *MockArticlesClient {
				return &MockArticlesClient{}
			},
			expectedErr: status.Errorf(codes.InvalidArgument, "missing form body"),
			expectedRes: nil,
		},
		{
			name:    "Scenario 3: Error in populating query parameters",
			context: context.Background(),
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles/feed", nil)
				req.Form = url.Values{"unknown": {"param"}}
				return req
			}(),
			mockClient: func() *MockArticlesClient {
				return &MockArticlesClient{}
			},
			expectedErr: status.Errorf(codes.InvalidArgument, "unable to populate query parameters"),
			expectedRes: nil,
		},
		{
			name:    "Scenario 4: gRPC client returns error",
			context: context.Background(),
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles/feed", nil)
				req.Form = url.Values{}
				return req
			}(),
			mockClient: func() *MockArticlesClient {
				mock := &MockArticlesClient{}

				return mock
			},
			expectedErr: status.Error(codes.Internal, "internal error"),
			expectedRes: nil,
		},
		{
			name:    "Scenario 5: Empty response from the gRPC client",
			context: context.Background(),
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/articles/feed", nil)
				req.Form = url.Values{}
				return req
			}(),
			mockClient: func() *MockArticlesClient {
				mock := &MockArticlesClient{}

				return mock
			},
			expectedErr: nil,
			expectedRes: &GetFeedArticlesResponse{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockClient := tc.mockClient()
			marshaler := &runtime.JSONPb{}
			pathParams := map[string]string{}

			res, _, err := request_Articles_GetFeedArticles_0(tc.context, marshaler, mockClient, tc.request, pathParams)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRes, res)

			if err != nil {
				t.Logf("Test case '%s' returned expected error: %v", tc.name, err)
			} else {
				t.Logf("Test case '%s' succeeding with expected response: %v", tc.name, res)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_CreateArticle_0_49003c8cfd
ROOST_METHOD_SIG_HASH=request_Articles_CreateArticle_0_108af334a7


 */
func Testrequest_Articles_CreateArticle_0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockArticlesClient(ctrl)

	type test struct {
		name          string
		reqBody       string
		mockSetup     func()
		expectedCode  codes.Code
		expectedError bool
	}

	tests := []test{
		{
			name:    "Successfully Create an Article",
			reqBody: `{"title":"Effective Go","content":"Go is an open source programming language."}`,
			mockSetup: func() {
				mockClient.EXPECT().CreateArticle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			expectedError: false,
		},
		{
			name:          "Handle Invalid Input Data",
			reqBody:       `{"title":123,"content":true}`,
			mockSetup:     func() {},
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:    "Client Fails to Create Article",
			reqBody: `{"title":"Failed Article","content":"This should fail."}`,
			mockSetup: func() {
				mockClient.EXPECT().CreateArticle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.Internal, "internal error"))
			},
			expectedCode:  codes.Internal,
			expectedError: true,
		},
		{
			name:          "Handle Empty Request Body",
			reqBody:       ``,
			mockSetup:     func() {},
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:    "Graceful Handling of EOF Error",
			reqBody: `{"title":"Valid Title","content":"Some content."}`,
			mockSetup: func() {
				mockClient.EXPECT().CreateArticle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			ctx := context.Background()
			marshaler := &runtime.JSONPb{}
			req := &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(tc.reqBody)),
			}

			resp, _, err := request_Articles_CreateArticle_0(ctx, marshaler, mockClient, req, nil)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error status %v, got %v", tc.expectedError, err)
			}

			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("expected grpc status error, got %v", err)
				} else if st.Code() != tc.expectedCode {
					t.Errorf("expected status code %v, got %v", tc.expectedCode, st.Code())
				}
			}

			if !tc.expectedError && resp == nil {
				t.Errorf("expected non-nil response on success scenario")
			}

			t.Logf("Test '%s' passed", tc.name)
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_CreateComment_0_bf83e78f76
ROOST_METHOD_SIG_HASH=request_Articles_CreateComment_0_213a1a5a4b


 */
func (m *mockArticlesClient) CreateArticle(ctx context.Context, req *CreateCommentRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

func (m *mockArticlesClient) CreateComment(ctx context.Context, req *CreateCommentRequest, opts ...grpc.CallOption) (proto.Message, error) {
	if m.createCommentFunc != nil {
		return m.createCommentFunc(ctx, req, opts...)
	}
	return nil, nil
}

func Testrequest_Articles_CreateComment_0(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		pathParams     map[string]string
		setupMock      func() mockArticlesClient
		expectedMsgNil bool
		expectErr      error
	}{
		{
			name:       "Successful Comment Creation",
			body:       `{"comment": {"body": "Great article!"}}`,
			pathParams: map[string]string{"slug": "test-slug"},
			setupMock: func() mockArticlesClient {
				return mockArticlesClient{
					createCommentFunc: func(ctx context.Context, req *CreateCommentRequest, opts ...grpc.CallOption) (proto.Message, error) {
						return &CreateCommentResponse{}, nil
					},
				}
			},
			expectedMsgNil: false,
			expectErr:      nil,
		},
		{
			name:           "Missing Path Parameter 'slug'",
			body:           `{"comment": {"body": "Great article!"}}`,
			pathParams:     map[string]string{},
			setupMock:      func() mockArticlesClient { return mockArticlesClient{} },
			expectedMsgNil: true,
			expectErr:      status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug"),
		},
		{
			name:           "Invalid Request Body",
			body:           `{"comment": "invalid"}`,
			pathParams:     map[string]string{"slug": "test-slug"},
			setupMock:      func() mockArticlesClient { return mockArticlesClient{} },
			expectedMsgNil: true,
			expectErr:      status.Errorf(codes.InvalidArgument, "unexpected end of JSON input"),
		},
		{
			name:           "Malformed Path Parameter 'slug'",
			body:           `{"comment": {"body": "Nice!"}}`,
			pathParams:     map[string]string{"slug": ""},
			setupMock:      func() mockArticlesClient { return mockArticlesClient{} },
			expectedMsgNil: true,
			expectErr:      status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", errors.New("invalid slug")),
		},
		{
			name:       "Client Fails to Create Comment",
			body:       `{"comment": {"body": "Great article!"}}`,
			pathParams: map[string]string{"slug": "test-slug"},
			setupMock: func() mockArticlesClient {
				return mockArticlesClient{
					createCommentFunc: func(ctx context.Context, req *CreateCommentRequest, opts ...grpc.CallOption) (proto.Message, error) {
						return nil, status.Errorf(codes.Internal, "internal error")
					},
				}
			},
			expectedMsgNil: true,
			expectErr:      status.Errorf(codes.Internal, "internal error"),
		},
		{
			name:       "EOF Behavior on Decoding",
			body:       "",
			pathParams: map[string]string{"slug": "test-slug"},
			setupMock: func() mockArticlesClient {
				return mockArticlesClient{
					createCommentFunc: func(ctx context.Context, req *CreateCommentRequest, opts ...grpc.CallOption) (proto.Message, error) {
						return &CreateCommentResponse{}, nil
					},
				}
			},
			expectedMsgNil: false,
			expectErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupMock()
			req, err := http.NewRequest(http.MethodPost, "/fake-url", bytes.NewReader([]byte(tt.body)))
			if err != nil {
				t.Fatalf("failed to create HTTP request: %v", err)
			}

			marshaler := &runtime.JSONPb{}
			resp, _, err := request_Articles_CreateComment_0(context.Background(), marshaler, &client, req, tt.pathParams)

			if !tt.expectedMsgNil && resp == nil {
				t.Errorf("expected non-nil response; got nil")
			}
			if tt.expectedMsgNil && resp != nil {
				t.Errorf("expected nil response; got non-nil")
			}

			if tt.expectErr != nil && err == nil {
				t.Errorf("expected error %v; got none", tt.expectErr)
			}
			if tt.expectErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectErr != nil && err != nil && tt.expectErr.Error() != err.Error() {
				t.Errorf("expected error %v; got %v", tt.expectErr, err)
			}

			t.Logf("Test case '%s' executed successfully", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_FavoriteArticle_0_b513096596
ROOST_METHOD_SIG_HASH=request_Articles_FavoriteArticle_0_f1000633c4


 */
func (m *MockArticlesClient) FavoriteArticle(ctx context.Context, req *FavoriteArticleRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return m.FavoriteResult[req.Slug], m.FavoriteError
}

func Testrequest_Articles_FavoriteArticle_0(t *testing.T) {
	marshaler := &runtime.JSONPb{}

	tests := []struct {
		name        string
		reqBody     string
		pathParams  map[string]string
		mockClient  *MockArticlesClient
		wantMessage proto.Message
		wantError   error
		wantMeta    runtime.ServerMetadata
	}{
		{
			name:    "Scenario 1: Valid Request Leading to Successful Favoriting of an Article",
			reqBody: `{"key":"value"}`,
			pathParams: map[string]string{
				"slug": "valid-slug",
			},
			mockClient: &MockArticlesClient{
				FavoriteResult: map[string]proto.Message{"valid-slug": proto.MessageV1(&FavoriteArticleRequest{})},
				FavoriteError:  nil,
			},
			wantMessage: proto.MessageV1(&FavoriteArticleRequest{}),
			wantError:   nil,
			wantMeta:    runtime.ServerMetadata{},
		},
		{
			name:        "Scenario 2: Missing 'slug' Parameter Leading to Error",
			reqBody:     `{"key":"value"}`,
			pathParams:  map[string]string{},
			mockClient:  &MockArticlesClient{},
			wantMessage: nil,
			wantError:   status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug"),
			wantMeta:    runtime.ServerMetadata{},
		},
		{
			name:        "Scenario 3: Invalid JSON Body in Request Causing Decoding Error",
			reqBody:     `{"invalid json"`,
			pathParams:  map[string]string{"slug": "valid-slug"},
			mockClient:  &MockArticlesClient{},
			wantMessage: nil,
			wantError:   status.Errorf(codes.InvalidArgument, "unexpected EOF"),
			wantMeta:    runtime.ServerMetadata{},
		},
		{
			name:        "Scenario 4: Type Mismatch in 'slug' Parameter Causing Error",
			reqBody:     `{"key":"value"}`,
			pathParams:  map[string]string{"slug": "not-a-string"},
			mockClient:  &MockArticlesClient{},
			wantMessage: nil,
			wantError:   status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", errors.New("some type mismatch error")),
			wantMeta:    runtime.ServerMetadata{},
		},
		{
			name:        "Scenario 5: Handling of Unexpected IO Error During Request Body Read",
			reqBody:     "",
			pathParams:  map[string]string{"slug": "valid-slug"},
			mockClient:  &MockArticlesClient{},
			wantMessage: nil,
			wantError:   status.Errorf(codes.InvalidArgument, "some io error"),
			wantMeta:    runtime.ServerMetadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Body: io.NopCloser(strings.NewReader(tt.reqBody)),
			}

			msg, meta, err := request_Articles_FavoriteArticle_0(context.Background(), marshaler, tt.mockClient, req, tt.pathParams)

			if !proto.Equal(msg, tt.wantMessage) {
				t.Errorf("expected message: %v, got: %v", tt.wantMessage, msg)
			}

			if status.Code(err) != status.Code(tt.wantError) {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}

			if !proto.Equal(&meta, &tt.wantMeta) {
				t.Errorf("expected metadata: %v, got: %v", tt.wantMeta, meta)
			}

			t.Logf("Test: %s - completed", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=request_Articles_UpdateArticle_0_79099719da
ROOST_METHOD_SIG_HASH=request_Articles_UpdateArticle_0_05fecbeaf1


 */
func Testrequest_Articles_UpdateArticle_0(t *testing.T) {
	type testCase struct {
		name           string
		body           string
		pathParams     map[string]string
		expectedError  string
		expectedOutput proto.Message
	}

	tests := []testCase{
		{
			name: "Scenario 1: Successful Update of an Article",
			body: `{"title": "Updated Title"}`,
			pathParams: map[string]string{
				"article.slug": "valid-slug",
			},
			expectedError:  "",
			expectedOutput: &Article{},
		},
		{
			name:           "Scenario 2: Missing Path Parameter",
			body:           `{"title": "Updated Title"}`,
			pathParams:     map[string]string{},
			expectedError:  "missing parameter article.slug",
			expectedOutput: nil,
		},
		{
			name: "Scenario 3: Invalid JSON Body",
			body: `{"title": "Updated Title",`,
			pathParams: map[string]string{
				"article.slug": "valid-slug",
			},
			expectedError:  "invalid character",
			expectedOutput: nil,
		},
		{
			name: "Scenario 4: Path Parameter Type Mismatch",
			body: `{"title": "Updated Title"}`,
			pathParams: map[string]string{
				"article.slug": "123",
			},
			expectedError:  "type mismatch",
			expectedOutput: nil,
		},
		{
			name: "Scenario 5: Client RPC Call Fails",
			body: `{"title": "Updated Title"}`,
			pathParams: map[string]string{
				"article.slug": "known-error",
			},
			expectedError:  "RPC call failed",
			expectedOutput: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &mockArticlesClient{}
			req := &http.Request{
				Body: io.NopCloser(strings.NewReader(test.body)),
			}
			marshaler := &runtime.JSONPb{}
			ctx := context.Background()

			msg, _, err := request_Articles_UpdateArticle_0(ctx, marshaler, client, req, test.pathParams)
			if err != nil && !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("expected error to contain %q, got %q", test.expectedError, err.Error())
			} else if err == nil && test.expectedError != "" {
				t.Errorf("expected error %q, got nil", test.expectedError)
			}

			if !proto.Equal(msg, test.expectedOutput) {
				t.Errorf("expected output %+v, got %+v", test.expectedOutput, msg)
			}

			t.Logf("Test %s: completed", test.name)
		})
	}
}

func (m *mockArticlesClient) UpdateArticle(ctx context.Context, req *UpdateArticleRequest, opts ...grpc.CallOption) (proto.Message, error) {
	if req.Slug == "known-error" {
		return nil, status.Errorf(codes.Internal, "RPC call failed")
	}
	return &Article{}, nil
}

/*
ROOST_METHOD_HASH=RegisterArticlesHandlerClient_42c9b8bd7f
ROOST_METHOD_SIG_HASH=RegisterArticlesHandlerClient_ac2c37fb8d


 */
func (m *mockArticlesClient) CreateArticle(ctx context.Context, in *CreateAritcleRequest, opts ...grpc.CallOption) (*CreateArticleResponse, error) {
	return &CreateArticleResponse{}, nil
}

func (m *mockArticlesClient) GetTags(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*TagsResponse, error) {
	return &TagsResponse{}, nil
}

func TestRegisterArticlesHandlerClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mux := runtime.NewServeMux()
	client := &mockArticlesClient{}

	type testCase struct {
		name       string
		method     string
		urlPattern string
		expect     bool
	}

	testCases := []testCase{
		{
			name:       "Scenario 1: Registration of Create Article Handler",
			method:     "POST",
			urlPattern: "/articles",
			expect:     true,
		},
		{
			name:       "Scenario 3: Handling of Update Article Request",
			method:     "PUT",
			urlPattern: "/articles/{article.slug}",
			expect:     true,
		},
		{
			name:       "Scenario 4: Non-standard HTTP Methods Handling",
			method:     "PATCH",
			urlPattern: "/articles",
			expect:     false,
		},
		{
			name:       "Scenario 5: Graceful Shutdown Context Handling",
			method:     "GET",
			urlPattern: "/tags",
			expect:     true,
		},
	}

	ctx := context.Background()
	err := RegisterArticlesHandlerClient(ctx, mux, client)
	assert.NoError(t, err, "RegisterArticlesHandlerClient should not return an error")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req, _ := http.NewRequest(tc.method, tc.urlPattern, bytes.NewReader([]byte("")))
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if tc.expect {
				assert.NotEqual(t, http.StatusNotFound, rr.Code, "Handler should be found")
				t.Logf("%s handler for %s method is successfully registered.", tc.name, tc.method)
			} else {
				assert.Equal(t, http.StatusNotFound, rr.Code, "Unexpected handler found registered.")
				t.Logf("%s method correctly not registered for unexpected HTTP method.", tc.method)
			}
		})
	}
}

func (m *mockArticlesClient) UpdateArticle(ctx context.Context, in *UpdateArticleRequest, opts ...grpc.CallOption) (*UpdateArticleResponse, error) {
	return &UpdateArticleResponse{}, nil
}

/*
ROOST_METHOD_HASH=RegisterArticlesHandlerServer_ce9ae7b704
ROOST_METHOD_SIG_HASH=RegisterArticlesHandlerServer_e6ac2d5c16


 */
func (s *MockArticlesServer) CreateArticle(ctx context.Context, req *CreateAritcleRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) CreateComment(ctx context.Context, req *CreateCommentRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) DeleteArticle(ctx context.Context, req *DeleteArticleRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) DeleteComment(ctx context.Context, req *DeleteCommentRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) FavoriteArticle(ctx context.Context, req *FavoriteArticleRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) GetArticle(ctx context.Context, req *GetArticleRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) GetArticles(ctx context.Context, req *GetArticlesRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) GetComments(ctx context.Context, req *GetCommentsRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) GetFeedArticles(ctx context.Context, req *GetFeedArticlesRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) GetTags(ctx context.Context, req *Empty) (proto.Message, error) {
	return nil, nil
}

func TestRegisterArticlesHandlerServer(t *testing.T) {
	type testCase struct {
		name        string
		mux         *runtime.ServeMux
		server      ArticlesServer
		context     context.Context
		expectedErr codes.Code
	}

	testCases := []testCase{
		{
			name:    "Success Registration",
			mux:     runtime.NewServeMux(),
			server:  &MockArticlesServer{},
			context: context.Background(),
		},
		{
			name:        "Nil Context Error",
			mux:         runtime.NewServeMux(),
			server:      &MockArticlesServer{},
			context:     nil,
			expectedErr: codes.Canceled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RegisterArticlesHandlerServer(tc.context, tc.mux, tc.server)

			if tc.expectedErr != codes.OK {
				if err == nil || status.Code(err) != tc.expectedErr {
					t.Errorf("expected error: %v, but got: %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				validateRegisteredRoutes(t, tc.mux)
			}
		})
	}

	t.Run("Concurrent Registration", func(t *testing.T) {
		var wg sync.WaitGroup
		t.Log("Testing concurrent registration of handlers")
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mux := runtime.NewServeMux()
				ctx := context.Background()
				server := &MockArticlesServer{}

				err := RegisterArticlesHandlerServer(ctx, mux, server)
				if err != nil {
					t.Error(err)
				}

				validateRegisteredRoutes(t, mux)
			}()
		}
		wg.Wait()
	})
}

func (s *MockArticlesServer) UnfavoriteArticle(ctx context.Context, req *UnfavoriteArticleRequest) (proto.Message, error) {
	return nil, nil
}

func (s *MockArticlesServer) UpdateArticle(ctx context.Context, req *UpdateArticleRequest) (proto.Message, error) {
	return nil, nil
}

func findHandlerForMethod(mux *runtime.ServeMux, method string, path string) (http.Handler, bool) {

	for _, h := range mux.handlers[method] {

		if h.op == path {
			return h, true
		}
	}
	return nil, false
}

func validateRegisteredRoutes(t *testing.T, mux *runtime.ServeMux) {
	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/articles"},
		{"GET", "/articles/feed"},
		{"GET", "/articles/{slug}"},
		{"GET", "/articles"},
		{"PUT", "/articles/{article.slug}"},
		{"DELETE", "/articles/{slug}"},
		{"POST", "/articles/{slug}/favorite"},
		{"DELETE", "/articles/{slug}/favorite"},
		{"GET", "/tags"},
		{"POST", "/articles/{slug}/comments"},
		{"GET", "/articles/{slug}/comments"},
		{"GET", "/articles/{slug}/comments/{id}"},
	}

	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		rr := httptest.NewRecorder()

		handler, ok := findHandlerForMethod(mux, route.method, route.path)
		if !ok {
			t.Errorf("Handler for %s not found", route.path)
			continue
		}
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status == http.StatusNotFound {
			t.Errorf("Handler for %s returned 404", route.path)
		}
	}
}

