package autowire

import (
	"reflect"
	"errors"
	"unsafe"
	"github.com/markdicksonjr/nibbler"
	"sort"
)

type Dependency struct {
	Parents		[]*Dependency
	Extension	*nibbler.Extension
}

// TODO: this will blow up if there's a cycle
func AutoWire(extensions *[]nibbler.Extension, logger nibbler.Logger) ([]nibbler.Extension, error) {
	treeMap := make(map[string]*Dependency)

	extensionInterfaceType := reflect.TypeOf(new(nibbler.Extension)).Elem()
	exts := *extensions

	// build a map of type name -> node
	for _, e := range exts {
		thisExt := e
		treeMap[reflect.TypeOf(e).String()] = &Dependency{Extension: &thisExt}
	}

	// go through the list again to assign fields and attach dependents to extensions
	for _, ext := range exts {
		extensionType := reflect.TypeOf(ext)
		extensionValue := reflect.ValueOf(ext).Elem()
		attributeCount := extensionValue.NumField()
		thisExt := treeMap[reflect.TypeOf(ext).String()]

		for i:=0; i<attributeCount; i++ {
			fieldTypeAssignable := extensionType.Elem().Field(i)
			fieldValue := extensionValue.Field(i)

			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Implements(extensionInterfaceType) {
				logger.Debug("autowiring " + fieldTypeAssignable.Name + " " + fieldTypeAssignable.Type.String() +
					" into " + extensionType.Elem().Name() + " " + extensionValue.Type().String())

				mapExt := treeMap[fieldTypeAssignable.Type.String()]

				if mapExt == nil {
					return nil, errors.New("could not autowire " + fieldValue.Type().Name() + " into " + extensionType.Name())
				}

				// if the value isn't set, populate it
				if fieldValue.IsNil() {
					unsafeExt := unsafe.Pointer(mapExt.Extension)
					ptr := reflect.NewAt(fieldValue.Type(), unsafeExt)
					fieldValue.Set(ptr.Elem())
				}

				thisExt.Parents = append(thisExt.Parents, mapExt)
				//mapExt.Parents = append(mapExt.Parents, thisExt)
			}
		}
	}

	return orderExtensions(treeMap), nil
}

// reorder extensions based on dependencies
func orderExtensions(treeMap map[string]*Dependency) []nibbler.Extension {

	// convert treemap to a slice of dependencies
	var dependencyList []*Dependency
	for _, v := range treeMap {
		dependencyList = append(dependencyList, v)
	}

	// use sort to sort the new slice
	sort.SliceStable(dependencyList, func(a, b int) bool {
		extA := dependencyList[a]
		extB := dependencyList[b]

		return !isDescendant(extA, extB)
	})

	// convert the slice of dependencies to a slice of extensions
	var sortedExtensions []nibbler.Extension
	for _, v := range dependencyList {
		sortedExtensions = append(sortedExtensions, *v.Extension)
	}

	return sortedExtensions
}

func isDescendant(candidateChild *Dependency, candidateAncestor *Dependency) bool {

	// the base case is that no parents are left
	if candidateChild.Parents == nil || len(candidateChild.Parents) == 0 {
		return false
	}

	// loop through parents
	for _, parent := range candidateChild.Parents {

		// if one is the candidate, we've proven the child to be a descendant
		if parent == candidateAncestor {
			return true
		}

		// if there are parents to traverse, run this function against that and the candidate child
		if isDescendant(parent, candidateAncestor) {
			return true
		}
	}

	// no ancestor was found
	return false
}