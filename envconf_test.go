package envconf

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("bind normal type", func(t *testing.T) {
		defenv := getenv

		getenv = func(key string) string {
			t.Helper()
			switch key {
			case "STR":
				return "str"
			case "INTEGER":
				return "123"
			case "BOOL":
				return "true"
			case "UINT":
				return "12"
			case "FLOAT":
				return "1.234"
			default:
				t.Fatalf("unkwnon envname %s", key)
			}
			return ""
		}

		defer func() {
			getenv = defenv
		}()

		type in struct {
			Str     string  `env:"STR"`
			Integer int     `env:"INTEGER"`
			Bool    bool    `env:"BOOL"`
			Uint    uint    `env:"UINT"`
			Float   float64 `env:"FLOAT"`
		}

		var i in
		if err := Load(&i); err != nil {
			t.Fatal(err)
		}

		expect := in{
			Str:     "str",
			Integer: 123,
			Bool:    true,
			Uint:    12,
			Float:   1.234,
		}

		if !reflect.DeepEqual(i, expect) {
			t.Errorf("must be same %+v and %+v", i, expect)
		}
	})

	t.Run("error when null env", func(t *testing.T) {
		defenv := getenv

		getenv = func(key string) string {
			t.Helper()
			switch key {
			case "STR":
				return "str"
			case "INTEGER":
				return "123"
			case "UINT":
				return "12"
			case "FLOAT":
				return "1.234"
			}
			return ""
		}

		defer func() {
			getenv = defenv
		}()

		type in struct {
			Str     string  `env:"STR"`
			Integer int     `env:"INTEGER"`
			Bool    bool    `env:"BOOL"`
			Uint    uint    `env:"UINT"`
			Float   float64 `env:"FLOAT"`
		}

		var i in
		err := Load(&i)
		if err == nil {
			t.Fatal("must be err")
		}

		if ok, _ := IsEnvironmentVariableNotFound(err); !ok {
			t.Errorf("should be environmentValueNotFound but %+v", err)
		}

	})

	t.Run("allow empty", func(t *testing.T) {
		defenv := getenv

		getenv = func(key string) string {
			t.Helper()
			switch key {
			case "STR":
				return "str"
			case "INTEGER":
				return "123"
			case "UINT":
				return "12"
			}
			return ""
		}

		defer func() {
			getenv = defenv
		}()

		type in struct {
			Str     string  `env:"STR"`
			Integer int     `env:"INTEGER"`
			Bool    bool    `env:"BOOL,allow-empty"`
			Uint    uint    `env:"UINT"`
			Float   float64 `env:"FLOAT,allow-empty"`
		}

		var i in
		if err := Load(&i); err != nil {
			t.Fatal(err)
		}

		expect := in{
			Str:     "str",
			Uint:    12,
			Integer: 123,
		}

		if !reflect.DeepEqual(i, expect) {
			t.Errorf("must be same %+v and %+v", i, expect)
		}

	})
}
