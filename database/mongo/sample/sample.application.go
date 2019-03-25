package main

import (
	"context"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/mongo"
	"log"
)

type Animal struct {
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"`
}

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	mongoExtension := mongo.Extension{}
	extensions := []nibbler.Extension{
		&mongoExtension,
	}

	// initialize the application
	app := nibbler.Application{}
	if err = app.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	collection := mongoExtension.Client.Database("baz").Collection("qux")

	animal := Animal{"Calvin", "Cow"}
	insertResult, err := collection.InsertOne(context.TODO(), &animal)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(insertResult.InsertedID)

	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
