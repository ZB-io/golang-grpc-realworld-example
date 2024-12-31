package handler

import (
		"context"
		"errors"
		"testing"
		"github.com/raahii/golang-grpc-realworld-example/auth"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"github.com/raahii/golang-grpc-realworld-example/proto"
		"github.com/raahii/golang-grpc-realworld-example/store"
		"github.com/rs/zerolog"
		"google.golang.org/grpc/codes"
		"google.golang.org/grpc/status"
		pb "github.com/raahii/golang-grpc-realworld-example/proto"
		"strconv"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/mock"
		"gorm.io/gorm"
)

type Handler struct {
	logger *zerolog.Logger
	us     *store.UserStore
	as     *store.ArticleStore
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}
type T struct {
	common
	isEnvSet bool
	context  *testContext
}T struct {
	common
	isEnvSet bool
	context  *testContext
}
type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}

type Article struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slug           string   `protobuf:"bytes,1,opt,name=slug,proto3" json:"slug,omitempty"`
	Title          string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description    string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Body           string   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	TagList        []string `protobuf:"bytes,5,rep,name=tagList,proto3" json:"tagList,omitempty"`
	CreatedAt      string   `protobuf:"bytes,6,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt      string   `protobuf:"bytes,7,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	Favorited      bool     `protobuf:"varint,8,opt,name=favorited,proto3" json:"favorited,omitempty"`
	FavoritesCount int32    `protobuf:"varint,9,opt,name=favoritesCount,proto3" json:"favoritesCount,omitempty"`
	Author         *Profile `protobuf:"bytes,10,opt,name=author,proto3" json:"author,omitempty"`
}

type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
}

type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}

type ArticleStore struct {
	db *gorm.DB
}

type UserStore struct {
	db *gorm.DB
}

type ArticleStore struct {
	db *gorm.DB
}
type UserStore struct {
	db *gorm.DB
}
type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
}
type MockUserStore struct {
	GetByUsernameFunc func(username string) (*model.User, error)
	GetByIDFunc       func(id uint) (*model.User, error)
	IsFollowingFunc   func(follower, followed *model.User) (bool, error)
}
type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}
type Call struct {
	Parent *Mock

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// Holds the caller info for the On() call
	callerInfo []string

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Amount of times this call has been called
	totalCalls int

	// Call to this method can be optional
	optional bool

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	RunFn func(Arguments)

	// PanicMsg holds msg to be used to mock panic on the function call
	//  if the PanicMsg is set to a non nil string the function call will panic
	// irrespective of other settings
	PanicMsg *string

	// Calls which must be satisfied before this call can be
	requires []*Call
}

type Call struct {
	Parent *Mock

	Method string

	Arguments Arguments

	ReturnArguments Arguments

	callerInfo []string

	Repeatability int

	totalCalls int

	optional bool

	WaitFor <-chan time.Time

	waitTime time.Duration

	RunFn func(Arguments)

	PanicMsg *string

	requires []*Call
}
type MockUserStore struct {
	mock.Mock
}
type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username  string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Bio       string `protobuf:"bytes,2,opt,name=bio,proto3" json:"bio,omitempty"`
	Image     string `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	Following bool   `protobuf:"varint,4,opt,name=following,proto3" json:"following,omitempty"`
}

/*
ROOST_METHOD_HASH=DeleteArticle_0347183038
ROOST_METHOD_SIG_HASH=DeleteArticle_b2585946c3


 */
func TestDeleteArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *store.MockArticleStore)
		userID         uint
		req            *proto.DeleteArticleRequest
		expectedError  error
		expectedStatus codes.Code
	}{
		{
			name: "Successfully Delete an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 1}}, nil)
				as.On("Delete", &model.Article{ID: 1, Author: model.User{ID: 1}}).Return(nil)
			},
			userID:        1,
			req:           &proto.DeleteArticleRequest{Slug: "1"},
			expectedError: nil,
		},
		{
			name:           "Attempt to Delete Article with Invalid Authentication",
			setupMocks:     func(us *store.MockUserStore, as *store.MockArticleStore) {},
			userID:         0,
			req:            &proto.DeleteArticleRequest{Slug: "1"},
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedStatus: codes.Unauthenticated,
		},
		{
			name: "Attempt to Delete Non-existent Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(999)).Return(nil, errors.New("article not found"))
			},
			userID:         1,
			req:            &proto.DeleteArticleRequest{Slug: "999"},
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Attempt to Delete Another User's Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 2}}, nil)
			},
			userID:         1,
			req:            &proto.DeleteArticleRequest{Slug: "1"},
			expectedError:  status.Error(codes.Unauthenticated, "forbidden"),
			expectedStatus: codes.Unauthenticated,
		},
		{
			name: "Handle Invalid Slug Format",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
			},
			userID:         1,
			req:            &proto.DeleteArticleRequest{Slug: "not-a-number"},
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Handle Database Error During Deletion",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 1}}, nil)
				as.On("Delete", &model.Article{ID: 1, Author: model.User{ID: 1}}).Return(errors.New("database error"))
			},
			userID:         1,
			req:            &proto.DeleteArticleRequest{Slug: "1"},
			expectedError:  status.Error(codes.Unauthenticated, "failed to delete article"),
			expectedStatus: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger := zerolog.Nop()
			us := &store.MockUserStore{}
			as := &store.MockArticleStore{}
			tt.setupMocks(us, as)

			h := &Handler{
				logger: &logger,
				us:     us,
				as:     as,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, tt.userID)

			_, err := h.DeleteArticle(ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					statusErr, ok := status.FromError(err)
					if !ok {
						t.Errorf("expected gRPC status error, got %v", err)
					} else if statusErr.Code() != tt.expectedStatus {
						t.Errorf("expected status code %v, got %v", tt.expectedStatus, statusErr.Code())
					}
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			us.AssertExpectations(t)
			as.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticle_8db60d3055
ROOST_METHOD_SIG_HASH=GetArticle_ea0095c9f8


 */
func TestGetArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockArticleStore, *store.MockUserStore)
		ctx            context.Context
		req            *pb.GetArticleRequest
		expectedResult *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully retrieve an article for an authenticated user",
			setupMocks: func(as *store.MockArticleStore, us *store.MockUserStore) {
				as.On("GetByID", uint(1)).Return(&model.Article{
					Slug:  "test-article",
					Title: "Test Article",
					Author: model.User{
						Username: "testuser",
					},
				}, nil)
				as.On("IsFavorited", &model.Article{}, &model.User{}).Return(true, nil)
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(true, nil)
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			req: &pb.GetArticleRequest{Slug: "1"},
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:      "test-article",
					Title:     "Test Article",
					Favorited: true,
					Author: &pb.Profile{
						Username:  "testuser",
						Following: true,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Retrieve an article for an unauthenticated user",
			setupMocks: func(as *store.MockArticleStore, us *store.MockUserStore) {
				as.On("GetByID", uint(1)).Return(&model.Article{
					Slug:  "test-article",
					Title: "Test Article",
					Author: model.User{
						Username: "testuser",
					},
				}, nil)
			},
			ctx: context.Background(),
			req: &pb.GetArticleRequest{Slug: "1"},
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:      "test-article",
					Title:     "Test Article",
					Favorited: false,
					Author: &pb.Profile{
						Username:  "testuser",
						Following: false,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Attempt to retrieve an article with an invalid slug",
			setupMocks:     func(as *store.MockArticleStore, us *store.MockUserStore) {},
			ctx:            context.Background(),
			req:            &pb.GetArticleRequest{Slug: "invalid"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Attempt to retrieve a non-existent article",
			setupMocks: func(as *store.MockArticleStore, us *store.MockUserStore) {
				as.On("GetByID", uint(999)).Return(nil, errors.New("article not found"))
			},
			ctx:            context.Background(),
			req:            &pb.GetArticleRequest{Slug: "999"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockArticleStore := &store.MockArticleStore{}
			mockUserStore := &store.MockUserStore{}

			tt.setupMocks(mockArticleStore, mockUserStore)

			h := &Handler{
				logger: zerolog.Nop(),
				as:     mockArticleStore,
				us:     mockUserStore,
			}

			result, err := h.GetArticle(tt.ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectedResult != nil {
				if result == nil {
					t.Error("Expected non-nil result, but got nil")
				} else if result.Article.Slug != tt.expectedResult.Article.Slug ||
					result.Article.Title != tt.expectedResult.Article.Title ||
					result.Article.Favorited != tt.expectedResult.Article.Favorited ||
					result.Article.Author.Username != tt.expectedResult.Article.Author.Username ||
					result.Article.Author.Following != tt.expectedResult.Article.Author.Following {
					t.Errorf("Expected result %v, but got %v", tt.expectedResult, result)
				}
			} else if result != nil {
				t.Errorf("Expected nil result, but got %v", result)
			}

			mockArticleStore.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=FavoriteArticle_29edacd2dc
ROOST_METHOD_SIG_HASH=FavoriteArticle_eb25e62ccd


 */
func TestFavoriteArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *store.MockArticleStore)
		req            *pb.FavoriteArticleRequest
		expectedResult *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Favorite an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Author: model.User{}}, nil)
				as.On("AddFavorite", mock.Anything, mock.Anything).Return(nil)
				us.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			req: &pb.FavoriteArticleRequest{Slug: "1"},
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:      "1",
					Favorited: true,
					Author:    &pb.Profile{},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Attempt to Favorite an Article with Unauthenticated User",
			setupMocks:     func(us *store.MockUserStore, as *store.MockArticleStore) {},
			req:            &pb.FavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Attempt to Favorite a Non-existent Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(999)).Return(nil, errors.New("article not found"))
			},
			req:            &pb.FavoriteArticleRequest{Slug: "999"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Handle Invalid Slug Format",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
			},
			req:            &pb.FavoriteArticleRequest{Slug: "not-a-number"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Handle Failure in Adding Favorite",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{}, nil)
				as.On("AddFavorite", mock.Anything, mock.Anything).Return(errors.New("failed to add favorite"))
			},
			req:            &pb.FavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "failed to add favorite"),
		},
		{
			name: "Handle Failure in Checking Following Status",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Author: model.User{}}, nil)
				as.On("AddFavorite", mock.Anything, mock.Anything).Return(nil)
				us.On("IsFollowing", mock.Anything, mock.Anything).Return(false, errors.New("failed to check following status"))
			},
			req:            &pb.FavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.NotFound, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.WithValue(context.Background(), auth.UserIDKey, uint(1))

			mockUserStore := new(store.MockUserStore)
			mockArticleStore := new(store.MockArticleStore)

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			result, err := h.FavoriteArticle(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Article.Slug, result.Article.Slug)
				assert.Equal(t, tt.expectedResult.Article.Favorited, result.Article.Favorited)
			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=UnfavoriteArticle_47bfda8100
ROOST_METHOD_SIG_HASH=UnfavoriteArticle_9043d547fd


 */
func TestUnfavoriteArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *store.MockArticleStore)
		req            *proto.UnfavoriteArticleRequest
		wantErr        bool
		expectedErrMsg string
		expectedCode   codes.Code
	}{
		{
			name: "Successful Unfavoriting of an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Author: model.User{}}, nil)
				as.On("DeleteFavorite", &model.Article{}, &model.User{}).Return(nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, nil)
			},
			req: &proto.UnfavoriteArticleRequest{Slug: "1"},
		},
		{
			name:           "Attempt to Unfavorite with Unauthenticated User",
			setupMocks:     func(us *store.MockUserStore, as *store.MockArticleStore) {},
			req:            &proto.UnfavoriteArticleRequest{Slug: "1"},
			wantErr:        true,
			expectedErrMsg: "unauthenticated",
			expectedCode:   codes.Unauthenticated,
		},
		{
			name: "Unfavorite with Invalid Article Slug",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "not-a-number"},
			wantErr:        true,
			expectedErrMsg: "invalid article id",
			expectedCode:   codes.InvalidArgument,
		},
		{
			name: "Unfavorite Non-Existent Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(999)).Return(nil, errors.New("article not found"))
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "999"},
			wantErr:        true,
			expectedErrMsg: "invalid article id",
			expectedCode:   codes.InvalidArgument,
		},
		{
			name: "Failure in Removing Favorite",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{}, nil)
				as.On("DeleteFavorite", &model.Article{}, &model.User{}).Return(errors.New("failed to remove favorite"))
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "1"},
			wantErr:        true,
			expectedErrMsg: "failed to remove favorite",
			expectedCode:   codes.InvalidArgument,
		},
		{
			name: "Error in Checking Following Status",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Author: model.User{}}, nil)
				as.On("DeleteFavorite", &model.Article{}, &model.User{}).Return(nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, errors.New("failed to check following status"))
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "1"},
			wantErr:        true,
			expectedErrMsg: "internal server error",
			expectedCode:   codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			us := &store.MockUserStore{}
			as := &store.MockArticleStore{}
			tt.setupMocks(us, as)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     us,
				as:     as,
			}

			ctx := context.Background()
			if tt.name != "Attempt to Unfavorite with Unauthenticated User" {
				ctx = auth.NewContextWithUserID(ctx, 1)
			}

			got, err := h.UnfavoriteArticle(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.UnfavoriteArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if status, ok := status.FromError(err); ok {
					if status.Code() != tt.expectedCode {
						t.Errorf("Handler.UnfavoriteArticle() error code = %v, want %v", status.Code(), tt.expectedCode)
					}
					if status.Message() != tt.expectedErrMsg {
						t.Errorf("Handler.UnfavoriteArticle() error message = %v, want %v", status.Message(), tt.expectedErrMsg)
					}
				} else {
					t.Errorf("Handler.UnfavoriteArticle() error is not a status error")
				}
			} else {
				if got == nil || got.Article == nil {
					t.Errorf("Handler.UnfavoriteArticle() returned nil response or article")
				}

			}

			us.AssertExpectations(t)
			as.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetFeedArticles_87ea56b889
ROOST_METHOD_SIG_HASH=GetFeedArticles_2be3462049


 */
func TestGetFeedArticles(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.UserStore, *store.ArticleStore)
		req            *proto.GetFeedArticlesRequest
		expectedResp   *proto.ArticlesResponse
		expectedErrMsg string
	}{
		{
			name: "Successful Retrieval of Feed Articles",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {
				us.GetByID = func(id uint) (*model.User, error) {
					return &model.User{Model: gorm.Model{ID: 1}}, nil
				}
				us.GetFollowingUserIDs = func(user *model.User) ([]uint, error) {
					return []uint{2, 3}, nil
				}
				as.GetFeedArticles = func(userIDs []uint, limit int64, offset int64) ([]model.Article, error) {
					return []model.Article{
						{Model: gorm.Model{ID: 1}, Title: "Test Article", Author: model.User{Model: gorm.Model{ID: 2}}},
					}, nil
				}
				as.IsFavorited = func(article *model.Article, user *model.User) (bool, error) {
					return false, nil
				}
				us.IsFollowing = func(follower *model.User, followed *model.User) (bool, error) {
					return true, nil
				}
			},
			req: &proto.GetFeedArticlesRequest{Limit: 20, Offset: 0},
			expectedResp: &proto.ArticlesResponse{
				Articles: []*proto.Article{
					{Title: "Test Article", Author: &proto.Profile{Following: true}},
				},
				ArticlesCount: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = auth.NewContext(ctx, 1)

			mockUserStore := &store.UserStore{}
			mockArticleStore := &store.ArticleStore{}

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			resp, err := h.GetFeedArticles(ctx, tt.req)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectedErrMsg)
				} else if status, ok := status.FromError(err); !ok || status.Message() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, status.Message())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Error("expected non-nil response, got nil")
				} else {
					if len(resp.Articles) != len(tt.expectedResp.Articles) {
						t.Errorf("expected %d articles, got %d", len(tt.expectedResp.Articles), len(resp.Articles))
					}
					if resp.ArticlesCount != tt.expectedResp.ArticlesCount {
						t.Errorf("expected ArticlesCount %d, got %d", tt.expectedResp.ArticlesCount, resp.ArticlesCount)
					}

				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_f87b10d80e
ROOST_METHOD_SIG_HASH=GetArticles_5d9fe7bf44


 */
func (m *MockArticleStore) GetArticles(tag, author string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	return m.GetArticlesFunc(tag, author, favoritedBy, limit, offset)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	return m.GetByIDFunc(id)
}

func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	return m.GetByUsernameFunc(username)
}

func (m *MockArticleStore) IsFavorited(article *model.Article, user *model.User) (bool, error) {
	return m.IsFavoritedFunc(article, user)
}

func (m *MockUserStore) IsFollowing(follower, followed *model.User) (bool, error) {
	return m.IsFollowingFunc(follower, followed)
}

func TestGetArticles(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore, *MockArticleStore)
		ctx            context.Context
		req            *proto.GetArticlesRequest
		expectedResp   *proto.ArticlesResponse
		expectedErrMsg string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUS := &MockUserStore{}
			mockAS := &MockArticleStore{}
			tt.setupMocks(mockUS, mockAS)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUS,
				as:     mockAS,
			}

			resp, err := h.GetArticles(tt.ctx, tt.req)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if status.Code(err) != codes.Aborted && status.Code(err) != codes.NotFound {
					t.Errorf("expected error code Aborted or NotFound, got %v", status.Code(err))
				} else if err.Error() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response, got nil")
				}
				if len(resp.Articles) != len(tt.expectedResp.Articles) {
					t.Errorf("expected %d articles, got %d", len(tt.expectedResp.Articles), len(resp.Articles))
				}
				if resp.ArticlesCount != tt.expectedResp.ArticlesCount {
					t.Errorf("expected ArticlesCount %d, got %d", tt.expectedResp.ArticlesCount, resp.ArticlesCount)
				}

			}
		})
	}
}


/*
ROOST_METHOD_HASH=CreateArticle_64372fa1a8
ROOST_METHOD_SIG_HASH=CreateArticle_ce1c125740


 */
func (m *MockArticleStore) Create(article *model.Article) error {
	args := m.Called(article)
	return args.Error(0)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) IsFollowing(follower, followed *model.User) (bool, error) {
	args := m.Called(follower, followed)
	return args.Bool(0), args.Error(1)
}

func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore, *MockArticleStore)
		input          *pb.CreateAritcleRequest
		expectedOutput *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Create an Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{Username: "testuser"}, nil)
				us.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
				as.On("Create", mock.AnythingOfType("*model.Article")).Return(nil)
			},
			input: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"test"},
				},
			},
			expectedOutput: &pb.ArticleResponse{
				Article: &pb.Article{
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"test"},
					Author: &pb.Profile{
						Username:  "testuser",
						Following: false,
					},
					Favorited: true,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			mockArticleStore := new(MockArticleStore)

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.New(nil),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, uint(1))

			got, err := h.CreateArticle(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.expectedOutput.Article.Title, got.Article.Title)
				assert.Equal(t, tt.expectedOutput.Article.Description, got.Article.Description)
				assert.Equal(t, tt.expectedOutput.Article.Body, got.Article.Body)
				assert.Equal(t, tt.expectedOutput.Article.TagList, got.Article.TagList)
				assert.Equal(t, tt.expectedOutput.Article.Author.Username, got.Article.Author.Username)
				assert.Equal(t, tt.expectedOutput.Article.Author.Following, got.Article.Author.Following)
				assert.Equal(t, tt.expectedOutput.Article.Favorited, got.Article.Favorited)
			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=UpdateArticle_c5b82e271b
ROOST_METHOD_SIG_HASH=UpdateArticle_f36cc09d87


 */
func TestUpdateArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *store.MockArticleStore)
		userID         uint
		req            *proto.UpdateArticleRequest
		expectedResult *proto.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Update an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 1}}, nil)
				as.On("Update", &model.Article{ID: 1, Author: model.User{ID: 1}, Title: "Updated Title", Description: "Updated Description", Body: "Updated Body"}).Return(nil)
				us.On("IsFollowing", &model.User{ID: 1}, &model.User{ID: 1}).Return(false, nil)
			},
			userID: 1,
			req: &proto.UpdateArticleRequest{
				Article: &proto.UpdateArticleRequest_Article{
					Slug:        "1",
					Title:       "Updated Title",
					Description: "Updated Description",
					Body:        "Updated Body",
				},
			},
			expectedResult: &proto.ArticleResponse{
				Article: &proto.Article{
					Slug:        "1",
					Title:       "Updated Title",
					Description: "Updated Description",
					Body:        "Updated Body",
					Author:      &proto.Profile{},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUserStore := &store.MockUserStore{}
			mockArticleStore := &store.MockArticleStore{}

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, tt.userID)

			result, err := h.UpdateArticle(ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Error("expected non-nil result, got nil")
				} else {
					if result.Article.Slug != tt.expectedResult.Article.Slug {
						t.Errorf("expected slug %s, got %s", tt.expectedResult.Article.Slug, result.Article.Slug)
					}
					if result.Article.Title != tt.expectedResult.Article.Title {
						t.Errorf("expected title %s, got %s", tt.expectedResult.Article.Title, result.Article.Title)
					}

				}
			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
		})
	}
}

