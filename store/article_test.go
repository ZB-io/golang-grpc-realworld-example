package article_test // or whatever package name is appropriate for your test file

import (
    "errors"
    "reflect"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/jinzhu/gorm"
    "github.com/raahii/golang-grpc-realworld-example/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type MockDB struct {
    mock.Mock
}

func (m *MockDB) Begin() *gorm.DB {
    args := m.Called()
    return args.Get(0).(*gorm.DB)
}

// Implement other necessary methods...

type MockAssociation struct {
    mock.Mock
}

func (m *MockAssociation) Append(values ...interface{}) error {
    args := m.Called(values)
    return args.Error(0)
}

// Implement other necessary methods...
