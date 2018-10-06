package nibbler

import (
	"reflect"
	"errors"
	"unsafe"
	"sort"
)

type dependency struct {
	parents			[]*dependency
	extension		*Extension
	typeName		string // mostly here for debugging, at the moment
}

// TODO: this will blow up if there's a cycle
func AutoWireExtensions(extensions *[]Extension, logger *Logger) ([]Extension, error) {

	// make a map to store dependency records by name
	treeMap := make(map[string]*dependency)

	// get the type of Extension, as it will be checked against often
	extensionInterfaceType := reflect.TypeOf(new(Extension)).Elem()

	// dereference extensions for ease of use
	extensionValues := *extensions

	// build a map of type name -> node
	for _, e := range extensionValues {
		thisExt := e
		typeName := reflect.TypeOf(e).String()
		treeMap[typeName] = &dependency{
			extension: &thisExt,
			typeName: typeName,
		}
	}

	// go through the list of extensions again to assign fields and attach dependents to extensions
	for extIndex, ext := range extensionValues {
		extensionType := reflect.TypeOf(ext)
		extensionValue := reflect.ValueOf(ext).Elem()
		fieldCount := extensionValue.NumField()
		thisExtensionDependency := treeMap[reflect.TypeOf(ext).String()]

		// loop through the fields for this extension
		for i:=0; i<fieldCount; i++ {
			fieldTypeAssignable := extensionType.Elem().Field(i)
			fieldValue := extensionValue.Field(i)

			// if we've encountered a pointer field
			if fieldValue.Kind() == reflect.Ptr {

				// if we've encountered a field that implements Extension
				if fieldValue.Type().Implements(extensionInterfaceType) {

					(*logger).Debug("autowiring " + fieldTypeAssignable.Name + " " + fieldTypeAssignable.Type.String() +
						" into " + extensionType.Elem().Name() + " " + extensionValue.Type().String())

					// TODO: this section looks repeated below

					// get the tree node by name
					mapExt := treeMap[fieldTypeAssignable.Type.String()]

					// if it's not found, something very bad happened
					if mapExt == nil {
						return nil, errors.New("could not autowire " + fieldValue.Type().Name() + " into " + extensionType.Name())
					}

					// if the value isn't set, populate it
					if fieldValue.IsNil() {
						unsafeExt := unsafe.Pointer(mapExt.extension)
						ptr := reflect.NewAt(fieldValue.Type(), unsafeExt)
						fieldValue.Set(ptr.Elem())
					}

					thisExtensionDependency.parents = append(thisExtensionDependency.parents, mapExt)
				} else {

					// look through all extensions to see if one of them implements the interface in question
					for compareIndex, compareExt := range extensionValues {

						// if the value is unset and not the one we're comparing against
						if compareIndex != extIndex && fieldValue.IsNil() {

							// the extension is assignable to the field
							compareExtensionType := reflect.TypeOf(compareExt)
							if fieldValue.Kind() == reflect.Interface &&  compareExtensionType.Implements(fieldValue.Type().Elem()) {

								(*logger).Debug("autowiring instance of " + compareExtensionType.String() +
									" as " + fieldTypeAssignable.Name + " " + fieldTypeAssignable.Type.String() +
									" into " + extensionType.Elem().Name() + " " + extensionValue.Type().String())

								// get the tree node by name
								mapExt := treeMap[compareExtensionType.String()]

								// if it's not found, something very bad happened
								if mapExt == nil {
									return nil, errors.New("could not autowire " + fieldValue.Type().Name() + " into " + extensionType.Name())
								}

								// if the value isn't set, populate it
								if fieldValue.IsNil() {
									unsafeExt := unsafe.Pointer(mapExt.extension)
									ptr := reflect.NewAt(fieldValue.Type(), unsafeExt)
									fieldValue.Set(ptr.Elem())
								}

								thisExtensionDependency.parents = append(thisExtensionDependency.parents, mapExt)
							}
						}
					}
				}
			}
		}
	}

	return orderExtensions(treeMap), nil
}

// reorder extensions based on dependencies
func orderExtensions(treeMap map[string]*dependency) []Extension {

	// convert treemap to a slice of dependencies
	var dependencyList []*dependency
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
	var sortedExtensions []Extension
	for _, v := range dependencyList {
		sortedExtensions = append(sortedExtensions, *v.extension)
	}

	return sortedExtensions
}

func isDescendant(candidateChild *dependency, candidateAncestor *dependency) bool {

	// the base case is that no parents are left
	if candidateChild.parents == nil || len(candidateChild.parents) == 0 {
		return false
	}

	// loop through parents
	for _, parent := range candidateChild.parents {

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