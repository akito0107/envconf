package main

import (
	"log"

	"github.com/akito0107/envconf"
)

// struct with `env` tagged literal
type Config struct {
	DBHost string `env:"DB_HOST"`
	DBPort int    `env:"DB_PORT,allow-empty"`
}

func main() {
	var conf Config

	// load environment variables with `Load` method
	if err := envconf.Load(&conf, &envconf.Option{}); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", conf) // 2019/03/01 16:57:57 {DBHost:localhost DBPort:0}

	conf2 := Config{DBPort: 12345} // passing default param

	// load environment variables with `Load` method
	if err := envconf.Load(&conf2, &envconf.Option{}); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", conf2) // 2019/03/01 16:57:57 {DBHost:localhost DBPort:12345}
}
