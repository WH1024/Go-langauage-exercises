package conf

import (
	"github.com/go-ini/ini"
	"os"
)

type Config struct {
	ClickConfig `ini:"clickhouse"`
}
type ClickConfig struct {
	Host             string `ini:"host"`
	Port             string `ini:"port"`
	Database         string `ini:"database"`
	DialTimeout      string `ini:"dial_timeout"`
	MaxExecutionTime string `ini:"max_execution_time"`
}

func GetConfig() (*Config, error) {
	dir, _ := os.Getwd()
	var config = new(Config)
	err := ini.MapTo(config, dir+"/conf/clickhouse.ini")
	if err != nil {
		return nil, err
	}
	return config, nil
}
