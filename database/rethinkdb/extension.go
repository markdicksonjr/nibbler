package rethinkdb

import (
	"github.com/markdicksonjr/nibbler"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"net/url"
)

type Extension struct {
	nibbler.NoOpExtension

	Url      string
	Username string
	Password string
	Session  *r.Session
}

func (s *Extension) Init(app *nibbler.Application) error {
	if len(s.Url) == 0 {
		s.Url = app.GetConfiguration().Raw.Get("rethinkdb", "url").String("")

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
		s.Password = app.GetConfiguration().Raw.Get("rethinkdb", "password").String("")

		if len(s.Password) == 0 {
			s.Password = app.GetConfiguration().Raw.Get("database", "password").String("")
		}
	}

	return s.Connect()
}

func (s *Extension) Connect() error {
	var err error

	if s.Username == "" {
		s.Session, err = r.Connect(r.ConnectOpts{
			Address: s.Url,
		})
	} else {
		s.Session, err = r.Connect(r.ConnectOpts{
			Address: s.Url,
			Username: s.Username,
			Password: s.Password,
		})
	}

	return err
}

func (s *Extension) CreateDbUser() error {
	return  r.DB("rethinkdb").Table("users").Insert(map[string]string{
		"id": s.Username,
		"password": s.Password,
	}).Exec(s.Session)

	// to grant access, do something like this:
	//err = r.DB("blog").Table("posts").Grant("john", map[string]bool{
	//	"read": true,
	//	"write": true,
	//}).Exec(session)
}
