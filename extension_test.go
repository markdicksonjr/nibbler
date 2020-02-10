package nibbler

import "testing"

func TestNoOpExtension_Init(t *testing.T) {
	e := NoOpExtension{}
	if err := e.Init(&Application{}); err != nil {
		t.Fatal("an error was returned from NoOpExtension.Init unexpectedly")
	}
}

func TestNoOpExtension_PostInit(t *testing.T) {
	e := NoOpExtension{}
	if err := e.PostInit(&Application{}); err != nil {
		t.Fatal("an error was returned from NoOpExtension.PostInit unexpectedly")
	}
}

func TestNoOpExtension_Destroy(t *testing.T) {
	e := NoOpExtension{}
	if err := e.Destroy(&Application{}); err != nil {
		t.Fatal("an error was returned from NoOpExtension.Destroy unexpectedly")
	}
}

func TestNoOpExtension_GetName(t *testing.T) {
	e := NoOpExtension{}
	if name := e.GetName(); name != "nameless" {
		t.Fatal("the wrong name was returned from the NoOpExtension")
	}
}


