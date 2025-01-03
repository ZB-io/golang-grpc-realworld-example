package store

import (
		"testing"
		"time"
		"github.com/DATA-DOG/go-sqlmock"
		"github.com/jinzhu/gorm"
		"github.com/stretchr/testify/assert"
		"database/sql"
		"regexp"
		sqlmock "github.com/DATA-DOG/go-sqlmock"
		"github.com/raahii/golang-grpc-realworld-example/model"
		_ "github.com/jinzhu/gorm/dialects/mysql"
		"errors"
)


/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
func TestNewUserStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		setupDB  func() (*gorm.DB, sqlmock.Sqlmock, error)
		validate func(*testing.T, *UserStore, *gorm.DB)
	}

	createMockDB := func() (*gorm.DB, sqlmock.Sqlmock, error) {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}
		gormDB, err := gorm.Open("mysql", sqlDB)
		if err != nil {
			return nil, nil, err
		}
		return gormDB, mock, nil
	}

	tests := []testCase{
		{
			name: "Scenario 1: Successfully Create New UserStore with Valid DB Connection",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return createMockDB()
			},
			validate: func(t *testing.T, us *UserStore, db *gorm.DB) {
				t.Log("Validating successful UserStore creation")
				assert.NotNil(t, us, "UserStore should not be nil")
				assert.Equal(t, db, us.db, "DB reference should match")
			},
		},
		{
			name: "Scenario 2: Create UserStore with Nil DB Connection",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return nil, nil, nil
			},
			validate: func(t *testing.T, us *UserStore, _ *gorm.DB) {
				t.Log("Validating UserStore creation with nil DB")
				assert.NotNil(t, us, "UserStore should not be nil even with nil DB")
				assert.Nil(t, us.db, "DB reference should be nil")
			},
		},
		{
			name: "Scenario 3: Verify DB Reference Integrity",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := createMockDB()
				if err != nil {
					return nil, nil, err
				}

				db = db.LogMode(true)
				return db, mock, nil
			},
			validate: func(t *testing.T, us *UserStore, db *gorm.DB) {
				t.Log("Validating DB reference integrity")
				assert.Equal(t, db, us.db, "DB reference should maintain integrity")
				isLogMode := db.GetLogger() != nil
				assert.True(t, isLogMode, "DB configuration should be preserved")
			},
		},
		{
			name: "Scenario 4: Multiple UserStore Instances Independence",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return createMockDB()
			},
			validate: func(t *testing.T, us *UserStore, db *gorm.DB) {
				t.Log("Validating multiple UserStore instances")

				db2, _, _ := createMockDB()
				us2 := NewUserStore(db2)

				assert.NotEqual(t, us.db, us2.db, "Different UserStore instances should have independent DB references")
				assert.Equal(t, db, us.db, "Original UserStore should maintain its DB reference")
			},
		},
		{
			name: "Scenario 5: UserStore Creation with Configured DB Connection",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := createMockDB()
				if err != nil {
					return nil, nil, err
				}

				db = db.LogMode(true)
				db.SetNowFuncOverride(func() time.Time {
					return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				})
				return db, mock, nil
			},
			validate: func(t *testing.T, us *UserStore, db *gorm.DB) {
				t.Log("Validating UserStore with configured DB")
				assert.Equal(t, db, us.db, "DB reference should maintain configurations")
				isLogMode := db.GetLogger() != nil
				assert.True(t, isLogMode, "LogMode configuration should be preserved")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, _, err := tc.setupDB()
			if err != nil {
				t.Fatalf("Failed to setup test DB: %v", err)
			}
			if db != nil {
				defer db.Close()
			}

			userStore := NewUserStore(db)

			tc.validate(t, userStore, db)
		})
	}
}


/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestUserStoreCreate(t *testing.T) {
	type testCase struct {
		name          string
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
		errorMessage  string
	}

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

	userStore := &UserStore{db: gormDB}

	tests := []testCase{
		{
			name: "Successful user creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "test@example.com", "password123", "Test bio", "test-image.jpg").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Duplicate username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&mysql.MySQLError{Number: 1062})
				mock.ExpectRollback()
			},
			expectedError: true,
			errorMessage:  "Duplicate entry",
		},
		{
			name: "Missing required fields",
			user: &model.User{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&mysql.MySQLError{Number: 1048})
				mock.ExpectRollback()
			},
			expectedError: true,
			errorMessage:  "Column cannot be null",
		},
		{
			name: "Database connection error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: true,
			errorMessage:  "sql: connection is already closed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(mock)

			err := userStore.Create(tc.user)

			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func TestUserStoreGetByID(t *testing.T) {
	type testCase struct {
		name          string
		userID        uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	validUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Bio:      "Test bio",
		Image:    "https://example.com/image.jpg",
	}

	tests := []testCase{
		{
			name:   "Successfully Retrieve User by Valid ID",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(validUser.ID, validUser.CreatedAt, validUser.UpdatedAt, nil, validUser.Username, validUser.Email, validUser.Password, validUser.Bio, validUser.Image)
				mock.ExpectQuery("^SELECT \\* FROM `users`").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser:  validUser,
			expectedError: nil,
		},
		{
			name:   "Attempt to Retrieve Non-existent User ID",
			userID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:   "Handle Database Connection Error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "Retrieve User with Zero ID",
			userID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users`").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
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

			gormDB.LogMode(false)

			tc.mockSetup(mock)

			userStore := &UserStore{db: gormDB}

			user, err := userStore.GetByID(tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				assert.Equal(t, tc.expectedUser.Bio, user.Bio)
				assert.Equal(t, tc.expectedUser.Image, user.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func TestUserStoreGetByEmail(t *testing.T) {
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

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name          string
		email         string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:  "Successful retrieval",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, time.Now(), time.Now(), nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:  "Special characters in email",
			email: "test+special@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, time.Now(), time.Now(), nil, "testuser", "test+special@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("test+special@example.com").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test+special@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			user, err := userStore.GetByEmail(tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
func TestUserStoreGetByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
		validateExtra func(*testing.T, *model.User, error)
	}{
		{
			name:     "Successfully retrieve existing user",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, "2024-01-01 00:00:00", "2024-01-01 00:00:00", nil, "testuser", "test@example.com", "hashedpassword", "Test bio", "image.jpg")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpassword",
				Bio:      "Test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty username",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Database connection error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:     "Special characters in username",
			username: "test@user#123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, "2024-01-01 00:00:00", "2024-01-01 00:00:00", nil, "test@user#123", "special@example.com", "hashedpassword", "Special bio", "image.jpg")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("test@user#123").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "test@user#123",
				Email:    "special@example.com",
				Password: "hashedpassword",
				Bio:      "Special bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
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

			tt.mockSetup(mock)

			userStore := &UserStore{db: gormDB}
			user, err := userStore.GetByUsername(tt.username)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			if tt.validateExtra != nil {
				tt.validateExtra(t, user, err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestUserStoreUpdate(t *testing.T) {
	type testCase struct {
		name    string
		user    *model.User
		mockDB  func(mock sqlmock.Sqlmock, user *model.User)
		wantErr bool
		errMsg  string
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
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WithArgs(
						sqlmock.AnyArg(),
						user.Username,
						user.Email,
						user.Password,
						user.Bio,
						user.Image,
						user.ID,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Database Error",
			user: baseUser,
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:    "Nil User Model",
			user:    nil,
			mockDB:  func(mock sqlmock.Sqlmock, user *model.User) {},
			wantErr: true,
			errMsg:  "invalid user model",
		},
		{
			name: "Empty Required Fields",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "",
				Email:    "",
				Password: "",
			},
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WillReturnError(errors.New("constraint violation"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "constraint violation",
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
			gormDB.LogMode(true)
			defer gormDB.Close()

			tc.mockDB(mock, tc.user)

			store := &UserStore{db: gormDB}
			err = store.Update(tc.user)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(t, err)
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
func TestUserStoreFollow(t *testing.T) {

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
		userA     *model.User
		userB     *model.User
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successful Follow",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Self Follow Attempt",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
			errMsg:  "cannot follow self",
		},
		{
			name:  "Nil User Parameters",
			userA: nil,
			userB: nil,
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Database Error",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			err := store.Follow(tt.userA, tt.userB)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func TestUserStoreIsFollowing(t *testing.T) {

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

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expected    bool
		expectError bool
	}{
		{
			name: "Valid Following Relationship",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "Valid Non-Following Relationship",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectError: false,
		},
		{
			name:  "Nil User A",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
			},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectError: false,
		},
		{
			name: "Nil User B",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			userB:       nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectError: false,
		},
		{
			name: "Database Error",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			expected:    false,
			expectError: true,
		},
		{
			name:        "Both Users Nil",
			userA:       nil,
			userB:       nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectError: false,
		},
		{
			name: "Same User Check",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectError: false,
		},
		{
			name: "Database Connection Lost",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("connection lost"))
			},
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			result, err := userStore.IsFollowing(tt.userA, tt.userB)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != tt.expected {
				t.Errorf("Expected result %v but got %v", tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUserStoreUnfollow(t *testing.T) {

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
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Unfollow",
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
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:  "Null User Parameters",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "primary key can't be nil",
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
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name: "Self-Unfollow Attempt",
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
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

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
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name      string
		user      *model.User
		mockSetup func(sqlmock.Sqlmock)
		expected  []uint
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successfully retrieve following user IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: []uint{2, 3, 4},
			wantErr:  false,
		},
		{
			name: "User with no followers",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: []uint{},
			wantErr:  false,
		},
		{
			name: "Database connection error",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expected: nil,
			wantErr:  true,
			errMsg:   "database connection error",
		},
		{
			name: "Invalid user ID",
			user: &model.User{Model: gorm.Model{ID: 0}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(0).
					WillReturnError(errors.New("invalid user ID"))
			},
			expected: nil,
			wantErr:  true,
			errMsg:   "invalid user ID",
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

			tt.mockSetup(mock)

			store := &UserStore{db: gormDB}
			got, err := store.GetFollowingUserIDs(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

