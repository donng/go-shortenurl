package main

import (
	"github.com/go-redis/redis/v7"
	"time"
)

const (
	URL_ID_KEY            = "next.url.id"          // 全局自增器
	SHORT_LINK_KEY        = "shortlink:%s:url"     // 短地址和地址的映射
	URL_HASH_KEY          = "urlhash:%s:url"       // 地址hash和短地址的映射
	SHORT_LINK_DETAIL_KEY = "shortlink:%s:detail"  // 短地址和详情的映射
)

type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	URL                 string
	CreatedAt           string
	ExpirationInMinutes time.Duration
}

func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisCli{c}
}
