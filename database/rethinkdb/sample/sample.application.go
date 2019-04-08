package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/rethinkdb"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
)

func main() {

	// allocate logger and configuration
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	rethinkExtension := rethinkdb.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&rethinkExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	if err := r.DB("testdb").Table("testtable").Insert(map[string]string{
		"id": "john",
		"password": "p455w0rd",
	}).Exec(rethinkExtension.Session); err != nil {
		log.Fatal(err)
	}

	res, err := r.DB("testdb").Table("testtable").Get("john").Run(rethinkExtension.Session)
	if err != nil {
		log.Fatal(err)
	}

	var result interface{}
	if err := res.One(result); err != nil {
		log.Fatal(err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
