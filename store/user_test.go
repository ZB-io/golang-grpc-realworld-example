package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/mock"
)

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9
*/
func TestNewUserStore(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920
*/
func TestCreate(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1
*/
func TestGetByID(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1
*/
func TestGetByEmail(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24
*/
func TestGetByUsername(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435
*/
func TestUpdate(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06
*/
func TestFollow(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c
*/
func TestIsFollowing(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55
*/
func TestUnfollow(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

func TestUnfollowDatabaseError(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

func TestUnfollowConcurrent(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7
*/
func TestGetFollowingUserIDs(t *testing.T) {
	// ... (rest of the function remains unchanged)
}

// MockGormDB is a mock of gorm.DB
type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

// MockAssociation is a mock of gorm.Association
type MockAssociation struct {
	mock.Mock
}

func (m *MockAssociation) Append(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

func (m *MockAssociation) Error() error {
	args := m.Called()
	return args.Error(0)
}

// MockUserStore is a mock of UserStore
type MockUserStore struct {
	db *MockGormDB
}
