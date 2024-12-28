package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)






func TestArticleStoreGetByID(t *testing.T) {

	scenarios := []struct {
		description  string
		mockSetup    func(mock sqlmock.Sqlmock)
		id           uint
		expectedErr  string
		expectedData *model.Article
	}{
		{
			description: "Successfully Retrieve an Article by ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(1, "Test Title")
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			id:           1,
			expectedErr:  "",
			expectedData: &model.Article{ID: 1, Title: "Test Title"},
		},
		{
			description: "Article Not Found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			id:           999,
			expectedErr:  gorm.ErrRecordNotFound.Error(),
			expectedData: nil,
		},
		{
			description: "Database Connection Failure",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(1).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			id:           1,
			expectedErr:  gorm.ErrInvalidSQL.Error(),
			expectedData: nil,
		},
		{
			description: "Preload Failure for Tags or Author",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
				mock.ExpectQuery("^SELECT (.+) FROM tags WHERE article_id = ?").
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			id:           1,
			expectedErr:  gorm.ErrRecordNotFound.Error(),
			expectedData: nil,
		},
		{
			description: "Edge Case with Minimum Possible ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(1, "Minimum ID Test Title")
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			id:           1,
			expectedErr:  "",
			expectedData: &model.Article{ID: 1, Title: "Minimum ID Test Title"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error opening a stub database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("error opening gorm db: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{db: gormDB}

			scenario.mockSetup(mock)

			article, err := store.GetByID(scenario.id)

			if scenario.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), scenario.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, scenario.expectedData, article)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Log("Test completed with description: ", scenario.description)
		})
	}
}
func (s *ArticleStore) GetByID(id uint) (*model.Article, error) {
	var m model.Article
	err := s.db.Preload("Tags").Preload("Author").Find(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}
