package redis

import (
	"github.com/go-redis/redis"
	"github.com/markdicksonjr/nibbler"
	"net/url"
)

type Extension struct {
	nibbler.NoOpExtension

	Url      string
	Password string
	DB       int
	Client   *redis.Client
}

func (s *Extension) Init(app *nibbler.Application) error {
	if len(s.Url) == 0 {
		s.Url = app.GetConfiguration().Raw.Get("redis", "url").String("")

		if len(s.Url) == 0 {
			s.Url = app.GetConfiguration().Raw.Get("rediscloud", "url").String("")
		}

		if len(s.Url) == 0 {
			s.Url = app.GetConfiguration().Raw.Get("database", "url").String("")
		}

		if len(s.Url) > 0 {
			parsedUrl, err := url.Parse(s.Url)
			if err != nil {
				return nil
			}

			if parsedUrl.Host != "" {
				s.Url = parsedUrl.Host

				if parsedUrl.User != nil {
					if pass, set := parsedUrl.User.Password(); set {
						s.Password = pass
					}
				}
			}
		}
	}

	if len(s.Password) == 0 {
		s.Password = app.GetConfiguration().Raw.Get("redis", "password").String("")

		if len(s.Password) == 0 {
			s.Password = app.GetConfiguration().Raw.Get("database", "password").String("")
		}
	}

	return s.Connect()
}

func (s *Extension) Connect() error {
	s.Client = redis.NewClient(&redis.Options{
		Addr:     s.Url,
		Password: s.Password,
		DB:       s.DB, // 0 = use default DB
	})

	_, err := s.Client.Ping().Result()

	return err
}
