package store

import (
	"fmt"
	"testing"
	
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreGetComments(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gdb, _ := gorm.Open("postgres", db)

	// Test data
	article := model.Article{
		Title: "Test article",
	}
	comment1 := model.Comment{
		Body: "First comment",
	}
	comment2 := model.Comment{
		Body: "Second comment",
	}
	allComments := []model.Comment{comment1, comment2}
	emptyComments := []model.Comment{}

	tests := []struct {
		name          string
		mockDbFunc    func() error
		article       model.Article
		expected      []model.Comment
		expectedError bool
	}{
		{
			name: "Retrieve Comments for Given Article",
			mockDbFunc: func() error {
				rows := sqlmock.
					NewRows([]string{"id", "body"}).
					AddRow(1, comment1.Body).
					AddRow(2, comment2.Body)

				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE (.+)$").WillReturnRows(rows)
				return nil
			},
			article:       article,
			expected:      allComments,
			expectedError: false,
		},
		{
			name: "Retrieve Comments for an Article with No Comments",
			mockDbFunc: func() error {
				rows := sqlmock.NewRows([]string{"id", "body"})
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE (.+)$").WillReturnRows(rows)
				return nil
			},
			article:       article,
			expected:      emptyComments,
			expectedError: false,
		},
		{
			name: "Retrieve Comments for Non-Existing Article",
			mockDbFunc: func() error {
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\"").WillReturnError(fmt.Errorf("database error"))
				return nil
			},
			article:       model.Article{Title: "Non-Existing Article"},
			expectedError: true,
		},
		{
			name: "Handle Database Connection Error",
			mockDbFunc: func() error {
				return fmt.Errorf("database connection error")
			},
			article:       article,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(`Scenario: ` + tt.name)
			
			store := &ArticleStore{db: gdb}

			// Mocking DB
			err := tt.mockDbFunc()

			if err != nil {
				if tt.expectedError {
					t.Logf(`Expected error: "%v"`, err)
				} else {
					t.Errorf(`Error should not have occurred: "%v"`, err)
				}
				return
			}

			result, err := store.GetComments(&tt.article)

			if err != nil && !tt.expectedError {
				t.Errorf(`Unexpected error: "%v"`, err)
				return
			}

			if tt.expectedError {
				t.Logf(`Expected error occurred: "%v"`, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf(`Count mismatch: got %v, expected %v`, len(result), len(tt.expected))
				return
			}

			for i, comment := range result {
				if comment.Body != tt.expected[i].Body {
					t.Errorf(`Body mismatch: got "%v", expected "%v"`, comment.Body, tt.expected[i].Body)
				}
			}

			t.Logf(`Success: expected result matches actual`)
		})
	}
}
