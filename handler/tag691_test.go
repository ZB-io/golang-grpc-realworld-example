package handler

import (
	"context"
	"testing"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockHandler is a struct used to simulate behavior for testing purposes.
type MockHandler struct {
	getTagsFunc func() ([]model.Tag, error)
}

func (m *MockHandler) GetTags() ([]model.Tag, error) {
	return m.getTagsFunc()
}

/*
ROOST_METHOD_HASH=GetTags_42221e4328
ROOST_METHOD_SIG_HASH=GetTags_52f72598a3
*/

func TestGetTags(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *Handler
		context   context.Context
		expected  *pb.TagsResponse
		expectErr error
	}{
		{
			name: "Success - Tags are fetched successfully",
			mockSetup: func() *Handler {
				return &Handler{as: &MockHandler{
					getTagsFunc: func() ([]model.Tag, error) {
						return []model.Tag{{Name: "go"}, {Name: "grpc"}}, nil
					},
				}}
			},
			context:   context.Background(),
			expected:  &pb.TagsResponse{Tags: []string{"go", "grpc"}},
			expectErr: nil,
		},
		{
			name: "Error - Internal server error when fetching tags fails",
			mockSetup: func() *Handler {
				return &Handler{as: &MockHandler{
					getTagsFunc: func() ([]model.Tag, error) {
						return nil, status.Error(codes.Aborted, "internal server error")
					},
				}}
			},
			context:   context.Background(),
			expected:  nil,
			expectErr: status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "Edge Case - No tags available",
			mockSetup: func() *Handler {
				return &Handler{as: &MockHandler{
					getTagsFunc: func() ([]model.Tag, error) {
						return []model.Tag{}, nil
					},
				}}
			},
			context:   context.Background(),
			expected:  &pb.TagsResponse{Tags: []string{}},
			expectErr: nil,
		},
		{
			name: "Error - Invalid context",
			mockSetup: func() *Handler {
				return &Handler{as: &MockHandler{
					getTagsFunc: func() ([]model.Tag, error) {
						return []model.Tag{}, nil
					},
				}}
			},
			expected:  nil,
			expectErr: context.Canceled,
		},
		{
			name: "Stress Test - Large number of tags",
			mockSetup: func() *Handler {
				return &Handler{as: &MockHandler{
					getTagsFunc: func() ([]model.Tag, error) {
						var tags []model.Tag
						for i := 0; i < 10000; i++ {
							tags = append(tags, model.Tag{Name: "tag" + strconv.Itoa(i)})
						}
						return tags, nil
					},
				}}
			},
			context: context.Background(),
			expected: func() *pb.TagsResponse {
				var tags []string
				for i := 0; i < 10000; i++ {
					tags = append(tags, "tag"+strconv.Itoa(i))
				}
				return &pb.TagsResponse{Tags: tags}
			}(),
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.mockSetup()

			if tt.name == "Error - Invalid context" {
				var cancel context.CancelFunc
				tt.context, cancel = context.WithCancel(context.Background())
				cancel()
			}

			got, err := handler.GetTags(tt.context, &pb.Empty{})

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			t.Log("Completed:", tt.name)
		})
	}
}
