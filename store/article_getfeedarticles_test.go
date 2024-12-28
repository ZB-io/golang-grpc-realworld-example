package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
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

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestArticleStoreGetFeedArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open("postgres", db)
	assert.NoError(t, err)

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		userIDs       []uint
		limit         int64
		offset        int64
		mockSetup     func()
		expectedError bool
		expectedCount int
	}{
		{
			name:    "Retrieve Articles for Empty User List",
			userIDs: []uint{},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}))
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:    "Retrieve Articles with Limit and Offset",
			userIDs: []uint{1},
			limit:   2,
			offset:  1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
					AddRow(2, "Title 2", "Description 2", "Body 2", 1).
					AddRow(3, "Title 3", "Description 3", "Body 3", 1)
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE user_id in (.+) LIMIT 2 OFFSET 1`).WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:    "Invalid User IDs",
			userIDs: []uint{999},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE user_id in (.+)`).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}))
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:    "Retrieve Articles Across Multiple Users",
			userIDs: []uint{1, 2},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
					AddRow(1, "Title 1", "Description 1", "Body 1", 1).
					AddRow(2, "Title 2", "Description 2", "Body 2", 2)
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE user_id in (.+)`).WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:    "Database Error Handling",
			userIDs: []uint{1},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE user_id in (.+)`).WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedError: true,
		},
		{
			name:    "Preloading Author Relation",
			userIDs: []uint{1},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
					AddRow(1, "Title 1", "Description 1", "Body 1", 1)
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE user_id in (.+)`).WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			articles, err := store.GetFeedArticles(tc.userIDs, tc.limit, tc.offset)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(articles))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %s", err)
			}
		})
	}
}
