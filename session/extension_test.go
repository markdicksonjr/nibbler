package session

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

func TestExtension_InitNeedsConnector(t *testing.T) {
	e := Extension{}
	if err := e.Init(&nibbler.Application{}); err == nil {
		t.Fatal("an error was not returned from init to extension without connector")
	} else if err.Error() != requiresConnectorError {
		t.Fatal("got the wrong error message from init to extension without connector")
	}
}

func TestExtension_GetName(t *testing.T) {
	e := Extension{}
	if e.GetName() != "session" {
		t.Fatal("the wrong name was returned by the extension")
	}
}


