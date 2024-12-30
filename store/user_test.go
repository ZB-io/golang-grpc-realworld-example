package user_test // or whatever package name you're using for your tests

import (
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    // other imports...
)

func TestGetByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn: db,
    }), &gorm.Config{})
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
    }

    s := &UserStore{db: gormDB}

    // Set up expectations
    mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`id` = ?").
        WithArgs(1).
        WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
            AddRow(1, "testuser", "test@example.com"))

    // Run the test
    user, err := s.GetByID(1)

    // Assert results
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, uint(1), user.ID)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "test@example.com", user.Email)

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
