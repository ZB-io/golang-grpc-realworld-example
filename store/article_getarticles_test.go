package store

import (
	"fmt"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)






type mockGormDB struct {
	*gorm.DB
	sqlmock.Sqlmock
}



func TestGetArticles(t *testing.T) {
	db, mock, err := newMockGormDB()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	store := ArticleStore{db: db}
	user := &model.User{ID: 1, Username: "testuser"}

	tests := []struct {
		name         string
		tagName      string
		username     string
		favoritedBy  *model.User
		limit        int64
		offset       int64
		expectErr    bool
		mockSetup    func(sqlmock.Sqlmock)
		verifyResult func([]model.Article) error
	}{
		{
			name:      "Retrieve Articles by Author Username",
			username:  "author1",
			expectErr: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles JOIN users").
					WithArgs("author1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "body"}).
						AddRow(1, "Article1", "Body1").
						AddRow(2, "Article2", "Body2"))
			},
			verifyResult: func(articles []model.Article) error {
				if len(articles) != 2 {
					return fmt.Errorf("expected 2 articles, got %d", len(articles))
				}
				for _, article := range articles {
					if article.Author.Username != "author1" {
						return fmt.Errorf("expected article by author1, got %s", article.Author.Username)
					}
				}
				return nil
			},
		},
		{
			name:      "Retrieve Articles Tagged with a Specific Tag",
			tagName:   "tag1",
			expectErr: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles JOIN").
					WithArgs("tag1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "body"}).
						AddRow(3, "Article3", "Body3"))
			},
			verifyResult: func(articles []model.Article) error {
				if len(articles) != 1 {
					return fmt.Errorf("expected 1 article, got %d", len(articles))
				}
				return nil
			},
		},
		{
			name:        "Retrieve Articles Favorited by a Specific User",
			favoritedBy: user,
			expectErr:   false,
			limit:       2,
			offset:      0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM favorite_articles WHERE user_id = ?").
					WithArgs(user.ID).
					WillReturnRows(sqlmock.NewRows([]string{"article_id"}).
						AddRow(1).
						AddRow(2))
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE id in").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "body"}).
						AddRow(1, "Article1", "Body1").
						AddRow(2, "Article2", "Body2"))
			},
			verifyResult: func(articles []model.Article) error {
				if len(articles) != 2 {
					return fmt.Errorf("expected 2 articles, got %d", len(articles))
				}
				return nil
			},
		},
		{
			name:      "Error Handling When Database Returns an Error",
			expectErr: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(fmt.Errorf("db error"))
			},
			verifyResult: func(articles []model.Article) error {
				if len(articles) != 0 {
					return fmt.Errorf("expected no articles due to error, got %d", len(articles))
				}
				return nil
			},
		},
		{
			name:      "Retrieve Articles with Specific Pagination",
			limit:     1,
			offset:    1,
			expectErr: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles LIMIT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "body"}).
						AddRow(2, "Article2", "Body2"))
			},
			verifyResult: func(articles []model.Article) error {
				if len(articles) != 1 {
					return fmt.Errorf("expected 1 article due to limit, got %d", len(articles))
				}
				if articles[0].ID != 2 {
					return fmt.Errorf("expected article ID 2, got %d", articles[0].ID)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := tt.verifyResult(articles); err != nil {
				t.Error(err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func newMockGormDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}
func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	d := s.db.Preload("Author")

	if username != "" {
		d = d.Joins("join users on articles.user_id = users.id").
			Where("users.username = ?", username)
	}

	if tagName != "" {
		d = d.Joins(
			"join article_tags on articles.id = article_tags.article_id "+
				"join tags on tags.id = article_tags.tag_id").
			Where("tags.name = ?", tagName)
	}

	if favoritedBy != nil {
		rows, err := s.db.Select("article_id").
			Table("favorite_articles").
			Where("user_id = ?", favoritedBy.ID).
			Offset(offset).Limit(limit).Rows()
		if err != nil {
			return []model.Article{}, err
		}
		defer rows.Close()

		var ids []uint
		for rows.Next() {
			var id uint
			rows.Scan(&id)
			ids = append(ids, id)
		}
		d = d.Where("id in (?)", ids)
	}

	d = d.Offset(offset).Limit(limit)

	var as []model.Article
	err := d.Find(&as).Error

	return as, err
}
