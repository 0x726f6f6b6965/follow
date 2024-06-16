package services

import (
	"context"
	"errors"
	"testing"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"github.com/0x726f6f6b6965/follow/mocks"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateUser(t *testing.T) {
	m := &mocks.MockSotrageUsers{}
	logger, _ := zap.NewDevelopment()
	ser := NewUserService(m, logger)
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
		_, err := ser.CreateUser(ctx, req)
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
