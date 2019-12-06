package main

import (
	"gopkg.in/ini.v1"
)

type Conf struct {
	Server *ServerConf
	Redis  *RedisConf
}

type ServerConf struct {
	HttpPort int
}

type RedisConf struct {
	Host     string
	Password string
	Port     int
	DB       int
}

func InitConfig() *Conf {
	cfg, err := ini.Load("conf.ini")
	if err != nil {
		panic(err)
	}

	server := new(ServerConf)
	redis := new(RedisConf)
	err = cfg.Section("server").MapTo(server)
	if err != nil {
		panic(err)
	}
	err = cfg.Section("redis").MapTo(redis)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Server: server,
		Redis:  redis,
	}
}
