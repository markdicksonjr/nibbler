package client

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Extension struct {
	nibbler.NoOpExtension
	conn *ssh.Client
	app  *nibbler.Application

	Username string
	Password string
	Host     string
	HostKey  *ssh.PublicKey

	Client *sftp.Client
}

func (s *Extension) Init(app *nibbler.Application) error {
	s.app = app
	return s.Connect()
}

func (s *Extension) Connect() error {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
		s.conn = nil
	}

	if s.Client != nil {
		if err := s.Client.Close(); err != nil {
			return err
		}
		s.Client = nil
	}

	if len(s.Username) == 0 {
		s.Username = s.app.GetConfiguration().Raw.Get("sftp", "client", "username").String("")
	}

	if len(s.Password) == 0 {
		s.Password = s.app.GetConfiguration().Raw.Get("sftp", "client", "password").String("")
	}

	if len(s.Host) == 0 {
		s.Host = s.app.GetConfiguration().Raw.Get("sftp", "client", "host").String("")
	}

	callback := ssh.InsecureIgnoreHostKey()

	if s.HostKey != nil {
		callback = ssh.FixedHostKey(*s.HostKey)
	}

	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: callback,
	}

	var err error
	s.conn, err = ssh.Dial("tcp", s.Host, config)
	if err != nil {
		return err
	}

	s.Client, err = sftp.NewClient(s.conn)
	return err
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
		s.conn = nil
	}

	if s.Client != nil {
		if err := s.Client.Close(); err != nil {
			return err
		}
		s.Client = nil
	}
	return nil
}
