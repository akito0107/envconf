# envconf

Tiny environment variable mapper for golang

## Getting Started

### Prerequisites
- Go 1.11+

### Installing
```
$ go get -u github.com/akito0107/envconf
```

### Examples
1. [Simple mapping case](./examples/simple)

set environment variables.
```sh
export DB_HOST=localhost
export DB_PORT=5432
```

and implement Go code like this 
```go
package main

import (
	"log"

	"github.com/akito0107/envconf"
)

// struct with `env` tagged literal
type Config struct {
    DBHost string `env:"DB_HOST"` // 
    DBPort int `env:"DB_PORT"`
}

func main() {
    var conf Config
    
    // load environment variables with `Load` method
    if err := envconf.Load(&conf, &envconf.Option{}); err != nil {
    	log.Fatal(err)
    }
    
    log.Printf("%+v\n", conf) // should be `{DBHost:localhost DBPort:5432}`
}
```

You can confirm that variables are mapped to struct via `envconf.Load`.
```
$ go run main.go
2019/03/01 16:46:14 {DBHost:localhost DBPort:5432}
```

2. [Blank environment variables](./examples/blank)
`envconf.Load`, in default behaviour, returns error when specified environment variables are blank.
If missing `DB_PORT` variable, returns error.

```sh
$ export DB_HOST=localhost
$ go run main.go
2019/03/01 16:49:56 init: relpaceEnv failed: environmentVariableNotFound Envname: DB_PORT
exit status 1
```

3. [Allow Empty](./examples/allowempty)
If you want to allow empty variables, you must specify `allow-empty` on `env` tag.
This feature convenient for setting default parameter.

```go
// struct with `env` tagged literal
type Config struct {
	DBHost string `env:"DB_HOST"`
	DBPort int    `env:"DB_PORT,allow-empty"` // set allow-empty
}

....
	
var conf Config

	// load environment variables with `Load` method
if err := envconf.Load(&conf, &envconf.Option{}); err != nil {
	log.Fatal(err) // should not be error
}

log.Printf("%+v\n", conf) // 2019/03/01 16:57:57 {DBHost:localhost DBPort:0}

conf2 := Config{DBPort: 12345} // passing default param

// load environment variables with `Load` method
if err := envconf.Load(&conf2, &envconf.Option{}); err != nil {
	log.Fatal(err)
}

log.Printf("%+v\n", conf2) // 2019/03/01 16:57:57 {DBHost:localhost DBPort:12345}
```

4. [with Dotenv](./examples/dotenv)

```sh
$ export DB_PORT=5432 # only set DB_PORT
$ go run main.go
2019/03/01 16:57:57 {DBHost:localhost DBPort:0}
2019/03/01 16:57:57 {DBHost:localhost DBPort:12345}
```



## License
This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
