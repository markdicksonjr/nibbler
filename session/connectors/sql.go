package connectors

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/wader/gormstore"
)

type SqlMemoryStoreConnector struct {
	SqlExtension 	*sql.Extension
	Secret     		string

	connectionInfo	string
	db				*gorm.DB
}

func (s SqlMemoryStoreConnector) Connect() (error, sessions.Store) {

	// if no DB provided, use sqlite3 memory DB
	if s.SqlExtension == nil {
		db, err := gorm.Open("sqlite3", ":memory:")

		if err != nil {
			return err, nil
		}

		s.db = db
	} else {
		s.db = s.SqlExtension.Db
	}

	if len(s.Secret) == 0 {
		return errors.New("sql connector requires secret"), nil
	}

	store := gormstore.NewOptions(s.db,
		gormstore.Options{},
		[]byte(s.Secret),
	)

	store.SessionOpts.MaxAge = s.MaxAge()

	return nil, store
}

func (s SqlMemoryStoreConnector) MaxAge() int {
	return 60 * 60 * 24 * 15 // 15 days
}
