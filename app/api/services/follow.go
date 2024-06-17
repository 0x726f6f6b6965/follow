package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

type followAPI struct {
	server pbFollow.FollowServiceServer
	group  singleflight.Group
}

const (
	FOLLOW_USER   = "follow_user"
	UNFOLLOW_USER = "unfollow_user"
	GET_FOLLOWERS = "get_followers"
	GET_FOLLOWING = "get_following"
	GET_FRIENDS   = "get_friends"
)

type FollowAPI interface {
	FollowUser(ctx *gin.Context)
	UnFollowUser(ctx *gin.Context)
	GetFollowers(ctx *gin.Context)
	GetFollowing(ctx *gin.Context)
	GetFriends(ctx *gin.Context)
}

func NewFollowAPI(server pbFollow.FollowServiceServer) FollowAPI {
	return &followAPI{server: server}
}

// FollowUser implements FollowAPI.
func (f *followAPI) FollowUser(ctx *gin.Context) {
	req := &pbFollow.FollowUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Following) || helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	dataChan := f.group.DoChan(fmt.Sprintf("%s-%s-%s", FOLLOW_USER, req.Username, req.Following),
		func() (interface{}, error) {
			return f.server.FollowUser(ctx, req)
		})
	select {
	case <-ctx.Done():
		ctx.JSON(http.StatusInternalServerError, MessageTimeout)
		return
	case res := <-dataChan:
		if res.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, MessageSuccess)
	}
}

// UnFollowUser implements FollowAPI.
func (f *followAPI) UnFollowUser(ctx *gin.Context) {
	req := &pbFollow.UnFollowUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Unfollow) || helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	dataChan := f.group.DoChan(fmt.Sprintf("%s-%s-%s", UNFOLLOW_USER, req.Username, req.Unfollow),
		func() (interface{}, error) {
			return f.server.UnFollowUser(ctx, req)
		})
	select {
	case <-ctx.Done():
		ctx.JSON(http.StatusInternalServerError, MessageTimeout)
		return
	case res := <-dataChan:
		if res.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, MessageSuccess)
	}
}

// GetFollowers implements FollowAPI.
func (f *followAPI) GetFollowers(ctx *gin.Context) {
	username := ctx.Param("username")
	if helper.IsEmpty(username) {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	sizeStr := ctx.DefaultQuery("size", "50")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	pageToken := ctx.DefaultQuery("page_token", "")

	req := &pbFollow.GetCommonRequest{
		Username:  username,
		Size:      int32(size),
		PageToken: pageToken,
	}

	dataChan := f.group.DoChan(fmt.Sprintf("%s-%s", GET_FOLLOWERS, req.Username),
		func() (interface{}, error) {
			return f.server.GetFollowers(ctx, req)
		})
	select {
	case <-ctx.Done():
		ctx.JSON(http.StatusInternalServerError, MessageTimeout)
		return
	case res := <-dataChan:
		if res.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
			return
		}
		resp, ok := res.Val.(*pbFollow.GetCommonResponse)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, MessageInvalidResponse)
			return
		}
		if len(resp.Usernames) == 0 {
			ctx.JSON(http.StatusOK, MessageEmpty)
			return
		}
		ctx.JSON(http.StatusOK, res.Val)
	}
}

// GetFollowing implements FollowAPI.
func (f *followAPI) GetFollowing(ctx *gin.Context) {
	username := ctx.Param("username")
	if helper.IsEmpty(username) {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	sizeStr := ctx.DefaultQuery("size", "50")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	pageToken := ctx.DefaultQuery("page_token", "")

	req := &pbFollow.GetCommonRequest{
		Username:  username,
		Size:      int32(size),
		PageToken: pageToken,
	}

	dataChan := f.group.DoChan(fmt.Sprintf("%s-%s", GET_FOLLOWING, req.Username),
		func() (interface{}, error) {
			return f.server.GetFollowing(ctx, req)
		})
	select {
	case <-ctx.Done():
		ctx.JSON(http.StatusInternalServerError, MessageTimeout)
		return
	case res := <-dataChan:
		if res.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
			return
		}
		resp, ok := res.Val.(*pbFollow.GetCommonResponse)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, MessageInvalidResponse)
			return
		}
		if len(resp.Usernames) == 0 {
			ctx.JSON(http.StatusOK, MessageEmpty)
			return
		}
		ctx.JSON(http.StatusOK, res.Val)
	}
}

// GetFriends implements FollowAPI.
func (f *followAPI) GetFriends(ctx *gin.Context) {
	username := ctx.Param("username")
	if helper.IsEmpty(username) {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	sizeStr := ctx.DefaultQuery("size", "50")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, MessageInvalidInput)
		return
	}
	pageToken := ctx.DefaultQuery("page_token", "")

	req := &pbFollow.GetCommonRequest{
		Username:  username,
		Size:      int32(size),
		PageToken: pageToken,
	}
	dataChan := f.group.DoChan(fmt.Sprintf("%s-%s", GET_FRIENDS, req.Username),
		func() (interface{}, error) {
			return f.server.GetFriends(ctx, req)
		})
	select {
	case <-ctx.Done():
		ctx.JSON(http.StatusInternalServerError, MessageTimeout)
		return
	case res := <-dataChan:
		if res.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
			return
		}
		resp, ok := res.Val.(*pbFollow.GetCommonResponse)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, MessageInvalidResponse)
			return
		}
		if len(resp.Usernames) == 0 {
			ctx.JSON(http.StatusOK, MessageEmpty)
			return
		}
		ctx.JSON(http.StatusOK, res.Val)
	}
}
