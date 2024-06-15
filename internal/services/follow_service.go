package services

import (
	"context"
	"errors"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/pagination"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

const DEFAUTL_SIZE = 50

type followService struct {
	pbFollow.UnsafeFollowServiceServer
	storage *user.SotrageUsers
}

// FollowUser implements v1.FollowServiceServer.
func (f *followService) FollowUser(ctx context.Context, req *pbFollow.FollowUserRequest) (*emptypb.Empty, error) {
	usersInfo, err := f.storage.GetUserInfo(req.Username, req.Following)
	if err != nil {
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 2 {
		return nil, ErrUserNotFound
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
		return nil, ErrUserNotFound
	}

	err = f.storage.SetFollowing(userId, targetId)
	if err != nil {
		return nil, errors.Join(ErrFollow, err)
	}
	return &emptypb.Empty{}, nil
}

// UnFollowUser implements v1.FollowServiceServer.
func (f *followService) UnFollowUser(ctx context.Context, req *pbFollow.UnFollowUserRequest) (*emptypb.Empty, error) {
	usersInfo, err := f.storage.GetUserInfo(req.Username, req.Following)
	if err != nil {
		return nil, errors.Join(ErrUserNotFound, err)
	}
	if len(usersInfo) != 2 {
		return nil, ErrUserNotFound
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
		return nil, ErrUserNotFound
	}
	err = f.storage.UnsetFollowing(userId, targetId)
	if err != nil {
		return nil, errors.Join(ErrFollow, err)
	}
	return &emptypb.Empty{}, nil
}

// GetFollowers implements v1.FollowServiceServer.
func (f *followService) GetFollowers(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	usersInfo, err := f.storage.GetUserInfo(req.Username, req.Username)
	if err != nil {
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
		if parsErr := token.DecodePageTokenStruct(req.PageToken); parsErr != nil {
			token.LastId = 0
			token.Size = DEFAUTL_SIZE
		}

	}
	followersInfo, err := f.storage.GetUserWithFollowers(userId, token.LastId, token.Size)
	if err != nil {
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

	resp := &pbFollow.GetCommonResponse{
		NextPageToken: nextToken,
		Usernames:     followers,
	}
	return resp, nil
}

// GetFollowing implements v1.FollowServiceServer.
func (f *followService) GetFollowing(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	usersInfo, err := f.storage.GetUserInfo(req.Username, req.Username)
	if err != nil {
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
		if parsErr := token.DecodePageTokenStruct(req.PageToken); parsErr != nil {
			token.LastId = 0
			token.Size = DEFAUTL_SIZE
		}

	}
	followingInfo, err := f.storage.GetUserWithFollowing(userId, token.LastId, token.Size)
	if err != nil {
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

	resp := &pbFollow.GetCommonResponse{
		NextPageToken: nextToken,
		Usernames:     following,
	}
	return resp, nil
}

// GetFriends implements v1.FollowServiceServer.
func (f *followService) GetFriends(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	usersInfo, err := f.storage.GetUserInfo(req.Username, req.Username)
	if err != nil {
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
		if parsErr := token.DecodePageTokenStruct(req.PageToken); parsErr != nil {
			token.LastId = 0
			token.Size = DEFAUTL_SIZE
		}

	}
	friendInfos, err := f.storage.GetUserWithFriends(userId, token.LastId, token.Size)
	if err != nil {
		return nil, err
	}
	friendIds := []int{}
	for _, info := range friendInfos {
		friendIds = append(friendIds, info.FollowerId)
		if info.FollowerId > token.LastId {
			token.LastId = info.FollowerId
		}
	}
	nextToken := token.String()

	resp := &pbFollow.GetCommonResponse{
		NextPageToken: nextToken,
	}
	friends, err := f.storage.GetUserInfoById(friendIds...)
	if err != nil {
		return nil, err
	}
	for _, friend := range friends {
		resp.Usernames = append(resp.Usernames, friend.Username)
	}

	return resp, nil
}

func NewFollowService(db *gorm.DB) pbFollow.FollowServiceServer {
	return &followService{
		storage: user.New(db),
	}
}
