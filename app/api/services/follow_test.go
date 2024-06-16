package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/0x726f6f6b6965/follow/mocks"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestFollowUser(t *testing.T) {
	api := ininFollowAPI()
	t.Run("multiple requsts", func(t *testing.T) {
		count := 0
		times := 100000
		mocks.FollowUserFunc = func(ctx context.Context, req *pbFollow.FollowUserRequest) (*emptypb.Empty, error) {
			count += 1
			return &emptypb.Empty{}, nil
		}
		wg := sync.WaitGroup{}
		errChan := make(chan int, times)
		for i := 0; i < times; i++ {
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				ctx := getTestGinContext(w)
				body := &pbFollow.FollowUserRequest{
					Username:  "abc",
					Following: "def",
				}
				mockJsonPost(ctx, body)

				api.FollowUser(ctx)
				if w.Code != http.StatusOK {
					errChan <- 1
				}
			}()
			wg.Add(1)
		}
		wg.Wait()
		close(errChan)
		errCount := 0
		for val := range errChan {
			errCount += val
		}
		assert.Less(t, count, times)
		assert.Zero(t, errCount)
	})
}

func TestUnFollowUser(t *testing.T) {
	api := ininFollowAPI()
	t.Run("multiple requsts", func(t *testing.T) {
		count := 0
		times := 100000
		mocks.UnFollowUserFunc = func(ctx context.Context, req *pbFollow.UnFollowUserRequest) (*emptypb.Empty, error) {
			count += 1
			return &emptypb.Empty{}, nil
		}
		wg := sync.WaitGroup{}
		errChan := make(chan int, times)
		for i := 0; i < times; i++ {
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				ctx := getTestGinContext(w)
				body := &pbFollow.UnFollowUserRequest{
					Username:  "abc",
					Following: "def",
				}
				mockJsonPost(ctx, body)

				api.UnFollowUser(ctx)
				if w.Code != http.StatusOK {
					errChan <- 1
				}
			}()
			wg.Add(1)
		}

		wg.Wait()
		close(errChan)
		errCount := 0
		for val := range errChan {
			errCount += val
		}
		assert.Less(t, count, times)
		assert.Zero(t, errCount)
	})
}

func TestGetFollowers(t *testing.T) {
	api := ininFollowAPI()
	t.Run("multiple requsts", func(t *testing.T) {
		count := 0
		times := 100000
		mocks.GetFollowersFunc = func(context.Context, *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
			count += 1
			resp := &pbFollow.GetCommonResponse{
				Usernames: []string{"a", "b", "c"},
			}
			return resp, nil
		}
		wg := sync.WaitGroup{}
		errChan := make(chan int, times)
		for i := 0; i < times; i++ {
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				ctx := getTestGinContext(w)
				body := &pbFollow.GetCommonRequest{
					Username: "abc",
				}
				mockJsonPost(ctx, body)

				api.GetFollowers(ctx)
				if w.Code != http.StatusOK {
					errChan <- 1
				}
			}()
			wg.Add(1)
		}

		wg.Wait()
		close(errChan)
		errCount := 0
		for val := range errChan {
			errCount += val
		}
		assert.Less(t, count, times)
		assert.Zero(t, errCount)
	})
}

func TestGetFollowing(t *testing.T) {
	api := ininFollowAPI()
	t.Run("multiple requsts", func(t *testing.T) {
		count := 0
		times := 100000
		mocks.GetFollowingFunc = func(context.Context, *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
			count += 1
			resp := &pbFollow.GetCommonResponse{
				Usernames: []string{"a", "b", "c"},
			}
			return resp, nil
		}
		wg := sync.WaitGroup{}
		errChan := make(chan int, times)
		for i := 0; i < times; i++ {
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				ctx := getTestGinContext(w)
				body := &pbFollow.GetCommonRequest{
					Username: "abc",
				}
				mockJsonPost(ctx, body)

				api.GetFollowing(ctx)
				if w.Code != http.StatusOK {
					errChan <- 1
				}
			}()
			wg.Add(1)
		}

		wg.Wait()
		close(errChan)
		errCount := 0
		for val := range errChan {
			errCount += val
		}
		assert.Less(t, count, times)
		assert.Zero(t, errCount)
	})
}

func TestGetFriends(t *testing.T) {
	api := ininFollowAPI()
	t.Run("multiple requsts", func(t *testing.T) {
		count := 0
		times := 100000
		mocks.GetFriendsFunc = func(context.Context, *pbFollow.GetCommonRequest) (*pbFollow.GetCommonResponse, error) {
			count += 1
			resp := &pbFollow.GetCommonResponse{
				Usernames: []string{"a", "b", "c"},
			}
			return resp, nil
		}
		wg := sync.WaitGroup{}
		errChan := make(chan int, times)
		for i := 0; i < times; i++ {
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				ctx := getTestGinContext(w)
				body := &pbFollow.GetCommonRequest{
					Username: "abc",
				}
				mockJsonPost(ctx, body)

				api.GetFriends(ctx)
				if w.Code != http.StatusOK {
					errChan <- 1
				}
			}()
			wg.Add(1)
		}

		wg.Wait()
		close(errChan)
		errCount := 0
		for val := range errChan {
			errCount += val
		}
		assert.Less(t, count, times)
		assert.Zero(t, errCount)
	})
}

func ininFollowAPI() FollowAPI {
	m := &mocks.MockFollowService{}
	return NewFollowAPI(m)
}

// mock gin context
func getTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx
}

func mockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
