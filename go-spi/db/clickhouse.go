package db

import (
	"fmt"
	"go-spi/common"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	dsn   string
	dbPtr *sqlx.DB
)

const (
	db_type = "clickhouse"
)

type Config struct {
	ClickHouse `ini:"clickhouse"`
}

type ClickHouse struct {
	Address  string `ini:"address"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

func init() {

	configPath := common.RootPath + "/conf/clickhouse.ini"
	config := &Config{}
	err := ini.MapTo(config, configPath)
	if err != nil {
		logrus.Fatalf("load clickhouse.ini failed, err: %v", err.Error())
	}
	dsn = fmt.Sprintf("%s://%s:%s@%s/%s", db_type, config.User, config.Password, config.Address, config.Database)
	fmt.Println("dsn连接: ", dsn)
	dbPtr, err = sqlx.Open(db_type, dsn)
	checkErr(err)
}

// GetDB
// 获取数据库连接实例
func GetDB() *sqlx.DB {
	check_sql_db_isopen()
	return dbPtr
}

// check_sql_db_isopen
// 检查数据库连接是否关闭，如果关闭了就重新打开
func check_sql_db_isopen() {
	err := dbPtr.Ping()
	if err != nil {
		dbPtr.Close()
		dbPtr, err = sqlx.Open(db_type, dsn)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		logrus.Println(err.Error())
	}
}
