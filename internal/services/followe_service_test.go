package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"github.com/0x726f6f6b6965/follow/mocks"
	"github.com/0x726f6f6b6965/follow/pkg/pagination"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestFollowUser(t *testing.T) {
	ser := initFollowService()
	userId := 3
	followingId := 5
	req := &pbFollow.FollowUserRequest{
		Username:  "test-user",
		Following: "test-following",
	}
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       userId,
					Username: "test-user",
				},
				{
					Id:       followingId,
					Username: "test-following",
				},
			}, nil
		}
		mocks.SetFollowingFunc = func(userId, targetId int) error {
			return nil
		}
		_, err := ser.FollowUser(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       followingId,
					Username: "test-following",
				},
			}, nil
		}
		_, err := ser.FollowUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("following user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       userId,
					Username: "test-user",
				},
			}, nil
		}
		_, err := ser.FollowUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("already follow", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       userId,
					Username: "test-user",
				},
				{
					Id:       followingId,
					Username: "test-following",
				},
			}, nil
		}
		mocks.SetFollowingFunc = func(userId, targetId int) error {
			return gorm.ErrDuplicatedKey
		}
		_, err := ser.FollowUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrSetFollow)
	})
}

func TestUnFollowUser(t *testing.T) {
	ser := initFollowService()
	userId := 3
	followingId := 5
	req := &pbFollow.UnFollowUserRequest{
		Username:  "test-user",
		Following: "test-following",
	}
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       userId,
					Username: "test-user",
				},
				{
					Id:       followingId,
					Username: "test-following",
				},
			}, nil
		}
		mocks.UnsetFollowingFunc = func(userId, targetId int) error {
			return nil
		}
		_, err := ser.UnFollowUser(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       followingId,
					Username: "test-following",
				},
			}, nil
		}
		_, err := ser.UnFollowUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("following user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       userId,
					Username: "test-user",
				},
			}, nil
		}
		_, err := ser.UnFollowUser(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestGetFollowers(t *testing.T) {
	ser := initFollowService()
	id := 3
	req := &pbFollow.GetCommonRequest{
		Username: "test-user",
		Size:     20,
	}
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		startFollowerId := 8
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       id,
					Username: "test-user",
				},
			}, nil
		}
		expect := []string{}
		mocks.GetUserWithFollowersFunc = func(userId int, lastId int, limit int) ([]models.User, error) {
			if lastId != 0 {
				t.Error("lastId should be 0")
			}
			if userId != id {
				t.Error("userId not match")
			}
			if limit != int(req.Size) {
				t.Error("limit not match")
			}
			result := []models.User{}
			for i := 0; i < int(req.Size); i++ {
				result = append(result, models.User{
					Id:       startFollowerId + i,
					Username: fmt.Sprintf("follower-%d", i),
				})
				expect = append(expect, fmt.Sprintf("follower-%d", i))
			}
			return result, nil
		}
		resp, err := ser.GetFollowers(ctx, req)
		assert.Nil(t, err)
		for _, follower := range resp.Usernames {
			assert.Contains(t, expect, follower)
		}
		next := &pagination.PageToken{}
		err = next.DecodePageTokenStruct(resp.NextPageToken)
		assert.Nil(t, err)
		assert.Equal(t, startFollowerId+int(req.Size)-1, next.LastId)
		assert.Equal(t, int(req.Size), next.Size)

		// with token
		req.PageToken = next.String()
		req.Size = 30
		expect = []string{}
		mocks.GetUserWithFollowersFunc = func(userId int, lastId int, limit int) ([]models.User, error) {
			if lastId != next.LastId {
				t.Error("lastId should be 0")
			}
			if userId != id {
				t.Error("userId not match")
			}
			if limit != int(req.Size) {
				t.Error("limit not match")
			}
			result := []models.User{}
			for i := 0; i < int(req.Size); i++ {
				result = append(result, models.User{
					Id:       next.LastId + i + 1,
					Username: fmt.Sprintf("follower-%d", next.LastId+i+1),
				})
				expect = append(expect, fmt.Sprintf("follower-%d", next.LastId+i+1))
			}
			return result, nil
		}
		resp, err = ser.GetFollowers(ctx, req)
		assert.Nil(t, err)
		for _, follower := range resp.Usernames {
			assert.Contains(t, expect, follower)
		}
		nextAndNext := &pagination.PageToken{}
		err = nextAndNext.DecodePageTokenStruct(resp.NextPageToken)
		assert.Nil(t, err)
		assert.Equal(t, next.LastId+int(req.Size), nextAndNext.LastId)
		assert.Equal(t, int(req.Size), nextAndNext.Size)
	})

	t.Run("user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{}, nil
		}
		_, err := ser.GetFollowers(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestGetFollowing(t *testing.T) {
	ser := initFollowService()
	id := 3
	req := &pbFollow.GetCommonRequest{
		Username: "test-user",
		Size:     20,
	}
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		startFollowerId := 8
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       id,
					Username: "test-user",
				},
			}, nil
		}
		expect := []string{}
		mocks.GetUserWithFollowingFunc = func(userId int, lastId int, limit int) ([]models.User, error) {
			if lastId != 0 {
				t.Error("lastId should be 0")
			}
			if userId != id {
				t.Error("userId not match")
			}
			if limit != int(req.Size) {
				t.Error("limit not match")
			}
			result := []models.User{}
			for i := 0; i < int(req.Size); i++ {
				result = append(result, models.User{
					Id:       startFollowerId + i,
					Username: fmt.Sprintf("follower-%d", i),
				})
				expect = append(expect, fmt.Sprintf("follower-%d", i))
			}
			return result, nil
		}
		resp, err := ser.GetFollowing(ctx, req)
		assert.Nil(t, err)
		for _, follower := range resp.Usernames {
			assert.Contains(t, expect, follower)
		}
		next := &pagination.PageToken{}
		err = next.DecodePageTokenStruct(resp.NextPageToken)
		assert.Nil(t, err)
		assert.Equal(t, startFollowerId+int(req.Size)-1, next.LastId)
		assert.Equal(t, int(req.Size), next.Size)

		// with token
		req.PageToken = next.String()
		req.Size = 30
		expect = []string{}
		mocks.GetUserWithFollowingFunc = func(userId int, lastId int, limit int) ([]models.User, error) {
			if lastId != next.LastId {
				t.Error("lastId should be 0")
			}
			if userId != id {
				t.Error("userId not match")
			}
			if limit != int(req.Size) {
				t.Error("limit not match")
			}
			result := []models.User{}
			for i := 0; i < int(req.Size); i++ {
				result = append(result, models.User{
					Id:       next.LastId + i + 1,
					Username: fmt.Sprintf("follower-%d", next.LastId+i+1),
				})
				expect = append(expect, fmt.Sprintf("follower-%d", next.LastId+i+1))
			}
			return result, nil
		}
		resp, err = ser.GetFollowing(ctx, req)
		assert.Nil(t, err)
		for _, follower := range resp.Usernames {
			assert.Contains(t, expect, follower)
		}
		nextAndNext := &pagination.PageToken{}
		err = nextAndNext.DecodePageTokenStruct(resp.NextPageToken)
		assert.Nil(t, err)
		assert.Equal(t, next.LastId+int(req.Size), nextAndNext.LastId)
		assert.Equal(t, int(req.Size), nextAndNext.Size)
	})

	t.Run("user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{}, nil
		}
		_, err := ser.GetFollowing(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestGetFriends(t *testing.T) {
	ser := initFollowService()
	id := 3
	req := &pbFollow.GetCommonRequest{
		Username: "test-user",
		Size:     20,
	}
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		startFollowerId := 8
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{
				{
					Id:       id,
					Username: "test-user",
				},
			}, nil
		}
		expect := map[int]bool{}
		mocks.GetUserWithFriendsFunc = func(userId int, lastId int, limit int) ([]models.Follower, error) {
			if lastId != 0 {
				t.Error("lastId should be 0")
			}
			if userId != id {
				t.Error("userId not match")
			}
			if limit != int(req.Size) {
				t.Error("limit not match")
			}
			result := []models.Follower{}
			for i := 0; i < int(req.Size); i++ {
				result = append(result, models.Follower{
					FollowingId: startFollowerId + i,
				})
				expect[startFollowerId+i] = true
			}
			return result, nil
		}
		mocks.GetUserInfoByIdFunc = func(ids ...int) ([]models.User, error) {
			result := []models.User{}
			for _, val := range ids {
				if !expect[val] {
					t.Errorf("id %d not in expect", val)
				}
				result = append(result, models.User{
					Id:       val,
					Username: fmt.Sprintf("follower-%d", val),
				})
			}
			return result, nil
		}
		resp, err := ser.GetFriends(ctx, req)
		assert.Nil(t, err)
		assert.Equal(t, len(expect), len(resp.Usernames))
		next := &pagination.PageToken{}
		err = next.DecodePageTokenStruct(resp.NextPageToken)
		assert.Nil(t, err)
		assert.Equal(t, startFollowerId+int(req.Size)-1, next.LastId)
		assert.Equal(t, int(req.Size), next.Size)

		// with token
		req.PageToken = next.String()
		req.Size = 30
		expect = map[int]bool{}
		mocks.GetUserWithFriendsFunc = func(userId int, lastId int, limit int) ([]models.Follower, error) {
			if lastId != next.LastId {
				t.Error("lastId should be 0")
			}
			if userId != id {
				t.Error("userId not match")
			}
			if limit != int(req.Size) {
				t.Error("limit not match")
			}

			result := []models.Follower{}
			for i := 0; i < int(req.Size); i++ {
				result = append(result, models.Follower{
					FollowingId: next.LastId + i + 1,
				})
				expect[next.LastId+i+1] = true
			}
			return result, nil
		}
		resp, err = ser.GetFriends(ctx, req)
		assert.Nil(t, err)
		assert.Equal(t, len(expect), len(resp.Usernames))
		nextAndNext := &pagination.PageToken{}
		err = nextAndNext.DecodePageTokenStruct(resp.NextPageToken)
		assert.Nil(t, err)
		assert.Equal(t, next.LastId+int(req.Size), nextAndNext.LastId)
		assert.Equal(t, int(req.Size), nextAndNext.Size)
	})

	t.Run("user not exist", func(t *testing.T) {
		mocks.GetUserInfoFunc = func(username ...string) ([]models.User, error) {
			return []models.User{}, nil
		}
		_, err := ser.GetFriends(ctx, req)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func initFollowService() pbFollow.FollowServiceServer {
	mockUsers := &mocks.MockSotrageUsers{}
	mockFollowers := &mocks.MockSotrageFollowers{}
	logger, _ := zap.NewDevelopment()
	followService := NewFollowService(mockUsers, mockFollowers, logger)
	return followService
}
