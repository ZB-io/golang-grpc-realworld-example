package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"database/sql"
	"regexp"
	"sync"
)

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestCreateUser(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		setupFn func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful user creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("Error 1062: Duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "Duplicate entry",
		},
		{
			name: "Missing required fields",
			user: &model.User{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("not null constraint violated"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "not null constraint violated",
		},
		{
			name: "Database connection error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("connection refused"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "connection refused",
		},
		{
			name: "Invalid email format",
			user: &model.User{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("invalid email format"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			if tt.setupFn != nil {
				tt.setupFn(mock)
			}

			store := &UserStore{
				db: gormDB,
			}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("User created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06


 */
func TestFollow(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &UserStore{db: gormDB}

	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		mockSQL func()
		wantErr bool
	}{
		{
			name: "Successful Follow",
			userA: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			userB: &model.User{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "Nil User A",
			userA:   nil,
			userB:   &model.User{},
			mockSQL: func() {},
			wantErr: true,
		},
		{
			name:    "Nil User B",
			userA:   &model.User{},
			userB:   nil,
			mockSQL: func() {},
			wantErr: true,
		},
		{
			name: "Self Follow",
			userA: &model.User{
				Model: gorm.Model{ID: 1},
			},
			userB: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database Error",
			userA: &model.User{
				Model: gorm.Model{ID: 1},
			},
			userB: &model.User{
				Model: gorm.Model{ID: 2},
			},
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSQL()

			err := store.Follow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if err != nil {
				t.Logf("Test '%s' failed with error: %v", tt.name, err)
			} else {
				t.Logf("Test '%s' passed successfully", tt.name)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func TestGetByEmail(t *testing.T) {

	createMockDB := func() (*sql.DB, sqlmock.Sqlmock, error) {
		return sqlmock.New()
	}

	tests := []struct {
		name          string
		email         string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:  "Successfully retrieve user by valid email",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \?\) ORDER BY "users"."id" ASC LIMIT 1`).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg"))
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Email:    "test@example.com",
				Username: "testuser",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:  "Non-existent email",
			email: "nonexistent@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \?\) ORDER BY "users"."id" ASC LIMIT 1`).
					WithArgs("nonexistent@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Empty email",
			email: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \?\) ORDER BY "users"."id" ASC LIMIT 1`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database connection error",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \?\) ORDER BY "users"."id" ASC LIMIT 1`).
					WithArgs("test@example.com").
					WillReturnError(errors.New("database connection error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name:  "Special characters in email",
			email: "test+special@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \?\) ORDER BY "users"."id" ASC LIMIT 1`).
					WithArgs("test+special@example.com").
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "testuser", "test+special@example.com", "hashedpass", "test bio", "image.jpg"))
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Email:    "test+special@example.com",
				Username: "testuser",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := createMockDB()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create gorm DB: %v", err)
			}
			defer gormDB.Close()

			tt.mockSetup(mock)

			store := &UserStore{db: gormDB}

			user, err := store.GetByEmail(tt.email)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func TestUserGetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mock    func(sqlmock.Sqlmock)
		want    *model.User
		wantErr error
	}{
		{
			name: "Successfully retrieve user by valid ID",
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg"))
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			wantErr: nil,
		},
		{
			name: "Non-existent user ID",
			id:   999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database connection error",
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: sql.ErrConnDone,
		},
		{
			name: "Zero ID value handling",
			id:   0,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Soft-deleted user retrieval",
			id:   2,
			mock: func(mock sqlmock.Sqlmock) {
				deletedAt := time.Now()
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(2, time.Now(), time.Now(), &deletedAt, "deleteduser", "deleted@example.com", "pass", "bio", "image.jpg"))
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database timeout",
			id:   3,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(3).
					WillReturnError(errors.New("database timeout"))
			},
			want:    nil,
			wantErr: errors.New("database timeout"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.mock(mock)

			store := &UserStore{db: gormDB}
			got, err := store.GetByID(tt.id)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			t.Logf("Test case '%s' completed", tt.name)
			if err != nil {
				t.Logf("Got error: %v", err)
			}
			if got != nil {
				t.Logf("Got user: %+v", got)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
func TestGetByUsername(t *testing.T) {

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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name      string
		username  string
		mockSetup func(sqlmock.Sqlmock)
		want      *model.User
		wantErr   error
	}{
		{
			name:     "Successfully retrieve existing user",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "password", "bio", "image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"testuser", "test@example.com", "hashedpass", "Test bio", "image.jpg",
				)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "Test bio",
				Image:    "image.jpg",
			},
			wantErr: nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty username",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Database connection error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			got, err := store.GetByUsername(tt.username)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Bio, got.Bio)
				assert.Equal(t, tt.want.Image, got.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Error: %v", tt.name, err)
		})
	}
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func TestIsFollowing(t *testing.T) {

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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		mockSetup   func(sqlmock.Sqlmock)
		want        bool
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Valid Users - User A Following User B",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid Users - User A Not Following User B",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
		},
		{
			name:  "Nil User A Parameter",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      false,
			wantErr:   false,
		},
		{
			name: "Nil User B Parameter",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB:     nil,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      false,
			wantErr:   false,
		},
		{
			name: "Database Error",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name:      "Both Users Are Nil",
			userA:     nil,
			userB:     nil,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      false,
			wantErr:   false,
		},
		{
			name: "Same User Reference",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Running test:", tt.name)

			tt.mockSetup(mock)

			got, err := store.IsFollowing(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.expectedErr.Error() {
				t.Errorf("IsFollowing() expected error = %v, got = %v", tt.expectedErr, err)
				return
			}

			if got != tt.want {
				t.Errorf("IsFollowing() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Log("Test completed successfully")
		})
	}
}

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
func TestNewUserStore(t *testing.T) {

	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name: "Scenario 1: Successfully Create New UserStore with Valid DB Connection",
			db: func() *gorm.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock DB: %v", err)
				}
				gormDB, err := gorm.Open("sqlite3", db)
				if err != nil {
					t.Fatalf("Failed to create GORM DB: %v", err)
				}
				return gormDB
			}(),
			wantNil:  false,
			scenario: "Valid DB Connection",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil DB Parameter",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB Connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.scenario)

			userStore := NewUserStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, userStore, "UserStore should be nil")
			} else {
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tt.db, userStore.db, "DB reference should match")
			}

			t.Log("Test completed successfully for scenario:", tt.scenario)
		})
	}

	t.Run("Scenario 3: Verify DB Reference Integrity", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		gormDB, err := gorm.Open("sqlite3", db)
		if err != nil {
			t.Fatalf("Failed to create GORM DB: %v", err)
		}

		userStore := NewUserStore(gormDB)
		assert.Equal(t, gormDB, userStore.db, "DB reference should be maintained")
	})

	t.Run("Scenario 4: Create Multiple UserStore Instances", func(t *testing.T) {
		db1, _, _ := sqlmock.New()
		db2, _, _ := sqlmock.New()
		gormDB1, _ := gorm.Open("sqlite3", db1)
		gormDB2, _ := gorm.Open("sqlite3", db2)

		userStore1 := NewUserStore(gormDB1)
		userStore2 := NewUserStore(gormDB2)

		assert.NotEqual(t, userStore1, userStore2, "UserStore instances should be independent")
	})

	t.Run("Scenario 7: Concurrent UserStore Creation", func(t *testing.T) {
		var wg sync.WaitGroup
		userStores := make([]*UserStore, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				db, _, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlite3", db)
				userStores[index] = NewUserStore(gormDB)
			}(i)
		}

		wg.Wait()

		for i := 0; i < 10; i++ {
			assert.NotNil(t, userStores[i], "Concurrent UserStore creation should succeed")
		}
	})
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUnfollow(t *testing.T) {

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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		setupMock   func()
		userA       *model.User
		userB       *model.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Unfollow",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			userA:       &model.User{Model: gorm.Model{ID: 1}, Username: "userA"},
			userB:       &model.User{Model: gorm.Model{ID: 2}, Username: "userB"},
			expectError: false,
		},
		{
			name: "Unfollow Non-Followed User",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 3).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			userA:       &model.User{Model: gorm.Model{ID: 1}, Username: "userA"},
			userB:       &model.User{Model: gorm.Model{ID: 3}, Username: "userC"},
			expectError: false,
		},
		{
			name:        "Invalid User Reference",
			userA:       &model.User{Model: gorm.Model{ID: 1}, Username: "userA"},
			userB:       nil,
			expectError: true,
			errorMsg:    "invalid user reference",
		},
		{
			name: "Database Connection Error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			userA:       &model.User{Model: gorm.Model{ID: 1}, Username: "userA"},
			userB:       &model.User{Model: gorm.Model{ID: 2}, Username: "userB"},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			err := store.Unfollow(tt.userA, tt.userB)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent Unfollow Operations", func(t *testing.T) {
		userA := &model.User{Model: gorm.Model{ID: 1}, Username: "userA"}
		userB := &model.User{Model: gorm.Model{ID: 2}, Username: "userB"}

		var wg sync.WaitGroup
		concurrentCalls := 5

		for i := 0; i < concurrentCalls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := store.Unfollow(userA, userB)
				assert.NoError(t, err)
			}()
		}
		wg.Wait()
	})

	t.Run("Unfollow After User Deletion", func(t *testing.T) {
		deletedAt := time.Now()
		userA := &model.User{Model: gorm.Model{ID: 1}, Username: "userA"}
		userB := &model.User{
			Model:    gorm.Model{ID: 2, DeletedAt: &deletedAt},
			Username: "userB",
		}

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `follows`").
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := store.Unfollow(userA, userB)
		assert.NoError(t, err)
	})
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestUserUpdate(t *testing.T) {

	type testCase struct {
		name          string
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	baseUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
	}

	tests := []testCase{
		{
			name: "Successful Update",
			user: baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WithArgs(
						sqlmock.AnyArg(),
						baseUser.Username,
						baseUser.Email,
						baseUser.Password,
						baseUser.Bio,
						baseUser.Image,
						baseUser.ID,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Non-Existent User",
			user: &model.User{
				Model: gorm.Model{ID: 999},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			user: baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Unique Constraint Violation",
			user: baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WillReturnError(errors.New("unique constraint violation"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("unique constraint violation"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

			store := &UserStore{db: gormDB}

			err = store.Update(tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestGetFollowingUserIDs(t *testing.T) {

	type testCase struct {
		name          string
		userID        uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError bool
	}

	tests := []testCase{
		{
			name:   "Successfully retrieve following user IDs",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: false,
		},
		{
			name:   "User with no followings",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: false,
		},
		{
			name:   "Database connection error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedIDs:   []uint{},
			expectedError: true,
		},
		{
			name:   "Invalid user ID",
			userID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(999).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &UserStore{db: gormDB}
			testUser := &model.User{
				Model: gorm.Model{
					ID:        tc.userID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			}

			ids, err := store.GetFollowingUserIDs(testUser)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedIDs, ids)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Expected IDs: %v, Got: %v", tc.name, tc.expectedIDs, ids)
		})
	}
}

