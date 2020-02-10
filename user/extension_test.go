package user

import (
	"github.com/markdicksonjr/nibbler"
	"testing"
)

func TestImplements(t *testing.T) {
	e := Extension{}
	var base nibbler.Extension = &e
	if base == nil {
		t.Fatal("base was nil")
	}
}

func TestExtension_Init(t *testing.T) {
	e := Extension{}
	if err := e.Init(&nibbler.Application{}); err == nil {
		t.Fatal("no error was given when it should have been - no persistence extension was provided")
	} else if err.Error() != noExtensionErrorMessage {
		t.Fatal("the error given by an init was not the expected value of " + noExtensionErrorMessage + " but was " + err.Error())
	}
}
