package proto

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
	"github.com/golang/mock/gomock"
	"time"
	"google.golang.org/grpc/metadata"
	"bytes"
	"fmt"
	"encoding/json"
	"io"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"io/ioutil"
	"strings"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/example/project/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/raahii/golang-grpc-realworld-example/proto/mocks"
)

var _ UsersClient = &MockUsersClient{}
type Empty struct{}
type MockUsersServer struct {
	mock.Mock
}
type SomeUserMessage struct {
	proto.Message
}
type UsersClient interface{}
type MockUsersServer struct {
	mockCtrl *gomock.Controller

	mock *MockProtoUsersServer
}
type UsersServer interface {
	UnfollowUser(ctx context.Context, req *UnfollowRequest) (proto.Message, error)
}
type mockProtoMessage struct{}
type mockUsersServer struct{}
type Empty struct{}
type MockUsersClient struct {
	mock.Mock
}
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

type CreateUserRequest struct {
	Name string
}
type errorReader struct{}
type mockUsersServer struct{}
type MockUsersServer struct {
	mock.Mock
}
type UsersServer interface {
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (proto.Message, error)
}
type mockDecoder struct {
	reader io.Reader
}
type mockMarshaler struct{}
type mockUsersServer struct{}
type MockUsersClient struct {
	mock.Mock
}
type ShowProfileRequest struct {
	Username string
}
type ShowProfileResponse struct{}
type UsersClient interface {
	ShowProfile(ctx context.Context, in *ShowProfileRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type CreateUserRequest struct {
	Name string
}
type FollowRequest struct {
	Username string
}
type mockUsersClient struct {
	unfollowUser func(ctx context.Context, in *UnfollowRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type FollowRequest struct {
	Username string
}
type FollowResponse struct {
	Success bool
}
type UsersServer interface {
	FollowUser(ctx context.Context, req *FollowRequest) (proto.Message, error)
}
type mockUsersServer struct{}
type CreateUserRequest struct {
	Name string
}
type CreateUserResponse struct {
	Id string
}
type UsersClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type mockUsersClient struct {
	unfollowUser func(ctx context.Context, in *UnfollowRequest, opts ...grpc.CallOption) (proto.Message, error)
}
type MockedUsersClient struct {
	mock.Mock
}
type MockUsersClient struct {
	mock.Mock
}
type UpdateUserResponse struct {
	Message string
}
type FollowRequest struct {
	Username string
}
type MockUsersClient struct {
	mock.Mock
}
type CreateUserRequest struct {
	Name string
}
type MockUsersServer struct {
	mock.Mock
}
type ShowProfileRequest struct {
	Username string
}
/*
ROOST_METHOD_HASH=local_request_Users_CurrentUser_0_fb61f27ac5
ROOST_METHOD_SIG_HASH=local_request_Users_CurrentUser_0_80f8e1f0b6


 */
func (m *MockUsersServer) CurrentUser(ctx context.Context, req *Empty) (proto.Message, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(proto.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func Testlocal_request_Users_CurrentUser_0(t *testing.T) {
	tests := []struct {
		name           string
		mockServerFunc func() *MockUsersServer
		ctx            context.Context
		expectedMsg    proto.Message
		expectedErr    error
	}{
		{
			name: "Successfully Retrieve Current User",
			mockServerFunc: func() *MockUsersServer {
				mockServer := &MockUsersServer{}
				mockUser := &SomeUserMessage{}
				mockServer.On("CurrentUser", mock.Anything, mock.Anything).Return(mockUser, nil)
				return mockServer
			},
			ctx:         context.Background(),
			expectedMsg: &SomeUserMessage{},
			expectedErr: nil,
		},
		{
			name: "Handle Server Error",
			mockServerFunc: func() *MockUsersServer {
				mockServer := &MockUsersServer{}
				mockServer.On("CurrentUser", mock.Anything, mock.Anything).Return(nil, errors.New("server error"))
				return mockServer
			},
			ctx:         context.Background(),
			expectedMsg: nil,
			expectedErr: errors.New("server error"),
		},
		{
			name: "Empty Context",
			mockServerFunc: func() *MockUsersServer {
				mockServer := &MockUsersServer{}
				mockUser := &SomeUserMessage{}
				mockServer.On("CurrentUser", context.TODO(), mock.Anything).Return(mockUser, nil)
				return mockServer
			},
			ctx:         context.TODO(),
			expectedMsg: &SomeUserMessage{},
			expectedErr: nil,
		},
		{
			name: "User Not Found",
			mockServerFunc: func() *MockUsersServer {
				mockServer := &MockUsersServer{}
				mockServer.On("CurrentUser", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user not found"))
				return mockServer
			},
			ctx:         context.Background(),
			expectedMsg: nil,
			expectedErr: status.Error(codes.NotFound, "user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := tt.mockServerFunc()
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			msg, _, err := local_request_Users_CurrentUser_0(tt.ctx, nil, mockServer, req, nil)

			assert.Equal(t, tt.expectedMsg, msg)
			assert.Equal(t, tt.expectedErr, err)

		})
	}
}

/*
ROOST_METHOD_HASH=RegisterUsersHandler_b6576ec644
ROOST_METHOD_SIG_HASH=RegisterUsersHandler_6b766bd753


 */
func TestRegisterUsersHandler(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		mux     *runtime.ServeMux
		conn    *grpc.ClientConn
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Scenario 1: Successful Registration of User Handler",
			ctx:     context.Background(),
			mux:     runtime.NewServeMux(),
			conn:    &grpc.ClientConn{},
			wantErr: false,
			errMsg:  "Handler registration should succeed with no errors",
		},
		{
			name:    "Scenario 2: Error Handling with Nil Context",
			ctx:     nil,
			mux:     runtime.NewServeMux(),
			conn:    &grpc.ClientConn{},
			wantErr: true,
			errMsg:  "Function should return an error due to nil context",
		},
		{
			name:    "Scenario 3: Error Handling with Nil ServeMux",
			ctx:     context.Background(),
			mux:     nil,
			conn:    &grpc.ClientConn{},
			wantErr: true,
			errMsg:  "Function should return an error due to nil ServeMux",
		},
		{
			name:    "Scenario 4: Error Handling with Nil ClientConn",
			ctx:     context.Background(),
			mux:     runtime.NewServeMux(),
			conn:    nil,
			wantErr: true,
			errMsg:  "Function should return an error due to nil ClientConn",
		},
		{
			name:    "Scenario 5: Checking Forward Response Options Integration",
			ctx:     context.Background(),
			mux:     runtime.NewServeMux(),
			conn:    &grpc.ClientConn{},
			wantErr: false,
			errMsg:  "Forward response options should be appropriately initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterUsersHandler(tt.ctx, tt.mux, tt.conn)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error: %v, wantErr: %v - %s", err, tt.wantErr, tt.errMsg)
			}

			if tt.name == "Scenario 5: Checking Forward Response Options Integration" && !tt.wantErr {
				if tt.mux == nil || tt.mux.GetForwardResponseOptions == nil || len(tt.mux.GetForwardResponseOptions()) == 0 {
					t.Errorf("Forward response options are not set as expected - %s", tt.errMsg)
				}
			}

			t.Logf("Test %s completed.", tt.name)
		})
	}

}

/*
ROOST_METHOD_HASH=local_request_Users_ShowProfile_0_7119d6fcac
ROOST_METHOD_SIG_HASH=local_request_Users_ShowProfile_0_59264256bc


 */
func Testlocal_request_Users_ShowProfile_0(t *testing.T) {
	tests := []struct {
		name         string
		pathParams   map[string]string
		setupMock    func(mock *MockProtoUsersServer)
		expectedMsg  proto.Message
		expectedErr  error
		expectedCode codes.Code
	}{
		{
			name: "Successfully show a user's profile",
			pathParams: map[string]string{
				"username": "testuser",
			},
			setupMock: func(mock *MockProtoUsersServer) {
				mock.EXPECT().
					ShowProfile(gomock.Any(), &ShowProfileRequest{Username: "testuser"}).
					Return(&ProfileResponse{}, nil)
			},
			expectedMsg:  &ProfileResponse{},
			expectedErr:  nil,
			expectedCode: codes.OK,
		},
		{
			name:         "Missing username parameter",
			pathParams:   map[string]string{},
			setupMock:    nil,
			expectedMsg:  nil,
			expectedErr:  status.Error(codes.InvalidArgument, "missing parameter username"),
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Invalid username parameter type",
			pathParams: map[string]string{
				"username": "123",
			},
			setupMock:    nil,
			expectedMsg:  nil,
			expectedErr:  status.Error(codes.InvalidArgument, "type mismatch, parameter: username, error: <specify error>"),
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Server error during ShowProfile call",
			pathParams: map[string]string{
				"username": "testuser",
			},
			setupMock: func(mock *MockProtoUsersServer) {
				mock.EXPECT().
					ShowProfile(gomock.Any(), &ShowProfileRequest{Username: "testuser"}).
					Return(nil, errors.New("internal server error"))
			},
			expectedMsg:  nil,
			expectedErr:  errors.New("internal server error"),
			expectedCode: codes.Internal,
		},
		{
			name: "Username parameter leads to non-existing profile",
			pathParams: map[string]string{
				"username": "unknownuser",
			},
			setupMock: func(mock *MockProtoUsersServer) {
				mock.EXPECT().
					ShowProfile(gomock.Any(), &ShowProfileRequest{Username: "unknownuser"}).
					Return(nil, status.Errorf(codes.NotFound, "profile not found"))
			},
			expectedMsg:  nil,
			expectedErr:  status.Error(codes.NotFound, "profile not found"),
			expectedCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockUsersServer(t)
			defer mockServer.mockCtrl.Finish()

			if tt.setupMock != nil {
				tt.setupMock(mockServer.mock)
			}

			req, err := http.NewRequest("GET", "http://example.com", nil)
			assert.NoError(t, err)

			msg, metadata, err := local_request_Users_ShowProfile_0(context.Background(), runtime.JSONPb{}, mockServer.mock, req, tt.pathParams)

			assert.Equal(t, tt.expectedMsg, msg)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newMockUsersServer(t *testing.T) *MockUsersServer {
	mockCtrl := gomock.NewController(t)
	mock := NewMockProtoUsersServer(mockCtrl)
	return &MockUsersServer{mockCtrl: mockCtrl, mock: mock}
}

/*
ROOST_METHOD_HASH=local_request_Users_UnfollowUser_0_2e77d364c1
ROOST_METHOD_SIG_HASH=local_request_Users_UnfollowUser_0_784b246e21


 */
func Testlocal_request_Users_UnfollowUser_0(t *testing.T) {
	type scenario struct {
		name        string
		username    string
		serverSetup func() *mockUsersServer
		expectError codes.Code
	}

	scenarios := []scenario{
		{
			name:     "Scenario 1: Successful Unfollow User",
			username: "validUser",
			serverSetup: func() *mockUsersServer {
				return &mockUsersServer{
					UnfollowUserFunc: func(ctx context.Context, req *UnfollowRequest) (proto.Message, error) {
						if req.Username == "validUser" {
							return &mockProtoMessage{}, nil
						}
						return nil, errors.New("user not found")
					},
				}
			},
			expectError: codes.OK,
		},
		{
			name:     "Scenario 2: Missing Username Parameter",
			username: "",
			serverSetup: func() *mockUsersServer {
				return &mockUsersServer{}
			},
			expectError: codes.InvalidArgument,
		},
		{
			name:     "Scenario 3: Invalid Username Type",
			username: "invalid!user*name",
			serverSetup: func() *mockUsersServer {
				return &mockUsersServer{}
			},
			expectError: codes.InvalidArgument,
		},
		{
			name:     "Scenario 4: Server Error During Unfollow Operation",
			username: "serverErrorUser",
			serverSetup: func() *mockUsersServer {
				return &mockUsersServer{
					UnfollowUserFunc: func(ctx context.Context, req *UnfollowRequest) (proto.Message, error) {
						return nil, status.Errorf(codes.Internal, "internal server error")
					},
				}
			},
			expectError: codes.Internal,
		},
		{
			name:     "Scenario 5: Network Context Cancellation",
			username: "cancellableUser",
			serverSetup: func() *mockUsersServer {
				return &mockUsersServer{
					UnfollowUserFunc: func(ctx context.Context, req *UnfollowRequest) (proto.Message, error) {
						time.Sleep(2 * time.Second)
						return nil, nil
					},
				}
			},
			expectError: codes.Canceled,
		},
		{
			name:     "Scenario 6: Unauthorized Access Attempt",
			username: "unauthorizedUser",
			serverSetup: func() *mockUsersServer {
				return &mockUsersServer{
					UnfollowUserFunc: func(ctx context.Context, req *UnfollowRequest) (proto.Message, error) {
						return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
					},
				}
			},
			expectError: codes.PermissionDenied,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			server := s.serverSetup()
			pathParams := make(map[string]string)
			if s.username != "" {
				pathParams["username"] = s.username
			}

			ctx := context.Background()
			if s.name == "Scenario 5: Network Context Cancellation" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				go func() {
					time.Sleep(1 * time.Second)
					cancel()
				}()
			}

			_, _, err := local_request_Users_UnfollowUser_0(ctx, &runtime.JSONPb{OrigName: true, EmitDefaults: true}, server, &http.Request{}, pathParams)

			if status.Code(err) != s.expectError {
				t.Errorf("expected error code %v, got %v, error: %v", s.expectError, status.Code(err), err)
			}

			t.Logf("Test scenario '%s' completed.\n", s.name)
		})
	}
}

func (m *mockUsersServer) UnfollowUser(ctx context.Context, req *UnfollowRequest) (proto.Message, error) {
	return m.UnfollowUserFunc(ctx, req)
}

/*
ROOST_METHOD_HASH=request_Users_CurrentUser_0_07909c811a
ROOST_METHOD_SIG_HASH=request_Users_CurrentUser_0_39d07ecd48


 */
func (m *MockUsersClient) CurrentUser(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error) {
	mdOut := &metadata.MD{}
	for _, opt := range opts {
		switch o := opt.(type) {
		case grpc.HeaderCallOption:
			mdOut = o.HeaderAddr
			*mdOut = metadata.Pairs("key", "value")
		}
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return &Empty{}, nil
}

func (m *MockUsersClient) CurrentUserError(ctx context.Context, in *Empty, opts ...grpc.CallOption) (proto.Message, error) {
	mdOut := &metadata.MD{}
	for _, opt := range opts {
		switch o := opt.(type) {
		case grpc.HeaderCallOption:
			mdOut = o.HeaderAddr
			*mdOut = metadata.Pairs("key", "value")
		}
	}
	return nil, status.Error(codes.Internal, "internal error")
}

func Testrequest_Users_CurrentUser_0(t *testing.T) {
	type args struct {
		ctx        context.Context
		marshaler  runtime.Marshaler
		client     UsersClient
		req        *http.Request
		pathParams map[string]string
	}
	type test struct {
		name       string
		args       args
		want       proto.Message
		wantErr    bool
		wantErrMsg string
	}

	mockSuccessClient := &MockUsersClient{}
	mockErrorClient := &MockUsersClient{}
	ctxSuccess := context.Background()
	ctxCanceled, cancel := context.WithCancel(context.Background())
	cancel()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()

	tests := []test{
		{
			name: "Scenario 1: Successful Retrieval of Current User",
			args: args{
				ctx:       ctxSuccess,
				client:    mockSuccessClient,
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
			},
			wantErr: false,
		},
		{
			name: "Scenario 2: Handling Server Error from CurrentUser Call",
			args: args{
				ctx:       ctxSuccess,
				client:    mockErrorClient,
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
			},
			wantErr:    true,
			wantErrMsg: "internal error",
		},
		{
			name: "Scenario 3: Invalid Context Handling",
			args: args{
				ctx:       ctxCanceled,
				client:    mockSuccessClient,
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
			},
			wantErr:    true,
			wantErrMsg: context.Canceled.Error(),
		},
		{
			name: "Scenario 4: Simulating Network Timeout",
			args: args{
				ctx:       ctxTimeout,
				client:    mockSuccessClient,
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
			},
			wantErr:    true,
			wantErrMsg: context.DeadlineExceeded.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := request_Users_CurrentUser_0(tt.args.ctx, tt.args.marshaler, tt.args.client, tt.args.req, tt.args.pathParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("request_Users_CurrentUser_0() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("Expected error message: %v, got: %v", tt.wantErrMsg, err.Error())
			}
			if !tt.wantErr {
				t.Logf("Received message: %v", got)
			} else {
				t.Logf("Received expected error: %v", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=RegisterUsersHandlerFromEndpoint_eb7c518bba
ROOST_METHOD_SIG_HASH=RegisterUsersHandlerFromEndpoint_1e496f2047


 */
func TestRegisterUsersHandlerFromEndpoint(t *testing.T) {
	type testScenario struct {
		description   string
		endpoint      string
		dialOptions   []grpc.DialOption
		mux           *runtime.ServeMux
		expectedError error
		validateFunc  func(*testing.T, *runtime.ServeMux, error)
	}

	mockGRPCServerEndpoint := "localhost:9090"

	tests := []testScenario{
		{
			description:  "Successful Connection and Registration",
			endpoint:     mockGRPCServerEndpoint,
			dialOptions:  []grpc.DialOption{grpc.WithInsecure()},
			mux:          new(runtime.ServeMux),
			validateFunc: validateSuccessfulRegistration,
		},
		{
			description: "Failed Connection to Endpoint",
			endpoint:    "invalid-endpoint",
			dialOptions: []grpc.DialOption{grpc.WithInsecure()},
			expectedError: status.Errorf(codes.Unavailable,
				"connection error"),
			mux:          new(runtime.ServeMux),
			validateFunc: validateConnectionError,
		},
		{
			description:  "Context Cancellation Handling",
			endpoint:     mockGRPCServerEndpoint,
			dialOptions:  []grpc.DialOption{grpc.WithInsecure()},
			mux:          new(runtime.ServeMux),
			validateFunc: validateContextCancellation,
		},
		{
			description:  "Connection Closure Failure",
			endpoint:     mockGRPCServerEndpoint,
			dialOptions:  []grpc.DialOption{grpc.WithInsecure()},
			mux:          new(runtime.ServeMux),
			validateFunc: validateClosureFailure,
		},
		{
			description:  "Checking Incorrect ServeMux Behavior",
			endpoint:     mockGRPCServerEndpoint,
			dialOptions:  []grpc.DialOption{grpc.WithInsecure()},
			mux:          new(runtime.ServeMux),
			validateFunc: validateIncorrectServeMuxBehavior,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Test output: %s\n", tt.description)

			err := RegisterUsersHandlerFromEndpoint(ctx, tt.mux, tt.endpoint, tt.dialOptions)

			if tt.validateFunc != nil {
				tt.validateFunc(t, tt.mux, err)
			}
		})
	}
}

func validateClosureFailure(t *testing.T, mux *runtime.ServeMux, err error) {
	if err != nil {
		t.Logf("Closure failure handled and logged: %v", err)
	} else {
		t.Log("Handled closure failure with no errors")
	}
}

func validateConnectionError(t *testing.T, mux *runtime.ServeMux, err error) {
	if err != nil && status.Code(err) != codes.Unavailable {
		t.Errorf("Expected unavailable error but got: %v", err)
	}
	t.Log("Correctly handled failed connection")
}

func validateContextCancellation(t *testing.T, mux *runtime.ServeMux, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Log("Managed context cancellation appropriately")
}

func validateIncorrectServeMuxBehavior(t *testing.T, mux *runtime.ServeMux, err error) {
	if err != nil {
		t.Logf("Handled incorrect ServeMux behavior with error: %v", err)
	} else {
		t.Log("Managed incorrect ServeMux behavior gracefully")
	}
}

func validateSuccessfulRegistration(t *testing.T, mux *runtime.ServeMux, err error) {
	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	}
	t.Log("Successful registration of handlers")
}

/*
ROOST_METHOD_HASH=local_request_Users_CreateUser_0_16c39336bf
ROOST_METHOD_SIG_HASH=local_request_Users_CreateUser_0_393813fbf0


 */
func (m *mockUsersServer) CreateUser(ctx context.Context, req *CreateUserRequest) (proto.Message, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, req)
	}
	return nil, status.Errorf(codes.Unimplemented, "CreateUser not implemented in mock")
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func Testlocal_request_Users_CreateUser_0(t *testing.T) {
	mockMarshaler := &runtime.JSONPb{}

	t.Run("Scenario 1: Successful User Creation", func(t *testing.T) {
		mockServer := &mockUsersServer{
			CreateUserFunc: func(ctx context.Context, req *CreateUserRequest) (proto.Message, error) {
				return &CreateUserResponse{UserId: "123"}, nil
			},
		}

		requestBody := `{"username": "testuser", "email": "test@example.com"}`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(requestBody))

		msg, _, err := local_request_Users_CreateUser_0(context.Background(), mockMarshaler, mockServer, req, nil)

		assert.NoError(t, err, "Expected no error during successful user creation")
		response, ok := msg.(*CreateUserResponse)
		assert.True(t, ok, "Expected message to be of type CreateUserResponse")
		assert.Equal(t, response.UserId, "123", "Expected UserID to match mock response")
	})

	t.Run("Scenario 2: Invalid JSON Request Body", func(t *testing.T) {
		mockServer := &mockUsersServer{}

		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid-json"))

		_, _, err := local_request_Users_CreateUser_0(context.Background(), mockMarshaler, mockServer, req, nil)

		assert.Error(t, err, "Expected error for invalid JSON")
		s, _ := status.FromError(err)
		assert.Equal(t, s.Code(), codes.InvalidArgument, "Expected InvalidArgument error for malformed JSON")
	})

	t.Run("Scenario 3: Server CreateUser Returns Error", func(t *testing.T) {
		mockServer := &mockUsersServer{
			CreateUserFunc: func(ctx context.Context, req *CreateUserRequest) (proto.Message, error) {
				return nil, status.Errorf(codes.Internal, "Internal Server Error")
			},
		}

		requestBody := `{"username": "testuser", "email": "test@example.com"}`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(requestBody))

		_, _, err := local_request_Users_CreateUser_0(context.Background(), mockMarshaler, mockServer, req, nil)

		assert.Error(t, err, "Expected error from server.CreateUser")
		s, _ := status.FromError(err)
		assert.Equal(t, s.Code(), codes.Internal, "Expected Internal error from server call")
	})

	t.Run("Scenario 4: IOReaderFactory Error Handling", func(t *testing.T) {
		mockServer := &mockUsersServer{}

		req, _ := http.NewRequest("POST", "/users", &errorReader{})

		_, _, err := local_request_Users_CreateUser_0(context.Background(), mockMarshaler, mockServer, req, nil)

		assert.Error(t, err, "Expected error from IOReaderFactory")
		s, _ := status.FromError(err)
		assert.Equal(t, s.Code(), codes.InvalidArgument, "Expected InvalidArgument error from IOReaderFactory")
	})

	t.Run("Scenario 5: EOF Encountered with Decoder", func(t *testing.T) {
		mockServer := &mockUsersServer{
			CreateUserFunc: func(ctx context.Context, req *CreateUserRequest) (proto.Message, error) {
				return &CreateUserResponse{UserId: "123"}, nil
			},
		}

		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(`{}`))

		msg, _, err := local_request_Users_CreateUser_0(context.Background(), mockMarshaler, mockServer, req, nil)

		assert.NoError(t, err, "Expect no error with empty JSON due to EOF")
		response, ok := msg.(*CreateUserResponse)
		assert.True(t, ok, "Expected message to be of type CreateUserResponse")
		assert.Equal(t, response.UserId, "123", "Expected UserID to match mock response for EOF")
	})
}

/*
ROOST_METHOD_HASH=local_request_Users_LoginUser_0_395d307a0e
ROOST_METHOD_SIG_HASH=local_request_Users_LoginUser_0_887c3218bd


 */
func (m *MockUsersServer) LoginUser(ctx context.Context, req *LoginUserRequest) (proto.Message, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(proto.Message), args.Error(1)
}

func Testlocal_request_Users_LoginUser_0(t *testing.T) {
	t.Run("Scenario 1: Successful User Login", func(t *testing.T) {
		server := new(MockUsersServer)
		marshaler := &runtime.JSONPb{}

		validReqBody := &LoginUserRequest{
			Username: "validuser",
			Password: "validpass",
		}
		bodyBytes, _ := json.Marshal(validReqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))

		expectedResponse := &LoginResponse{UserId: "123", Token: "token"}
		server.On("LoginUser", mock.Anything, validReqBody).Return(expectedResponse, nil)

		resp, _, err := local_request_Users_LoginUser_0(context.Background(), marshaler, server, req, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Errorf("expected a response, got nil")
		}
		t.Log("Scenario 1: Successful login executed correctly with valid response and no errors.")
	})

	t.Run("Scenario 2: Invalid Request Body", func(t *testing.T) {
		server := new(MockUsersServer)
		marshaler := &runtime.JSONPb{}

		req, _ := http.NewRequest("POST", "/login", strings.NewReader("{invalid-json}"))

		resp, _, err := local_request_Users_LoginUser_0(context.Background(), marshaler, server, req, nil)

		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("expected InvalidArgument error, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Scenario 2: Invalid request body handled with proper error response.")
	})

	t.Run("Scenario 3: Empty Request Body", func(t *testing.T) {
		server := new(MockUsersServer)
		marshaler := &runtime.JSONPb{}

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(""))

		_, _, err := local_request_Users_LoginUser_0(context.Background(), marshaler, server, req, nil)

		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("expected InvalidArgument error, got %v", err)
		}
		t.Log("Scenario 3: Empty request body error handled correctly.")
	})

	t.Run("Scenario 4: Successful Login No Metadata", func(t *testing.T) {
		server := new(MockUsersServer)
		marshaler := &runtime.JSONPb{}

		reqData := &LoginUserRequest{
			Username: "validuser",
			Password: "validpass",
		}
		bodyBytes, _ := json.Marshal(reqData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))

		expectedResponse := &LoginResponse{UserId: "123", Token: "token"}
		server.On("LoginUser", mock.Anything, reqData).Return(expectedResponse, nil)

		resp, _, err := local_request_Users_LoginUser_0(context.Background(), marshaler, server, req, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Errorf("expected a response, got nil")
		}
		t.Log("Scenario 4: Function operated without any metadata successfully.")
	})

	t.Run("Scenario 5: Error During Decoding", func(t *testing.T) {
		server := new(MockUsersServer)
		marshaler := &runtime.JSONPb{}

		utilities.IOReaderFactory = func(r io.Reader) (func() io.Reader, error) {
			return nil, errors.New("decoding failed")
		}

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(""))

		_, _, err := local_request_Users_LoginUser_0(context.Background(), marshaler, server, req, nil)

		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("expected InvalidArgument error, got %v", err)
		}
		t.Log("Scenario 5: Decoding error properly handled.")
	})

	t.Run("Scenario 6: Authentication Failure", func(t *testing.T) {
		server := new(MockUsersServer)
		marshaler := &runtime.JSONPb{}

		reqData := &LoginUserRequest{
			Username: "invaliduser",
			Password: "invalidpass",
		}
		bodyBytes, _ := json.Marshal(reqData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))

		server.On("LoginUser", mock.Anything, reqData).Return(nil, status.Error(codes.Unauthenticated, "authentication failed"))

		resp, _, err := local_request_Users_LoginUser_0(context.Background(), marshaler, server, req, nil)

		if status.Code(err) != codes.Unauthenticated {
			t.Errorf("expected Unauthenticated error, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
		t.Log("Scenario 6: Authentication failure handled correctly.")
	})
}

/*
ROOST_METHOD_HASH=local_request_Users_UpdateUser_0_98df809638
ROOST_METHOD_SIG_HASH=local_request_Users_UpdateUser_0_c03fd099fb


 */
func (d *mockDecoder) Decode(v interface{}) error {
	data, _ := ioutil.ReadAll(d.reader)
	if strings.TrimSpace(string(data)) == "invalid" {
		return errors.New("decoding error")
	}
	*v.(*UpdateUserRequest) = UpdateUserRequest{}
	return nil
}

func (m *mockMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &mockDecoder{reader: r}
}

func Testlocal_request_Users_UpdateUser_0(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (*http.Request, func() io.Reader)
		server    UsersServer
		expectErr bool
		errCode   codes.Code
	}{
		{
			name: "Successful Update of User Information",
			setup: func() (*http.Request, func() io.Reader) {
				req, _ := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader("{}")))
				return req, func() io.Reader { return bytes.NewReader([]byte("{}")) }
			},
			server:    &mockUsersServer{shouldFail: false},
			expectErr: false,
		},
		{
			name: "Handling Invalid HTTP Request Body",
			setup: func() (*http.Request, func() io.Reader) {
				return nil, nil
			},
			server:    &mockUsersServer{shouldFail: false},
			expectErr: true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "Marshaler Fails to Decode Request",
			setup: func() (*http.Request, func() io.Reader) {
				req, _ := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader("invalid")))
				return req, func() io.Reader { return bytes.NewReader([]byte("invalid")) }
			},
			server:    &mockUsersServer{shouldFail: false},
			expectErr: true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "Error from UsersServer UpdateUser",
			setup: func() (*http.Request, func() io.Reader) {
				req, _ := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader("{}")))
				return req, func() io.Reader { return bytes.NewReader([]byte("{}")) }
			},
			server:    &mockUsersServer{shouldFail: true},
			expectErr: true,
			errCode:   codes.Internal,
		},
		{
			name: "Handling Empty Map of Path Parameters",
			setup: func() (*http.Request, func() io.Reader) {
				req, _ := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader("{}")))
				return req, func() io.Reader { return bytes.NewReader([]byte("{}")) }
			},
			server:    &mockUsersServer{shouldFail: false},
			expectErr: false,
		},
		{
			name: "Network Context Cancellation During Execution",
			setup: func() (*http.Request, func() io.Reader) {
				req, _ := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader("{}")))
				return req, func() io.Reader { return bytes.NewReader([]byte("{}")) }
			},
			server:    &mockUsersServer{shouldFail: false},
			expectErr: true,
			errCode:   codes.Canceled,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, reader := tt.setup()

			originalIOReaderFactory := utilities.IOReaderFactory
			utilities.IOReaderFactory = func(_ io.Reader) (func() io.Reader, error) {
				return reader, nil
			}
			defer func() { utilities.IOReaderFactory = originalIOReaderFactory }()

			mockMarshaler := &mockMarshaler{}
			msg, _, err := local_request_Users_UpdateUser_0(ctx, mockMarshaler, tt.server, req, map[string]string{})

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err != nil && status.Code(err) != tt.errCode {
				t.Errorf("expected error code: %v, got: %v", tt.errCode, status.Code(err))
			}

			if err == nil && msg == nil {
				t.Error("expected non-nil message on success, got nil")
			}

			t.Logf("Test finished. Scenario: %s", tt.name)
		})
	}
}

func (m *mockUsersServer) UpdateUser(ctx context.Context, req *UpdateUserRequest) (proto.Message, error) {
	if m.shouldFail {
		return nil, status.Error(codes.Internal, "update error")
	}
	return &UpdateUserResponse{}, nil
}

/*
ROOST_METHOD_HASH=request_Users_ShowProfile_0_8a0ef792d8
ROOST_METHOD_SIG_HASH=request_Users_ShowProfile_0_a5b74123a3


 */
func (m *MockUsersClient) ShowProfile(ctx context.Context, in *ShowProfileRequest, opts ...grpc.CallOption) (proto.Message, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(proto.Message), args.Error(1)
}

func Testrequest_Users_ShowProfile_0_InvalidUsernameType(t *testing.T) {

	req, err := http.NewRequest("GET", "/profiles/invalidType", nil)
	assert.NoError(t, err)

	pathParams := map[string]string{"username": "invalidType"}
	ctx := context.Background()
	marshaler := &runtime.JSONPb{}

	oldString := runtime.String
	runtime.String = func(val string) (string, error) {
		return "", errors.New("type mismatch")
	}
	defer func() { runtime.String = oldString }()

	response, _, err := request_Users_ShowProfile_0(ctx, marshaler, new(MockUsersClient), req, pathParams)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.True(t, strings.Contains(err.Error(), "type mismatch, parameter: username"))
	t.Log("Scenario 3: Invalid Username Type Conversion passed.")
}

func Testrequest_Users_ShowProfile_0_MetadataPropagation(t *testing.T) {

	mockClient := new(MockUsersClient)
	expectedHeader := metadata.MD{"header-key": "header-value"}
	expectedTrailer := metadata.MD{"trailer-key": "trailer-value"}

	mockResponse := &ShowProfileResponse{}
	mockClient.On("ShowProfile", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil).Run(func(args mock.Arguments) {
		header := args.Get(2).(*metadata.MD)
		trailer := args.Get(3).(*metadata.MD)
		*header = expectedHeader
		*trailer = expectedTrailer
	})

	req, err := http.NewRequest("GET", "/profiles/validUsername", nil)
	assert.NoError(t, err)

	pathParams := map[string]string{"username": "validUsername"}
	ctx := context.Background()
	marshaler := &runtime.JSONPb{}

	response, respMetadata, err := request_Users_ShowProfile_0(ctx, marshaler, mockClient, req, pathParams)

	assert.NoError(t, err)
	assert.Equal(t, mockResponse, response)
	assert.Equal(t, expectedHeader, respMetadata.HeaderMD)
	assert.Equal(t, expectedTrailer, respMetadata.TrailerMD)
	t.Log("Scenario 5: Correct Metadata Propagation passed.")
}

func Testrequest_Users_ShowProfile_0_MissingUsername(t *testing.T) {

	req, err := http.NewRequest("GET", "/profiles/", nil)
	assert.NoError(t, err)

	pathParams := map[string]string{}
	ctx := context.Background()
	marshaler := &runtime.JSONPb{}

	response, _, err := request_Users_ShowProfile_0(ctx, marshaler, new(MockUsersClient), req, pathParams)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.True(t, strings.Contains(err.Error(), "missing parameter"))
	t.Log("Scenario 2: Missing Username Path Parameter passed.")
}

func Testrequest_Users_ShowProfile_0_ShowProfileError(t *testing.T) {

	mockClient := new(MockUsersClient)
	mockError := status.Errorf(codes.Internal, "internal error")
	mockClient.On("ShowProfile", mock.Anything, mock.Anything).Return(nil, mockError)

	req, err := http.NewRequest("GET", "/profiles/validUsername", nil)
	assert.NoError(t, err)

	pathParams := map[string]string{"username": "validUsername"}
	ctx := context.Background()
	marshaler := &runtime.JSONPb{}

	response, _, err := request_Users_ShowProfile_0(ctx, marshaler, mockClient, req, pathParams)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, mockError, err)
	t.Log("Scenario 4: Handling Error from the ShowProfile Call passed.")
}

func Testrequest_Users_ShowProfile_0_Success(t *testing.T) {

	mockClient := new(MockUsersClient)
	expectedResponse := &ShowProfileResponse{}
	mockClient.On("ShowProfile", mock.Anything, &ShowProfileRequest{Username: "validUsername"}).Return(expectedResponse, nil)

	req, err := http.NewRequest("GET", "/profiles/validUsername", nil)
	assert.NoError(t, err)

	pathParams := map[string]string{"username": "validUsername"}
	ctx := context.Background()
	marshaler := &runtime.JSONPb{}

	response, metadata, err := request_Users_ShowProfile_0(ctx, marshaler, mockClient, req, pathParams)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
	mockClient.AssertExpectations(t)
	t.Log("Scenario 1: Successful Profile Request for Valid Username passed.")
}

/*
ROOST_METHOD_HASH=request_Users_UnfollowUser_0_bad3e79511
ROOST_METHOD_SIG_HASH=request_Users_UnfollowUser_0_f1a5538d33


 */
func (m *mockUsersClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

func (m *mockUsersClient) CurrentUser(ctx context.Context, in *GetCurrentUserRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

func (m *mockUsersClient) FollowUser(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

func (m *mockUsersClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

func (m *mockUsersClient) ShowProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

func Testrequest_Users_UnfollowUser_0(t *testing.T) {
	type args struct {
		ctx        context.Context
		marshaler  runtime.Marshaler
		client     UsersClient
		req        *http.Request
		pathParams map[string]string
	}

	tests := []struct {
		name       string
		args       args
		wantErr    codes.Code
		setupMocks func() UsersClient
	}{
		{
			name: "Scenario 1: Successful Unfollow when all parameters are Correct",
			args: args{
				ctx:       context.Background(),
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
				pathParams: map[string]string{
					"username": "validUser",
				},
			},
			setupMocks: func() UsersClient {
				return &mockUsersClient{
					unfollowUser: func(ctx context.Context, in *UnfollowRequest, opts ...grpc.CallOption) (proto.Message, error) {
						assert.Equal(t, "validUser", in.Username)
						return &UnfollowResponse{}, nil
					},
				}
			},
			wantErr: codes.OK,
		},
		{
			name: "Scenario 2: Missing Username Parameter",
			args: args{
				ctx:        context.Background(),
				marshaler:  &runtime.JSONPb{},
				req:        &http.Request{},
				pathParams: map[string]string{},
			},
			setupMocks: func() UsersClient {
				return nil
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "Scenario 3: Type Mismatch for Username Parameter",
			args: args{
				ctx:       context.Background(),
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
				pathParams: map[string]string{
					"username": "[]",
				},
			},
			setupMocks: func() UsersClient {
				return nil
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "Scenario 4: UnfollowUser Client Method Fails",
			args: args{
				ctx:       context.Background(),
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
				pathParams: map[string]string{
					"username": "someUser",
				},
			},
			setupMocks: func() UsersClient {
				return &mockUsersClient{
					unfollowUser: func(ctx context.Context, in *UnfollowRequest, opts ...grpc.CallOption) (proto.Message, error) {
						return nil, errors.New("method failure")
					},
				}
			},
			wantErr: codes.Unknown,
		},
		{
			name: "Scenario 5: Invalid Path Parameters",
			args: args{
				ctx:       context.Background(),
				marshaler: &runtime.JSONPb{},
				req:       &http.Request{},
				pathParams: map[string]string{
					"wrongkey": "value",
				},
			},
			setupMocks: func() UsersClient {
				return nil
			},
			wantErr: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			fmt.Fprintf(output, "Executing test: %s\n", tt.name)
			client := tt.setupMocks()
			msg, md, err := request_Users_UnfollowUser_0(tt.args.ctx, tt.args.marshaler, client, tt.args.req, tt.args.pathParams)

			if err != nil {
				st, _ := status.FromError(err)
				if st.Code() != tt.wantErr {
					t.Errorf("expected error code %v, got %v", tt.wantErr, st.Code())
				} else {
					t.Logf("%s passed, error: %v", tt.name, err)
				}
			} else {
				t.Fatalf("Expected error, got success: msg=%v, md=%v", msg, md)
			}
		})
	}
}

func (m *mockUsersClient) UnfollowUser(ctx context.Context, in *UnfollowRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return m.unfollowUser(ctx, in, opts...)
}

func (m *mockUsersClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (proto.Message, error) {
	return nil, nil
}

/*
ROOST_METHOD_HASH=local_request_Users_FollowUser_0_5c53888ab3
ROOST_METHOD_SIG_HASH=local_request_Users_FollowUser_0_093a00b2f9


 */
func (m *mockUsersServer) FollowUser(ctx context.Context, req *FollowRequest) (proto.Message, error) {
	if req.Username == "errorUser" {
		return nil, status.Errorf(codes.Internal, "mock server error")
	}
	return &FollowResponse{Success: true}, nil
}

func Testlocal_request_Users_FollowUser_0(t *testing.T) {
	type args struct {
		ctx        context.Context
		marshaler  runtime.Marshaler
		server     UsersServer
		req        *http.Request
		pathParams map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		errCode    codes.Code
		mockDBFunc func() (sqlmock.Sqlmock, *testing.T)
	}{
		{
			name: "Successful Follow User Request",
			args: args{
				ctx:       context.TODO(),
				marshaler: &runtime.JSONPb{},
				server:    &mockUsersServer{},
				req: &http.Request{
					Body: ioutil.NopCloser(strings.NewReader(`{"username": "validUser"}`)),
				},
				pathParams: map[string]string{"username": "validUser"},
			},
			wantErr: false,
		},
		{
			name: "Missing Username Parameter",
			args: args{
				ctx:       context.TODO(),
				marshaler: &runtime.JSONPb{},
				server:    &mockUsersServer{},
				req: &http.Request{
					Body: ioutil.NopCloser(strings.NewReader(`{"dummy": "data"}`)),
				},
				pathParams: map[string]string{},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Invalid Username Type",
			args: args{
				ctx:       context.TODO(),
				marshaler: &runtime.JSONPb{},
				server:    &mockUsersServer{},
				req: &http.Request{
					Body: ioutil.NopCloser(strings.NewReader(`{"username": 123}`)),
				},
				pathParams: map[string]string{"username": "123"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Decode Failure from Request Body",
			args: args{
				ctx:       context.TODO(),
				marshaler: &runtime.JSONPb{},
				server:    &mockUsersServer{},
				req: &http.Request{
					Body: ioutil.NopCloser(strings.NewReader(`{invalid_json}`)),
				},
				pathParams: map[string]string{"username": "validUser"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Simulated Server Error during FollowUser Call",
			args: args{
				ctx:       context.TODO(),
				marshaler: &runtime.JSONPb{},
				server:    &mockUsersServer{},
				req: &http.Request{
					Body: ioutil.NopCloser(strings.NewReader(`{"username": "errorUser"}`)),
				},
				pathParams: map[string]string{"username": "errorUser"},
			},
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.mockDBFunc == nil {
				tt.mockDBFunc = func() (sqlmock.Sqlmock, *testing.T) {
					return nil, t
				}
			}
			mock, _ := tt.mockDBFunc()

			if mock != nil {
				defer mock.ExpectationsWereMet()
			}

			got, _, err := local_request_Users_FollowUser_0(tt.args.ctx, tt.args.marshaler, tt.args.server, tt.args.req, tt.args.pathParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("local_request_Users_FollowUser_0() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && status.Code(err) != tt.errCode {
				t.Errorf("Expected error code %v, got %v", tt.errCode, status.Code(err))
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Expected valid proto.Message but got nil")
			}
			t.Logf("Success: Test case %s passed", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=request_Users_CreateUser_0_597d6c6d58
ROOST_METHOD_SIG_HASH=request_Users_CreateUser_0_aa01923cc2


 */
func (m *mockUsersClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (proto.Message, error) {
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}
	if in.GetName() == "error" {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &CreateUserResponse{Id: "12345"}, nil
}

func (req *CreateUserRequest) GetName() string {
	return req.Name
}

func Testrequest_Users_CreateUser_0(t *testing.T) {
	tests := []struct {
		name             string
		reqBody          string
		mockError        error
		expectedError    codes.Code
		expectedResponse *CreateUserResponse
	}{
		{
			name:             "Successful User Creation",
			reqBody:          `{"name": "Valid User"}`,
			expectedError:    codes.OK,
			expectedResponse: &CreateUserResponse{Id: "12345"},
		},
		{
			name:          "Invalid JSON in Request Body",
			reqBody:       `{"name": "Invalid User"}`,
			expectedError: codes.InvalidArgument,
		},
		{
			name:          "Empty Request Body",
			reqBody:       ``,
			expectedError: codes.InvalidArgument,
		},
		{
			name:          "Client Error During User Creation",
			reqBody:       `{"name": "error"}`,
			expectedError: codes.Internal,
		},
		{
			name:          "Network Issues During Client Call",
			reqBody:       `{"name": "Valid User"}`,
			mockError:     status.Errorf(codes.Unavailable, "network issue"),
			expectedError: codes.Unavailable,
		},
		{
			name:          "Context Cancellation",
			reqBody:       `{"name": "Valid User"}`,
			mockError:     context.Canceled,
			expectedError: codes.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running scenario: %s", tt.name)
			ctx := context.Background()

			if tt.expectedError == codes.Canceled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			req := &http.Request{Body: ioutil.NopCloser(strings.NewReader(tt.reqBody))}
			marshaler := &runtime.JSONPb{}

			client := &mockUsersClient{}

			resp, _, err := request_Users_CreateUser_0(ctx, marshaler, client, req, nil)

			if tt.expectedError != codes.OK {
				assert.Error(t, err)
				if status.Code(err) != tt.expectedError {
					t.Errorf("Expected error code %v, got %v", tt.expectedError, status.Code(err))
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			}

			if tt.name == "Network Issues During Client Call" || tt.name == "Client Error During User Creation" {
				assert.Equal(t, tt.mockError, err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Users_LoginUser_0_333f425f0e
ROOST_METHOD_SIG_HASH=request_Users_LoginUser_0_be66a6b8e2


 */
func (m *MockedUsersClient) LoginUser(ctx context.Context, in *proto.LoginUserRequest, opts ...grpc.CallOption) (proto.Message, error) {
	args := m.Called(ctx, in)
	if response, ok := args.Get(0).(proto.Message); ok {
		return response, args.Error(1)
	}
	return nil, args.Error(1)
}

func Testrequest_Users_LoginUser_0(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    io.Reader
		mockSetup      func(client *MockedUsersClient)
		expectedError  codes.Code
		expectedResult proto.Message
	}{
		{
			name:        "Successful User Login Request",
			requestBody: bytes.NewBuffer([]byte(`{"username":"validuser","password":"validpass"}`)),
			mockSetup: func(client *MockedUsersClient) {
				client.On("LoginUser", mock.Anything, &proto.LoginUserRequest{Username: "validuser", Password: "validpass"}).
					Return(&proto.LoginUserResponse{Success: true}, nil)
			},
			expectedError:  codes.OK,
			expectedResult: &proto.LoginUserResponse{Success: true},
		},
		{
			name:          "Invalid Request Body",
			requestBody:   bytes.NewBuffer([]byte(`{"invalid-json"}`)),
			mockSetup:     func(client *MockedUsersClient) {},
			expectedError: codes.InvalidArgument,
		},
		{
			name:        "Unsuccessful User Login Due to Client Error",
			requestBody: bytes.NewBuffer([]byte(`{"username":"wronguser","password":"wrongpass"}`)),
			mockSetup: func(client *MockedUsersClient) {
				client.On("LoginUser", mock.Anything, &proto.LoginUserRequest{Username: "wronguser", Password: "wrongpass"}).
					Return(nil, status.Errorf(codes.Unauthenticated, "credentials are invalid"))
			},
			expectedError: codes.Unauthenticated,
		},
		{
			name:          "Empty Request Body",
			requestBody:   bytes.NewBuffer([]byte(``)),
			mockSetup:     func(client *MockedUsersClient) {},
			expectedError: codes.InvalidArgument,
		},
		{
			name:        "Network Error During User Login",
			requestBody: bytes.NewBuffer([]byte(`{"username":"user","password":"pass"}`)),
			mockSetup: func(client *MockedUsersClient) {
				client.On("LoginUser", mock.Anything, &proto.LoginUserRequest{Username: "user", Password: "pass"}).
					Return(nil, status.Errorf(codes.Unavailable, "network error"))
			},
			expectedError: codes.Unavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(MockedUsersClient)
			tt.mockSetup(client)

			req, err := http.NewRequest("POST", "/users/login", tt.requestBody)
			if err != nil {
				t.Fatalf("failed to create HTTP request: %v", err)
			}

			marshaler := &runtime.JSONPb{}
			response, _, err := proto.Request_Users_LoginUser_0(context.Background(), marshaler, client, req, nil)

			if response == nil && tt.expectedResult != nil {
				t.Errorf("expected result: %v, got: nil", tt.expectedResult)
			}

			if response != nil && !proto.Equal(response, tt.expectedResult) {
				t.Errorf("expected result: %v, got: %v", tt.expectedResult, response)
			}

			if err != nil {
				code := status.Code(err)
				if code != tt.expectedError {
					t.Errorf("expected error code: %v, got: %v", tt.expectedError, code)
				}
				t.Logf("received expected error: %v", err)
			} else if tt.expectedError != codes.OK {
				t.Errorf("expected error code: %v, got: nil", tt.expectedError)
			} else {
				t.Log("no error raised as expected")
			}
		})
	}
}

/*
ROOST_METHOD_HASH=request_Users_UpdateUser_0_a1eede20f0
ROOST_METHOD_SIG_HASH=request_Users_UpdateUser_0_3ed4d48da0


 */
func Testrequest_Users_UpdateUser_0(t *testing.T) {
	type testScenario struct {
		name      string
		body      string
		setup     func(*MockUsersClient)
		expected  proto.Message
		expectErr bool
	}

	scenarios := []testScenario{
		{
			name: "Successful User Update",
			body: `{"username":"john_doe", "email":"john@example.com"}`,
			setup: func(client *MockUsersClient) {
				msg := &UpdateUserResponse{Message: "User updated successfully"}
				client.On("UpdateUser", mock.Anything, mock.Anything).Return(msg, nil)
			},
			expected:  &UpdateUserResponse{Message: "User updated successfully"},
			expectErr: false,
		},
		{
			name:      "Invalid JSON Body Parsing Error",
			body:      `{"username":}`,
			setup:     func(client *MockUsersClient) {},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Client Error on Update Request",
			body: `{"username":"john_doe", "email":"john@example.com"}`,
			setup: func(client *MockUsersClient) {
				client.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, status.Error(codes.Internal, "client error"))
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Empty Request Body Handling",
			body:      `{}`,
			setup:     func(client *MockUsersClient) {},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			client := new(MockUsersClient)
			scenario.setup(client)

			req := &http.Request{
				Body: io.NopCloser(bytes.NewReader([]byte(scenario.body))),
			}
			marshaler := &runtime.JSONPb{}
			pathParams := map[string]string{}

			result, _, err := request_Users_UpdateUser_0(context.Background(), marshaler, client, req, pathParams)

			if scenario.expectErr {
				assert.Error(t, err, "Expected an error but none was returned")
			} else {
				assert.NoError(t, err, "Expected no error but an error was returned")
				assert.Equal(t, scenario.expected, result)
			}

			client.AssertExpectations(t)
		})
	}
}

func (m *MockUsersClient) UpdateUser(ctx context.Context, req *UpdateUserRequest, opts ...interface{}) (proto.Message, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(proto.Message), args.Error(1)
}

/*
ROOST_METHOD_HASH=request_Users_FollowUser_0_0f14b9ad3e
ROOST_METHOD_SIG_HASH=request_Users_FollowUser_0_98a8b098a0


 */
func (m *MockUsersClient) FollowUser(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (proto.Message, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(proto.Message), args.Error(1)
}

func Testrequest_Users_FollowUser_0(t *testing.T) {
	t.Run("Table-driven tests", func(t *testing.T) {
		tests := []struct {
			name         string
			mockResponse proto.Message
			mockError    error
			reqBody      string
			pathParams   map[string]string
			expectedErr  string
		}{
			{
				name:         "Valid Input and Successful Follow",
				mockResponse: &empty.Empty{},
				mockError:    nil,
				reqBody:      `{"some_key": "some_value"}`,
				pathParams:   map[string]string{"username": "validUser"},
				expectedErr:  "",
			},
			{
				name:         "Missing Username Parameter",
				mockResponse: nil,
				mockError:    status.Errorf(codes.InvalidArgument, "missing parameter %s", "username"),
				reqBody:      `{"some_key": "some_value"}`,
				pathParams:   map[string]string{},
				expectedErr:  "missing parameter username",
			},
			{
				name:         "Invalid JSON Body",
				mockResponse: nil,
				mockError:    status.Errorf(codes.InvalidArgument, "invalid json"),
				reqBody:      `{invalid json}`,
				pathParams:   map[string]string{"username": "validUser"},
				expectedErr:  "invalid json",
			},
			{
				name:         "Type Mismatch for Username Parameter",
				mockResponse: nil,
				mockError:    status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "username", "conversion failed"),
				reqBody:      `{"some_key": "some_value"}`,
				pathParams:   map[string]string{"username": "123"},
				expectedErr:  "type mismatch, parameter: username, error: conversion failed",
			},
			{
				name:         "GRPC Client Error Response",
				mockResponse: nil,
				mockError:    errors.New("grpc error"),
				reqBody:      `{"some_key": "some_value"}`,
				pathParams:   map[string]string{"username": "validUser"},
				expectedErr:  "grpc error",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockClient := new(MockUsersClient)
				ctx := context.Background()
				req := &http.Request{
					Body: ioutil.NopCloser(bytes.NewReader([]byte(tt.reqBody))),
				}
				marshaler := &runtime.JSONPb{}

				mockClient.On("FollowUser", ctx, mock.AnythingOfType("*proto.FollowRequest")).Return(tt.mockResponse, tt.mockError)

				resp, _, err := request_Users_FollowUser_0(ctx, marshaler, mockClient, req, tt.pathParams)

				if tt.expectedErr != "" {
					if err == nil || err.Error() != tt.expectedErr {
						t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
					}
				} else if err != nil {
					t.Errorf("Did not expect an error, but got %v", err)
				} else if resp == nil {
					t.Error("Expected a non-nil response")
				}

			})
		}
	})
}

/*
ROOST_METHOD_HASH=RegisterUsersHandlerClient_99d4372219
ROOST_METHOD_SIG_HASH=RegisterUsersHandlerClient_8f1226dca0


 */
func TestRegisterUsersHandlerClient(t *testing.T) {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	mockClient := new(mocks.UsersClient)

	tests := []struct {
		name       string
		method     string
		url        string
		body       io.Reader
		setupMock  func()
		expectFunc func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Scenario 1: Successful User Registration",
			method: "POST",
			url:    "/users",
			body:   bytes.NewBuffer([]byte(`{"user":{"email":"test@example.com","password":"testpass"}}`)),
			setupMock: func() {
				mockClient.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&CreateUserResponse{User: &User{Email: "test@example.com"}}, nil).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "test@example.com")
				mockClient.AssertExpectations(t)
			},
		},
		{
			name:   "Scenario 2: Handle User Login with Invalid Credentials",
			method: "POST",
			url:    "/users/login",
			body:   bytes.NewBuffer([]byte(`{"user":{"email":"fail@example.com","password":"wrongpass"}}`)),
			setupMock: func() {
				mockClient.On("LoginUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "Invalid credentials")
				mockClient.AssertExpectations(t)
			},
		},
		{
			name:   "Scenario 3: Successful Retrieval of Current User",
			method: "GET",
			url:    "/user",
			body:   nil,
			setupMock: func() {
				mockClient.On("CurrentUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&CurrentUserResponse{User: &User{Email: "current@example.com"}}, nil).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "current@example.com")
				mockClient.AssertExpectations(t)
			},
		},
		{
			name:   "Scenario 4: Update User Profile with Invalid Data",
			method: "PUT",
			url:    "/user",
			body:   bytes.NewBuffer([]byte(`{"user":{"invalid":"data"}}`)),
			setupMock: func() {
				mockClient.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.InvalidArgument, "Invalid data")).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "Invalid data")
				mockClient.AssertExpectations(t)
			},
		},
		{
			name:   "Scenario 5: Show User Profile Successfully",
			method: "GET",
			url:    "/profiles/testuser",
			body:   nil,
			setupMock: func() {
				mockClient.On("ShowProfile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&ShowProfileResponse{Profile: &Profile{Username: "testuser"}}, nil).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "testuser")
				mockClient.AssertExpectations(t)
			},
		},
		{
			name:   "Scenario 6: Handle Follow User Action",
			method: "POST",
			url:    "/profiles/testuser/follow",
			body:   nil,
			setupMock: func() {
				mockClient.On("FollowUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&FollowUserResponse{}, nil).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				mockClient.AssertExpectations(t)
			},
		},
		{
			name:   "Scenario 7: Error When User Unfollow Fails",
			method: "DELETE",
			url:    "/profiles/testuser/follow",
			body:   nil,
			setupMock: func() {
				mockClient.On("UnfollowUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.NotFound, "User not found")).Once()
			},
			expectFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "User not found")
				mockClient.AssertExpectations(t)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(tt.method, tt.url, tt.body)

			RegisterUsersHandlerClient(ctx, mux, mockClient)
			mux.ServeHTTP(recorder, request)

			tt.expectFunc(t, recorder)
		})
	}
}

/*
ROOST_METHOD_HASH=RegisterUsersHandlerServer_8960cae25c
ROOST_METHOD_SIG_HASH=RegisterUsersHandlerServer_3cfc375408


 */
func (m *MockUsersServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*UserResponse), args.Error(1)
}

func (m *MockUsersServer) LoginUser(ctx context.Context, req *LoginUserRequest) (*User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUsersServer) ShowProfile(ctx context.Context, req *ShowProfileRequest) (*Profile, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*Profile), args.Error(1)
}

func TestRegisterUsersHandlerServer(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		body       io.Reader
		setupMock  func(*MockUsersServer)
		assertResp func(*testing.T, *http.Response)
	}{
		{
			name:   "Successful Registration of User Handlers",
			method: "POST",
			url:    "/users/login",
			body:   bytes.NewBuffer([]byte(`{"email":"test@example.com", "password":"password123"}`)),
			setupMock: func(s *MockUsersServer) {
				s.On("LoginUser", mock.Anything, mock.AnythingOfType("*proto.LoginUserRequest")).
					Return(&User{}, nil)
			},
			assertResp: func(t *testing.T, resp *http.Response) {
				if resp.StatusCode != http.StatusOK {
					t.Errorf("expected status code %v, got %v", http.StatusOK, resp.StatusCode)
				}
			},
		},
		{
			name:   "Login User Handler with Invalid Credentials",
			method: "POST",
			url:    "/users/login",
			body:   bytes.NewBuffer([]byte(`{"email":"invalid@example.com", "password":"wrong"}`)),
			setupMock: func(s *MockUsersServer) {
				s.On("LoginUser", mock.Anything, mock.AnythingOfType("*proto.LoginUserRequest")).
					Return(nil, status.Error(codes.Unauthenticated, "invalid credentials"))
			},
			assertResp: func(t *testing.T, resp *http.Response) {
				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected status code %v, got %v", http.StatusUnauthorized, resp.StatusCode)
				}
			},
		},
	}

	mux := runtime.NewServeMux()
	server := new(MockUsersServer)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(server)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, tt.body)

			RegisterUsersHandlerServer(context.Background(), mux, server)

			mux.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			tt.assertResp(t, resp)
		})
	}
}

