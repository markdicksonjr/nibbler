package autowire

import (
	"testing"
	"github.com/markdicksonjr/nibbler"
	"reflect"
	"log"
)

type A struct {
	nibbler.NoOpExtension
}

type B struct {
	nibbler.NoOpExtension
}

type C struct {
	nibbler.NoOpExtension
}

type A1 struct {
	nibbler.NoOpExtension
	A *A
}

type B1 struct {
	nibbler.NoOpExtension
	B *B
}

type AB struct {
	nibbler.NoOpExtension
	A *A
	B *B
}

type BC struct {
	nibbler.NoOpExtension
	B *B
	C *C
}

func TestAutoWire(t *testing.T) {
	exts := []nibbler.Extension{
		&A{},
		&A1{},
		&B1{},
		&AB{},
		&B{},
		&C{},
		&BC{},
	}
	exts, err := AutoWire(&exts, nibbler.DefaultLogger{})

	if err != nil {
		t.Fail()
	}

	for _, v := range exts {
		log.Println(reflect.TypeOf(v).String())
	}

	aIndex := IndexOfType(exts, "*autowire.A")
	a1Index := IndexOfType(exts, "*autowire.A1")
	abIndex := IndexOfType(exts, "*autowire.AB")

	if aIndex == -1 || a1Index == -1 || aIndex > a1Index {
		t.Fatal("A at index", aIndex, "is not in correct index relative to A1 at index", a1Index)
	}

	if aIndex == -1 || abIndex == -1 || aIndex > abIndex {
		t.Fatal("A at index", aIndex, "is not in correct index relative to AB at index", abIndex)
	}
}

func IndexOfType(exts []nibbler.Extension, typeName string) int {
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