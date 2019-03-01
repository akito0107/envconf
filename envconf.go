package envconf

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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

//go:generate generr -t environmentVariableNotFound -t environmentVariableNotFound -i -c -u
type environmentVariableNotFound interface {
	EnvironmentVariableNotFound() (envname string)
}

type option struct {
	UseDotEnv       bool
	DotEnvNameAlias string
}

type Option func(*option)

func UseDotEnv() Option {
	return func(opt *option) {
		opt.UseDotEnv = true
	}
}

func DotEnvNameAlias(filename string) Option {
	return func(opt *option) {
		opt.DotEnvNameAlias = filename
	}
}

func Load(i interface{}, opts ...Option) error {
	option := &option{}
	for _, o := range opts {
		o(option)
	}

	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return &NotPointerType{}
	}

	if option.UseDotEnv {
		fname := ".env"
		if option.DotEnvNameAlias != "" {
			fname = option.DotEnvNameAlias
		}
		loadEnv(fname)
	}

	if err := replaceEnv(i); err != nil {
		return errors.Wrap(err, "init: relpaceEnv failed")
	}

	return nil
}

func loadEnv(fname string) error {
	load := os.Getenv("ENVCONF_LOAD_DOTFILE")
	if load == "disable" {
		return nil
	}

	if err := godotenv.Overload(fname); err != nil {
		return errors.Wrap(err, "init: godotenv Load failed")
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
		env := tag.Get("env")
		if env == "" {
			continue
		}

		var allowEmpty bool

		strs := strings.Split(env, ",")
		envname := strs[0]
		if len(strs) == 2 && strs[1] == "allow-empty" {
			allowEmpty = true
		}

		value := getenv(envname)
		if value == "" {
			if allowEmpty {
				continue
			}
			return &EnvironmentVariableNotFound{Envname: envname}
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
