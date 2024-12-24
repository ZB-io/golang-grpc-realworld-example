package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
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
func TestArticleStoreGetTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error during mock database creation: %v", err)
	}
	gormDB, gormErr := gorm.Open("postgres", db)
	if gormErr != nil {
		t.Fatalf("failed opening gorm database: %v", gormErr)
	}

	articleStore := &store.ArticleStore{db: gormDB}

	tests := []struct {
		name     string
		setup    func()
		expected []model.Tag
		wantErr  bool
	}{
		{
			name: "Retrieve All Tags Successfully",
			setup: func() {
				mock.ExpectQuery("SELECT (.+) FROM tags").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).AddRow(1, time.Now(), time.Now(), nil, "tag1").AddRow(2, time.Now(), time.Now(), nil, "tag2"))
			},
			expected: []model.Tag{
				{
					Model: gorm.Model{ID: 1},
					Name:  "tag1",
				},
				{
					Model: gorm.Model{ID: 2},
					Name:  "tag2",
				},
			},
			wantErr: false,
		},
		{
			name: "Database is Empty",
			setup: func() {
				mock.ExpectQuery("SELECT (.+) FROM tags").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}))
			},
			expected: []model.Tag{},
			wantErr:  false,
		},
		{
			name: "An Error Occurs During Retrieval",
			setup: func() {
				mock.ExpectQuery("SELECT (.+) FROM tags").WillReturnError(errors.New("failed to retrieve tags"))
			},
			wantErr: true,
		},
		{
			name: "Database has Duplicate Entries",
			setup: func() {
				mock.ExpectQuery("SELECT (.+) FROM tags").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).AddRow(1, time.Now(), time.Now(), nil, "tag1").AddRow(2, time.Now(), time.Now(), nil, "tag1"))
			},
			expected: []model.Tag{
				{
					Model: gorm.Model{ID: 1},
					Name:  "tag1",
				},
				{
					Model: gorm.Model{ID: 2},
					Name:  "tag1",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			actual, err := articleStore.GetTags()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %v", err)
		}
	}
}
