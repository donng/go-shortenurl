package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pilu/go-base62"
	"time"
)

const (
	URL_ID_KEY            = "next.url.id"         // 全局自增器
	SHORT_LINK_KEY        = "shortlink:%s:url"    // 短地址和原地址的映射
	URL_HASH_KEY          = "urlhash:%s:url"      // 地址hash和短地址的映射
	SHORT_LINK_DETAIL_KEY = "shortlink:%s:detail" // 短地址和详情的映射
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

func (r *RedisCli) Shorten(url string, exp int) (string, error) {
	h := toSHA1(url)
	d, err := r.Cli.Get(fmt.Sprintf(URL_HASH_KEY, h)).Result()
	if err == redis.Nil {
		// not exist, nothing to do
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// expiration, noting to do
		} else {
			return d, nil
		}
	}

	// increase the global counter
	err = r.Cli.Incr(URL_ID_KEY).Err()
	if err != nil {
		return "", err
	}

	id, err := r.Cli.Get(URL_ID_KEY).Int()
	if err != nil {
		return "", err
	}
	// encode global counter to base62
	// URL_ID_KEY 是单纯的数字，需要编码成字母数字的短地址形式
	eid := base62.Encode(id)

	// 存储短地址和原地址的映射
	err = r.Cli.Set(fmt.Sprintf(SHORT_LINK_KEY, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}

	// 存储哈希值和短地址的映射
	err = r.Cli.Set(fmt.Sprintf(URL_HASH_KEY, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(&URLDetail{
		URL:                 url,
		CreatedAt:           time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})
	if err != nil {
		return "", err
	}

	// 存储短地址和详情的映射
	err = r.Cli.Set(fmt.Sprintf(SHORT_LINK_DETAIL_KEY, eid), detail, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}

	return eid, nil
}

func (r *RedisCli) ShortLinkInfo(eid string) (interface{}, error) {
	detail, err := r.Cli.Get(fmt.Sprintf(SHORT_LINK_DETAIL_KEY, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{
			Code: 404,
			Err:  errors.New("unknown short URL"),
		}
	} else if err != nil {
		return "", err
	}

	return detail, nil
}

func (r *RedisCli) UnShorten(eid string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(SHORT_LINK_KEY, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{
			Code: 404,
			Err:  errors.New("unknown short URL"),
		}
	} else if err != nil {
		return "", err
	}

	return url, nil
}

func toSHA1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
