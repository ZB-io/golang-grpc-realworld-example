package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
)

// TestArticleStoreGetTags unit tests the GetTags function in ArticleStore.
func TestArticleStoreGetTags(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Scenario 1: Valid retrieval of tags",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(
					[]string{
						"id", "tag",
					},
				).AddRow(
					1, "article",
				)
				mock.ExpectQuery("SELECT (.+) FROM").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{{Model: gorm.Model{ID: 1}, Tag: "article"}},
		},
		{
			name: "Scenario 2: Handling of empty database",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "tag"})
				mock.ExpectQuery("SELECT (.+) FROM").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{},
		},
		{
			name: "Scenario 3: Handling of database retrieval errors",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error %v was not expected when opening stub database connection", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error %v was not expected when opening gorm database", err)
			}

			articleStore := ArticleStore{gdb}

			tags, err := articleStore.GetTags()

			if tt.expectedError != nil && err != tt.expectedError {
				t.Errorf("expected error %v, but got %v", tt.expectedError, err)
			}

			if len(tags) != len(tt.expectedTags) {
				t.Errorf("expected %v tags, but got %v", len(tt.expectedTags), len(tags))
			}

			for i, tag := range tags {
				if tag != tt.expectedTags[i] {
					t.Errorf("expected tag %v, but got %v", tt.expectedTags[i], tag)
				}
			}
		})
	}
}
