package store

import (
	"database/sql"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful Article Creation",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Empty Fields",
			article: &model.Article{
				Title:       "",
				Description: "",
				Body:        "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Duplicate Article",
			article: &model.Article{
				Title:       "Duplicate Title",
				Description: "Duplicate Description",
				Body:        "Duplicate Body",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Article with Tags",
			article: &model.Article{
				Title:       "Article with Tags",
				Description: "Description with Tags",
				Body:        "Body with Tags",
				Tags:        []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `article_tags`").
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

			err := store.Create(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Logf("Article created successfully: %+v", tt.article)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

