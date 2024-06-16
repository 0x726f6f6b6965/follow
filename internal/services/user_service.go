package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0x726f6f6b6965/follow/internal/helper"
	"github.com/0x726f6f6b6965/follow/internal/storage/cache"
	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	boom "github.com/tylertreat/BoomFilters"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var (
	testSalt *helper.Salt
)

type userService struct {
	pbUser.UnimplementedUserServiceServer
	logger       *zap.Logger
	userStorage  user.SotrageUsers
	cacheTTL     time.Duration
	cacheStorage cache.Cache
	filter       *boom.CountingBloomFilter
}

func NewUserService(storage user.SotrageUsers, cacheStorage cache.Cache, cacheTTL time.Duration, filter *boom.CountingBloomFilter, logger *zap.Logger) pbUser.UserServiceServer {
	return &userService{
		userStorage:  storage,
		cacheStorage: cacheStorage,
		cacheTTL:     cacheTTL,
		filter:       filter,
		logger:       logger,
	}
}

// CreateUser: create a new user account.
func (s *userService) CreateUser(ctx context.Context, req *pbUser.CreateUserRequest) (*emptypb.Empty, error) {
	if helper.IsEmpty(req.Username) || helper.IsEmpty(req.Password) {
		return nil, ErrInvalidInput
	}
	// check if user already exists
	cacheInfo, _ := s.cacheStorage.Get(ctx, fmt.Sprintf(UserExistKey, req.Username))
	if !helper.IsEmpty(cacheInfo) {
		return nil, ErrUserExists
	}

	info, err := s.userStorage.GetUserInfo(req.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("CreateUser: error getting user",
				zap.Any("request", req),
				zap.Error(err))
			return nil, err
		}
	} else if len(info) != 0 {
		return nil, ErrUserExists
	}

	// salt pwd
	var salt *helper.Salt
	if testSalt == nil {
		salt = helper.CreateNewSalt()
	} else {
		salt = testSalt
	}
	pwd, err := salt.SaltInput(req.Password)
	if err != nil {
		s.logger.Error("CreateUser: error hashing password",
			zap.Any("request", req),
			zap.Error(err))
		return nil, errors.Join(ErrSalt, err)
	}

	//insert data
	insertData := &models.User{
		Username: req.Username,
		Password: pwd,
		Salt:     salt.SaltString,
	}

	err = s.userStorage.CreateUser(insertData)
	if err != nil {
		s.logger.Error("CreateUser: error inserting user",
			zap.Any("request", req),
			zap.Error(err))
		return nil, err
	}

	_ = s.filter.TestAndRemove([]byte(req.Username))
	_ = s.cacheStorage.Set(ctx, fmt.Sprintf(UserExistKey, req.Username), insertData.Id, s.cacheTTL)

	return &emptypb.Empty{}, nil
}
