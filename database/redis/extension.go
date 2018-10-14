package redis

import (
	_ "github.com/go-redis/redis"
	"github.com/markdicksonjr/nibbler"
	"github.com/go-redis/redis"
)

type Extension struct {
	nibbler.NoOpExtension

	Url			string
	Password	string
	DB			int
	Client		*redis.Client
}

func (s *Extension) Init(app *nibbler.Application) error {
	if len(s.Url) == 0 {
		s.Url = (*app.GetConfiguration().Raw).Get("redis", "url").String("")

		if len(s.Url) == 0 {
			s.Url = (*app.GetConfiguration().Raw).Get("database", "url").String("")
		}
	}

	if len(s.Password) == 0 {
		s.Password = (*app.GetConfiguration().Raw).Get("redis", "password").String("")

		if len(s.Password) == 0 {
			s.Password = (*app.GetConfiguration().Raw).Get("database", "password").String("")
		}
	}

	s.Client = redis.NewClient(&redis.Options{
		Addr:     s.Url,
		Password: s.Password,
		DB:       s.DB,			// 0 = use default DB
	})

	_, err := s.Client.Ping().Result()

	return err
}
