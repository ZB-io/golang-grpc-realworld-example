package store

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9
*/
func TestNewUserStore(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	store := NewUserStore(gormDB)
	assert.NotNil(t, store)
}

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920
*/
func TestUserStoreCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	user := &model.User{
		Email:    "test@example.com",
		Username: "testuser",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(user.Email, user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewUserStore(gormDB)
	err = store.Create(user)
	assert.NoError(t, err)
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1
*/
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	userID := uint(1)
	rows := sqlmock.NewRows([]string{"id", "email", "username"}).
		AddRow(userID, "test@example.com", "testuser")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(userID).
		WillReturnRows(rows)

	store := NewUserStore(gormDB)
	user, err := store.GetByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1
*/
func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"id", "email", "username"}).
		AddRow(1, email, "testuser")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(email).
		WillReturnRows(rows)

	store := NewUserStore(gormDB)
	user, err := store.GetByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24
*/
func TestUserStoreGetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	username := "testuser"
	rows := sqlmock.NewRows([]string{"id", "email", "username"}).
		AddRow(1, "test@example.com", username)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(username).
		WillReturnRows(rows)

	store := NewUserStore(gormDB)
	user, err := store.GetByUsername(username)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435
*/
func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	user := &model.User{
		ID:       1,
		Email:    "updated@example.com",
		Username: "updateduser",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewUserStore(gormDB)
	err = store.Update(user)
	assert.NoError(t, err)
}

/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06
*/
func TestFollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	followerID := uint(1)
	followingID := uint(2)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `follows`")).
		WithArgs(followerID, followingID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewUserStore(gormDB)
	err = store.Follow(followerID, followingID)
	assert.NoError(t, err)
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c
*/
func TestUserStoreIsFollowing(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	followerID := uint(1)
	followingID := uint(2)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `follows`")).
		WithArgs(followerID, followingID).
		WillReturnRows(rows)

	store := NewUserStore(gormDB)
	following, err := store.IsFollowing(followerID, followingID)
	assert.NoError(t, err)
	assert.True(t, following)
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55
*/
func TestUserStoreUnfollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	followerID := uint(1)
	followingID := uint(2)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `follows`")).
		WithArgs(followerID, followingID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewUserStore(gormDB)
	err = store.Unfollow(followerID, followingID)
	assert.NoError(t, err)
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7
*/
func TestGetFollowingUserIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	assert.NoError(t, err)

	userID := uint(1)
	rows := sqlmock.NewRows([]string{"following_id"}).
		AddRow(2).
		AddRow(3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT following_id FROM `follows`")).
		WithArgs(userID).
		WillReturnRows(rows)

	store := NewUserStore(gormDB)
	ids, err := store.GetFollowingUserIDs(userID)
	assert.NoError(t, err)
	assert.Len(t, ids, 2)
}
