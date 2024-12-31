package handler

import (
		"context"
		"errors"
		"testing"
		"github.com/raahii/golang-grpc-realworld-example/auth"
		"github.com/raahii/golang-grpc-realworld-example/model"
		pb "github.com/raahii/golang-grpc-realworld-example/proto"
		"github.com/raahii/golang-grpc-realworld-example/store"
		"github.com/rs/zerolog"
		"google.golang.org/grpc/codes"
		"google.golang.org/grpc/status"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/mock"
		"github.com/raahii/golang-grpc-realworld-example/proto"
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

type T struct {
	common
	isEnvSet bool
	context  *testContext
}T struct {
	common
	isEnvSet bool
	context  *testContext
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
type mockArticleStore struct {
	mock.Mock
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

type MockArticleStore struct {
	mock.Mock
}
type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}
type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
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
		req            *pb.DeleteArticleRequest
		expectedResult *pb.Empty
		expectedError  error
	}{
		{
			name: "Successfully delete an article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 1}}, nil)
				as.On("Delete", &model.Article{ID: 1, Author: model.User{ID: 1}}).Return(nil)
			},
			userID:         1,
			req:            &pb.DeleteArticleRequest{Slug: "1"},
			expectedResult: &pb.Empty{},
			expectedError:  nil,
		},
		{
			name:           "Attempt to delete an article without authentication",
			setupMocks:     func(us *store.MockUserStore, as *store.MockArticleStore) {},
			userID:         0,
			req:            &pb.DeleteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Attempt to delete a non-existent article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(999)).Return(nil, errors.New("article not found"))
			},
			userID:         1,
			req:            &pb.DeleteArticleRequest{Slug: "999"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Attempt to delete another user's article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 2}}, nil)
			},
			userID:         1,
			req:            &pb.DeleteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "forbidden"),
		},
		{
			name: "Handle invalid slug format",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
			},
			userID:         1,
			req:            &pb.DeleteArticleRequest{Slug: "not-an-integer"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Handle database error during article deletion",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 1}}, nil)
				as.On("Delete", &model.Article{ID: 1, Author: model.User{ID: 1}}).Return(errors.New("database error"))
			},
			userID:         1,
			req:            &pb.DeleteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "failed to delete article"),
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
			if tt.userID != 0 {
				ctx = auth.NewContext(ctx, tt.userID)
			}

			result, err := h.DeleteArticle(ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected non-nil result, got nil")
				}
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
func (m *mockArticleStore) GetByID(id uint) (*model.Article, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *mockAuth) GetUserID(ctx context.Context) (uint, error) {
	if ctx.Value("userID") != nil {
		return ctx.Value("userID").(uint), nil
	}
	return 0, errors.New("user not authenticated")
}

func (m *mockArticleStore) IsFavorited(article *model.Article, user *model.User) (bool, error) {
	args := m.Called(article, user)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserStore) IsFollowing(follower *model.User, followed *model.User) (bool, error) {
	args := m.Called(follower, followed)
	return args.Bool(0), args.Error(1)
}

func TestGetArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*mockArticleStore, *mockUserStore)
		ctx            context.Context
		req            *pb.GetArticleRequest
		expectedResp   *pb.ArticleResponse
		expectedErrMsg string
	}{
		{
			name: "Successfully retrieve an article for an authenticated user",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				article := &model.Article{Slug: "1", Title: "Test Article", Author: model.User{Username: "author"}}
				as.On("GetByID", uint(1)).Return(article, nil)
				as.On("IsFavorited", article, mock.AnythingOfType("*model.User")).Return(true, nil)
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				us.On("IsFollowing", mock.AnythingOfType("*model.User"), &article.Author).Return(true, nil)
			},
			ctx: context.WithValue(context.Background(), "userID", uint(1)),
			req: &pb.GetArticleRequest{Slug: "1"},
			expectedResp: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:      "1",
					Title:     "Test Article",
					Favorited: true,
					Author: &pb.Profile{
						Username:  "author",
						Following: true,
					},
				},
			},
		},
		{
			name: "Retrieve an article for an unauthenticated user",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				article := &model.Article{Slug: "1", Title: "Test Article", Author: model.User{Username: "author"}}
				as.On("GetByID", uint(1)).Return(article, nil)
			},
			ctx: context.Background(),
			req: &pb.GetArticleRequest{Slug: "1"},
			expectedResp: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:      "1",
					Title:     "Test Article",
					Favorited: false,
					Author: &pb.Profile{
						Username:  "author",
						Following: false,
					},
				},
			},
		},
		{
			name:           "Attempt to retrieve an article with an invalid slug",
			setupMocks:     func(as *mockArticleStore, us *mockUserStore) {},
			ctx:            context.Background(),
			req:            &pb.GetArticleRequest{Slug: "invalid"},
			expectedErrMsg: "invalid article id",
		},
		{
			name: "Attempt to retrieve a non-existent article",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				as.On("GetByID", uint(999)).Return((*model.Article)(nil), errors.New("article not found"))
			},
			ctx:            context.Background(),
			req:            &pb.GetArticleRequest{Slug: "999"},
			expectedErrMsg: "invalid article id",
		},
		{
			name: "Retrieve an article with a valid token but non-existent user",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				article := &model.Article{Slug: "1", Title: "Test Article", Author: model.User{Username: "author"}}
				as.On("GetByID", uint(1)).Return(article, nil)
				us.On("GetByID", uint(1)).Return((*model.User)(nil), errors.New("user not found"))
			},
			ctx:            context.WithValue(context.Background(), "userID", uint(1)),
			req:            &pb.GetArticleRequest{Slug: "1"},
			expectedErrMsg: "token is valid but the user not found",
		},
		{
			name: "Handle error when checking if article is favorited",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				article := &model.Article{Slug: "1", Title: "Test Article", Author: model.User{Username: "author"}}
				as.On("GetByID", uint(1)).Return(article, nil)
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("IsFavorited", article, mock.AnythingOfType("*model.User")).Return(false, errors.New("database error"))
			},
			ctx:            context.WithValue(context.Background(), "userID", uint(1)),
			req:            &pb.GetArticleRequest{Slug: "1"},
			expectedErrMsg: "internal server error",
		},
		{
			name: "Handle error when checking if user is following the author",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				article := &model.Article{Slug: "1", Title: "Test Article", Author: model.User{Username: "author"}}
				as.On("GetByID", uint(1)).Return(article, nil)
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("IsFavorited", article, mock.AnythingOfType("*model.User")).Return(true, nil)
				us.On("IsFollowing", mock.AnythingOfType("*model.User"), &article.Author).Return(false, errors.New("database error"))
			},
			ctx:            context.WithValue(context.Background(), "userID", uint(1)),
			req:            &pb.GetArticleRequest{Slug: "1"},
			expectedErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAS := new(mockArticleStore)
			mockUS := new(mockUserStore)
			tt.setupMocks(mockAS, mockUS)

			h := &Handler{
				logger: zerolog.Nop(),
				as:     mockAS,
				us:     mockUS,
				auth:   &mockAuth{},
			}

			resp, err := h.GetArticle(tt.ctx, tt.req)

			if tt.expectedErrMsg != "" {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Contains(t, st.Message(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}

			mockAS.AssertExpectations(t)
			mockUS.AssertExpectations(t)
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
			name: "Successful Favoriting of an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Author: model.User{}}, nil)
				as.On("AddFavorite", &model.Article{}, &model.User{}).Return(nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, nil)
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
			name: "Favoriting with Invalid Article Slug",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
			},
			req:            &pb.FavoriteArticleRequest{Slug: "not-a-number"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
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
			name: "Error During Favoriting Process",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{}, nil)
				as.On("AddFavorite", &model.Article{}, &model.User{}).Return(errors.New("failed to add favorite"))
			},
			req:            &pb.FavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "failed to add favorite"),
		},
		{
			name: "Error Retrieving Following Status",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Author: model.User{}}, nil)
				as.On("AddFavorite", &model.Article{}, &model.User{}).Return(nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, errors.New("failed to get following status"))
			},
			req:            &pb.FavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.NotFound, "internal server error"),
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

			ctx := context.Background()
			if tt.name != "Attempt to Favorite an Article with Unauthenticated User" {
				ctx = auth.NewContext(ctx, 1)
			}

			result, err := h.FavoriteArticle(ctx, tt.req)

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
				} else {
					if result.Article.Slug != tt.expectedResult.Article.Slug {
						t.Errorf("Expected slug %s, but got %s", tt.expectedResult.Article.Slug, result.Article.Slug)
					}
					if result.Article.Favorited != tt.expectedResult.Article.Favorited {
						t.Errorf("Expected favorited %v, but got %v", tt.expectedResult.Article.Favorited, result.Article.Favorited)
					}
				}
			} else if result != nil {
				t.Errorf("Expected nil result, but got %v", result)
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
		expectedResult *proto.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Unfavorite an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{FavoritesCount: 1}, nil)
				as.On("DeleteFavorite", &model.Article{}, &model.User{}).Return(nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, nil)
			},
			req: &proto.UnfavoriteArticleRequest{Slug: "1"},
			expectedResult: &proto.ArticleResponse{
				Article: &proto.Article{
					Slug:           "1",
					Favorited:      false,
					FavoritesCount: 0,
					Author:         &proto.Profile{},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Attempt to Unfavorite with Unauthenticated User",
			setupMocks:     func(us *store.MockUserStore, as *store.MockArticleStore) {},
			req:            &proto.UnfavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Attempt to Unfavorite a Non-existent Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(999)).Return(nil, errors.New("article not found"))
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "999"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Invalid Slug Format",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "not-a-number"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Failure to Remove Favorite",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{}, nil)
				as.On("DeleteFavorite", &model.Article{}, &model.User{}).Return(errors.New("failed to remove favorite"))
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "failed to remove favorite"),
		},
		{
			name: "Failure to Get Following Status",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{}, nil)
				as.On("DeleteFavorite", &model.Article{}, &model.User{}).Return(nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, errors.New("failed to get following status"))
			},
			req:            &proto.UnfavoriteArticleRequest{Slug: "1"},
			expectedResult: nil,
			expectedError:  status.Error(codes.NotFound, "internal server error"),
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

			ctx := context.Background()
			if tt.name != "Attempt to Unfavorite with Unauthenticated User" {
				ctx = auth.NewContext(ctx, 1)
			}

			result, err := h.UnfavoriteArticle(ctx, tt.req)

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
				} else {
					if result.Article.Slug != tt.expectedResult.Article.Slug {
						t.Errorf("Expected slug %s, but got %s", tt.expectedResult.Article.Slug, result.Article.Slug)
					}
					if result.Article.Favorited != tt.expectedResult.Article.Favorited {
						t.Errorf("Expected favorited %v, but got %v", tt.expectedResult.Article.Favorited, result.Article.Favorited)
					}
					if result.Article.FavoritesCount != tt.expectedResult.Article.FavoritesCount {
						t.Errorf("Expected favoritesCount %d, but got %d", tt.expectedResult.Article.FavoritesCount, result.Article.FavoritesCount)
					}
				}
			} else if result != nil {
				t.Error("Expected nil result, but got non-nil")
			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
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
		setupMocks     func(*store.MockUserStore, *store.MockArticleStore)
		req            *proto.GetFeedArticlesRequest
		expectedResp   *proto.ArticlesResponse
		expectedErrMsg string
		expectedCode   codes.Code
	}{
		{
			name: "Successful Retrieval of Feed Articles",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{Model: model.Model{ID: 1}}, nil)
				us.On("GetFollowingUserIDs", &model.User{Model: model.Model{ID: 1}}).Return([]uint{2, 3}, nil)
				as.On("GetFeedArticles", []uint{2, 3}, int64(10), int64(0)).Return([]model.Article{
					{Model: model.Model{ID: 1}, Title: "Article 1", Author: model.User{Model: model.Model{ID: 2}}},
					{Model: model.Model{ID: 2}, Title: "Article 2", Author: model.User{Model: model.Model{ID: 3}}},
				}, nil)
				as.On("IsFavorited", &model.Article{Model: model.Model{ID: 1}}, &model.User{Model: model.Model{ID: 1}}).Return(false, nil)
				as.On("IsFavorited", &model.Article{Model: model.Model{ID: 2}}, &model.User{Model: model.Model{ID: 1}}).Return(false, nil)
				us.On("IsFollowing", &model.User{Model: model.Model{ID: 1}}, &model.User{Model: model.Model{ID: 2}}).Return(true, nil)
				us.On("IsFollowing", &model.User{Model: model.Model{ID: 1}}, &model.User{Model: model.Model{ID: 3}}).Return(true, nil)
			},
			req: &proto.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			expectedResp: &proto.ArticlesResponse{
				Articles: []*proto.Article{
					{Title: "Article 1", Author: &proto.Profile{Following: true}},
					{Title: "Article 2", Author: &proto.Profile{Following: true}},
				},
				ArticlesCount: 2,
			},
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

			ctx := context.Background()
			if tt.expectedCode != codes.Unauthenticated {
				ctx = auth.NewContext(ctx, 1)
			}

			resp, err := h.GetFeedArticles(ctx, tt.req)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Errorf("expected error code %v, got %v", tt.expectedCode, st.Code())
				}
				if st.Message() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, st.Message())
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedResp != nil {
				if resp == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(resp.Articles) != len(tt.expectedResp.Articles) {
					t.Errorf("expected %d articles, got %d", len(tt.expectedResp.Articles), len(resp.Articles))
				}
				if resp.ArticlesCount != tt.expectedResp.ArticlesCount {
					t.Errorf("expected ArticlesCount %d, got %d", tt.expectedResp.ArticlesCount, resp.ArticlesCount)
				}

			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetArticles_f87b10d80e
ROOST_METHOD_SIG_HASH=GetArticles_5d9fe7bf44


 */
func (m *MockArticleStore) GetArticles(tag, author string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	args := m.Called(tag, author, favoritedBy, limit, offset)
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockArticleStore) IsFavorited(article *model.Article, user *model.User) (bool, error) {
	args := m.Called(article, user)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStore) IsFollowing(follower, followed *model.User) (bool, error) {
	args := m.Called(follower, followed)
	return args.Bool(0), args.Error(1)
}

func TestGetArticles(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockArticleStore, *MockUserStore)
		req            *pb.GetArticlesRequest
		ctx            context.Context
		expectedResp   *pb.ArticlesResponse
		expectedErrMsg string
	}{
		{
			name: "Successful retrieval of articles with default limit",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				articles := make([]model.Article, 20)
				for i := range articles {
					articles[i] = model.Article{Title: "Article " + string(rune(i))}
				}
				mas.On("GetArticles", "", "", (*model.User)(nil), int64(20), int64(0)).Return(articles, nil)
				mas.On("IsFavorited", mock.Anything, mock.Anything).Return(false, nil)
				mus.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			req:          &pb.GetArticlesRequest{},
			ctx:          context.Background(),
			expectedResp: &pb.ArticlesResponse{Articles: make([]*pb.Article, 20), ArticlesCount: 20},
		},
		{
			name: "Retrieval of articles with specified tag",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				articles := []model.Article{{Title: "Tagged Article"}}
				mas.On("GetArticles", "test-tag", "", (*model.User)(nil), int64(20), int64(0)).Return(articles, nil)
				mas.On("IsFavorited", mock.Anything, mock.Anything).Return(false, nil)
				mus.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			req:          &pb.GetArticlesRequest{Tag: "test-tag"},
			ctx:          context.Background(),
			expectedResp: &pb.ArticlesResponse{Articles: []*pb.Article{{Title: "Tagged Article"}}, ArticlesCount: 1},
		},
		{
			name: "Handling of non-existent author in request",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				mus.On("GetByUsername", "non-existent").Return((*model.User)(nil), errors.New("user not found"))
				mas.On("GetArticles", "", "non-existent", (*model.User)(nil), int64(20), int64(0)).Return([]model.Article{}, nil)
			},
			req:          &pb.GetArticlesRequest{Author: "non-existent"},
			ctx:          context.Background(),
			expectedResp: &pb.ArticlesResponse{Articles: []*pb.Article{}, ArticlesCount: 0},
		},
		{
			name: "Error handling for database failure",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				mas.On("GetArticles", "", "", (*model.User)(nil), int64(20), int64(0)).Return([]model.Article{}, errors.New("database error"))
			},
			req:            &pb.GetArticlesRequest{},
			ctx:            context.Background(),
			expectedErrMsg: "internal server error",
		},
		{
			name: "Retrieval of articles with favorited filter",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				favoritedUser := &model.User{Username: "favorited-user"}
				mus.On("GetByUsername", "favorited-user").Return(favoritedUser, nil)
				articles := []model.Article{{Title: "Favorited Article"}}
				mas.On("GetArticles", "", "", favoritedUser, int64(20), int64(0)).Return(articles, nil)
				mas.On("IsFavorited", mock.Anything, mock.Anything).Return(true, nil)
				mus.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			req:          &pb.GetArticlesRequest{Favorited: "favorited-user"},
			ctx:          context.Background(),
			expectedResp: &pb.ArticlesResponse{Articles: []*pb.Article{{Title: "Favorited Article", Favorited: true}}, ArticlesCount: 1},
		},
		{
			name: "Handling of authenticated user context",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				currentUser := &model.User{Username: "current-user"}
				mus.On("GetByID", uint(1)).Return(currentUser, nil)
				articles := []model.Article{{Title: "Auth User Article", Author: model.User{Username: "author"}}}
				mas.On("GetArticles", "", "", (*model.User)(nil), int64(20), int64(0)).Return(articles, nil)
				mas.On("IsFavorited", mock.Anything, currentUser).Return(true, nil)
				mus.On("IsFollowing", currentUser, mock.Anything).Return(true, nil)
			},
			req: &pb.GetArticlesRequest{},
			ctx: auth.NewContext(context.Background(), 1),
			expectedResp: &pb.ArticlesResponse{
				Articles: []*pb.Article{
					{
						Title:     "Auth User Article",
						Favorited: true,
						Author:    &pb.Profile{Username: "author", Following: true},
					},
				},
				ArticlesCount: 1,
			},
		},
		{
			name: "Pagination with offset and limit",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				articles := []model.Article{{Title: "Paginated Article 1"}, {Title: "Paginated Article 2"}}
				mas.On("GetArticles", "", "", (*model.User)(nil), int64(2), int64(5)).Return(articles, nil)
				mas.On("IsFavorited", mock.Anything, mock.Anything).Return(false, nil)
				mus.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			req:          &pb.GetArticlesRequest{Limit: 2, Offset: 5},
			ctx:          context.Background(),
			expectedResp: &pb.ArticlesResponse{Articles: []*pb.Article{{Title: "Paginated Article 1"}, {Title: "Paginated Article 2"}}, ArticlesCount: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleStore := new(MockArticleStore)
			mockUserStore := new(MockUserStore)
			tt.setupMocks(mockArticleStore, mockUserStore)

			h := &Handler{
				logger: zerolog.New(zerolog.NewTestWriter(t)),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			resp, err := h.GetArticles(tt.ctx, tt.req)

			if tt.expectedErrMsg != "" {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Aborted, st.Code())
				assert.Contains(t, st.Message(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp.ArticlesCount, resp.ArticlesCount)
				assert.Len(t, resp.Articles, len(tt.expectedResp.Articles))
				for i, article := range resp.Articles {
					assert.Equal(t, tt.expectedResp.Articles[i].Title, article.Title)
					assert.Equal(t, tt.expectedResp.Articles[i].Favorited, article.Favorited)
					if tt.expectedResp.Articles[i].Author != nil {
						assert.Equal(t, tt.expectedResp.Articles[i].Author.Username, article.Author.Username)
						assert.Equal(t, tt.expectedResp.Articles[i].Author.Following, article.Author.Following)
					}
				}
			}

			mockArticleStore.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=CreateArticle_64372fa1a8
ROOST_METHOD_SIG_HASH=CreateArticle_ce1c125740


 */
func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *store.MockArticleStore)
		setupAuth      func(context.Context) context.Context
		input          *pb.CreateAritcleRequest
		expectedOutput *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Create an Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				us.On("IsFollowing", &model.User{}, &model.User{}).Return(false, nil)
				as.On("Create", &model.Article{}).Return(nil)
			},
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, "user_id", uint(1))
			},
			input: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectedOutput: &pb.ArticleResponse{
				Article: &pb.Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
					Favorited:   true,
					Author:      &pb.Profile{},
				},
			},
			expectedError: nil,
		},
		{
			name:       "Attempt to Create an Article with Unauthenticated User",
			setupMocks: func(*store.MockUserStore, *store.MockArticleStore) {},
			setupAuth: func(ctx context.Context) context.Context {
				return ctx
			},
			input:          &pb.CreateAritcleRequest{},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Attempt to Create an Article with Non-existent User",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(nil, errors.New("user not found"))
			},
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, "user_id", uint(1))
			},
			input:          &pb.CreateAritcleRequest{},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "user not found"),
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

			ctx := tt.setupAuth(context.Background())

			got, err := h.CreateArticle(ctx, tt.input)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectedOutput != nil {
				if got == nil {
					t.Error("expected non-nil output, got nil")
				} else {
					if got.Article.Title != tt.expectedOutput.Article.Title {
						t.Errorf("expected title %s, got %s", tt.expectedOutput.Article.Title, got.Article.Title)
					}

				}
			} else if got != nil {
				t.Errorf("expected nil output, got %v", got)
			}

			us.AssertExpectations(t)
			as.AssertExpectations(t)
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
		{
			name: "Attempt to Update Another User's Article",
			setupMocks: func(us *store.MockUserStore, as *store.MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{ID: 1, Author: model.User{ID: 2}}, nil)
			},
			userID: 1,
			req: &proto.UpdateArticleRequest{
				Article: &proto.UpdateArticleRequest_Article{
					Slug: "1",
				},
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "forbidden"),
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
				} else {
					if result.Article.Slug != tt.expectedResult.Article.Slug {
						t.Errorf("Expected slug %s, but got %s", tt.expectedResult.Article.Slug, result.Article.Slug)
					}
					if result.Article.Title != tt.expectedResult.Article.Title {
						t.Errorf("Expected title %s, but got %s", tt.expectedResult.Article.Title, result.Article.Title)
					}

				}
			} else if result != nil {
				t.Errorf("Expected nil result, but got %v", result)
			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
		})
	}
}

