package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"github.com/0x726f6f6b6965/follow/mocks"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	"github.com/stretchr/testify/assert"
	boom "github.com/tylertreat/BoomFilters"
	"go.uber.org/zap"
)

func TestCreateUser(t *testing.T) {
	mUser := &mocks.MockSotrageUsers{}
	mCache := &mocks.MockSotrageCache{}
	ttl := time.Minute
	filter := boom.NewDefaultCountingBloomFilter(100, 0.1)
	logger, _ := zap.NewDevelopment()
	ser := NewUserService(mUser, mCache, ttl, filter, logger)
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		req := &pbUser.CreateUserRequest{
			Username: "abc",
			Password: "123456",
		}
		testSalt = helper.CreateNewSalt()
		pwd, _ := testSalt.SaltInput(req.Password)
		mocks.GetUserInfoFunc = func(usernames ...string) ([]models.User, error) {
			return []models.User{}, nil
		}
		mocks.CreateUserFunc = func(user *models.User) error {
			if user.Username != req.Username || user.Password != pwd {
				return errors.New("user info not match")
			}
			return nil
		}
		// cache
		mocks.GetFunc = func(ctx context.Context, key string) (string, error) {
			if key != fmt.Sprintf(UserExistKey, req.Username) {
				t.Error("cache get key not match")
			}
			return "", nil
		}

		mocks.DeleteFunc = func(ctx context.Context, key string) error {
			if key != fmt.Sprintf(UserNotExistKey, req.Username) {
				t.Error("cache delete key not match")
			}
			return nil
		}

		mocks.SetFunc = func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
			if key != fmt.Sprintf(UserExistKey, req.Username) {
				t.Error("cache set key not match")
			}
			return nil
		}
		_, err := ser.CreateUser(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("username exist", func(t *testing.T) {
		req := &pbUser.CreateUserRequest{
			Username: "abc",
			Password: "123456",
		}
		mocks.GetUserInfoFunc = func(usernames ...string) ([]models.User, error) {
			return []models.User{
				{Id: 3, Username: "abc"},
			}, nil
		}
		// cache
		mocks.GetFunc = func(ctx context.Context, key string) (string, error) {
			if key != fmt.Sprintf(UserExistKey, req.Username) {
				t.Error("cache get key not match")
			}
			return "", nil
		}
		_, err := ser.CreateUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserExists)

		// get user from cache
		mocks.GetFunc = func(ctx context.Context, key string) (string, error) {
			if key != fmt.Sprintf(UserExistKey, req.Username) {
				t.Error("cache get key not match")
			}
			return "3", nil
		}
		_, err = ser.CreateUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserExists)
	})

	t.Run("invalid input", func(t *testing.T) {
		req := &pbUser.CreateUserRequest{
			Username: "",
			Password: "123456",
		}
		_, err := ser.CreateUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrInvalidInput)

		req = &pbUser.CreateUserRequest{
			Username: "123",
			Password: "",
		}
		_, err = ser.CreateUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})
}
