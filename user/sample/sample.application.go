package main

import (
	"log"
	//"github.com/markdicksonjr/nibbler/elasticsearch"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/user"
	NibUserSql "github.com/markdicksonjr/nibbler/user/sql"
	//NibUserElastic "github.com/markdicksonjr/nibbler/user/elastic"
	"github.com/markdicksonjr/nibbler"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadApplicationConfiguration(nil)

	// prepare models for initialization
	var models []interface{}
	models = append(models, user.User{})

	// allocate an SQL controller, providing a sql extension
	sqlController := NibUserSql.SqlPersistenceController{
		SqlExtension: &sql.Extension{
			Models: models,
		},
	}

	// allocate an ES controller, providing an ES extension
	//elasticController := NibUserElastic.ElasticPersistenceController{
	//	ElasticExtension: &elasticsearch.Extension{},
	//}

	// allocate our user extension, providing the SQL controller
	userExtension := user.Extension{
		PersistenceController:  &sqlController, // &elasticController,
	}

	// prepare extension(s) for initialization
	extensions := []nibbler.Extension{
		sqlController.SqlExtension,
		//elasticController.ElasticExtension,
		&userExtension,
	}

	// initialize the application context
	appContext := nibbler.Application{}
	err = appContext.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	// create a test user
	emailVal := "test@example.com"
	_, errCreate := userExtension.Create(&user.User{
		Email: &emailVal,
	})

	if errCreate != nil {
		log.Fatal(errCreate.Error())
	}

	uV, errFind := userExtension.GetUserByEmail(emailVal)

	if errFind != nil {
		log.Fatal(errFind.Error())
	}

	log.Println(uV)

	firstName := "testfirst"
	lastName := "testlast"
	uV.FirstName = &firstName
	uV.LastName = &lastName
	err = userExtension.Update(uV)

	if err != nil {
		log.Fatal(err.Error())
	}

	// start the app
	err = appContext.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
