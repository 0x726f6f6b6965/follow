package user

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetUserWithFollowers(t *testing.T) {
	db, cleanup, mock, err := initSQL()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	storage := New(db)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectQuery("SELECT u.id, u.username FROM t_followers f INNER JOIN t_user u ON f.follower_id = u.id WHERE f.following_id = $1 and f.follower_id > $2 ORDER BY f.follower_id asc LIMIT $3").
		WithArgs(1, 0, 5).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "username"}).
			AddRow(2, "user-2").AddRow(3, "user-3").
			AddRow(4, "user-4").AddRow(5, "user-5").
			AddRow(6, "user-6"))
	result, err := storage.GetUserWithFollowers(1, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	expect := []int{2, 3, 4, 5, 6}
	for _, follower := range result {
		assert.Contains(t, expect, follower.Id)
	}
}

func TestGetUserWithFollowing(t *testing.T) {
	db, cleanup, mock, err := initSQL()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	storage := New(db)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectQuery("SELECT u.id, u.username FROM t_followers f INNER JOIN t_user u ON f.following_id = u.id WHERE f.follower_id = $1 and f.following_id > $2 ORDER BY f.following_id asc LIMIT $3").
		WithArgs(1, 0, 5).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "username"}).
			AddRow(2, "user-2").AddRow(3, "user-3").
			AddRow(4, "user-4").AddRow(5, "user-5").
			AddRow(6, "user-6"))
	result, err := storage.GetUserWithFollowing(1, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	expect := []int{2, 3, 4, 5, 6}
	for _, follower := range result {
		assert.Contains(t, expect, follower.Id)
	}
}

func TestGetUserWithFriends(t *testing.T) {
	db, cleanup, mock, err := initSQL()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	storage := New(db)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectQuery("SELECT f.following_id FROM t_followers f INNER JOIN t_followers f2 ON f.follower_id = f2.following_id and f.following_id = f2.follower_id WHERE f.follower_id < f.following_id and f.follower_id = $1 and f.following_id > $2 ORDER BY f.following_id asc LIMIT $3").
		WithArgs(1, 0, 5).
		WillReturnRows(sqlmock.NewRows(
			[]string{"following_id"}).
			AddRow(2).AddRow(3).
			AddRow(4).AddRow(5).
			AddRow(6))
	result, err := storage.GetUserWithFriends(1, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	expect := []int{2, 3, 4, 5, 6}
	for _, follower := range result {
		assert.Contains(t, expect, follower.FollowingId)
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
