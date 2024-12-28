package store

import (
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)





func TestArticleStoreGetTags(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedTags  []model.Tag
		expectedError bool
	}{
		{
			name: "Retrieve a List of Tags Successfully",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ID", "name"}).
					AddRow(1, "tag1").
					AddRow(2, "tag2")
				mock.ExpectQuery("^SELECT \\* FROM .tags.$").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
			},
			expectedError: false,
		},
		{
			name: "Return an Empty List When No Tags Exist",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ID", "name"})
				mock.ExpectQuery("^SELECT \\* FROM .tags.$").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: false,
		},
		{
			name: "Handle Database Error Gracefully",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM .tags.$").WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedTags:  []model.Tag{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening gorm DB", err)
			}

			store := &ArticleStore{db: gormDB}

			tt.mockSetup(mock)

			tags, err := store.GetTags()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTags, tags)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Log("Test completed for scenario:", tt.name)
		})
	}
}
func (s *ArticleStore) GetTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return tags, err
	}
	return tags, nil
}
