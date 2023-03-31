package services

import (
	"log"
	"os"
	"strconv"
)

type Configs struct {
	DatabaseUri   string
	RedisUri      string
	RedisPassword string
	AppEnv        string
	AppFork       bool
	BindAddr      string
	UseTls        bool
	TlsCertPath   string
	TlsKeyPath    string
	RedisDb       int
	JwtSecret     string
	BcryptSalt    int
}

var (
	instance        *Configs = nil
	requiredEnvVars []string = []string{"REDIS_URI", "DATABASE_URI", "JWT_SECRET"}
)

func GenerateConfigsFromEnv() {
	for _, v := range requiredEnvVars {
		if _, ok := os.LookupEnv(v); !ok {
			log.Panicf("environment: variable %s is not set", v)
		}
	}
	instance = &Configs{}

	instance.DatabaseUri = os.Getenv("DATABASE_URI")
	instance.RedisUri = os.Getenv("REDIS_URI")
	instance.RedisPassword = os.Getenv("REDIS_PASSWORD")
	instance.JwtSecret = os.Getenv("JWT_SECRET")

	if os.Getenv("BCRYPT_SALT") == "" {
		instance.BcryptSalt = 12
	} else {
		v, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
		if err != nil {
			log.Panicf("environment: variable BCRYPT_SALT must be an integer")
		}
		instance.BcryptSalt = v
	}

	if os.Getenv("APP_ENV") == "" {
		instance.AppEnv = "development"
	} else {
		instance.AppEnv = os.Getenv("APP_ENV")
	}

	if os.Getenv("APP_FORK") == "" {
		instance.AppFork = false
	} else {
		instance.AppFork = os.Getenv("APP_FORK") == "true"
	}

	if os.Getenv("BIND_ADDR") == "" {
		instance.BindAddr = ":3333"
	} else {
		instance.BindAddr = os.Getenv("BIND_ADDR")
	}

	if os.Getenv("REDIS_DB") == "" {
		instance.RedisDb = 0
	} else {
		v, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			log.Panicf("environment: variable REDIS_DB must be an integer")
		}
		instance.RedisDb = v
	}

	if os.Getenv("USE_TLS") == "true" {
		if os.Getenv("TLS_CERT_PATH") == "" || os.Getenv("TLS_KEY_PATH") == "" {
			log.Panicln("environment: TLS_CERT_PATH and TLS_KEY_PATH must be set if USE_TLS is true")
		}
		instance.TlsCertPath = os.Getenv("TLS_CERT_PATH")
		instance.TlsKeyPath = os.Getenv("TLS_KEY_PATH")
		instance.UseTls = true
	} else {
		instance.UseTls = false
	}
}

func ConfigProvider() *Configs {
	if instance == nil {
		instance = &Configs{}
	}

	return instance
}
