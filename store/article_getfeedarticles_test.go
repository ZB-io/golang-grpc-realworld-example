package store_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func openTestDb() (*gorment.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, nil, err
	}

	return gdb, mock, nil
}

func TestArticleStoreGetFeedArticles(t *testing.T) {
	db, mock, err := openTestDb()
	if err != nil {
		t.Fatalf("error during opening database connection %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"users\".*").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	tests := []struct {
		name    string
		userIDs []uint
		limit   int64
		offset  int64
		wantErr bool
		err     error
	}{
		{
			name:    "Normal Operation",
			userIDs: []uint{1},
			limit:   5,
			offset:  0,
			wantErr: false,
		},
		{
			name:    "Empty user IDs array",
			userIDs: []uint{},
			limit:   5,
			offset:  0,
			wantErr: false,
		},
        {
			name:    "No articles matching user IDs",
			userIDs: []uint{555},
			limit:   5,
			offset:  0,
			wantErr: false,
		},
        {
			name:    "Database fetch error",
			userIDs: []uint{1},
			limit:   5,
			offset:  0,
			wantErr: true,
			err:     errors.New("Fetch error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

            if err != nil && err != tt.err {
				t.Errorf("GetFeedArticles() error = %v, expected %v", err, tt.err)
				return
			}

			if tt.wantErr && tt.err.Error() != err.Error() {
				t.Errorf("GetFeedArticles() error = %v, wantErr %v", err, tt.err)
				return
			}

            if len(got) != len(tt.articles) {
				t.Errorf("GetFeedArticles() got %v, want %v", len(got), len(tt.articles))
			};
		})
	}
}
