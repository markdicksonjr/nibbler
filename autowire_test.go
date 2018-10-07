package nibbler

import (
	"testing"
	"reflect"
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

type D interface {
	Extension
}

type D0 struct {
	NoOpExtension
}

type E struct {
	NoOpExtension
	D *D
}

type F struct {
	NoOpExtension
	D D
}

func TestAutoWireExtensions(t *testing.T) {
	var logger Logger = DefaultLogger{}

	a1 := A1{}
	b1 := B1{}
	ab := AB{}
	bc := BC{}

	exts := []Extension{
		&A{},
		&a1,
		&b1,
		&ab,
		&B{},
		&C{},
		&bc,
	}
	exts, err := AutoWireExtensions(&exts, &logger)

	if err != nil {
		t.Fail()
	}

	aIndex := IndexOfType(exts, "*nibbler.A")
	a1Index := IndexOfType(exts, "*nibbler.A1")
	abIndex := IndexOfType(exts, "*nibbler.AB")

	if aIndex == -1 || a1Index == -1 || aIndex > a1Index {
		t.Fatal("A at index", aIndex, "is not in correct index relative to A1 at index", a1Index)
	}

	if a1.A == nil {
		t.Fatal("A1.A not wired")
	}

	if aIndex == -1 || abIndex == -1 || aIndex > abIndex {
		t.Fatal("A at index", aIndex, "is not in correct index relative to AB at index", abIndex)
	}

	if ab.A == nil {
		t.Fatal("AB.A not wired")
	}

	if bc.B == nil {
		t.Fatal("BC.B not wired")
	}

	if bc.C == nil {
		t.Fatal("BC.C not wired")
	}
}

func TestAutoWireExtensionsForInterfaces(t *testing.T) {
	var logger Logger = DefaultLogger{}

	e := E{}
	f := F{}
	exts := []Extension{
		&e,
		&D0{},
		&f,
	}
	exts, err := AutoWireExtensions(&exts, &logger)

	if err != nil {
		t.Fail()
	}

	fIndex := IndexOfType(exts, "*nibbler.F")
	eIndex := IndexOfType(exts, "*nibbler.E")
	dIndex := IndexOfType(exts, "*nibbler.D0")

	if eIndex == -1 || dIndex == -1 || eIndex < dIndex {
		t.Fatal("E at index", eIndex, "is not in correct index relative to D0 at index", dIndex)
	}

	if e.D == nil {
		t.Fatal("E.D not wired")
	}

	if f.D == nil {
		t.Fatal("F.D not wired")
	}
	if fIndex == -1 || dIndex == -1 || fIndex < dIndex {
		t.Fatal("F at index", fIndex, "is not in correct index relative to D0 at index", dIndex)
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