package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)






func TestGetByID(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		id        uint
		wantUser  *model.User
		wantErr   error
	}{
		{
			name: "Successful Retrieval of User by Valid ID",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(1).WillReturnRows(rows)
			},
			id:       1,
			wantUser: &model.User{ID: 1, Username: "testuser", Email: "test@example.com"},
			wantErr:  nil,
		},
		{
			name: "User Not Found in Database",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(2).WillReturnError(gorm.ErrRecordNotFound)
			},
			id:       2,
			wantUser: nil,
			wantErr:  gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(3).WillReturnError(errors.New("connection error"))
			},
			id:       3,
			wantUser: nil,
			wantErr:  errors.New("connection error"),
		},
		{
			name: "Retrieval with Invalid ID Input (Zero ID)",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(0).WillReturnError(errors.New("invalid query condition"))
			},
			id:       0,
			wantUser: nil,
			wantErr:  errors.New("invalid query condition"),
		},
		{
			name: "Retrieval with Invalid ID Input (Negative ID)",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(-1).WillReturnError(errors.New("invalid query condition"))
			},
			id:       -1,
			wantUser: nil,
			wantErr:  errors.New("invalid query condition"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			if err != nil {
				t.Fatalf("Failed to initialize GORM database: %v", err)
			}

			userStore := &UserStore{db: gdb}

			gotUser, gotErr := userStore.GetByID(tt.id)

			if tt.wantUser != nil && gotUser != nil {
				if *tt.wantUser != *gotUser {
					t.Errorf("got user %v, want %v", gotUser, tt.wantUser)
				}
			} else if tt.wantUser != gotUser {
				t.Errorf("got user %v, want %v", gotUser, tt.wantUser)
			}

			if (gotErr != nil && tt.wantErr == nil) || (gotErr == nil && tt.wantErr != nil) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			} else if gotErr != nil && gotErr.Error() != tt.wantErr.Error() {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

