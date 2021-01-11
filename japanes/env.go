package main

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/subosito/gotenv"
)

// MyEnv 環境變數 可從系統中或.env檔案取得
type MyEnv struct {
	HTTPListenHost string `env:"HTTP_LISTEN_HOST" json:"HTTP_LISTEN_HOST" envDefault:""`
	HTTPListenPort int    `env:"HTTP_LISTEN_PORT" json:"HTTP_LISTEN_PORT" envDefault:"80"`

	MysqlReadHost     string `env:"MYSQL_READ_HOST" json:"MYSQL_READ_HOST"`
	MysqlReadPort     int    `env:"MYSQL_READ_PORT" json:"MYSQL_READ_PORT"`
	MysqlReadUser     string `env:"MYSQL_READ_USER" json:"MYSQL_READ_USER"`
	MysqlReadPassword string `env:"MYSQL_READ_PASSWORD" json:"MYSQL_READ_PASSWORD"`
	MysqlReadDatabase string `env:"MYSQL_READ_DATABASE" json:"MYSQL_READ_DATABASE"`

	MysqlWriteHost     string `env:"MYSQL_WRITE_HOST" json:"MYSQL_WRITE_HOST"`
	MysqlWritePort     int    `env:"MYSQL_WRITE_PORT" json:"MYSQL_WRITE_PORT"`
	MysqlWriteUser     string `env:"MYSQL_WRITE_USER" json:"MYSQL_WRITE_USER"`
	MysqlWritePassword string `env:"MYSQL_WRITE_PASSWORD" json:"MYSQL_WRITE_PASSWORD"`
	MysqlWriteDatabase string `env:"MYSQL_WRITE_DATABASE" json:"MYSQL_WRITE_DATABASE"`

	RedisHost     string `env:"REDIS_HOST" json:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT" json:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD" json:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" json:"REDIS_DB"`
}

// Init Init
func (e *MyEnv) Init(envFilePath string) error {
	// 若有使用 .env.xxx 設定檔案，使用 gotenv 將所有設定載入後，服務就可以從os.getenv()取得設定資料
	if envFilePath != "" {
		envLoadErr := gotenv.Load(envFilePath)
		if envLoadErr != nil {
			return envLoadErr
		}
	} else {
		fmt.Println("is null")
	}

	// 透過 caarlos0/env 將 os.getenv() 可取得的設定資料轉入 Env struct
	err := env.Parse(e)
	if err != nil {
		return err
	}

	if e.HTTPListenHost == "" {
		e.HTTPListenHost = "0.0.0.0"
	}

	return nil
}
