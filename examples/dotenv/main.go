package main

import (
	"log"

	"github.com/akito0107/envconf"
)

// struct with `env` tagged literal
type Config struct {
	DBHost string `env:"DB_HOST"`
	DBPort int    `env:"DB_PORT"`
	DBUser string `env:"DB_USER"`
}

func main() {
	var conf Config

	// load environment variables and dotenv variables with `Load` method
	if err := envconf.Load(&conf, envconf.UseDotEnv()); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", conf)
}
