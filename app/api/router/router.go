package router

import (
	"github.com/0x726f6f6b6965/follow/app/api/services"
	"github.com/gin-gonic/gin"
)

type router struct {
	userAPI   services.UserAPI
	followAPI services.FollowAPI
}
type Router interface {
	RegisterRoutes(server *gin.Engine)
}

func NewRouter(userAPI services.UserAPI, followAPI services.FollowAPI) Router {
	return &router{
		userAPI:   userAPI,
		followAPI: followAPI,
	}
}

func (r *router) RegisterRoutes(server *gin.Engine) {
	r.registerUserRouter(server.Group("/v1/user/"))
	r.registerFollowRouter(server.Group("/v1/follow/"))
}

func (r *router) registerUserRouter(group *gin.RouterGroup) {
	group.POST("/register", r.userAPI.CreateUser)
}

func (r *router) registerFollowRouter(group *gin.RouterGroup) {
	group.POST("/follow", r.followAPI.FollowUser)
	group.POST("/unfollow", r.followAPI.UnFollowUser)

	list := group.Group("/list")

	list.POST("/followers", r.followAPI.GetFollowers)
	list.POST("/following", r.followAPI.GetFollowing)
	list.POST("/friends", r.followAPI.GetFriends)
}
