package mocks

import (
	"context"

	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	// FollowUser: follow a user.
	FollowUserFunc func(context.Context, *pbFollow.FollowUserRequest) (*emptypb.Empty, error)
	// UnFollowUser: unfollow a user.
	UnFollowUserFunc func(context.Context, *pbFollow.UnFollowUserRequest) (*emptypb.Empty, error)
	// GetFollowers: get followers list.
	GetFollowersFunc func(context.Context, *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error)
	// GetFollowing: get following list.
	GetFollowingFunc func(context.Context, *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error)
	// GetFriends: get friends list.
	GetFriendsFunc func(context.Context, *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error)
)

type MockFollowService struct {
	pbFollow.UnimplementedFollowServiceServer
}

// FollowUser: follow a user.
func (m *MockFollowService) FollowUser(ctx context.Context, req *pbFollow.FollowUserRequest) (*emptypb.Empty, error) {
	return FollowUserFunc(ctx, req)
}

// UnFollowUser: unfollow a user.
func (m *MockFollowService) UnFollowUser(ctx context.Context, req *pbFollow.UnFollowUserRequest) (*emptypb.Empty, error) {
	return UnFollowUserFunc(ctx, req)
}

// GetFollowers: get followers list.
func (m *MockFollowService) GetFollowers(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	return GetFollowersFunc(ctx, req)
}

// GetFollowing: get following list.
func (m *MockFollowService) GetFollowing(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	return GetFollowingFunc(ctx, req)
}

// GetFriends: get friends list.
func (m *MockFollowService) GetFriends(ctx context.Context, req *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
	return GetFriendsFunc(ctx, req)
}
