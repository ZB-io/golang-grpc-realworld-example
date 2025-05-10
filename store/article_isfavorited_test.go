package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestIsFavorited(t *testing.T) {
	tests := []struct {
		name          string
		article       *model.Article
		user          *model.User
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedFav   bool
		expectedError bool
	}{
		{
			name: "Scenario 1: Article and User Exist and User Has Favorited the Article",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(*) FROM favorite_articles").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedFav:   true,
			expectedError: false,
		},
		{
			name: "Scenario 2: Article and User Exist but User Has Not Favorited the Article",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			user: &model.User{
				Model: gorm.Model{ID: 2},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(*) FROM favorite_articles").
					WithArgs(2, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedFav:   false,
			expectedError: false,
		},
		{
			name:    "Scenario 3: Article Object is Nil",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 3},
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedFav:   false,
			expectedError: false,
		},
		{
			name: "Scenario 4: User Object is Nil",
			article: &model.Article{
				Model: gorm.Model{ID: 4},
			},
			user:          nil,
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedFav:   false,
			expectedError: false,
		},
		{
			name: "Scenario 5: Database Connection Error",
			article: &model.Article{
				Model: gorm.Model{ID: 5},
			},
			user: &model.User{
				Model: gorm.Model{ID: 5},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(*) FROM favorite_articles").
					WithArgs(5, 5).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectedFav:   false,
			expectedError: true,
		},
		{
			name: "Scenario 6: Article ID or User ID is Zero",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			user: &model.User{
				Model: gorm.Model{ID: 0},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count(*) FROM favorite_articles").
					WithArgs(0, 0).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedFav:   false,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm DB connection: %v", err)
			}

			store := &ArticleStore{db: gormDB}

			tt.mockSetup(mock)

			favorited, err := store.IsFavorited(tt.article, tt.user)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if favorited != tt.expectedFav {
				t.Errorf("expected favored: %v, got: %v", tt.expectedFav, favorited)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
