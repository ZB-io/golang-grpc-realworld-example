package store

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {

	oldIn := os.Stdin
	oldOut := os.Stdout
	defer func() {
		os.Stdin = oldIn
		os.Stdout = oldOut
	}()
	var stdout, stdin bytes.Buffer
	os.Stdin = io.Reader(&stdin)
	os.Stdout = io.Writer(&stdout)

	scenarios := []struct {
		description string
		setup       func(mock sqlmock.Sqlmock, article *model.Article)
		article     model.Article
		shouldError bool
	}{
		{
			description: "Article Creation with Valid Parameters",
			setup: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `articles`").WithArgs(article.Title, article.Description, article.Body).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article:     model.Article{Title: "Test Title", Description: "Test Description", Body: "Test Body"},
			shouldError: false,
		},
		{
			description: "Article Creation with Invalid Parameters",
			setup:       func(mock sqlmock.Sqlmock, article *model.Article) {},
			article:     model.Article{Body: "Test Body"},
			shouldError: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open the DB connection: %s", err)
			}

			s.setup(mock, &s.article)

			store := ArticleStore{db: gormDB}
			e := store.Create(&s.article)

			if !s.shouldError {
				assert.NoError(t, e, fmt.Sprintf("%s: error was not expected, but got %v", s.description, e))
			} else {
				assert.Error(t, e, fmt.Sprintf("%s: error was expected", s.description))
			}
		})
	}
}

func (s *UserStore) Create(m *model.User) error {
	return s.db.Create(m).Error
}

