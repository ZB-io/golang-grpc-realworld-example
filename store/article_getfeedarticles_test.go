// assume the current package as "store"
package store

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

const (
	Limit  = 5
	Offset = 0
)

// assuming these exist somewhere
var validUserIDs = []uint{1, 2, 3, 4, 5}
var invalidUserIDs = []uint{10, 20, 30, 40, 50}

// This is the test function for GetFeedArticles implemented as table-driven tests
func TestArticleStoreGetFeedArticles(t *testing.T) {
	// Mock DB using sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	// Define test scenarios
	tests := []struct {
		name         string
		userIDs      []uint
		limit        int64
		offset       int64
		setup        func()
		expectedErr  error
		expectedResp []model.Article
	}{
		{
			name:         "Scenario 1: Valid UserIDs and valid limit and offset values",
			userIDs:      validUserIDs,
			limit:        Limit,
			offset:       Offset,
			expectedResp: []model.Article{},
			setup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `articles`  WHERE (`user_id` in (?)) LIMIT ? OFFSET ?")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))
			},
		},
		{
			name:    "Scenario 2: Invalid UserIDs",
			userIDs: invalidUserIDs,
			limit:   Limit,
			offset:  Offset,
			setup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `articles`  WHERE (`user_id` in (?)) LIMIT ? OFFSET ?")).
					WillReturnError(errors.New("record not found"))
			},
			expectedErr: errors.New("record not found"),
		},
		{
			name:    "Scenario 3: Edge Case of Zero Limit and Offset",
			userIDs: validUserIDs,
			limit:   0,
			offset:  0,
			setup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `articles`  WHERE (`user_id` in (?)) LIMIT ? OFFSET ?")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))
			},
			expectedResp: []model.Article{},
		},
		{
			name:    "Scenario 4: Error While Fetching the Articles from DB",
			userIDs: validUserIDs,
			limit:   Limit,
			offset:  Offset,
			setup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `articles`  WHERE (`user_id` in (?)) LIMIT ? OFFSET ?")).
					WillReturnError(errors.New("DB error"))
			},
			expectedErr: errors.New("DB error"),
		},
	}

	// Run test scenarios
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			as := &ArticleStore{db: gdb}
			resp, err := as.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResp, resp)
			if err != nil {
				t.Log("Failure reason: ", err.Error())
			}
			if err == nil {
				t.Log("Success: correct articles were returned")
			}
		})
	}
}
