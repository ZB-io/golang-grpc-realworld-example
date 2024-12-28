package handler

import (
	"context"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type MockHandler struct {
	getTagsFunc func() ([]Tag, error)
}

func (m *MockHandler) GetTags() ([]Tag, error) {
	return m.getTagsFunc()
}
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
				handler := &Handler{as: &MockHandler{
					getTagsFunc: func() ([]Tag, error) {
						return []Tag{{Name: "go"}, {Name: "grpc"}}, nil
					},
				}}
				return handler
			},
			context:   context.Background(),
			expected:  &pb.TagsResponse{Tags: []string{"go", "grpc"}},
			expectErr: nil,
		},
		{
			name: "Error - Internal server error when fetching tags fails",
			mockSetup: func() *Handler {
				handler := &Handler{as: &MockHandler{
					getTagsFunc: func() ([]Tag, error) {
						return nil, sqlmock.ErrCancelled
					},
				}}
				return handler
			},
			context:   context.Background(),
			expected:  nil,
			expectErr: status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "Edge Case - No tags available",
			mockSetup: func() *Handler {
				handler := &Handler{as: &MockHandler{
					getTagsFunc: func() ([]Tag, error) {
						return []Tag{}, nil
					},
				}}
				return handler
			},
			context:   context.Background(),
			expected:  &pb.TagsResponse{Tags: []string{}},
			expectErr: nil,
		},
		{
			name: "Error - Invalid context",
			mockSetup: func() *Handler {
				handler := &Handler{as: &MockHandler{
					getTagsFunc: func() ([]Tag, error) {
						return []Tag{}, nil
					},
				}}
				return handler
			},
			context:   context.Background(),
			expected:  nil,
			expectErr: context.Canceled,
		},
		{
			name: "Stress Test - Large number of tags",
			mockSetup: func() *Handler {
				handler := &Handler{as: &MockHandler{
					getTagsFunc: func() ([]Tag, error) {
						var tags []Tag
						for i := 0; i < 10000; i++ {
							tags = append(tags, Tag{Name: "tag" + string(i)})
						}
						return tags, nil
					},
				}}
				return handler
			},
			context: context.Background(),
			expected: func() *pb.TagsResponse {
				var tags []string
				for i := 0; i < 10000; i++ {
					tags = append(tags, "tag"+string(i))
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
				tt.context, cancel = context.WithCancel(tt.context)
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


func (m *MockHandler) GetTags() ([]Tag, error) {
	return m.getTagsFunc()
}
