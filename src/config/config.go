package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	_APP_PORT   int
	_APP_DOMAIN string

	_DB_HOST string
	_DB_PORT int
	_DB_NAME string
	_DB_USER string
	_DB_PASS string
	_DB_TIME string
}

var config Config

func New() (*Config, error) {
	if config == (Config{}) {
		// APPLICATION

		appPortStr, exists := os.LookupEnv("APP_PORT")
		if exists != true {
			return nil, errors.New("APP_PORT enviroment variable is not set")
		}
		APP_PORT, err := strconv.Atoi(appPortStr)
		if err != nil {
			return nil, errors.New("incorrect APP_PORT value")
		}
		config._APP_PORT = APP_PORT

		APP_DOMAIN, exists := os.LookupEnv("APP_DOMAIN")
		if exists != true {
			return nil, errors.New("APP_DOMAIN enviroment variable is not set")
		}
		_, err = url.ParseRequestURI(APP_DOMAIN)
		if err != nil {
			return nil, errors.New("APP_DOMAIN enviroment variable is not a valid URL")
		}
		config._APP_DOMAIN = APP_DOMAIN

		// DATABASE
		DB_HOST, exists := os.LookupEnv("DB_HOST")
		if exists != true {
			return nil, errors.New("DB_HOST enviroment variable is not set")
		}
		config._DB_HOST = DB_HOST

		dbPortStr, exists := os.LookupEnv("DB_PORT")
		if exists != true {
			return nil, errors.New("DB_PORT enviroment variable is not set")
		}
		DB_PORT, err := strconv.Atoi(dbPortStr)
		if err != nil {
			return nil, errors.New("incorrect DB_PORT value")
		}
		config._DB_PORT = DB_PORT

		DB_NAME, exists := os.LookupEnv("DB_NAME")
		if exists != true {
			return nil, errors.New("DB_NAME enviroment variable is not set")
		}
		config._DB_NAME = DB_NAME

		DB_USER, exists := os.LookupEnv("DB_USER")
		if exists != true {
			return nil, errors.New("DB_USER enviroment variable is not set")
		}
		config._DB_USER = DB_USER

		DB_PASS, exists := os.LookupEnv("DB_PASS")
		if exists != true {
			return nil, errors.New("DB_PASS enviroment variable is not set")
		}
		config._DB_PASS = DB_PASS

		DB_TIME, exists := os.LookupEnv("DB_TIME")
		if exists != true {
			return nil, errors.New("DB_TIME enviroment variable is not set")
		}
		config._DB_TIME = DB_TIME
	}
	fmt.Println("enviroment variables are ok...")

	return &config, nil
}

func (c *Config) APP_DOMAIN() string {
	return c._APP_DOMAIN
}

func (c *Config) GetDbUrl() string {
	return c._DB_USER + ":" + c._DB_PASS + "@tcp(" + c._DB_HOST + ":" + strconv.Itoa(c._DB_PORT) + ")/" + c._DB_NAME
}

func (c *Config) GetServingAddress() string {
	return ":" + strconv.Itoa(c._APP_PORT)
}
