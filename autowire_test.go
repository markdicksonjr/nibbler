package nibbler

import (
	"testing"
	"reflect"
	"log"
)

type A struct {
	NoOpExtension
}

type B struct {
	NoOpExtension
}

type C struct {
	NoOpExtension
}

type A1 struct {
	NoOpExtension
	A *A
}

type B1 struct {
	NoOpExtension
	B *B
}

type AB struct {
	NoOpExtension
	A *A
	B *B
}

type BC struct {
	NoOpExtension
	B *B
	C *C
}

func TestAutoWireExtensions(t *testing.T) {
	var logger Logger = DefaultLogger{}

	exts := []Extension{
		&A{},
		&A1{},
		&B1{},
		&AB{},
		&B{},
		&C{},
		&BC{},
	}
	exts, err := AutoWireExtensions(&exts, &logger)

	if err != nil {
		t.Fail()
	}

	for _, v := range exts {
		log.Println(reflect.TypeOf(v).String())
	}

	aIndex := IndexOfType(exts, "*nibbler.A")
	a1Index := IndexOfType(exts, "*nibbler.A1")
	abIndex := IndexOfType(exts, "*nibbler.AB")

	if aIndex == -1 || a1Index == -1 || aIndex > a1Index {
		t.Fatal("A at index", aIndex, "is not in correct index relative to A1 at index", a1Index)
	}

	if aIndex == -1 || abIndex == -1 || aIndex > abIndex {
		t.Fatal("A at index", aIndex, "is not in correct index relative to AB at index", abIndex)
	}
}

func IndexOfType(exts []Extension, typeName string) int {
	return SliceIndex(len(exts), func(i int) bool {
		return reflect.TypeOf(exts[i]).String() == typeName
	})
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}