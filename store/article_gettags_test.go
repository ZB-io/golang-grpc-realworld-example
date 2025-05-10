package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"strconv"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestGetTags(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when creating the mock database", err)
	}
	defer mockDB.Close()
	gormDB, err := gorm.Open("sqlite3", mockDB)
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}
	defer gormDB.Close()

	articleStore := ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedTags  []model.Tag
		expectedError bool
	}{
		{
			name: "Successfully Retrieve Tag List",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test").
					AddRow(2, "example")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{Name: "test"},
				{Name: "example"},
			},
			expectedError: false,
		},
		{
			name: "No Tags in Database",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: false,
		},
		{
			name: "Database Retrieval Error",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedTags:  nil,
			expectedError: true,
		},
		{
			name: "Duplicate Tag Names",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "duplicate").
					AddRow(2, "duplicate")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{Name: "duplicate"},
				{Name: "duplicate"},
			},
			expectedError: false,
		},
		{
			name: "Performance with Large Data Set",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})

				for i := 1; i <= 1000; i++ {
					rows.AddRow(i, "tag"+strconv.Itoa(i))
				}
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)

				expected := make([]model.Tag, 1000)
				for i := range expected {
					expected[i].Name = "tag" + strconv.Itoa(i+1)
				}
				tests[4].expectedTags = expected
			},
			expectedTags:  nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			tags, err := articleStore.GetTags()
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, tags, "The tags slice should be nil on error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTags, tags, "Tags fetched do not match expected result")
			}
		})
	}
}
