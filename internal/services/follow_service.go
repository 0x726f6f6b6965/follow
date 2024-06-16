package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	"github.com/0x726f6f6b6965/follow/internal/storage/follower"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/pagination"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

const DEFAUTL_SIZE = 50

type followService struct {
	pbFollow.UnsafeFollowServiceServer
	logger          *zap.Logger
	userStorage     user.SotrageUsers
	followerStorage follower.SotrageFollowers
}

// FollowUser follow a user.
func (f *followService) FollowUser(ctx context.Context, req *pbFollow.FollowUserRequest) (*emptypb.Empty, error) {
	usersInfo, err := f.userStorage.GetUserInfo(req.Username, req.Following)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Error("get user info failed", zap.Any("request", req), zap.Error(err))
		}
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 2 {
		return nil, errors.Join(ErrUserNotFound,
			fmt.Errorf("get not enough users: %v", usersInfo))
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
	usersInfo, err := f.userStorage.GetUserInfo(req.Username, req.Following)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Error("get user info failed", zap.Any("request", req), zap.Error(err))
		}
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 2 {
		return nil, errors.Join(ErrUserNotFound,
			fmt.Errorf("get not enough users: %v", usersInfo))
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
	usersInfo, err := f.userStorage.GetUserInfo(req.Username, req.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Error("get user info failed", zap.Any("request", req), zap.Error(err))
		}
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 1 {
		return nil, ErrUserNotFound
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
	if lastId == token.LastId {
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
	usersInfo, err := f.userStorage.GetUserInfo(req.Username, req.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Error("get user info failed", zap.Any("request", req), zap.Error(err))
		}
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 1 {
		return nil, ErrUserNotFound
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
	if lastId == token.LastId {
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
	usersInfo, err := f.userStorage.GetUserInfo(req.Username, req.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Error("get user info failed", zap.Any("request", req), zap.Error(err))
		}
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 1 {
		return nil, ErrUserNotFound
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
	if lastId == token.LastId {
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

func NewFollowService(storageUser user.SotrageUsers, storageFollower follower.SotrageFollowers, logger *zap.Logger) pbFollow.FollowServiceServer {
	return &followService{
		userStorage:     storageUser,
		followerStorage: storageFollower,
		logger:          logger,
	}
}
