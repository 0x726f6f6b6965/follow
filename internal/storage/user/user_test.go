package user

import (
	"testing"
	"time"

	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, cleanup, mock, err := initSQL()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	storage := New(db)
	userId := 3
	data := &models.User{
		Username: "test-user",
		Password: "pwd",
		Salt:     "salt",
	}
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"t_user\" (\"username\",\"password\",\"salt\") VALUES ($1,$2,$3) RETURNING \"create_time\",\"update_time\",\"id\"").
		WithArgs(data.Username, data.Password, data.Salt).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "create_time", "update_time"}).
			AddRow(userId, time.Now(), time.Now()))
	mock.ExpectCommit()
	err = storage.CreateUser(data)
	if err != nil {
		t.Fatal(err)
	}
}

func initSQL() (*gorm.DB, func() error, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, func() error { return nil }, nil, err
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, db.Close, mock, err
	}
	return gormdb, db.Close, mock, nil
}
