package envconf

import (
	"os"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

//go:generate generr -t notPointerType -i -c -u
type notPointerType interface {
	NotPointerType()
}

//go:generate generr -t unsupportedType -i -c -u
type unsupportedType interface {
	UnsupportedType() (propName string, typeName string)
}

//go:generate generr -t environmentValueNotFound -t environmentValueNotFound -i -c -u
type environmentValueNotFound interface {
	EnvironmentValueNotFound() (envname string)
}

func Init(i interface{}) error {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return &NotPointerType{}
	}
	if err := replaceEnv(i); err != nil {
		return err
	}

	return nil
}

var getenv = func(key string) string {
	return os.Getenv(key)
}

func replaceEnv(i interface{}) error {
	interfacevalue := reflect.Indirect(reflect.ValueOf(i)).Interface()
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(interfacevalue)

	if t.Kind() != reflect.Struct {
		return &UnsupportedType{}
	}

	for j := 0; j < t.NumField(); j++ {
		sf := t.Field(j)
		tag := sf.Tag
		envname := tag.Get("env")
		if envname == "" {
			continue
		}
		value := getenv(envname)
		if value == "" {
			return &EnvironmentValueNotFound{Envname: envname}
		}
		f := v.Elem().Field(j)

		switch sf.Type.Kind() {
		case reflect.Int, reflect.Int64:
			val, err := strconv.Atoi(value)
			if err != nil {
				return errors.Wrapf(err, "typename: %s, envname: %s, value :%s, can't cast to int", sf.Type.Name(), envname, value)
			}
			f.SetInt(int64(val))
		case reflect.Uint, reflect.Uint64:
			val, err := strconv.Atoi(value)
			if err != nil {
				return errors.Wrapf(err, "typename: %s, envname: %s, value :%s, can't cast to uint", sf.Type.Name(), envname, value)
			}
			f.SetUint(uint64(val))
		case reflect.String:
			f.SetString(value)
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				return errors.Wrapf(err, "typename: %s, envname: %s, value :%s, can't cast to bool", sf.Type.Name(), envname, value)
			}
			f.SetBool(b)
		case reflect.Float32, reflect.Float64:
			fl, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return errors.Wrapf(err, "typename: %s, envname: %s, value :%s, can't cast to float", sf.Type.Name(), envname, value)
			}
			f.SetFloat(fl)
		default:
			return &UnsupportedType{TypeName: sf.Type.Name(), PropName: sf.Name}
		}
	}

	return nil
}
