package store

import (
	"database/sql"
	"testing"
	
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStoreGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gdb, err := gorm.Open("postgres", db)
	require.NoError(t, err)

	store := UserStore{db: gdb}
	email := "test@example.com"
	
	tests := []struct {
		name     string
		mock     func()
		wantErr  bool
	}{
		{
			name: "Successful retrieval of user by email",
			mock: func() {
				rows := sqlmock.
					NewRows([]string{"id", "email", "password"}).
					AddRow(1, email, "password")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
					WithArgs(email).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "User email not found in the database",
			mock: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
					WithArgs(email).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "Database connection error during execution",
			mock: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
					WithArgs(email).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := store.GetByEmail(email)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
