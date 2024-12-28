package handler

import (
	"context"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/store"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}

type ArticleStore struct {
	db *gorm.DB
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
	context  *testContext // For running tests and subtests.
}
func TestHandlerGetTags(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	logger := zerolog.New(nil)
	as := &store.ArticleStore{
		DB: mockDB,
	}
	handler := Handler{
		logger: &logger,
		as:     as,
	}

	tests := []struct {
		name      string
		mockSetup func()
		expected  *pb.TagsResponse
		code      codes.Code
	}{
		{
			name: "Retrieve Tags Successfully",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM tags").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).
						AddRow("Go").
						AddRow("Programming"))
			},
			expected: &pb.TagsResponse{
				Tags: []string{"Go", "Programming"},
			},
			code: codes.OK,
		},
		{
			name: "No Tags Available",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM tags").
					WillReturnRows(sqlmock.NewRows([]string{"name"}))
			},
			expected: &pb.TagsResponse{
				Tags: []string{},
			},
			code: codes.OK,
		},
		{
			name: "Article Store Returns an Error",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM tags").WillReturnError(status.Error(codes.Aborted, "database error"))
			},
			expected: nil,
			code:     codes.Aborted,
		},
		{
			name: "Large Number of Tags",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"name"})
				for i := 0; i < 1000; i++ {
					rows.AddRow("Tag" + string(i))
				}
				mock.ExpectQuery("SELECT \\* FROM tags").WillReturnRows(rows)
			},
			expected: &pb.TagsResponse{
				Tags: func() []string {
					tags := make([]string, 1000)
					for i := 0; i < 1000; i++ {
						tags[i] = "Tag" + string(i)
					}
					return tags
				}(),
			},
			code: codes.OK,
		},
		{
			name: "Duplicate Tag Names",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM tags").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).
						AddRow("Go").
						AddRow("Go"))
			},
			expected: &pb.TagsResponse{
				Tags: []string{"Go", "Go"},
			},
			code: codes.OK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			resp, err := handler.GetTags(context.Background(), &pb.Empty{})
			if tc.code == codes.OK {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, resp)
				t.Logf("Test '%s' succeeded as expected.", tc.name)
			} else {
				require.Error(t, err)
				statusErr, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.code, statusErr.Code())
				t.Logf("Test '%s' failed as expected with error: %v.", tc.name, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
