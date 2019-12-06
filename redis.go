package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pilu/go-base62"
)

const (
	URLIdKey           = "next.url.id"         // 全局自增器
	ShortLinkKey       = "shortlink:%s:url"    // 短地址和原地址的映射
	URLHashKey         = "urlhash:%s:url"      // 地址hash和短地址的映射
	ShortLinkDetailKey = "shortlink:%s:detail" // 短地址和详情的映射
)

type RedisClient struct {
	cli *redis.Client
}

type URLDetail struct {
	URL                 string
	CreatedAt           string
	ExpirationInMinutes time.Duration
}

func NewRedisClient(r *Redis) *RedisClient {
	addr := fmt.Sprintf("%s:%d", r.Host, r.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: r.Password,
		DB:       r.DB,
	})
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisClient{client}
}

func (r *RedisClient) Shorten(url string, exp int64) (string, error) {
	urlHash := toSHA1(url)

	d, err := r.cli.Get(fmt.Sprintf(URLHashKey, urlHash)).Result()
	if err == redis.Nil {
		// URLHashKey not exist, nothing to do
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// expiration, noting to do
		} else {
			return d, nil
		}
	}

	// 1. increase the global counter
	id, err := r.cli.Incr(URLIdKey).Result()
	if err != nil {
		return "", err
	}

	// encode global counter to base62
	// notice: int64->int is not safe where int is 32 bits
	encodeId := base62.Encode(int(id))

	// 2. save short link and origin url
	err = r.cli.Set(fmt.Sprintf(ShortLinkKey, encodeId), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	// 3. save hash value and short link
	err = r.cli.Set(fmt.Sprintf(URLHashKey, urlHash), encodeId, time.Minute*time.Duration(exp)).Err()
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

	// 4. save short link and link detail
	err = r.cli.Set(fmt.Sprintf(ShortLinkDetailKey, encodeId), detail, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}

	return encodeId, nil
}

func (r *RedisClient) UnShorten(encodeId string) (string, error) {
	url, err := r.cli.Get(fmt.Sprintf(ShortLinkKey, encodeId)).Result()
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

func (r *RedisClient) ShortLinkInfo(encodeId string) (interface{}, error) {
	detail, err := r.cli.Get(fmt.Sprintf(ShortLinkDetailKey, encodeId)).Result()
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

func toSHA1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
