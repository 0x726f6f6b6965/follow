package services

import (
	"net/http"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

type userAPI struct {
	server pbUser.UserServiceServer
	group  singleflight.Group
}

type UserAPI interface {
	CreateUser(ctx *gin.Context)
}

func NewUserAPI(server pbUser.UserServiceServer) UserAPI {
	return &userAPI{
		server: server,
	}
}

func (u *userAPI) CreateUser(ctx *gin.Context) {
	req := pbUser.CreateUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if helper.IsEmpty(req.Password) || helper.IsEmpty(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	dataChan := u.group.DoChan(req.Username, func() (interface{}, error) {
		return u.server.CreateUser(ctx, &req)
	})
	select {
	case <-ctx.Done():
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrTimeout.Error()})
		return
	case res := <-dataChan:
		if res.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, res.Val)
	}
}
