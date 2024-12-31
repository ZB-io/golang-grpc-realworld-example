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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserStore struct {
	GetByIDFunc     func(uint) (*model.User, error)
	IsFollowingFunc func(*model.User, *model.User) (bool, error)
}

type MockArticleStore struct {
	GetByIDFunc        func(uint) (*model.Article, error)
	DeleteFunc          func(*model.Article) error
	CreateFunc          func(*model.Article) error
	UpdateFunc          func(*model.Article) error
	IsFavoritedFunc     func(*model.Article, *model.User) (bool, error)
	AddFavoriteFunc     func(*model.Article, *model.User) error
	DeleteFavoriteFunc  func(*model.Article, *model.User) error
	GetArticlesFunc     func(string, string, *model.User, int64, int64) ([]model.Article, error)
	GetFeedArticlesFunc func([]uint, int64, int64) ([]model.Article, error)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	return m.GetByIDFunc(id)
}

func (m *MockUserStore) IsFollowing(follower, followed *model.User) (bool, error) {
	return m.IsFollowingFunc(follower, followed)
}

func (m *MockArticleStore) GetByID(id uint) (*model.Article, error) {
	return m.GetByIDFunc(id)
}

func (m *MockArticleStore) Delete(article *model.Article) error {
	return m.DeleteFunc(article)
}

func (m *MockArticleStore) Create(article *model.Article) error {
	return m.CreateFunc(article)
}

func (m *MockArticleStore) Update(article *model.Article) error {
	return m.UpdateFunc(article)
}

func (m *MockArticleStore) IsFavorited(article *model.Article, user *model.User) (bool, error) {
	return m.IsFavoritedFunc(article, user)
}

func (m *MockArticleStore) AddFavorite(article *model.Article, user *model.User) error {
	return m.AddFavoriteFunc(article, user)
}

func (m *MockArticleStore) DeleteFavorite(article *model.Article, user *model.User) error {
	return m.DeleteFavoriteFunc(article, user)
}

func (m *MockArticleStore) GetArticles(tag, author string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	return m.GetArticlesFunc(tag, author, favoritedBy, limit, offset)
}

func (m *MockArticleStore) GetFeedArticles(userIDs []uint, limit, offset int64) ([]model.Article, error) {
	return m.GetFeedArticlesFunc(userIDs, limit, offset)
}

/*
ROOST_METHOD_HASH=DeleteArticle_0347183038
ROOST_METHOD_SIG_HASH=DeleteArticle_b2585946c3
*/
func TestDeleteArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore, *MockArticleStore)
		userID         uint
		req            *pb.DeleteArticleRequest
		expectedError  error
		expectedStatus codes.Code
	}{
		{
			name: "Successfully Delete an Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				as.GetByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1, Author: model.User{ID: 1}}, nil
				}
				as.DeleteFunc = func(article *model.Article) error {
					return nil
				}
			},
			userID:         1,
			req:            &pb.DeleteArticleRequest{Slug: "1"},
			expectedError:  nil,
			expectedStatus: codes.OK,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &MockUserStore{}
			as := &MockArticleStore{}
			tt.setupMocks(us, as)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     us,
				as:     as,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, tt.userID)

			_, err := h.DeleteArticle(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, statusErr.Code())
			} else {
				assert.NoError(t, err)
			}
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
		setupMocks     func(*MockArticleStore, *MockUserStore)
		req            *pb.GetArticleRequest
		ctx            context.Context
		expectedResult *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully retrieve an article for an authenticated user",
			setupMocks: func(as *MockArticleStore, us *MockUserStore) {
				as.GetByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{
						Slug:  "test-article",
						Title: "Test Article",
						Author: model.User{
							Username: "testauthor",
						},
					}, nil
				}
				as.IsFavoritedFunc = func(a *model.Article, u *model.User) (bool, error) {
					return true, nil
				}
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{Username: "testuser"}, nil
				}
				us.IsFollowingFunc = func(a, b *model.User) (bool, error) {
					return true, nil
				}
			},
			req: &pb.GetArticleRequest{Slug: "1"},
			ctx: context.WithValue(context.Background(), auth.UserIDKey, uint(3)),
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:      "test-article",
					Title:     "Test Article",
					Favorited: true,
					Author: &pb.Profile{
						Username:  "testauthor",
						Following: true,
					},
				},
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := &MockArticleStore{}
			us := &MockUserStore{}
			tt.setupMocks(as, us)

			h := &Handler{
				logger: zerolog.Nop(),
				as:     as,
				us:     us,
			}

			result, err := h.GetArticle(tt.ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Article.Slug, result.Article.Slug)
				assert.Equal(t, tt.expectedResult.Article.Title, result.Article.Title)
				assert.Equal(t, tt.expectedResult.Article.Favorited, result.Article.Favorited)
				assert.Equal(t, tt.expectedResult.Article.Author.Username, result.Article.Author.Username)
				assert.Equal(t, tt.expectedResult.Article.Author.Following, result.Article.Author.Following)
			}
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
		setupMocks     func(*MockUserStore, *MockArticleStore)
		contextUserID  uint
		inputSlug      string
		expectedResult *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successful Favorite Article Operation",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				as.GetByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1, Author: model.User{}}, nil
				}
				as.AddFavoriteFunc = func(a *model.Article, u *model.User) error {
					return nil
				}
				us.IsFollowingFunc = func(a, b *model.User) (bool, error) {
					return false, nil
				}
			},
			contextUserID: 1,
			inputSlug:     "1",
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Favorited: true,
					Author:    &pb.Profile{},
				},
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			us := &MockUserStore{}
			as := &MockArticleStore{}
			tt.setupMocks(us, as)

			h := &Handler{
				logger: &logger,
				us:     us,
				as:     as,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, tt.contextUserID)

			req := &pb.FavoriteArticleRequest{
				Slug: tt.inputSlug,
			}

			result, err := h.FavoriteArticle(ctx, req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Article.Favorited, result.Article.Favorited)
			}
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
		setupMocks     func(*MockUserStore, *MockArticleStore)
		userID         uint
		req            *pb.UnfavoriteArticleRequest
		expectedResult *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Unfavorite an Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				as.GetByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1, Favorited: true, FavoritesCount: 1}, nil
				}
				as.DeleteFavoriteFunc = func(article *model.Article, user *model.User) error {
					return nil
				}
				us.IsFollowingFunc = func(follower, followed *model.User) (bool, error) {
					return false, nil
				}
			},
			userID: 1,
			req:    &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:           "1",
					Favorited:      false,
					FavoritesCount: 0,
					Author:         &pb.Profile{},
				},
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &MockUserStore{}
			mockArticleStore := &MockArticleStore{}
			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, tt.userID)

			result, err := h.UnfavoriteArticle(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Article.Slug, result.Article.Slug)
				assert.Equal(t, tt.expectedResult.Article.Favorited, result.Article.Favorited)
				assert.Equal(t, tt.expectedResult.Article.FavoritesCount, result.Article.FavoritesCount)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetFeedArticles_87ea56b889
ROOST_METHOD_SIG_HASH=GetFeedArticles_2be3462049
*/
func TestGetFeedArticles(t *testing.T) {
	tests := []struct {
		name            string
		setupMocks      func(*MockUserStore, *MockArticleStore)
		req             *pb.GetFeedArticlesRequest
		expectedResp    *pb.ArticlesResponse
		expectedErrCode codes.Code
	}{
		{
			name: "Successfully retrieve feed articles",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				us.IsFollowingFunc = func(follower, followed *model.User) (bool, error) {
					return true, nil
				}
				as.GetFeedArticlesFunc = func(userIDs []uint, limit, offset int64) ([]model.Article, error) {
					return []model.Article{
						{ID: 1, Title: "Article 1", Author: model.User{ID: 2}},
						{ID: 2, Title: "Article 2", Author: model.User{ID: 3}},
					}, nil
				}
				as.IsFavoritedFunc = func(article *model.Article, user *model.User) (bool, error) {
					return false, nil
				}
			},
			req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			expectedResp: &pb.ArticlesResponse{
				Articles: []*pb.Article{
					{Title: "Article 1", Author: &pb.Profile{Following: true}},
					{Title: "Article 2", Author: &pb.Profile{Following: true}},
				},
				ArticlesCount: 2,
			},
			expectedErrCode: codes.OK,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &MockUserStore{}
			as := &MockArticleStore{}

			tt.setupMocks(us, as)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     us,
				as:     as,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, uint(1))

			resp, err := h.GetFeedArticles(ctx, tt.req)

			if tt.expectedErrCode != codes.OK {
				assert.Error(t, err)
				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedErrCode, statusErr.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(tt.expectedResp.Articles), len(resp.Articles))
				assert.Equal(t, tt.expectedResp.ArticlesCount, resp.ArticlesCount)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetArticles_f87b10d80e
ROOST_METHOD_SIG_HASH=GetArticles_5d9fe7bf44
*/
func TestGetArticles(t *testing.T) {
	tests := []struct {
		name           string
		req            *pb.GetArticlesRequest
		setupMocks     func(*MockArticleStore, *MockUserStore)
		expectedResult *pb.ArticlesResponse
		expectedError  error
	}{
		{
			name: "Successful retrieval of articles with default limit",
			req:  &pb.GetArticlesRequest{},
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				articles := make([]model.Article, 20)
				for i := range articles {
					articles[i] = model.Article{Title: "Article " + string(rune(i))}
				}
				mas.GetArticlesFunc = func(tag, author string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
					return articles, nil
				}
				mas.IsFavoritedFunc = func(article *model.Article, user *model.User) (bool, error) {
					return false, nil
				}
				mus.IsFollowingFunc = func(follower, followed *model.User) (bool, error) {
					return false, nil
				}
			},
			expectedResult: &pb.ArticlesResponse{
				Articles:      make([]*pb.Article, 20),
				ArticlesCount: 20,
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mas := new(MockArticleStore)
			mus := new(MockUserStore)

			tt.setupMocks(mas, mus)

			h := &Handler{
				logger: zerolog.New(nil),
				us:     mus,
				as:     mas,
			}

			ctx := context.Background()
			result, err := h.GetArticles(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.ArticlesCount, result.ArticlesCount)
				assert.Len(t, result.Articles, int(tt.expectedResult.ArticlesCount))
			}
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
		setupMocks     func(*MockUserStore, *MockArticleStore)
		input          *pb.CreateAritcleRequest
		expectedOutput *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Create an Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.GetByIDFunc = func(uint) (*model.User, error) {
					return &model.User{Username: "testuser"}, nil
				}
				us.IsFollowingFunc = func(*model.User, *model.User) (bool, error) {
					return false, nil
				}
				as.CreateFunc = func(*model.Article) error {
					return nil
				}
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
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &MockUserStore{}
			mockArticleStore := &MockArticleStore{}

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.New(zerolog.NewTestWriter(t)),
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
			}
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
		setupMocks     func(*MockUserStore, *MockArticleStore)
		userID         uint
		req            *pb.UpdateArticleRequest
		expectedResult *pb.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Update an Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				as.GetByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1, Author: model.User{ID: 1}}, nil
				}
				as.UpdateFunc = func(article *model.Article) error {
					return nil
				}
				us.IsFollowingFunc = func(a, b *model.User) (bool, error) {
					return false, nil
				}
			},
			userID: 1,
			req: &pb.UpdateArticleRequest{
				Article: &pb.UpdateArticleRequest_Article{
					Slug:        "1",
					Title:       "Updated Title",
					Description: "Updated Description",
					Body:        "Updated Body",
				},
			},
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:        "1",
					Title:       "Updated Title",
					Description: "Updated Description",
					Body:        "Updated Body",
					Author:      &pb.Profile{},
					Favorited:   true,
				},
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &MockUserStore{}
			mockArticleStore := &MockArticleStore{}

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			ctx := context.WithValue(context.Background(), auth.UserIDKey, tt.userID)

			result, err := h.UpdateArticle(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Article.Slug, result.Article.Slug)
				assert.Equal(t, tt.expectedResult.Article.Title, result.Article.Title)
				assert.Equal(t, tt.expectedResult.Article.Description, result.Article.Description)
				assert.Equal(t, tt.expectedResult.Article.Body, result.Article.Body)
			}
		})
	}
}
