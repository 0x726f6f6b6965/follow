package services

import (
	"net/http"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/gin-gonic/gin"
)

type followAPI struct {
	server pbFollow.FollowServiceServer
}

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := f.server.FollowUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UnFollowUser implements FollowAPI.
func (f *followAPI) UnFollowUser(ctx *gin.Context) {
	req := &pbFollow.UnFollowUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Following) || helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := f.server.UnFollowUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetFollowers implements FollowAPI.
func (f *followAPI) GetFollowers(ctx *gin.Context) {
	req := &pbFollow.GetCommonRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := f.server.GetFollowers(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetFollowing implements FollowAPI.
func (f *followAPI) GetFollowing(ctx *gin.Context) {
	req := &pbFollow.GetCommonRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := f.server.GetFollowing(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetFriends implements FollowAPI.
func (f *followAPI) GetFriends(ctx *gin.Context) {
	req := &pbFollow.GetCommonRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := f.server.GetFriends(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
