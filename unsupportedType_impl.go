// Code generated by "generr"; DO NOT EDIT.
package envconf

import (
	"fmt"

	"github.com/pkg/errors"
)

func IsUnsupportedType(err error) (bool, string, string) {
	err = errors.Cause(err)
	var propName string
	var typeName string
	if e, ok := err.(unsupportedType); ok {
		propName, typeName = e.UnsupportedType()
		return true, propName, typeName
	}
	return false, propName, typeName
}

type UnsupportedType struct {
	PropName string
	TypeName string
}

func (e *UnsupportedType) UnsupportedType() (string, string) {
	return e.PropName, e.TypeName
}
func (e *UnsupportedType) Error() string {
	return fmt.Sprintf("unsupportedType PropName: %v TypeName: %v", e.PropName, e.TypeName)
}
