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
		if err := Init(&i); err != nil {
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

}
