package store

import (
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func TestArticleStoreGetCommentByID(t *testing.T) {
	type testScenario struct {
		description string
		setup       func() (*gorm.DB, sqlmock.Sqlmock)
		id          uint
		expected    *model.Comment
		expectErr   bool
	}

	mockComment := &model.Comment{
		ID:        1,
		Body:      "Sample comment",
		ArticleID: 1,
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []testScenario{
		{
			description: "Successfully Retrieve an Existing Comment by ID",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("unexpected error when opening a stub database connection: %v", err)
				}
				gormDB, err := gorm.Open(postgres.New(postgres.Config{
					Conn: db,
				}), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				if err != nil {
					t.Fatalf("failed to open gorm db: %v", err)
				}

				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE "comments"\."id" = \$1`).
					WithArgs(mockComment.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "article_id", "author_id", "created_at", "updated_at"}).
						AddRow(mockComment.ID, mockComment.Body, mockComment.ArticleID, mockComment.AuthorID, mockComment.CreatedAt, mockComment.UpdatedAt))

				return gormDB, mock
			},
			id:        mockComment.ID,
			expected:  mockComment,
			expectErr: false,
		},
		{
			description: "Return Error for Non-existent Comment ID",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("unexpected error when opening a stub database connection: %v", err)
				}
				gormDB, err := gorm.Open(postgres.New(postgres.Config{
					Conn: db,
				}), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				if err != nil {
					t.Fatalf("failed to open gorm db: %v", err)
				}

				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE "comments"\."id" = \$1`).
					WithArgs(99).
					WillReturnError(gorm.ErrRecordNotFound)

				return gormDB, mock
			},
			id:        99,
			expected:  nil,
			expectErr: true,
		},
		{
			description: "Return Error When Database Connection Fails",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("unexpected error when opening a stub database connection: %v", err)
				}
				gormDB, err := gorm.Open(postgres.New(postgres.Config{
					Conn: db,
				}), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				if err != nil {
					t.Fatalf("failed to open gorm db: %v", err)
				}

				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE "comments"\."id" = \$1`).
					WithArgs(1).
					WillReturnError(gorm.ErrInvalidTransaction)

				return gormDB, mock
			},
			id:        1,
			expected:  nil,
			expectErr: true,
		},
		{
			description: "Retrieve Comment with Minimum (Boundary) Valid ID",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("unexpected error when opening a stub database connection: %v", err)
				}
				gormDB, err := gorm.Open(postgres.New(postgres.Config{
					Conn: db,
				}), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				if err != nil {
					t.Fatalf("failed to open gorm db: %v", err)
				}

				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE "comments"\."id" = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "body", "article_id", "author_id", "created_at", "updated_at"}).
						AddRow(1, "Boundary comment", 1, 1, mockComment.CreatedAt, mockComment.UpdatedAt))

				return gormDB, mock
			},
			id:        1,
			expected:  mockComment,
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			gormDB, mock := tc.setup()
			db, _ := gormDB.DB()
			defer db.Close()

			store := &ArticleStore{db: gormDB}
			comment, err := store.GetCommentByID(tc.id)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tc.expected, comment)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

