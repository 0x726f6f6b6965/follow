package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/0x726f6f6b6965/follow/app/api"
	"github.com/0x726f6f6b6965/follow/app/follow"
	"github.com/0x726f6f6b6965/follow/app/psql"
	"github.com/0x726f6f6b6965/follow/app/rds"
	"github.com/0x726f6f6b6965/follow/app/user"
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	boom "github.com/tylertreat/BoomFilters"
	"gopkg.in/yaml.v3"
)

const (
	COUNTING_BLOOM_FILTER_N               = 10000
	COUNTING_BLOOM_FILTER_FP_Rate float64 = 0.01
)

func main() {
	godotenv.Load()
	path := os.Getenv("CONFIG")
	cfg := new(config.AppConfig)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read yaml error", err)
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("unmarshal yaml error", err)
		return
	}
	db, dbCleanup, err := psql.InitDBService(cfg)
	if err != nil {
		log.Fatal("init db error", err)
		return
	}
	defer dbCleanup()
	rdb, rdbCleanup, err := rds.InitRdsService(cfg)
	if err != nil {
		log.Fatal("init redis error", err)
		return
	}
	defer rdbCleanup()

	filter := boom.NewDefaultCountingBloomFilter(COUNTING_BLOOM_FILTER_N, COUNTING_BLOOM_FILTER_FP_Rate)

	userServer, uServerCleanup, err := user.InitUserService(cfg, db, rdb, filter)
	if err != nil {
		log.Fatal("init user server error", err)
		return
	}
	defer uServerCleanup()

	followServer, fServerCleanup, err := follow.InitFollowService(cfg, db, rdb, filter)
	if err != nil {
		log.Fatal("init follow server error", err)
		return
	}
	defer fServerCleanup()

	router, err := api.InitRouter(userServer, followServer)
	if err != nil {
		log.Fatal("init gin router error", err)
		return
	}

	engine := initGin(cfg)
	router.RegisterRoutes(engine)

	engine.Run(fmt.Sprintf(":%d", cfg.HttpPort))
}

func initGin(cfg *config.AppConfig) *gin.Engine {
	gin.SetMode(func() string {
		if cfg.Env == "dev" {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "Service internal exception!",
		})
	}))
	return engine
}
