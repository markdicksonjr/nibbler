package session

import (
	"bytes"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/markdicksonjr/nibbler"
	"net/http"
	"net/http/httptest"
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

func TestExtension_GetName_DefaultValue(t *testing.T) {
	e := Extension{}
	if e.GetName() != "session" {
		t.Fatal("the wrong name was returned by the extension")
	}
}

func TestExtension_EnforceLogin_NoCaller(t *testing.T) {
	e := Extension{
		StoreConnector: &MockStoreConnector{
			Store: &MockStore{Session: &sessions.Session{
				Values: map[interface{}]interface{}{
					"user": nil,
				},
			}},
		},
	}

	if err := e.Init(&nibbler.Application{
		Logger: nibbler.SilentLogger{},
	}); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", bytes.NewBufferString("{}"))
	resWriter := httptest.NewRecorder()

	isAuthed := false
	e.EnforceLoggedIn(func(writer http.ResponseWriter, request *http.Request) {
		isAuthed = true
	})(resWriter, req)

	if isAuthed {
		t.Fatal("not supposed to be able to reach resource")
	}
}

func TestExtension_EnforceLogin_CallerAllowed(t *testing.T) {
	e := Extension{
		StoreConnector: &MockStoreConnector{
			Store: &MockStore{Session: &sessions.Session{
				Values: map[interface{}]interface{}{
					"user": "{\"name\": \"bob\"}",
				},
			}},
		},
	}

	e.Init(&nibbler.Application{
		Logger: nibbler.SilentLogger{},
	})
	req := httptest.NewRequest("GET", "/", bytes.NewBufferString("{}"))
	resWriter := httptest.NewRecorder()

	isAuthed := false
	e.EnforceLoggedIn(func(writer http.ResponseWriter, request *http.Request) {
		isAuthed = true
	})(resWriter, req)

	if !isAuthed {
		t.Fatal("supposed to be able to reach resource")
	}
}

func TestExtension_EnforceLogin_GetCallerFailPropagates(t *testing.T) {
	e := Extension{
		StoreConnector: &MockStoreConnector{
			Store: &MockStoreFailedGet{
				ErrOnGet: errors.New("some error"),
				MockStore: MockStore{Session: &sessions.Session{
					Values: nil,
				}},
			},
		},
	}

	e.Init(&nibbler.Application{
		Logger: nibbler.SilentLogger{},
	})
	req := httptest.NewRequest("GET", "/", bytes.NewBufferString("{}"))
	resWriter := httptest.NewRecorder()

	isAuthed := false
	e.EnforceLoggedIn(func(writer http.ResponseWriter, request *http.Request) {
		isAuthed = true
	})(resWriter, req)

	if isAuthed {
		t.Fatal("not supposed to be able to reach resource")
	}

	// TODO: validate status code, and error
}
