package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	"github.com/0x726f6f6b6965/follow/internal/storage/cache"
	"github.com/0x726f6f6b6965/follow/internal/storage/follower"
	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/pagination"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	boom "github.com/tylertreat/BoomFilters"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type followService struct {
	pbFollow.UnimplementedFollowServiceServer
	logger          *zap.Logger
	userStorage     user.SotrageUsers
	followerStorage follower.SotrageFollowers
	cacheTTL        time.Duration
	cacheStorage    cache.Cache
	filter          *boom.CountingBloomFilter
}

func NewFollowService(storageUser user.SotrageUsers, storageFollower follower.SotrageFollowers, cacheStorage cache.Cache, cacheTTL time.Duration, filter *boom.CountingBloomFilter, logger *zap.Logger) pbFollow.FollowServiceServer {
	return &followService{
		userStorage:     storageUser,
		followerStorage: storageFollower,
		cacheStorage:    cacheStorage,
		cacheTTL:        cacheTTL,
		filter:          filter,
		logger:          logger,
	}

}

// FollowUser follow a user.
func (f *followService) FollowUser(ctx context.Context, req *pbFollow.FollowUserRequest) (*emptypb.Empty, error) {
	usersInfo, err := f.getUserInfo(ctx, req.Username, req.Following)
	if err != nil {
		return nil, err
	}
	var (
		userId   int
		targetId int
	)
	for _, info := range usersInfo {
		if info.Username == req.Username {
			userId = info.Id
		} else if info.Username == req.Following {
			targetId = info.Id
		}
	}
	if userId == 0 || targetId == 0 {
		return nil, errors.Join(ErrUserNotFound,
			fmt.Errorf("id not found, user: %d, following: %d", userId, targetId))
	}

	err = f.followerStorage.SetFollowing(userId, targetId)
	if err != nil {
		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			f.logger.Error("set user followeing failed", zap.Any("request", req), zap.Error(err))
		}
		return nil, errors.Join(ErrSetFollow, err)
	}
	return &emptypb.Empty{}, nil
}

// UnFollowUser unfollow a user.
func (f *followService) UnFollowUser(ctx context.Context, req *pbFollow.UnFollowUserRequest) (*emptypb.Empty, error) {
	usersInfo, err := f.getUserInfo(ctx, req.Username, req.Following)
	if err != nil {
		return nil, err
	}

	var (
		userId   int
		targetId int
	)
	for _, info := range usersInfo {
		if info.Username == req.Username {
			userId = info.Id
		} else if info.Username == req.Following {
			targetId = info.Id
		}
	}
	if userId == 0 || targetId == 0 {
		return nil, errors.Join(ErrUserNotFound,
			fmt.Errorf("id not found, user: %d, following: %d", userId, targetId))
	}
	err = f.followerStorage.UnsetFollowing(userId, targetId)
	if err != nil {
		f.logger.Error("set user unfolloweing failed", zap.Any("request", req), zap.Error(err))
		return nil, errors.Join(ErrSetFollow, err)
	}
	return &emptypb.Empty{}, nil
}

// GetFollowers get followers list.
func (f *followService) GetFollowers(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	usersInfo, err := f.getUserInfo(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	userId := usersInfo[0].Id
	token := &pagination.PageToken{
		LastId: 0,
		Size:   DEFAUTL_SIZE,
	}
	if !helper.IsEmpty(req.PageToken) {
		_ = token.DecodePageTokenStruct(req.PageToken)
	}
	if req.Size > 0 {
		token.Size = int(req.Size)
	}
	lastId := token.LastId
	followersInfo, err := f.followerStorage.GetUserWithFollowers(userId, token.LastId, token.Size)
	if err != nil {
		f.logger.Error("get follower info failed", zap.Any("request", req), zap.Error(err))
		return nil, err
	}
	followers := []string{}
	for _, info := range followersInfo {
		followers = append(followers, info.Username)
		if info.Id > token.LastId {
			token.LastId = info.Id
		}
	}
	nextToken := token.String()
	if lastId == token.LastId || len(followers) < token.Size {
		nextToken = ""
	}

	resp := &pbFollow.GetCommonResponse{
		NextPageToken: nextToken,
		Usernames:     followers,
	}
	return resp, nil
}

// GetFollowing get following list.
func (f *followService) GetFollowing(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	usersInfo, err := f.getUserInfo(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	userId := usersInfo[0].Id
	token := &pagination.PageToken{
		LastId: 0,
		Size:   DEFAUTL_SIZE,
	}
	if !helper.IsEmpty(req.PageToken) {
		_ = token.DecodePageTokenStruct(req.PageToken)
	}
	if req.Size > 0 {
		token.Size = int(req.Size)
	}
	lastId := token.LastId
	followingInfo, err := f.followerStorage.GetUserWithFollowing(userId, token.LastId, token.Size)
	if err != nil {
		f.logger.Error("get following info failed", zap.Any("request", req), zap.Error(err))
		return nil, err
	}
	following := []string{}
	for _, info := range followingInfo {
		following = append(following, info.Username)
		if info.Id > token.LastId {
			token.LastId = info.Id
		}
	}
	nextToken := token.String()
	if lastId == token.LastId || len(following) < token.Size {
		nextToken = ""
	}

	resp := &pbFollow.GetCommonResponse{
		NextPageToken: nextToken,
		Usernames:     following,
	}
	return resp, nil
}

// GetFriends get friends list.
func (f *followService) GetFriends(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	usersInfo, err := f.getUserInfo(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	userId := usersInfo[0].Id
	token := &pagination.PageToken{
		LastId: 0,
		Size:   DEFAUTL_SIZE,
	}
	if !helper.IsEmpty(req.PageToken) {
		_ = token.DecodePageTokenStruct(req.PageToken)
	}
	if req.Size > 0 {
		token.Size = int(req.Size)
	}
	lastId := token.LastId
	friendInfos, err := f.followerStorage.GetUserWithFriends(userId, token.LastId, token.Size)
	if err != nil {
		f.logger.Error("get friends info failed", zap.Any("request", req), zap.Error(err))
		return nil, err
	}
	friendIds := []int{}
	for _, info := range friendInfos {
		friendIds = append(friendIds, info.FollowingId)
		if info.FollowingId > token.LastId {
			token.LastId = info.FollowingId
		}
	}
	nextToken := token.String()
	if lastId == token.LastId || len(friendInfos) < token.Size {
		nextToken = ""
	}
	resp := &pbFollow.GetCommonResponse{
		NextPageToken: nextToken,
	}
	friends, err := f.userStorage.GetUserInfoById(friendIds...)
	if err != nil {
		f.logger.Error("get user info failed", zap.Any("request", req), zap.Error(err))
		return nil, err
	}
	for _, friend := range friends {
		resp.Usernames = append(resp.Usernames, friend.Username)
	}

	return resp, nil
}

func (f *followService) getUserInfo(ctx context.Context, usernames ...string) ([]models.User, error) {
	noCache := []string{}
	result := []models.User{}
	for _, name := range usernames {
		// for not exist user
		if f.filter.Test([]byte(name)) {
			return nil, ErrUserNotFound
		}
		cacheInfo, _ := f.cacheStorage.Get(ctx, fmt.Sprintf(UserExistKey, name))
		if helper.IsEmpty(cacheInfo) {
			noCache = append(noCache, name)
		} else {
			id, parseErr := strconv.Atoi(cacheInfo)
			if parseErr != nil {
				f.logger.Error("parse cache failed", zap.String("cache", cacheInfo), zap.Error(parseErr))
				noCache = append(noCache, name)
				continue
			}
			result = append(result, models.User{
				Username: name,
				Id:       id,
			})
		}
	}
	// get from db
	if len(noCache) > 0 {
		userInfos, err := f.userStorage.GetUserInfo(noCache...)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				f.logger.Error("get user info failed", zap.Strings("username", usernames), zap.Error(err))
			}
			return nil, errors.Join(ErrUserNotFound, err)
		}
		result = append(result, userInfos...)
		for _, info := range userInfos {
			_ = f.cacheStorage.Set(ctx, fmt.Sprintf(UserExistKey, info.Username), info.Id, f.cacheTTL)
		}
	}
	if len(usernames) != len(result) {
		set := map[string]bool{}
		for _, info := range result {
			set[info.Username] = true
		}
		for _, name := range usernames {
			if !set[name] {
				f.filter.Add([]byte(name))
			}
		}
		return nil, errors.Join(ErrUserNotFound,
			fmt.Errorf("get not enough users: %v", result))
	}
	return result, nil
}
