package config

import (
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

type Config struct {
	conf.Version
	Build string
	Web   struct {
		APIHost         string        `conf:"default:0.0.0.0:4000"`
		DebugHost       string        `conf:"default:0.0.0.0:5000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
		ShutdownTimeout time.Duration `conf:"default:20s"`
	}
	Auth struct {
		Algorithm  string `conf:"default:HS256"`
		SigningKey string `conf:"env:JWT_SIGNING_KEY"`
	}

	DB struct {
		User         string `conf:"default:postgres,env:DB_USER"`
		Password     string `conf:"default:postgres,mask,env:SQL_PWD"`
		Host         string `conf:"default:db,env:SQL_HOST"`
		DBName       string `conf:"default:bookmarks,env:SQL_DB_NAME"`
		MaxIdleConns int    `conf:"default:2"`
		MaxOpenConns int    `conf:"default:0"`
		DisableTLS   bool   `conf:"default:true"`
	}
}

func Initialize() (Config, error) {
	cfg := Config{}

	cfg.Build = getEnvOrDefault("API_BUILD", "dev")

	cfg.Version = conf.Version{
		SVN:  cfg.Build,
		Desc: "code created by Marksthought",
	}

	//setup WEB
	cfg.Web.APIHost = getEnvOrDefault("", "0.0.0.0:4000")
	cfg.Web.DebugHost = getEnvOrDefault("", "0.0.0.0:5000")
	cfg.Web.ReadTimeout = time.Second * 5
	cfg.Web.WriteTimeout = time.Second * 10
	cfg.Web.IdleTimeout = time.Second * 120
	cfg.Web.ShutdownTimeout = time.Second * 20

	//setup auth
	cfg.Auth.Algorithm = getEnvOrDefault("", "HS256")
	cfg.Auth.SigningKey = getEnvOrDefault("JWT_SIGNING_KEY", "")

	//setup db
	cfg.DB.User = getEnvOrDefault("DB_USER", "postgres")
	cfg.DB.Password = getEnvOrDefault("SQL_PWD", "pg")
	cfg.DB.Host = getEnvOrDefault("SQL_HOST", "empty")
	cfg.DB.DBName = getEnvOrDefault("SQL_DB_NAME", "bookmarks")

	return cfg, nil
}

func getEnvOrDefault(envVar string, defEnvVar string) (newEnvVar string) {
	if newEnvVar = os.Getenv(envVar); len(newEnvVar) == 0 {
		return defEnvVar
	} else {
		return newEnvVar
	}
}
