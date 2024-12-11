package store

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type MockSql struct {
	DB      *sql.DB
	Mock    sqlmock.Sqlmock
	Request *http.Request
}
/*
type UserStore struct {
	db *gorm.DB
}
type DB struct {
	sync.RWMutex
	Value			interface{}
	Error			error
	RowsAffected		int64
	db			SQLCommon
	blockGlobalUpdate	bool
	logMode			logModeValue
	logger			logger
	search			*search
	values			sync.Map
	parent			*DB
	callbacks		*Callback
	dialect			Dialect
	singularTable		bool
	nowFuncOverride		func() time.Time
}// single db
// function to be used to override the creating of a new timestamp


type User struct {
	gorm.Model
	Username		string		`gorm:"unique_index;not null"`
	Email			string		`gorm:"unique_index;not null"`
	Password		string		`gorm:"not null"`
	Bio			string		`gorm:"not null"`
	Image			string		`gorm:"not null"`
	Follows			[]User		`gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles	[]Article	`gorm:"many2many:favorite_articles;"`
}


ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920

func (s *UserStore) Create(m *model.User) error {
	return s.db.Create(m).Error
}
*/

func TestUserStoreCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		dbMock  func() (sqlmock.Sqlmock, *gorm.DB)
		wantErr bool
	}{
		{
			name: "Valid User",
			user: &model.User{Username: "Test", Email: "test@test.com"},
			dbMock: func() (sqlmock.Sqlmock, *gorm.DB) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
				gormdb, _ := gorm.Open("postgres", db)
				return mock, gormdb
			},
			wantErr: false,
		},
		{
			name: "Invalid User",
			user: &model.User{Username: "", Email: ""},
			dbMock: func() (sqlmock.Sqlmock, *gorm.DB) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
				mock.ExpectExec("^INSERT INTO `users` .*").WillReturnError(errors.New("cannot create user"))
				gormdb, _ := gorm.Open("postgres", db)
				return mock, gormdb
			},
			wantErr: true,
		},
		{
			name: "Nil User",
			user: nil,
			dbMock: func() (sqlmock.Sqlmock, *gorm.DB) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
				mock.ExpectExec("^INSERT INTO `users` .*").WillReturnError(gorm.ErrRecordNotFound)
				gormdb, _ := gorm.Open("postgres", db)
				return mock, gormdb
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, gormDB := tt.dbMock()
			defer func() {
				gormDB.Close()
			}()

			us := &UserStore{db: gormDB}

			err := us.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

