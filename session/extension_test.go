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


