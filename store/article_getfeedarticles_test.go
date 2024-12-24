package store

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreGetFeedArticles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE EXTENSION").WillReturnResult(sqlmock.NewResult(1, 1))

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
	}
	s := &ArticleStore{db: gdb}

	tests := []struct {
		name    string
		s       *ArticleStore
		userIDs []uint
		limit   int64
		offset  int64
		mock    func()
		want    []model.Article
		wantErr bool
	}{
		{
			name:    "success",
			s:       s,
			userIDs: []uint{1, 2, 3},
			limit:   5,
			offset:  0,
			mock: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"deleted_at\" IS NULL AND \\(user_id in \\(\\?\\,\\?\\,\\?\\)\\) ORDER BY \"articles\".\"id\" ASC LIMIT 5").WithArgs(1, 2, 3).WillReturnRows(sqlmock.NewRows(nil))
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "invalid userIDs",
			s:       s,
			userIDs: nil,
			limit:   5,
			offset:  0,
			mock:    func() {},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid limit",
			s:       s,
			userIDs: []uint{1, 2, 3},
			limit:   -5,
			offset:  0,
			mock:    func() {},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "DB error",
			s:       s,
			userIDs: []uint{1, 2, 3},
			limit:   5,
			offset:  0,
			mock: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"deleted_at\" IS NULL AND \\(user_id in \\(\\?\\,\\?\\,\\?\\)\\) ORDER BY \"articles\".\"id\" ASC LIMIT 5").WithArgs(1, 2, 3).WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "zero offset",
			s:       s,
			userIDs: []uint{1, 2, 3},
			limit:   5,
			offset:  0,
			mock: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"deleted_at\" IS NULL AND \\(user_id in \\(\\?\\,\\?\\,\\?\\)\\) ORDER BY \"articles\".\"id\" ASC LIMIT 5").WithArgs(1, 2, 3).WillReturnRows(sqlmock.NewRows(nil))
			},
			want:    []model.Article{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetFeedArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetFeedArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}
